package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/petrejonn/naytife/internal/api"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/observability"
	"github.com/petrejonn/naytife/internal/services"
	"go.uber.org/zap"
)

type TemplateHandler struct {
	repository          db.Repository
	RetryClient         *retryablehttp.Client
	StoreDeployerClient *services.StoreDeployerClient
}

func NewTemplateHandler(repo db.Repository) *TemplateHandler {
	return &TemplateHandler{
		repository: repo,
	}
}

// NOTE: Store data update functionality is now handled by proxy handlers
// (see proxy.handlers.go ProxyUpdateStoreData) which proxy to store-deployer service

// @Summary      Get deployment status
// @Description  Get the current deployment status of a store
// @Tags         deployment
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Success      200  {object}  models.SuccessResponse{data=models.DeploymentStatus}
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /shops/{shop_id}/deployment-status [get]
func (h *TemplateHandler) GetDeploymentStatus(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id")
	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	// Get shop details
	shop, err := h.repository.GetShop(c.Context(), shopID)
	if err != nil {
		zap.L().Warn("GetDeploymentStatus: shop not found", zap.Int64("shop_id", shopID), zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusNotFound, "Shop not found", nil)
	}

	// Get latest deployment from database (webhook-updated source of truth)
	deployment, err := h.repository.GetLatestDeploymentByShop(c.Context(), shopID)
	if err != nil {
		zap.L().Warn("GetDeploymentStatus: no deployment found for shop", zap.Int64("shop_id", shopID), zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusNotFound, "Shop has no deployment", nil)
	}

	// Build response using actual database state
	response := models.DeploymentStatus{
		ShopID:          fmt.Sprintf("%d", shopID),
		Subdomain:       shop.Subdomain,
		Status:          deployment.Status,
		TemplateName:    deployment.TemplateName,
		TemplateVersion: deployment.TemplateVersion,
		DeploymentID:    fmt.Sprintf("%d", deployment.DeploymentID),
		Message:         "",
		ProductionURL:   fmt.Sprintf("https://%s.naytife.com", shop.Subdomain),
	}

	// Add timestamps if available
	if deployment.CompletedAt.Valid {
		response.LastDeployedAt = &deployment.CompletedAt.Time
	}
	if deployment.StartedAt.Valid {
		response.LastUpdateAt = &deployment.StartedAt.Time
	}

	// Add error message if deployment failed
	if deployment.Message != nil {
		response.Message = *deployment.Message
	}

	zap.L().Info("GetDeploymentStatus: retrieved deployment status",
		zap.Int64("shop_id", shopID),
		zap.Int64("deployment_id", deployment.DeploymentID),
		zap.String("status", deployment.Status))

	return api.SuccessResponse(c, fiber.StatusOK, response, "Deployment status retrieved successfully")
}

// @Summary      Get current template
// @Description  Get the current deployed template and version for a shop
// @Tags         templates
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Success      200  {object}  models.SuccessResponse{data=models.CurrentTemplate}
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/template/current [get]
func (h *TemplateHandler) GetCurrentTemplate(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id")
	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	// Verify shop exists
	_, err = h.repository.GetShop(c.Context(), shopID)
	if err != nil {
		zap.L().Warn("GetCurrentTemplate: shop not found", zap.Int64("shop_id", shopID), zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusNotFound, "Shop not found", nil)
	}

	// Get current deployed template
	templateInfo, err := h.repository.GetShopCurrentTemplate(c.Context(), shopID)
	if err != nil {
		zap.L().Warn("GetCurrentTemplate: no deployment found for shop", zap.Int64("shop_id", shopID), zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusNotFound, "Shop has no active deployment", nil)
	}

	// Convert pgtype.Timestamptz to *time.Time
	var deployedAt *time.Time
	if templateInfo.CompletedAt.Valid {
		deployedAt = &templateInfo.CompletedAt.Time
	}

	response := models.CurrentTemplate{
		ShopID:          fmt.Sprintf("%d", shopID),
		TemplateName:    templateInfo.TemplateName,
		TemplateVersion: templateInfo.TemplateVersion,
		DeployedAt:      deployedAt,
		Status:          templateInfo.Status,
	}

	return api.SuccessResponse(c, fiber.StatusOK, response, "Current template fetched successfully")
}

// @Summary      Update template to latest version
// @Description  Update a shop's template to the latest available version (prevents downgrades)
// @Tags         templates
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Success      202  {object}  models.SuccessResponse{data=models.TemplateUpdateResponse}
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      409  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/template/update-latest [post]
func (h *TemplateHandler) UpdateToLatestTemplate(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id")
	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	// Verify shop exists
	shop, err := h.repository.GetShop(c.Context(), shopID)
	if err != nil {
		zap.L().Warn("UpdateToLatestTemplate: shop not found", zap.Int64("shop_id", shopID), zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusNotFound, "Shop not found", nil)
	}

	// Get current deployed template
	currentTemplate, err := h.repository.GetShopCurrentTemplate(c.Context(), shopID)
	if err != nil {
		zap.L().Warn("UpdateToLatestTemplate: no deployment found for shop", zap.Int64("shop_id", shopID), zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusNotFound, "Shop has no active deployment", nil)
	}

	// Fetch latest template version from template-registry
	latestVersion, err := h.fetchLatestTemplateVersionFromService(c.Context(), currentTemplate.TemplateName)
	if err != nil {
		zap.L().Error("UpdateToLatestTemplate: failed to fetch latest template version", zap.Int64("shop_id", shopID), zap.String("template", currentTemplate.TemplateName), zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch latest template version", nil)
	}

	// Check if already on latest version
	if currentTemplate.TemplateVersion == latestVersion.Version {
		return api.SuccessResponse(c, fiber.StatusOK, models.TemplateUpdateResponse{
			ShopID:           fmt.Sprintf("%d", shopID),
			CurrentVersion:   currentTemplate.TemplateVersion,
			TargetVersion:    latestVersion.Version,
			IsUpdateRequired: false,
			Message:          "Already on latest version",
		}, "Shop is already using the latest template version")
	}

	// Compare versions to prevent downgrades (simple semantic versioning check)
	// Extract version numbers for comparison (e.g., "1.2.3" -> [1, 2, 3])
	currentParts := parseSemanticVersion(currentTemplate.TemplateVersion)
	latestParts := parseSemanticVersion(latestVersion.Version)

	if compareVersions(currentParts, latestParts) > 0 {
		// Current version is newer than latest available - this shouldn't happen but prevent downgrade
		zap.L().Warn("UpdateToLatestTemplate: attempted downgrade detected", zap.Int64("shop_id", shopID), zap.String("current", currentTemplate.TemplateVersion), zap.String("latest", latestVersion.Version))
		return api.ErrorResponse(c, fiber.StatusConflict, "Cannot downgrade template version", nil)
	}

	// Create deployment record (status = 'deploying')
	startedAt := pgtype.Timestamptz{Time: time.Now(), Valid: true}
	deployment, err := h.repository.CreateDeployment(c.Context(), db.CreateDeploymentParams{
		ShopID:          shopID,
		TemplateName:    currentTemplate.TemplateName,
		TemplateVersion: latestVersion.Version,
		Status:          "deploying",
		DeploymentType:  "template_update",
		Message:         nil,
		StartedAt:       startedAt,
	})
	if err != nil {
		zap.L().Error("UpdateToLatestTemplate: failed to create deployment record",
			zap.Int64("shop_id", shopID),
			zap.String("subdomain", shop.Subdomain),
			zap.String("template", currentTemplate.TemplateName),
			zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create deployment record", nil)
	}

	// Trigger async redeployment via store-deployer client
	go func(shopID int64, deploymentID int64, subdomain, templateName, templateVersion string) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		ctx, finish := observability.StartSpan(ctx, "asyncTemplateUpdate", "store-deployer", "POST", "redeploy")
		defer finish(0, nil)

		// Trigger deployment with store-deployer client
		if err := h.StoreDeployerClient.Deploy(ctx, shopID, deploymentID, subdomain, templateName); err != nil {
			errMsg := fmt.Sprintf("failed to trigger template update: %v", err)
			_ = h.repository.UpdateDeploymentStatus(ctx, db.UpdateDeploymentStatusParams{
				DeploymentID: deploymentID,
				Status:       "failed",
				Message:      &errMsg,
			})

			zap.L().Error("UpdateToLatestTemplate: failed to trigger store-deployer",
				zap.Int64("shop_id", shopID),
				zap.String("subdomain", subdomain),
				zap.String("template", templateName),
				zap.String("version", templateVersion),
				zap.Error(err))
			return
		}

		zap.L().Info("UpdateToLatestTemplate: template update initiated, waiting for store-deployer callback",
			zap.Int64("shop_id", shopID),
			zap.String("subdomain", subdomain),
			zap.String("template", templateName),
			zap.String("version", templateVersion))
	}(shopID, deployment.DeploymentID, shop.Subdomain, currentTemplate.TemplateName, latestVersion.Version)

	response := models.TemplateUpdateResponse{
		ShopID:           fmt.Sprintf("%d", shopID),
		CurrentVersion:   currentTemplate.TemplateVersion,
		TargetVersion:    latestVersion.Version,
		IsUpdateRequired: true,
		DeploymentID:     fmt.Sprintf("%d", deployment.DeploymentID),
		Status:           "deploying",
		Message:          "Template update initiated",
	}

	return api.SuccessResponse(c, fiber.StatusAccepted, response, "Template update initiated successfully")
}

func (h *TemplateHandler) fetchLatestTemplateVersionFromService(ctx context.Context, templateName string) (*models.TemplateVersion, error) {
	serviceURL := getServiceURL("template-registry", "8002")
	// Accept caller ctx to preserve cancellation/tracing - use context directly without creating new timeout
	ctx, finish := observability.StartSpan(ctx, "fetchLatestTemplateVersionFromService", "template-registry", http.MethodGet, fmt.Sprintf("%s/templates/%s/latest", serviceURL, templateName))
	defer finish(0, nil)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/templates/%s/latest", serviceURL, templateName), nil)
	if err != nil {
		zap.L().Error("failed to create request for template-registry", zap.Error(err), zap.String("template", templateName))
		return nil, err
	}
	observability.InjectTraceHeaders(ctx, req)
	observability.EnsureRequestID(req)

	var resp *http.Response
	if h.RetryClient != nil {
		resp, err = h.RetryClient.StandardClient().Do(req)
	} else {
		resp, err = http.DefaultClient.Do(req)
	}
	if err != nil {
		zap.L().Error("failed to fetch latest template version from service", zap.Error(err), zap.String("template", templateName), zap.String("url", fmt.Sprintf("%s/templates/%s/latest", serviceURL, templateName)))
		return nil, err
	}
	defer resp.Body.Close()

	// Log response status for debugging
	zap.L().Debug("template-registry response", zap.Int("status", resp.StatusCode), zap.String("template", templateName))

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		zap.L().Warn("template-registry returned non-OK status", zap.Int("status", resp.StatusCode), zap.String("template", templateName), zap.String("body", string(body)))
		return nil, fmt.Errorf("template-registry service returned status %d for template %s", resp.StatusCode, templateName)
	}

	// Decode into wrapper response struct (template-registry returns {status, version})
	var wrapper models.LatestTemplateVersionResponse
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		zap.L().Error("failed to decode template-registry response", zap.Error(err), zap.String("template", templateName))
		return nil, fmt.Errorf("failed to decode template-registry response: %w", err)
	}

	// Validate response structure
	if wrapper.Status != "success" {
		zap.L().Warn("unexpected status from template-registry", zap.String("status", wrapper.Status), zap.String("template", templateName))
		return nil, fmt.Errorf("template-registry returned status: %s", wrapper.Status)
	}

	return &wrapper.Version, nil
}

// Version comparison utilities for semantic versioning

// parseSemanticVersion converts a version string like "1.2.3" to [1, 2, 3]
func parseSemanticVersion(versionStr string) []int {
	parts := strings.Split(versionStr, ".")
	result := make([]int, 0, len(parts))

	for _, part := range parts {
		// Extract only the numeric part
		var numStr string
		for _, c := range part {
			if c >= '0' && c <= '9' {
				numStr += string(c)
			} else {
				break
			}
		}
		if numStr == "" {
			result = append(result, 0)
		} else {
			if v, err := strconv.Atoi(numStr); err == nil {
				result = append(result, v)
			} else {
				result = append(result, 0)
			}
		}
	}

	// Pad with zeros to ensure consistent length
	for len(result) < 3 {
		result = append(result, 0)
	}

	return result[:3]
}

// compareVersions compares two semantic version arrays
// Returns: positive if v1 > v2, 0 if equal, negative if v1 < v2
func compareVersions(v1, v2 []int) int {
	for i := 0; i < len(v1) && i < len(v2); i++ {
		if v1[i] > v2[i] {
			return 1
		} else if v1[i] < v2[i] {
			return -1
		}
	}
	return 0
}

// Service URL configuration
var serviceConfig = map[string]struct {
	envVar      string
	defaultPort string
}{
	"template-registry": {"TEMPLATE_REGISTRY_URL", "8002"},
	"store-deployer":    {"STORE_DEPLOYER_URL", "8001"},
}

func getServiceURL(serviceName, defaultPort string) string {
	// Try environment variable first
	if config, exists := serviceConfig[serviceName]; exists {
		if url := os.Getenv(config.envVar); url != "" {
			return url
		}
		// Use default port from config
		return fmt.Sprintf("http://%s:%s", serviceName, config.defaultPort)
	}

	// Fallback to provided default port
	return fmt.Sprintf("http://%s:%s", serviceName, defaultPort)
}

func jsonPayload(data []byte) *bytes.Reader {
	return bytes.NewReader(data)
}

// CompleteDeploymentCallback handles deployment completion notifications from store-deployer
// @Summary      Complete deployment callback
// @Description  Webhook endpoint called by store-deployer when deployment completes
// @Tags         deployments-internal
// @Accept       json
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        callback body object{deployment_id=int64,subdomain=string,status=string,message=string,completed_at=string} true "Deployment completion callback"
// @Success      200  {object}  models.SuccessResponse "Deployment status updated successfully"
// @Failure      400  {object}  models.ErrorResponse "Invalid request"
// @Failure      404  {object}  models.ErrorResponse "Deployment not found"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Router       /internal/deployments/{shop_id}/complete [post]
func (h *TemplateHandler) CompleteDeploymentCallback(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id")
	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
	if err != nil {
		zap.L().Warn("CompleteDeploymentCallback: invalid shop_id", zap.String("shop_id", shopIDStr))
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	// Parse callback payload
	var callback struct {
		DeploymentID int64  `json:"deployment_id"`
		Subdomain    string `json:"subdomain"`
		Status       string `json:"status"`
		Message      string `json:"message"`
		CompletedAt  string `json:"completed_at"`
	}

	if err := c.BodyParser(&callback); err != nil {
		zap.L().Warn("CompleteDeploymentCallback: failed to parse request body",
			zap.Int64("shop_id", shopID),
			zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	// Validate required fields
	if callback.DeploymentID == 0 || callback.Status == "" {
		zap.L().Warn("CompleteDeploymentCallback: missing required fields",
			zap.Int64("shop_id", shopID),
			zap.Int64("deployment_id", callback.DeploymentID),
			zap.String("status", callback.Status))
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Missing required fields", nil)
	}

	// Update deployment status in database
	var messagePtr *string
	if callback.Message != "" {
		messagePtr = &callback.Message
	}

	err = h.repository.UpdateDeploymentStatus(c.Context(), db.UpdateDeploymentStatusParams{
		DeploymentID: callback.DeploymentID,
		Status:       callback.Status,
		Message:      messagePtr,
	})
	if err != nil {
		zap.L().Error("CompleteDeploymentCallback: failed to update deployment status",
			zap.Int64("shop_id", shopID),
			zap.Int64("deployment_id", callback.DeploymentID),
			zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update deployment status", nil)
	}

	// If deployment succeeded, also update shop's last_deployment_id and complete it
	if callback.Status == "deployed" {
		err = h.repository.CompleteDeployment(c.Context(), db.CompleteDeploymentParams{
			DeploymentID: callback.DeploymentID,
			Status:       "deployed",
			Message:      messagePtr,
		})
		if err != nil {
			zap.L().Warn("CompleteDeploymentCallback: failed to mark deployment as completed",
				zap.Int64("shop_id", shopID),
				zap.Int64("deployment_id", callback.DeploymentID),
				zap.Error(err))
		}

		// Update shop's last_deployment_id
		err = h.repository.UpdateShopLastDeployment(c.Context(), db.UpdateShopLastDeploymentParams{
			ShopID:           shopID,
			LastDeploymentID: &callback.DeploymentID,
		})
		if err != nil {
			zap.L().Warn("CompleteDeploymentCallback: failed to update shop last deployment",
				zap.Int64("shop_id", shopID),
				zap.Int64("deployment_id", callback.DeploymentID),
				zap.Error(err))
		}
	}

	zap.L().Info("CompleteDeploymentCallback: deployment callback processed successfully",
		zap.Int64("shop_id", shopID),
		zap.Int64("deployment_id", callback.DeploymentID),
		zap.String("status", callback.Status),
		zap.String("subdomain", callback.Subdomain))

	return api.SuccessResponse(c, fiber.StatusOK, map[string]interface{}{
		"deployment_id": callback.DeploymentID,
		"status":        callback.Status,
	}, "Deployment status updated successfully")
}
