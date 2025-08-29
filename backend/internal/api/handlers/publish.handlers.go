package handlers

// import (
// 	"encoding/json"
// 	"fmt"
// 	"strconv"
// 	"time"

// 	"github.com/go-redis/redis/v8"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/petrejonn/naytife/internal/api"
// 	"github.com/petrejonn/naytife/internal/api/models"
// 	"github.com/petrejonn/naytife/internal/db"
// )

// // PublishHandler handles publishing operations
// type PublishHandler struct {
// 	Repository  db.Repository
// 	RedisClient *redis.Client
// }

// // NewPublishHandler creates a new publish handler
// func NewPublishHandler(repo db.Repository, redisClient *redis.Client) *PublishHandler {
// 	return &PublishHandler{
// 		Repository:  repo,
// 		RedisClient: redisClient,
// 	}
// }

// // BuildJob represents a simplified build job that relies on GraphQL queries
// type BuildJob struct {
// 	SiteName     string         `json:"site_name"`
// 	TemplateName string         `json:"template_name"`
// 	ShopID       int64          `json:"shop_id"`
// 	Subdomain    string         `json:"subdomain"`
// 	Changes      []ChangeRecord `json:"changes"`
// 	Priority     int            `json:"priority"`
// 	Timestamp    time.Time      `json:"timestamp"`
// 	DataOnly     bool           `json:"data_only"` // true for JSON-only updates, false for full rebuilds
// 	DataType     string         `json:"data_type"` // "shop", "products", "all" - determines which data files to update
// }

// // ChangeRecord represents a tracked change
// type ChangeRecord struct {
// 	ID          string    `json:"id"`
// 	Type        string    `json:"type"`
// 	Entity      string    `json:"entity"`
// 	Timestamp   time.Time `json:"timestamp"`
// 	Description string    `json:"description"`
// }

// // PublishRequest represents the request body for publishing
// type PublishRequest struct {
// 	TemplateName string         `json:"template_name" validate:"required"`
// 	Changes      []ChangeRecord `json:"changes"`
// }

// // PublishResponse represents the response from publishing
// type PublishResponse struct {
// 	JobID     string    `json:"job_id"`
// 	Status    string    `json:"status"`
// 	Timestamp time.Time `json:"timestamp"`
// 	QueuePos  int       `json:"queue_position"`
// }

// // TriggerSiteBuild triggers a site build for a shop
// // @Summary Trigger site build
// // @Description Queue a build job for the shop's static site
// // @Tags publish
// // @Accept json
// // @Produce json
// // @Param shop_id path string true "Shop ID"
// // @Param request body PublishRequest true "Publish request"
// // @Success 200 {object} models.SuccessResponse{data=PublishResponse} "Build job queued successfully"
// // @Failure 400 {object} models.ErrorResponse "Invalid request"
// // @Failure 404 {object} models.ErrorResponse "Shop not found"
// // @Failure 500 {object} models.ErrorResponse "Failed to queue build job"
// // @Security OAuth2AccessCode
// // @Router /shops/{shop_id}/publish [post]
// func (h *PublishHandler) TriggerSiteBuild(c *fiber.Ctx) error {
// 	shopIDStr := c.Params("shop_id")
// 	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
// 	if err != nil {
// 		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
// 	}

// 	var req PublishRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
// 	}

// 	// Validate request
// 	validator := &models.XValidator{}
// 	if errs := validator.Validate(&req); len(errs) > 0 {
// 		errMsgs := models.FormatValidationErrors(errs)
// 		return api.ErrorResponse(c, fiber.StatusBadRequest, errMsgs, nil)
// 	}

// 	// Fetch shop to get subdomain for the build
// 	shop, err := h.Repository.GetShop(c.Context(), shopID)
// 	if err != nil {
// 		return api.ErrorResponse(c, fiber.StatusNotFound, "Shop not found", nil)
// 	}

// 	// Create simplified build job - cloud-build will query data via GraphQL
// 	jobID := fmt.Sprintf("build_%s_%d", shop.Subdomain, time.Now().Unix())

// 	// Determine if this should be a data-only build and what type of data
// 	dataOnly := false
// 	dataType := "all" // Default to updating all data

// 	// Check if all changes are data-only changes and determine specific data type
// 	if len(req.Changes) > 0 {
// 		allDataOnlyChanges := true
// 		shopChanges := false
// 		productChanges := false

// 		for _, change := range req.Changes {
// 			switch change.Type {
// 			case "product_create", "product_update", "product_delete", "inventory_update":
// 				// These changes only require product data updates
// 				productChanges = true
// 				continue
// 			case "shop_update":
// 				// These changes only require shop data updates
// 				shopChanges = true
// 				continue
// 			case "category_update", "template_update":
// 				// These changes require full rebuilds
// 				allDataOnlyChanges = false
// 				break
// 			default:
// 				// Unknown change types default to full rebuild for safety
// 				allDataOnlyChanges = false
// 				break
// 			}
// 		}

// 		if allDataOnlyChanges {
// 			dataOnly = true
// 			// Determine specific data type if only one type was changed
// 			if shopChanges && !productChanges {
// 				dataType = "shop"
// 			} else if productChanges && !shopChanges {
// 				dataType = "products"
// 			}
// 			// If both changed, keep "all"
// 		}
// 	}

// 	buildJob := BuildJob{
// 		SiteName:     shop.Subdomain,
// 		TemplateName: req.TemplateName,
// 		ShopID:       shopID,
// 		Subdomain:    shop.Subdomain,
// 		Changes:      req.Changes,
// 		Priority:     1, // Normal priority
// 		Timestamp:    time.Now(),
// 		DataOnly:     dataOnly,
// 		DataType:     dataType,
// 	}

// 	// Serialize job to JSON
// 	jobJSON, err := json.Marshal(buildJob)
// 	if err != nil {
// 		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to serialize build job", nil)
// 	}

// 	// Queue the job in Redis
// 	queueName := "build-queue"
// 	err = h.RedisClient.RPush(c.Context(), queueName, jobJSON).Err()
// 	if err != nil {
// 		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to queue build job", nil)
// 	}

// 	// Get queue position
// 	queuePos, err := h.RedisClient.LLen(c.Context(), queueName).Result()
// 	if err != nil {
// 		queuePos = 0 // Default if we can't get queue length
// 	}

// 	// Prepare response
// 	response := PublishResponse{
// 		JobID:     jobID,
// 		Status:    "queued",
// 		Timestamp: time.Now(),
// 		QueuePos:  int(queuePos),
// 	}

// 	return api.SuccessResponse(c, fiber.StatusOK, response, "Build job queued successfully")
// }

// // GetPublishStatus gets the status of a publish job
// // @Summary Get publish status
// // @Description Get the current status of a publish job
// // @Tags publish
// // @Produce json
// // @Param shop_id path string true "Shop ID"
// // @Param job_id query string true "Job ID"
// // @Success 200 {object} models.SuccessResponse{data=map[string]interface{}} "Job status retrieved"
// // @Failure 400 {object} models.ErrorResponse "Invalid request"
// // @Failure 404 {object} models.ErrorResponse "Job not found"
// // @Security OAuth2AccessCode
// // @Router /shops/{shop_id}/publish/status [get]
// func (h *PublishHandler) GetPublishStatus(c *fiber.Ctx) error {
// 	shopIDStr := c.Params("shop_id")
// 	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
// 	if err != nil {
// 		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
// 	}

// 	jobID := c.Query("job_id")
// 	if jobID == "" {
// 		return api.ErrorResponse(c, fiber.StatusBadRequest, "job_id is required", nil)
// 	}

// 	// For now, return a simple status
// 	// TODO: Implement proper job tracking in Redis or database
// 	status := map[string]interface{}{
// 		"shop_id":   shopID,
// 		"job_id":    jobID,
// 		"status":    "processing", // Could be: queued, processing, completed, failed
// 		"progress":  75,           // Percentage
// 		"message":   "Building site...",
// 		"timestamp": time.Now(),
// 	}

// 	return api.SuccessResponse(c, fiber.StatusOK, status, "Job status retrieved")
// }

// // GetPublishHistory gets the publish history for a shop
// // @Summary Get publish history
// // @Description Get the publish history for a shop
// // @Tags publish
// // @Produce json
// // @Param shop_id path string true "Shop ID"
// // @Param limit query int false "Limit" default(10)
// // @Param offset query int false "Offset" default(0)
// // @Success 200 {object} models.SuccessResponse{data=[]map[string]interface{}} "Publish history retrieved"
// // @Failure 400 {object} models.ErrorResponse "Invalid request"
// // @Security OAuth2AccessCode
// // @Router /shops/{shop_id}/publish/history [get]
// func (h *PublishHandler) GetPublishHistory(c *fiber.Ctx) error {
// 	shopIDStr := c.Params("shop_id")
// 	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
// 	if err != nil {
// 		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
// 	}

// 	// Parse pagination
// 	limit, _ := strconv.Atoi(c.Query("limit", "10"))
// 	offset, _ := strconv.Atoi(c.Query("offset", "0"))

// 	// For now, return mock data
// 	// TODO: Implement proper publish history tracking
// 	history := []map[string]interface{}{
// 		{
// 			"id":           fmt.Sprintf("build_shop%d_1703123456", shopID),
// 			"shop_id":      shopID,
// 			"status":       "completed",
// 			"timestamp":    time.Now().Add(-1 * time.Hour),
// 			"changes":      3,
// 			"duration":     "2m 34s",
// 			"triggered_by": "admin@example.com",
// 		},
// 		{
// 			"id":           fmt.Sprintf("build_shop%d_1703123400", shopID),
// 			"shop_id":      shopID,
// 			"status":       "completed",
// 			"timestamp":    time.Now().Add(-3 * time.Hour),
// 			"changes":      1,
// 			"duration":     "1m 45s",
// 			"triggered_by": "admin@example.com",
// 		},
// 	}

// 	// Apply pagination
// 	start := offset
// 	end := offset + limit
// 	if start >= len(history) {
// 		history = []map[string]interface{}{}
// 	} else {
// 		if end > len(history) {
// 			end = len(history)
// 		}
// 		history = history[start:end]
// 	}

// 	return api.SuccessResponse(c, fiber.StatusOK, history, "Publish history retrieved")
// }
