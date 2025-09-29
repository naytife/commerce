package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
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

// @Summary      Get deployment status (internal method, not routed)
// @Description  Get the deployment status for a shop (internal method, actual endpoint is proxied)
// @Tags         deployments-internal
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Success      200  {object}  models.SuccessResponse{data=models.DeploymentStatus}
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
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

	// Pass the incoming request context into the helper to preserve cancellation and tracing.
	status, err := h.fetchDeploymentStatusFromService(c.Context(), shop.Subdomain)
	if err != nil {
		zap.L().Error("GetDeploymentStatus: failed to fetch deployment status", zap.Int64("shop_id", shopID), zap.Error(err))
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch deployment status", nil)
	}

	return api.SuccessResponse(c, fiber.StatusOK, status, "Deployment status fetched successfully")
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
