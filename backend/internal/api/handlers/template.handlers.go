package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"
)

type TemplateHandler struct {
	repository db.Repository
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
	templates, err := h.fetchTemplatesFromService()
	if err != nil {
		log.Printf("Failed to fetch templates: %v", err)
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

	versions, err := h.fetchTemplateVersionsFromService(templateName)
	if err != nil {
		log.Printf("Failed to fetch template versions: %v", err)
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

	latest, err := h.fetchLatestTemplateVersionFromService(templateName)
	if err != nil {
		log.Printf("Failed to fetch latest template version: %v", err)
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

	response, err := h.triggerTemplateBuild(req)
	if err != nil {
		log.Printf("Failed to trigger template build: %v", err)
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
		return api.ErrorResponse(c, fiber.StatusNotFound, "Shop not found", nil)
	}

	status, err := h.fetchDeploymentStatusFromService(shop.Subdomain)
	if err != nil {
		log.Printf("Failed to fetch deployment status: %v", err)
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch deployment status", nil)
	}

	return api.SuccessResponse(c, fiber.StatusOK, status, "Deployment status fetched successfully")
}

// Service integration methods

func (h *TemplateHandler) fetchTemplatesFromService() ([]models.Template, error) {
	serviceURL := getServiceURL("store-deployer", "9003")
	resp, err := http.Get(fmt.Sprintf("%s/templates", serviceURL))
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

func (h *TemplateHandler) fetchTemplateVersionsFromService(templateName string) ([]models.TemplateVersion, error) {
	serviceURL := getServiceURL("template-registry", "9001")
	resp, err := http.Get(fmt.Sprintf("%s/versions/%s", serviceURL, templateName))
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

func (h *TemplateHandler) fetchLatestTemplateVersionFromService(templateName string) (*models.TemplateVersion, error) {
	serviceURL := getServiceURL("template-registry", "9001")
	resp, err := http.Get(fmt.Sprintf("%s/latest/%s", serviceURL, templateName))
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

func (h *TemplateHandler) triggerTemplateBuild(req models.TemplateBuildRequest) (*models.BuildResponse, error) {
	serviceURL := getServiceURL("template-registry", "9001")

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/build", serviceURL),
		"application/json",
		jsonPayload(payload),
	)
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

func (h *TemplateHandler) triggerStoreDeployment(req models.StoreDeploymentRequest) (*models.DeploymentResponse, error) {
	serviceURL := getServiceURL("store-deployer", "9003")

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/deploy", serviceURL),
		"application/json",
		jsonPayload(payload),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("service returned status %d", resp.StatusCode)
	}

	var result models.DeploymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// NOTE: Data update functionality is now handled by proxy handlers which proxy to store-deployer service

func (h *TemplateHandler) fetchDeploymentStatusFromService(subdomain string) (*models.DeploymentStatus, error) {
	serviceURL := getServiceURL("store-deployer", "9003")
	resp, err := http.Get(fmt.Sprintf("%s/status/%s", serviceURL, subdomain))
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
	"template-registry": {"TEMPLATE_REGISTRY_URL", "9001"},
	"store-deployer":    {"STORE_DEPLOYER_URL", "9003"},
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
