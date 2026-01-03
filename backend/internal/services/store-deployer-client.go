package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/petrejonn/naytife/internal/observability"
	"go.uber.org/zap"
)

type StoreDeployerClient struct {
	BaseURL    string
	HTTPClient *retryablehttp.Client
}

func NewStoreDeployerClient(client *retryablehttp.Client) *StoreDeployerClient {
	baseURL := os.Getenv("STORE_DEPLOYER_URL")
	if baseURL == "" {
		baseURL = "http://store-deployer:8001"
	}
	return &StoreDeployerClient{
		BaseURL:    baseURL,
		HTTPClient: client,
	}
}

func (s *StoreDeployerClient) UpdateData(ctx context.Context, subdomain string, shopID int64, dataType string) error {
	updateReq := map[string]interface{}{
		"shop_id":   fmt.Sprintf("%d", shopID),
		"data_type": dataType,
	}

	reqBody, err := json.Marshal(updateReq)
	if err != nil {
		return fmt.Errorf("marshal update request: %w", err)
	}

	url := fmt.Sprintf("%s/update-data/%s", s.BaseURL, subdomain)
	req, err := retryablehttp.NewRequest("POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("create update request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")

	observability.InjectTraceHeaders(ctx, req.Request)
	observability.EnsureRequestID(req.Request)

	start := time.Now()
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("call store-deployer: %w", err)
	}
	defer resp.Body.Close()

	observability.RecordServiceRequest("store-deployer", "POST", url, resp.StatusCode, time.Since(start))

	if resp.StatusCode != 200 {
		return fmt.Errorf("store-deployer returned status %d", resp.StatusCode)
	}

	zap.L().Info("store-deployer update succeeded",
		zap.Int64("shop_id", shopID),
		zap.String("subdomain", subdomain),
		zap.String("data_type", dataType))

	return nil
}

func (s *StoreDeployerClient) Cleanup(ctx context.Context, subdomain string, shopID int64) error {
	cleanupReq := map[string]interface{}{
		"shop_id": fmt.Sprintf("%d", shopID),
	}
	reqBody, err := json.Marshal(cleanupReq)
	if err != nil {
		return fmt.Errorf("failed to marshal cleanup request: %w", err)
	}

	url := fmt.Sprintf("%s/cleanup/%s", s.BaseURL, subdomain)
	req, err := retryablehttp.NewRequest("DELETE", url, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create cleanup request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")

	observability.InjectTraceHeaders(ctx, req.Request)
	observability.EnsureRequestID(req.Request)

	start := time.Now()
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call store-deployer cleanup: %w", err)
	}
	defer resp.Body.Close()

	observability.RecordServiceRequest("store-deployer", "DELETE", url, resp.StatusCode, time.Since(start))

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("store-deployer cleanup failed: %d %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *StoreDeployerClient) Deploy(ctx context.Context, shopID int64, subdomain, templateName string) error {
	deploymentReq := map[string]interface{}{
		"shop_id":       fmt.Sprintf("%d", shopID),
		"subdomain":     subdomain,
		"template_name": templateName,
		"version":       "",
		"data_override": map[string]string{},
	}
	reqBody, err := json.Marshal(deploymentReq)
	if err != nil {
		return fmt.Errorf("failed to marshal deployment request: %w", err)
	}

	url := fmt.Sprintf("%s/deploy", c.BaseURL)
	req, err := retryablehttp.NewRequest("POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create deployment request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")

	observability.InjectTraceHeaders(ctx, req.Request)
	observability.EnsureRequestID(req.Request)

	start := time.Now()
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call store-deployer deploy: %w", err)
	}
	defer resp.Body.Close()

	observability.RecordServiceRequest("store-deployer", "POST", url, resp.StatusCode, time.Since(start))

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("store-deployer deploy failed: %d %s", resp.StatusCode, string(body))
	}

	return nil
}

// NotifyDeploymentComplete notifies the backend when deployment completes
// Implements 3-attempt retry logic without exponential backoff
func (c *StoreDeployerClient) NotifyDeploymentComplete(ctx context.Context, shopID int64, deploymentID int64, subdomain, status, message string) error {
	const maxAttempts = 3
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// Skip backoff on first attempt
		if attempt > 1 {
			time.Sleep(time.Duration(attempt-1) * time.Second)
		}

		payload := map[string]interface{}{
			"deployment_id": deploymentID,
			"shop_id":       shopID,
			"subdomain":     subdomain,
			"status":        status,
			"message":       message,
			"completed_at":  time.Now().UTC().Format(time.RFC3339),
		}

		reqBody, err := json.Marshal(payload)
		if err != nil {
			zap.L().Warn("NotifyDeploymentComplete: failed to marshal notification payload",
				zap.Int64("deployment_id", deploymentID),
				zap.Error(err))
			return fmt.Errorf("marshal notification: %w", err)
		}

		url := fmt.Sprintf("%s/api/v1/internal/deployments/%d/complete", c.BaseURL, shopID)
		req, err := retryablehttp.NewRequest("POST", url, bytes.NewReader(reqBody))
		if err != nil {
			zap.L().Warn("NotifyDeploymentComplete: failed to create notification request",
				zap.Int64("deployment_id", deploymentID),
				zap.Error(err))
			lastErr = fmt.Errorf("create notification request: %w", err)
			continue
		}

		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/json")
		observability.InjectTraceHeaders(ctx, req.Request)
		observability.EnsureRequestID(req.Request)

		start := time.Now()
		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			zap.L().Warn("NotifyDeploymentComplete: failed to call deployment completion webhook",
				zap.Int64("deployment_id", deploymentID),
				zap.Int("attempt", attempt),
				zap.Error(err))
			lastErr = fmt.Errorf("webhook call attempt %d: %w", attempt, err)
			continue
		}
		defer resp.Body.Close()

		observability.RecordServiceRequest("backend", "POST", url, resp.StatusCode, time.Since(start))

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted {
			zap.L().Info("NotifyDeploymentComplete: deployment completion notification succeeded",
				zap.Int64("deployment_id", deploymentID),
				zap.String("subdomain", subdomain),
				zap.String("status", status),
				zap.Int("attempt", attempt))
			return nil
		}

		body, _ := io.ReadAll(resp.Body)
		zap.L().Warn("NotifyDeploymentComplete: webhook returned non-success status",
			zap.Int64("deployment_id", deploymentID),
			zap.Int("status_code", resp.StatusCode),
			zap.String("response_body", string(body)),
			zap.Int("attempt", attempt))
		lastErr = fmt.Errorf("webhook returned %d: %s", resp.StatusCode, string(body))
	}

	zap.L().Error("NotifyDeploymentComplete: failed after 3 attempts",
		zap.Int64("deployment_id", deploymentID),
		zap.String("subdomain", subdomain),
		zap.Error(lastErr))
	return lastErr
}
