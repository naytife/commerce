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

// DeploymentStatusResponse represents the status response from store-deployer
type DeploymentStatusResponse struct {
	Status          string `json:"status"`
	Message         string `json:"message"`
	Subdomain       string `json:"subdomain"`
	TemplateName    string `json:"template_name"`
	TemplateVersion string `json:"template_version"`
}

// WaitForDeploymentCompletion polls the store-deployer service until deployment completes
// It uses exponential backoff with a maximum timeout of 60 seconds
// Returns nil on success, error if deployment fails or times out
func (c *StoreDeployerClient) WaitForDeploymentCompletion(ctx context.Context, subdomain string) error {
	maxWaitTime := 60 * time.Second
	initialInterval := 2 * time.Second
	maxInterval := 5 * time.Second

	url := fmt.Sprintf("%s/status/%s", c.BaseURL, subdomain)
	interval := initialInterval
	startTime := time.Now()

	for {
		// Check timeout
		if time.Since(startTime) > maxWaitTime {
			zap.L().Warn("WaitForDeploymentCompletion: timeout waiting for deployment",
				zap.String("subdomain", subdomain),
				zap.Duration("max_wait", maxWaitTime))
			return fmt.Errorf("deployment polling timeout after %s", maxWaitTime)
		}

		// Poll status
		req, err := retryablehttp.NewRequest("GET", url, nil)
		if err != nil {
			zap.L().Warn("WaitForDeploymentCompletion: failed to create status request",
				zap.String("subdomain", subdomain),
				zap.Error(err))
			// Don't fail immediately, retry
			time.Sleep(interval)
			continue
		}

		req = req.WithContext(ctx)
		observability.InjectTraceHeaders(ctx, req.Request)
		observability.EnsureRequestID(req.Request)

		start := time.Now()
		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			zap.L().Warn("WaitForDeploymentCompletion: failed to call store-deployer status",
				zap.String("subdomain", subdomain),
				zap.Error(err))
			// Service might be temporarily unavailable, retry
			time.Sleep(interval)
			if interval < maxInterval {
				interval = time.Duration(float64(interval) * 1.5)
				if interval > maxInterval {
					interval = maxInterval
				}
			}
			continue
		}

		observability.RecordServiceRequest("store-deployer", "GET", url, resp.StatusCode, time.Since(start))

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			zap.L().Warn("WaitForDeploymentCompletion: non-200 status from store-deployer",
				zap.String("subdomain", subdomain),
				zap.Int("status_code", resp.StatusCode))
			// Service error, retry
			time.Sleep(interval)
			if interval < maxInterval {
				interval = time.Duration(float64(interval) * 1.5)
				if interval > maxInterval {
					interval = maxInterval
				}
			}
			continue
		}

		// Parse response
		var statusResp DeploymentStatusResponse
		if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
			resp.Body.Close()
			zap.L().Warn("WaitForDeploymentCompletion: failed to decode status response",
				zap.String("subdomain", subdomain),
				zap.Error(err))
			// Bad response, retry
			time.Sleep(interval)
			continue
		}
		resp.Body.Close()

		// Check deployment status
		switch statusResp.Status {
		case "deployed", "success":
			zap.L().Info("WaitForDeploymentCompletion: deployment completed successfully",
				zap.String("subdomain", subdomain),
				zap.String("template", statusResp.TemplateName),
				zap.String("version", statusResp.TemplateVersion),
				zap.Duration("total_wait", time.Since(startTime)))
			return nil

		case "failed", "error":
			errMsg := fmt.Sprintf("deployment failed: %s", statusResp.Message)
			zap.L().Warn("WaitForDeploymentCompletion: deployment failed",
				zap.String("subdomain", subdomain),
				zap.String("message", statusResp.Message))
			return fmt.Errorf(errMsg)

		case "deploying", "pending":
			// Still deploying, wait and retry
			zap.L().Debug("WaitForDeploymentCompletion: deployment in progress",
				zap.String("subdomain", subdomain),
				zap.String("status", statusResp.Status),
				zap.Duration("elapsed", time.Since(startTime)))
			time.Sleep(interval)
			if interval < maxInterval {
				interval = time.Duration(float64(interval) * 1.5)
				if interval > maxInterval {
					interval = maxInterval
				}
			}
			continue

		default:
			zap.L().Warn("WaitForDeploymentCompletion: unknown deployment status",
				zap.String("subdomain", subdomain),
				zap.String("status", statusResp.Status))
			time.Sleep(interval)
			if interval < maxInterval {
				interval = time.Duration(float64(interval) * 1.5)
				if interval > maxInterval {
					interval = maxInterval
				}
			}
			continue
		}
	}
}
