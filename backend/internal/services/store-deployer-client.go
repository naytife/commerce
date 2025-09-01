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
		baseURL = "http://store-deployer:9003"
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
