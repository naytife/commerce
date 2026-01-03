package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/petrejonn/naytife/internal/api"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/observability"
	"go.uber.org/zap"
)

type TemplateHandler struct {
	repository  db.Repository
	RetryClient *retryablehttp.Client
}

func NewTemplateHandler(repo db.Repository) *TemplateHandler {
	return &TemplateHandler{
		repository: repo,
	}
}

// Template management endpoints

// @Summary      List available templates
// @Description  Get all available templates from the template system (internal method, not routed)
// @Tags         templates-internal
// @Produce      json
// @Success      200  {object}  models.SuccessResponse{data=[]models.Template}
// @Failure      500  {object}  models.ErrorResponse
func (h *TemplateHandler) ListTemplates(c *fiber.Ctx) error {
	// Pass the incoming request context into the helper to preserve cancellation and tracing.
	templates, err := h.fetchTemplatesFromService(c.Context())
	if err != nil {
		zap.L().Error("ListTemplates: failed to fetch templates", zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch templates", nil)
	}

	return api.SuccessResponse(c, fiber.StatusOK, templates, "Templates fetched successfully")
}

// @Summary      Get template versions
// @Description  Get all versions for a specific template
// @Tags         templates
// @Produce      json
// @Param        template_name path string true "Template name"
// @Success      200  {object}  models.SuccessResponse{data=[]models.TemplateVersion}
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Security     OAuth2AccessCode
// @Router       /templates/{template_name}/versions [get]
func (h *TemplateHandler) GetTemplateVersions(c *fiber.Ctx) error {
	templateName := c.Params("template_name")
	if templateName == "" {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Template name is required", nil)
	}

	// Pass the incoming request context into the helper to preserve cancellation and tracing.
	versions, err := h.fetchTemplateVersionsFromService(c.Context(), templateName)
	if err != nil {
		zap.L().Error("GetTemplateVersions: failed to fetch template versions", zap.String("template", templateName), zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch template versions", nil)
	}

	return api.SuccessResponse(c, fiber.StatusOK, versions, "Template versions fetched successfully")
}

// @Summary      Get latest template version
// @Description  Get the latest version for a specific template
// @Tags         templates
// @Produce      json
// @Param        template_name path string true "Template name"
// @Success      200  {object}  models.SuccessResponse{data=models.TemplateVersion}
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Security     OAuth2AccessCode
// @Router       /templates/{template_name}/latest [get]
func (h *TemplateHandler) GetLatestTemplateVersion(c *fiber.Ctx) error {
	templateName := c.Params("template_name")
	if templateName == "" {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Template name is required", nil)
	}

	// Pass the incoming request context into the helper to preserve cancellation and tracing.
	latest, err := h.fetchLatestTemplateVersionFromService(c.Context(), templateName)
	if err != nil {
		zap.L().Error("GetLatestTemplateVersion: failed to fetch latest template version", zap.String("template", templateName), zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch latest template version", nil)
	}

	return api.SuccessResponse(c, fiber.StatusOK, latest, "Latest template version fetched successfully")
}

// @Summary      Build template
// @Description  Trigger a template build
// @Tags         templates
// @Accept       json
// @Produce      json
// @Param        request body models.TemplateBuildRequest true "Build request"
// @Success      202  {object}  models.SuccessResponse{data=models.BuildResponse}
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Security     OAuth2AccessCode
// @Router       /templates/build [post]
func (h *TemplateHandler) BuildTemplate(c *fiber.Ctx) error {
	var req models.TemplateBuildRequest
	if err := c.BodyParser(&req); err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	if req.TemplateName == "" {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Template name is required", nil)
	}

	// Pass the incoming request context into the helper to preserve cancellation and tracing.
	response, err := h.triggerTemplateBuild(c.Context(), req)
	if err != nil {
		zap.L().Error("BuildTemplate: failed to trigger template build", zap.String("template", req.TemplateName), zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to trigger template build", nil)
	}

	return api.SuccessResponse(c, fiber.StatusAccepted, response, "Template build initiated successfully")
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

	// Trigger async redeployment with new template version
	deploymentResp, err := h.triggerStoreRedeployment(c.Context(), shop.Subdomain, currentTemplate.TemplateName, latestVersion.Version)
	if err != nil {
		zap.L().Error("UpdateToLatestTemplate: failed to trigger redeployment", zap.Int64("shop_id", shopID), zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to trigger template update", nil)
	}

	response := models.TemplateUpdateResponse{
		ShopID:           fmt.Sprintf("%d", shopID),
		CurrentVersion:   currentTemplate.TemplateVersion,
		TargetVersion:    latestVersion.Version,
		IsUpdateRequired: true,
		DeploymentID:     deploymentResp.DeploymentID,
		Status:           deploymentResp.Status,
		Message:          "Template update initiated",
	}

	return api.SuccessResponse(c, fiber.StatusAccepted, response, "Template update initiated successfully")
}

// Service integration methods

func (h *TemplateHandler) fetchTemplatesFromService(ctx context.Context) ([]models.Template, error) {
	serviceURL := getServiceURL("store-deployer", "8001")
	// TODO: This helper should accept a caller-provided ctx so cancellation and tracing propagate.
	// If caller passed a background/TODO context, we still create a short timeout as a safeguard.
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	ctx, finish := observability.StartSpan(ctx, "fetchTemplatesFromService", "store-deployer", http.MethodGet, fmt.Sprintf("%s/templates", serviceURL))
	defer finish(0, nil)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/templates", serviceURL), nil)
	if err != nil {
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
		return nil, fmt.Errorf("failed to connect to template service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("template service returned status %d", resp.StatusCode)
	}

	var result struct {
		Templates []models.Template `json:"templates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode template response: %w", err)
	}

	return result.Templates, nil
}

// Template version and build service methods

func (h *TemplateHandler) fetchTemplateVersionsFromService(ctx context.Context, templateName string) ([]models.TemplateVersion, error) {
	serviceURL := getServiceURL("template-registry", "8002")
	// TODO: Accept caller ctx to preserve cancellation/tracing.
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	ctx, finish := observability.StartSpan(ctx, "fetchTemplateVersionsFromService", "template-registry", http.MethodGet, fmt.Sprintf("%s/versions/%s", serviceURL, templateName))
	defer finish(0, nil)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/versions/%s", serviceURL, templateName), nil)
	if err != nil {
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
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("service returned status %d", resp.StatusCode)
	}

	var result struct {
		Versions []models.TemplateVersion `json:"versions"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Versions, nil
}

func (h *TemplateHandler) fetchLatestTemplateVersionFromService(ctx context.Context, templateName string) (*models.TemplateVersion, error) {
	serviceURL := getServiceURL("template-registry", "8002")
	// TODO: Accept caller ctx to preserve cancellation/tracing.
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	ctx, finish := observability.StartSpan(ctx, "fetchLatestTemplateVersionFromService", "template-registry", http.MethodGet, fmt.Sprintf("%s/latest/%s", serviceURL, templateName))
	defer finish(0, nil)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/latest/%s", serviceURL, templateName), nil)
	if err != nil {
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
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("service returned status %d", resp.StatusCode)
	}

	var result models.TemplateVersion
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (h *TemplateHandler) triggerTemplateBuild(ctx context.Context, req models.TemplateBuildRequest) (*models.BuildResponse, error) {
	serviceURL := getServiceURL("template-registry", "8002")

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	// TODO: Accept caller ctx to preserve cancellation/tracing.
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	ctx, finish := observability.StartSpan(ctx, "triggerTemplateBuild", "template-registry", http.MethodPost, fmt.Sprintf("%s/build", serviceURL))
	defer finish(0, nil)
	reqHttp, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/build", serviceURL), jsonPayload(payload))
	if err != nil {
		return nil, err
	}
	reqHttp.Header.Set("Content-Type", "application/json")
	observability.InjectTraceHeaders(ctx, reqHttp)
	observability.EnsureRequestID(reqHttp)

	var resp *http.Response
	if h.RetryClient != nil {
		resp, err = h.RetryClient.StandardClient().Do(reqHttp)
	} else {
		resp, err = http.DefaultClient.Do(reqHttp)
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("service returned status %d", resp.StatusCode)
	}

	var result models.BuildResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// triggerStoreDeployment removed: store deployment is handled by proxy handlers (store-deployer service).

// NOTE: Data update functionality is now handled by proxy handlers which proxy to store-deployer service

func (h *TemplateHandler) fetchDeploymentStatusFromService(ctx context.Context, subdomain string) (*models.DeploymentStatus, error) {
	serviceURL := getServiceURL("store-deployer", "8001")
	// TODO: Accept caller ctx to preserve cancellation/tracing.
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	ctx, finish := observability.StartSpan(ctx, "fetchDeploymentStatusFromService", "store-deployer", http.MethodGet, fmt.Sprintf("%s/status/%s", serviceURL, subdomain))
	defer finish(0, nil)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/status/%s", serviceURL, subdomain), nil)
	if err != nil {
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
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("service returned status %d", resp.StatusCode)
	}

	var result models.DeploymentStatus
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// triggerStoreRedeployment triggers an async redeployment of a store with a new template version
func (h *TemplateHandler) triggerStoreRedeployment(ctx context.Context, subdomain, templateName, templateVersion string) (*models.DeploymentResponse, error) {
	serviceURL := getServiceURL("store-deployer", "8001")

	redeployReq := models.StoreRedeploymentRequest{
		Subdomain:    subdomain,
		TemplateName: templateName,
		Version:      templateVersion,
	}

	payload, err := json.Marshal(redeployReq)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	ctx, finish := observability.StartSpan(ctx, "triggerStoreRedeployment", "store-deployer", http.MethodPost, fmt.Sprintf("%s/redeploy", serviceURL))
	defer finish(0, nil)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/redeploy", serviceURL), jsonPayload(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	observability.InjectTraceHeaders(ctx, req)
	observability.EnsureRequestID(req)

	var resp *http.Response
	if h.RetryClient != nil {
		resp, err = h.RetryClient.StandardClient().Do(req)
	} else {
		resp, err = http.DefaultClient.Do(req)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to store-deployer service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("store-deployer service returned status %d", resp.StatusCode)
	}

	var result models.DeploymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode deployment response: %w", err)
	}

	return &result, nil
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
