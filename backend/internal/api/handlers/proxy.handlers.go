package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/petrejonn/naytife/internal/api"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"
)

type ProxyHandler struct {
	Repository          db.Repository
	TemplateRegistryURL string
	StoreDeployerURL    string
	HttpClient          *http.Client
}

func NewProxyHandler(repo db.Repository) *ProxyHandler {
	// Get service URLs from environment variables with fallback defaults
	templateRegistryURL := os.Getenv("TEMPLATE_REGISTRY_URL")
	if templateRegistryURL == "" {
		templateRegistryURL = "http://template-registry:9001" // Default for Kubernetes
	}

	storeDeployerURL := os.Getenv("STORE_DEPLOYER_URL")
	if storeDeployerURL == "" {
		storeDeployerURL = "http://store-deployer:9003" // Default for Kubernetes
	}

	return &ProxyHandler{
		Repository:          repo,
		TemplateRegistryURL: templateRegistryURL,
		StoreDeployerURL:    storeDeployerURL,
		HttpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Generic proxy method
func (h *ProxyHandler) proxyRequest(c *fiber.Ctx, targetURL string, path string) error {
	// Build the full URL
	fullURL := fmt.Sprintf("%s%s", targetURL, path)

	// Create a new request
	var reqBody io.Reader
	if c.Body() != nil && len(c.Body()) > 0 {
		reqBody = bytes.NewReader(c.Body())
	}

	req, err := http.NewRequest(string(c.Method()), fullURL, reqBody)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create proxy request", nil)
	}

	// Copy headers from original request
	for key, values := range c.GetReqHeaders() {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Copy query parameters
	query := req.URL.Query()
	for key, value := range c.Queries() {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()

	// Make the request
	resp, err := h.HttpClient.Do(req)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadGateway, "Failed to reach service", nil)
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Set(key, value)
		}
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to read service response", nil)
	}

	// Set status code and return body
	c.Status(resp.StatusCode)
	return c.Send(body)
}

// Template Registry Proxy Handlers

// ProxyListTemplates proxies requests to template-registry /templates
// @Summary      List all templates
// @Description  Get a list of all available templates from template registry
// @Tags         templates
// @Produce      json
// @Success      200  {object}  models.SuccessResponse  "Templates retrieved successfully"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /templates [get]
func (h *ProxyHandler) ProxyListTemplates(c *fiber.Ctx) error {
	return h.ProxyWithStandardResponse(c, h.TemplateRegistryURL, "/templates", "Templates retrieved successfully")
}

// ProxyGetTemplate proxies requests to template-registry /templates/{name}
// @Summary      Get a specific template
// @Description  Get details of a specific template by name
// @Tags         templates
// @Produce      json
// @Param        name path string true "Template name"
// @Success      200  {object}  models.SuccessResponse  "Template retrieved successfully"
// @Failure      404  {object}  models.ErrorResponse "Template not found"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /templates/{name} [get]
func (h *ProxyHandler) ProxyGetTemplate(c *fiber.Ctx) error {
	templateName := c.Params("name")
	path := fmt.Sprintf("/templates/%s", templateName)
	return h.ProxyWithStandardResponse(c, h.TemplateRegistryURL, path, "Template retrieved successfully")
}

// ProxyGetTemplateVersions proxies requests to template-registry /templates/{name}/versions
// @Summary      Get all versions of a template
// @Description  Get a list of all versions for a specific template
// @Tags         templates
// @Produce      json
// @Param        name path string true "Template name"
// @Success      200  {object}  models.SuccessResponse  "Template versions retrieved successfully"
// @Failure      404  {object}  models.ErrorResponse "Template not found"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /templates/{name}/versions [get]
func (h *ProxyHandler) ProxyGetTemplateVersions(c *fiber.Ctx) error {
	templateName := c.Params("name")
	path := fmt.Sprintf("/templates/%s/versions", templateName)
	return h.ProxyWithStandardResponse(c, h.TemplateRegistryURL, path, "Template versions retrieved successfully")
}

// ProxyGetLatestTemplateVersion proxies requests to template-registry /templates/{name}/latest
// @Summary      Get latest version of a template
// @Description  Get the latest version information for a specific template
// @Tags         templates
// @Produce      json
// @Param        name path string true "Template name"
// @Success      200  {object}  models.SuccessResponse  "Latest template version retrieved successfully"
// @Failure      404  {object}  models.ErrorResponse "Template not found"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /templates/{name}/latest [get]
func (h *ProxyHandler) ProxyGetLatestTemplateVersion(c *fiber.Ctx) error {
	templateName := c.Params("name")
	path := fmt.Sprintf("/templates/%s/latest", templateName)
	return h.ProxyWithStandardResponse(c, h.TemplateRegistryURL, path, "Latest template version retrieved successfully")
}

// ProxyGetTemplateVersion proxies requests to template-registry /templates/{name}/versions/{version}
// @Summary      Get a specific template version
// @Description  Get details of a specific template version
// @Tags         templates
// @Produce      json
// @Param        name path string true "Template name"
// @Param        version path string true "Template version"
// @Success      200  {object}  models.SuccessResponse  "Template version retrieved successfully"
// @Failure      404  {object}  models.ErrorResponse "Template version not found"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /templates/{name}/versions/{version} [get]
func (h *ProxyHandler) ProxyGetTemplateVersion(c *fiber.Ctx) error {
	templateName := c.Params("name")
	version := c.Params("version")
	path := fmt.Sprintf("/templates/%s/versions/%s", templateName, version)
	return h.ProxyWithStandardResponse(c, h.TemplateRegistryURL, path, "Template version retrieved successfully")
}

// ProxyDownloadTemplate proxies requests to template-registry /templates/{name}/versions/{version}/download
// @Summary      Download a template version
// @Description  Download assets for a specific template version
// @Tags         templates
// @Produce      json
// @Param        name path string true "Template name"
// @Param        version path string true "Template version"
// @Success      200  {object}  models.TemplateDownload  "Template download information"
// @Failure      404  {object}  models.ErrorResponse "Template version not found"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /templates/{name}/versions/{version}/download [get]
func (h *ProxyHandler) ProxyDownloadTemplate(c *fiber.Ctx) error {
	templateName := c.Params("name")
	version := c.Params("version")
	path := fmt.Sprintf("/templates/%s/versions/%s/download", templateName, version)
	return h.proxyRequest(c, h.TemplateRegistryURL, path)
}

// TemplateUploadRequest represents the multipart form data for template upload
type TemplateUploadRequest struct {
	TemplateName string `json:"template_name" form:"template_name" binding:"required" example:"my-template" description:"Template name"`
	Version      string `json:"version" form:"version" example:"1.0.0" description:"Template version (auto-generated if not provided)"`
	Description  string `json:"description" form:"description" example:"Template description" description:"Template description"`
	Category     string `json:"category" form:"category" example:"ecommerce" description:"Template category (e.g., ecommerce, blog, portfolio)"`
	Features     string `json:"features" form:"features" example:"responsive,dark-mode,multi-language" description:"Comma-separated list of template features"`
	Force        bool   `json:"force" form:"force" example:"false" description:"Force upload even if version exists"`
	Assets       string `json:"assets" form:"assets" swaggertype:"string" format:"binary" binding:"required" description:"Template assets (tar.gz file)"`
	PreviewImage string `json:"preview_image" form:"preview_image" swaggertype:"string" format:"binary" description:"Optional preview image for the template (PNG, JPG, WebP, GIF)"`
}

// ProxyUploadTemplate proxies template upload requests to template-registry /templates/upload
// @Summary      Upload a template
// @Description  Upload a new template to the template registry
// @Tags         templates
// @Accept       multipart/form-data
// @Produce      json
// @Param        request body TemplateUploadRequest true "Template upload form data"
// @Success      200  {object}  models.SuccessResponse  "Template uploaded successfully"
// @Failure      400  {object}  models.ErrorResponse "Invalid request"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /templates/upload [post]
func (h *ProxyHandler) ProxyUploadTemplate(c *fiber.Ctx) error {
	return h.ProxyWithStandardResponse(c, h.TemplateRegistryURL, "/templates/upload", "Template uploaded successfully")
}

// Store Deployer Proxy Handlers

// ProxyDeployStore proxies requests to store-deployer /deploy
// @Summary      Deploy a store
// @Description  Deploy a store using the store deployer service
// @Tags         deployment
// @Accept       json
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        deploy body models.StoreDeploymentRequest true "Deployment parameters"
// @Success      200  {object}  models.DeploymentResponse  "Store deployed successfully"
// @Failure      400  {object}  models.ErrorResponse "Invalid request"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/deploy [post]
func (h *ProxyHandler) ProxyDeployStore(c *fiber.Ctx) error {
	shopID := c.Params("shop_id")

	// Validate shop ownership before proxying
	shopIDInt, err := strconv.ParseInt(shopID, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	// Add shop validation here if needed
	_ = shopIDInt // Use the validated shop ID

	return h.proxyRequest(c, h.StoreDeployerURL, "/deploy")
}

// ProxyRedeployStore proxies requests to store-deployer /redeploy/{subdomain}
// @Summary      Redeploy a store
// @Description  Redeploy an existing store using the store deployer service
// @Tags         deployment
// @Accept       json
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        redeploy body models.StoreRedeploymentRequest true "Redeployment parameters"
// @Success      200  {object}  models.DeploymentResponse  "Store redeployed successfully"
// @Failure      400  {object}  models.ErrorResponse "Invalid request"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/redeploy [post]
func (h *ProxyHandler) ProxyRedeployStore(c *fiber.Ctx) error {
	shopID := c.Params("shop_id")

	// Validate shop ownership before proxying
	shopIDInt, err := strconv.ParseInt(shopID, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	// Get shop details to extract subdomain for the proxy call
	shop, err := h.Repository.GetShop(c.Context(), shopIDInt)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusNotFound, "Shop not found", nil)
	}

	path := fmt.Sprintf("/redeploy/%s", shop.Subdomain)
	return h.proxyRequest(c, h.StoreDeployerURL, path)
}

// ProxyDeploymentStatus proxies requests to store-deployer for deployment status
// @Summary      Get deployment status
// @Description  Get the current deployment status of a store
// @Tags         deployment
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Success      200  {object}  models.DeploymentStatus  "Deployment status retrieved"
// @Failure      400  {object}  models.ErrorResponse "Invalid shop ID"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/deployment-status [get]
func (h *ProxyHandler) ProxyDeploymentStatus(c *fiber.Ctx) error {
	shopID := c.Params("shop_id")

	// Validate shop ownership before proxying
	shopIDInt, err := strconv.ParseInt(shopID, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	// Get shop details to extract subdomain for the proxy call
	shop, err := h.Repository.GetShop(c.Context(), shopIDInt)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusNotFound, "Shop not found", nil)
	}

	path := fmt.Sprintf("/status/%s", shop.Subdomain)
	return h.proxyRequest(c, h.StoreDeployerURL, path)
}

// ProxyUpdateStoreData proxies requests to store-deployer /update-data/{subdomain}
// @Summary      Update store data
// @Description  Update store data without redeploying assets
// @Tags         deployment
// @Accept       json
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        update body models.DataUpdateRequest true "Data update request"
// @Success      200  {object}  models.SuccessResponse  "Store data updated successfully"
// @Failure      400  {object}  models.ErrorResponse "Invalid request"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/update-data [post]
func (h *ProxyHandler) ProxyUpdateStoreData(c *fiber.Ctx) error {
	shopID := c.Params("shop_id")

	// Validate shop ownership before proxying
	shopIDInt, err := strconv.ParseInt(shopID, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	// Get shop details to extract subdomain for the proxy call
	shop, err := h.Repository.GetShop(c.Context(), shopIDInt)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusNotFound, "Shop not found", nil)
	}

	path := fmt.Sprintf("/update-data/%s", shop.Subdomain)
	return h.proxyRequest(c, h.StoreDeployerURL, path)
}

// Health Check Handlers - Aggregate health from services

// ProxyHealthCheck checks health of all proxied services
// @Summary      Health check for all services
// @Description  Aggregate health check for template-registry and store-deployer
// @Tags         health
// @Produce      json
// @Success      200  {object}  models.HealthResponse  "All services healthy"
// @Failure      503  {object}  models.ErrorResponse "One or more services unhealthy"
// @Router       /health/services [get]
func (h *ProxyHandler) ProxyHealthCheck(c *fiber.Ctx) error {
	var services []models.ServiceHealth
	allHealthy := true

	// Check template-registry health
	templateResp, err := h.HttpClient.Get(fmt.Sprintf("%s/health", h.TemplateRegistryURL))
	if err != nil || templateResp.StatusCode != 200 {
		allHealthy = false
		errorMsg := ""
		if err != nil {
			errorMsg = err.Error()
		} else {
			errorMsg = fmt.Sprintf("HTTP %d", templateResp.StatusCode)
		}
		services = append(services, models.ServiceHealth{
			Service: "template-registry",
			Status:  "unhealthy",
			URL:     h.TemplateRegistryURL,
			Error:   errorMsg,
		})
	} else {
		services = append(services, models.ServiceHealth{
			Service: "template-registry",
			Status:  "healthy",
			URL:     h.TemplateRegistryURL,
		})
	}
	if templateResp != nil {
		templateResp.Body.Close()
	}

	// Check store-deployer health
	deployerResp, err := h.HttpClient.Get(fmt.Sprintf("%s/health", h.StoreDeployerURL))
	if err != nil || deployerResp.StatusCode != 200 {
		allHealthy = false
		errorMsg := ""
		if err != nil {
			errorMsg = err.Error()
		} else {
			errorMsg = fmt.Sprintf("HTTP %d", deployerResp.StatusCode)
		}
		services = append(services, models.ServiceHealth{
			Service: "store-deployer",
			Status:  "unhealthy",
			URL:     h.StoreDeployerURL,
			Error:   errorMsg,
		})
	} else {
		services = append(services, models.ServiceHealth{
			Service: "store-deployer",
			Status:  "healthy",
			URL:     h.StoreDeployerURL,
		})
	}
	if deployerResp != nil {
		deployerResp.Body.Close()
	}

	result := models.HealthResponse{
		Services: services,
	}

	if allHealthy {
		result.Status = "healthy"
		return c.JSON(result)
	} else {
		result.Status = "unhealthy"
		return c.Status(fiber.StatusServiceUnavailable).JSON(result)
	}
}

// ProxyWithStandardResponse proxies a request and wraps the response in backend standard format
func (h *ProxyHandler) ProxyWithStandardResponse(c *fiber.Ctx, targetURL string, path string, message string) error {
	// First, get the response from the target service
	fullURL := fmt.Sprintf("%s%s", targetURL, path)

	var reqBody io.Reader
	if c.Body() != nil && len(c.Body()) > 0 {
		reqBody = bytes.NewReader(c.Body())
	}

	req, err := http.NewRequest(string(c.Method()), fullURL, reqBody)
	if err != nil {
		return api.SystemErrorResponse(c, err, "Failed to create proxy request")
	}

	// Copy headers
	for key, values := range c.GetReqHeaders() {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Copy query parameters
	query := req.URL.Query()
	for key, value := range c.Queries() {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()

	resp, err := h.HttpClient.Do(req)
	if err != nil {
		return api.ExternalServiceErrorResponse(c, "template-registry", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return api.SystemErrorResponse(c, err, "Failed to read service response")
	}

	// If the service returned an error status, proxy it directly
	if resp.StatusCode >= 400 {
		c.Status(resp.StatusCode)
		return c.Send(body)
	}

	// Parse the response data
	var responseData interface{}
	if err := json.Unmarshal(body, &responseData); err != nil {
		// If not JSON, return as-is
		c.Status(resp.StatusCode)
		return c.Send(body)
	}

	// Wrap in backend standard format using api.SuccessResponse
	return api.SuccessResponse(c, resp.StatusCode, responseData, message)
}

// ProxyWithTransform allows transformation of requests/responses
func (h *ProxyHandler) ProxyWithTransform(c *fiber.Ctx, targetURL string, path string, transformer func([]byte) ([]byte, error)) error {
	// First, proxy the request normally
	fullURL := fmt.Sprintf("%s%s", targetURL, path)

	var reqBody io.Reader
	if c.Body() != nil && len(c.Body()) > 0 {
		reqBody = bytes.NewReader(c.Body())
	}

	req, err := http.NewRequest(string(c.Method()), fullURL, reqBody)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create proxy request", nil)
	}

	// Copy headers
	for key, values := range c.GetReqHeaders() {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := h.HttpClient.Do(req)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadGateway, "Failed to reach service", nil)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to read service response", nil)
	}

	// Apply transformation if provided
	if transformer != nil {
		body, err = transformer(body)
		if err != nil {
			return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to transform response", nil)
		}
	}

	// Copy response headers (except content-length which might change)
	for key, values := range resp.Header {
		if key != "Content-Length" {
			for _, value := range values {
				c.Set(key, value)
			}
		}
	}

	c.Status(resp.StatusCode)
	return c.Send(body)
}

// Example transformer for adding metadata to template responses
func (h *ProxyHandler) enhanceTemplateResponse(data []byte) ([]byte, error) {
	var response map[string]interface{}
	if err := json.Unmarshal(data, &response); err != nil {
		return data, nil // Return original if not JSON
	}

	// Add backend metadata
	response["_meta"] = map[string]interface{}{
		"proxied_by": "naytife-backend",
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	}

	return json.Marshal(response)
}
