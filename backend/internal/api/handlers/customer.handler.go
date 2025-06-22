package handlers

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/petrejonn/naytife/internal/api"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"
)

// UpsertCustomer creates or updates a customer
// @Summary      Create or update a customer
// @Description
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param        customer body models.RegisterCustomerParams true "Customer object that needs to be created or updated"
// @Success      200  {object}   models.SuccessResponse{data=models.CustomerResponse} "Customer created or updated successfully"
// @Router       /auth/register-customer [post]
func (h *Handler) UpsertCustomer(c *fiber.Ctx) error {
	var param models.RegisterCustomerParams
	if err := c.BodyParser(&param); err != nil {
		return api.BusinessLogicErrorResponse(c, "Failed to parse request body")
	}

	if err := api.ValidateRequest(c, &param); err != nil {
		return err
	}

	objDB, err := h.Repository.UpsertCustomer(c.Context(), db.UpsertCustomerParams{
		ShopID:         param.ShopID,
		Email:          *param.Email,
		Name:           param.Name,
		Locale:         param.Locale,
		ProfilePicture: param.ProfilePicture,
		VerifiedEmail:  param.VerifiedEmail,
		AuthProvider:   param.AuthProvider,
		AuthProviderID: param.AuthProviderID,
	})
	if err != nil {
		return api.SystemErrorResponse(c, err, "Failed to create or update customer")
	}
	resp := models.CustomerResponse{
		CustomerID:     objDB.ShopCustomerID,
		ShopID:         objDB.ShopID,
		Email:          &objDB.Email,
		Name:           objDB.Name,
		ProfilePicture: objDB.ProfilePicture,
		CreatedAt:      objDB.CreatedAt,
		LastLogin:      objDB.LastLogin,
		VerifiedEmail:  objDB.VerifiedEmail,
		AuthProvider:   objDB.AuthProvider,
		AuthProviderID: objDB.AuthProviderID,
	}
	return api.SuccessResponse(c, fiber.StatusOK, resp, "Customer created or updated successfully")
}

// GetCustomerByEmail fetches a customer by email for a specific shop
// @Summary      Fetch a customer by email for a specific shop
// @Description
// @Tags         customer
// @Produce      json
// @Param        subdomain query string true "Shop subdomain"
// @Param        email query string true "Customer email"
// @Success      200  {object}   models.SuccessResponse{data=models.CustomerResponse} "Customer fetched successfully"
// @Router       /shops/customerinfo [get]
func (h *Handler) GetCustomerByEmail(c *fiber.Ctx) error {
	subdomain := c.Query("subdomain")
	email := c.Query("email")

	// Fetch the customer using the email and the retrieved shop_id
	customer, err := h.Repository.GetCustomerByEmail(c.Context(), db.GetCustomerByEmailParams{
		Email:     email,
		Subdomain: subdomain,
	})
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get customer", nil)
	}

	resp := models.CustomerResponse{
		CustomerID:     customer.ShopCustomerID,
		Email:          &customer.Email,
		Name:           customer.Name,
		ProfilePicture: customer.ProfilePicture,
		CreatedAt:      customer.CreatedAt,
		LastLogin:      customer.LastLogin,
		VerifiedEmail:  customer.VerifiedEmail,
		AuthProvider:   customer.AuthProvider,
		AuthProviderID: customer.AuthProviderID,
	}
	return api.SuccessResponse(c, fiber.StatusOK, resp, "Customer fetched successfully")
}

// GetCustomers fetches all customers for a shop with pagination
// @Summary      Fetch all customers for a shop
// @Description  Get paginated list of customers for a specific shop
// @Tags         customer
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        limit query int false "Limit" default(20)
// @Param        offset query int false "Offset" default(0)
// @Success      200  {object}   models.SuccessResponse{data=models.CustomerListResponse} "Customers fetched successfully"
// @Failure      400  {object}   models.ErrorResponse "Bad request"
// @Failure      500  {object}   models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/customers [get]
func (h *Handler) GetCustomers(c *fiber.Ctx) error {
	shopID, err := api.ParseIDParameter(c, "shop_id", "Shop")
	if err != nil {
		return err
	}

	limit, offset, err := api.ParsePaginationParams(c)
	if err != nil {
		return api.BusinessLogicErrorResponse(c, "Invalid pagination parameters")
	}

	customers, err := h.Repository.GetCustomers(c.Context(), db.GetCustomersParams{
		ShopID: shopID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return api.SystemErrorResponse(c, err, "Failed to fetch customers")
	}

	totalCount, err := h.Repository.GetCustomersCount(c.Context(), shopID)
	if err != nil {
		return api.SystemErrorResponse(c, err, "Failed to get customer count")
	}

	// Convert to response format
	customerResponses := make([]models.CustomerResponse, len(customers))
	for i, customer := range customers {
		customerResponses[i] = models.CustomerResponse{
			CustomerID:     customer.ShopCustomerID,
			ShopID:         customer.ShopID,
			Email:          &customer.Email,
			Name:           customer.Name,
			Locale:         customer.Locale,
			ProfilePicture: customer.ProfilePicture,
			CreatedAt:      customer.CreatedAt,
			LastLogin:      customer.LastLogin,
			VerifiedEmail:  customer.VerifiedEmail,
			AuthProvider:   customer.AuthProvider,
			AuthProviderID: customer.AuthProviderID,
		}
	}

	page := (offset / limit) + 1
	return api.PaginatedSuccessResponse(c, fiber.StatusOK, customerResponses, totalCount, page, limit, "Customers fetched successfully")
}

// SearchCustomers searches customers by name or email
// @Summary      Search customers by name or email
// @Description  Search customers in a shop by name or email
// @Tags         customer
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        q query string true "Search query"
// @Param        limit query int false "Limit" default(20)
// @Param        offset query int false "Offset" default(0)
// @Success      200  {object}   models.SuccessResponse{data=models.CustomerListResponse} "Customers found successfully"
// @Failure      400  {object}   models.ErrorResponse "Bad request"
// @Failure      500  {object}   models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/customers/search [get]
func (h *Handler) SearchCustomers(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	query := c.Query("q", "")
	if query == "" {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Search query is required", nil)
	}

	limitStr := c.Query("limit", "20")
	offsetStr := c.Query("offset", "0")

	limit, _ := strconv.ParseInt(limitStr, 10, 32)
	offset, _ := strconv.ParseInt(offsetStr, 10, 32)

	customers, err := h.Repository.SearchCustomers(c.Context(), db.SearchCustomersParams{
		ShopID: shopID,
		Lower:  "%" + strings.ToLower(query) + "%",
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to search customers", nil)
	}

	// Convert to response format
	customerResponses := make([]models.CustomerResponse, len(customers))
	for i, customer := range customers {
		customerResponses[i] = models.CustomerResponse{
			CustomerID:     customer.ShopCustomerID,
			ShopID:         customer.ShopID,
			Email:          &customer.Email,
			Name:           customer.Name,
			Locale:         customer.Locale,
			ProfilePicture: customer.ProfilePicture,
			CreatedAt:      customer.CreatedAt,
			LastLogin:      customer.LastLogin,
			VerifiedEmail:  customer.VerifiedEmail,
			AuthProvider:   customer.AuthProvider,
			AuthProviderID: customer.AuthProviderID,
		}
	}

	response := models.CustomerListResponse{
		Customers: customerResponses,
		Total:     int64(len(customerResponses)), // For search results, we return the current count
		Page:      int(offset/limit) + 1,
		Limit:     int(limit),
	}

	return api.SuccessResponse(c, fiber.StatusOK, response, "Customers found successfully")
}

// GetCustomerById fetches a specific customer by ID
// @Summary      Fetch a customer by ID
// @Description  Get detailed customer information by customer ID
// @Tags         customer
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        customer_id path string true "Customer ID"
// @Success      200  {object}   models.SuccessResponse{data=models.CustomerResponse} "Customer fetched successfully"
// @Failure      400  {object}   models.ErrorResponse "Bad request"
// @Failure      404  {object}   models.ErrorResponse "Customer not found"
// @Failure      500  {object}   models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/customers/{customer_id} [get]
func (h *Handler) GetCustomerById(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	customerIDStr := c.Params("customer_id")
	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid customer ID", nil)
	}

	customer, err := h.Repository.GetCustomerById(c.Context(), db.GetCustomerByIdParams{
		ShopCustomerID: customerID,
		ShopID:         shopID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Customer not found", nil)
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch customer", nil)
	}

	response := models.CustomerResponse{
		CustomerID:     customer.ShopCustomerID,
		ShopID:         customer.ShopID,
		Email:          &customer.Email,
		Name:           customer.Name,
		Locale:         customer.Locale,
		ProfilePicture: customer.ProfilePicture,
		CreatedAt:      customer.CreatedAt,
		LastLogin:      customer.LastLogin,
		VerifiedEmail:  customer.VerifiedEmail,
		AuthProvider:   customer.AuthProvider,
		AuthProviderID: customer.AuthProviderID,
	}

	return api.SuccessResponse(c, fiber.StatusOK, response, "Customer fetched successfully")
}

// UpdateCustomer updates customer information
// @Summary      Update customer information
// @Description  Update customer details
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        customer_id path string true "Customer ID"
// @Param        customer body models.UpdateCustomerParams true "Customer update parameters"
// @Success      200  {object}   models.SuccessResponse{data=models.CustomerResponse} "Customer updated successfully"
// @Failure      400  {object}   models.ErrorResponse "Bad request"
// @Failure      404  {object}   models.ErrorResponse "Customer not found"
// @Failure      500  {object}   models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/customers/{customer_id} [put]
func (h *Handler) UpdateCustomer(c *fiber.Ctx) error {
	shopID, err := api.ParseIDParameter(c, "shop_id", "Shop")
	if err != nil {
		return err
	}

	customerIDStr := c.Params("customer_id")
	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		return api.BusinessLogicErrorResponse(c, "Invalid customer ID")
	}

	var param models.UpdateCustomerParams
	if err := c.BodyParser(&param); err != nil {
		return api.BusinessLogicErrorResponse(c, "Failed to parse request body")
	}

	if err := api.ValidateRequest(c, &param); err != nil {
		return err
	}

	updatedCustomer, err := h.Repository.UpdateCustomer(c.Context(), db.UpdateCustomerParams{
		ShopCustomerID: customerID,
		ShopID:         shopID,
		Name:           param.Name,
		Locale:         param.Locale,
		ProfilePicture: param.ProfilePicture,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.NotFoundErrorResponse(c, "Customer")
		}
		return api.SystemErrorResponse(c, err, "Failed to update customer")
	}

	response := models.CustomerResponse{
		CustomerID:     updatedCustomer.ShopCustomerID,
		ShopID:         updatedCustomer.ShopID,
		Email:          &updatedCustomer.Email,
		Name:           updatedCustomer.Name,
		Locale:         updatedCustomer.Locale,
		ProfilePicture: updatedCustomer.ProfilePicture,
		CreatedAt:      updatedCustomer.CreatedAt,
		LastLogin:      updatedCustomer.LastLogin,
		VerifiedEmail:  updatedCustomer.VerifiedEmail,
		AuthProvider:   updatedCustomer.AuthProvider,
		AuthProviderID: updatedCustomer.AuthProviderID,
	}

	return api.SuccessResponse(c, fiber.StatusOK, response, "Customer updated successfully")
}

// DeleteCustomer deletes a customer
// @Summary      Delete a customer
// @Description  Delete a customer from the shop
// @Tags         customer
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        customer_id path string true "Customer ID"
// @Success      200  {object}   models.SuccessResponse "Customer deleted successfully"
// @Failure      400  {object}   models.ErrorResponse "Bad request"
// @Failure      404  {object}   models.ErrorResponse "Customer not found"
// @Failure      500  {object}   models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/customers/{customer_id} [delete]
func (h *Handler) DeleteCustomer(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	customerIDStr := c.Params("customer_id")
	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid customer ID", nil)
	}

	err = h.Repository.DeleteCustomer(c.Context(), db.DeleteCustomerParams{
		ShopCustomerID: customerID,
		ShopID:         shopID,
	})
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete customer", nil)
	}

	return api.SuccessResponse(c, fiber.StatusOK, nil, "Customer deleted successfully")
}

// GetCustomerOrders fetches order history for a specific customer
// @Summary      Fetch customer order history
// @Description  Get paginated order history for a specific customer
// @Tags         customer
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        customer_id path string true "Customer ID"
// @Param        limit query int false "Limit" default(20)
// @Param        offset query int false "Offset" default(0)
// @Success      200  {object}   models.SuccessResponse{data=models.CustomerOrdersResponse} "Customer orders fetched successfully"
// @Failure      400  {object}   models.ErrorResponse "Bad request"
// @Failure      404  {object}   models.ErrorResponse "Customer not found"
// @Failure      500  {object}   models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/customers/{customer_id}/orders [get]
func (h *Handler) GetCustomerOrders(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id", "0")
	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	customerIDStr := c.Params("customer_id")
	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid customer ID", nil)
	}

	// First get customer to get their email
	customer, err := h.Repository.GetCustomerById(c.Context(), db.GetCustomerByIdParams{
		ShopCustomerID: customerID,
		ShopID:         shopID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Customer not found", nil)
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch customer", nil)
	}

	limitStr := c.Query("limit", "20")
	offsetStr := c.Query("offset", "0")

	limit, _ := strconv.ParseInt(limitStr, 10, 32)
	offset, _ := strconv.ParseInt(offsetStr, 10, 32)

	orders, err := h.Repository.GetCustomerOrders(c.Context(), db.GetCustomerOrdersParams{
		CustomerEmail: &customer.Email,
		ShopID:        shopID,
		Limit:         int32(limit),
		Offset:        int32(offset),
	})
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch customer orders", nil)
	}

	// Group orders and their items
	orderMap := make(map[int64]*models.OrderResponse)
	for _, orderRow := range orders {
		if _, exists := orderMap[orderRow.OrderID]; !exists {
			orderMap[orderRow.OrderID] = &models.OrderResponse{
				OrderID:         orderRow.OrderID,
				Status:          string(orderRow.Status),
				Amount:          models.NumericToFloat64(orderRow.Amount),
				Discount:        models.NumericToFloat64(orderRow.Discount),
				ShippingCost:    models.NumericToFloat64(orderRow.ShippingCost),
				Tax:             models.NumericToFloat64(orderRow.Tax),
				ShippingAddress: orderRow.ShippingAddress,
				PaymentMethod:   string(orderRow.PaymentMethod),
				PaymentStatus:   string(orderRow.PaymentStatus),
				ShippingMethod:  orderRow.ShippingMethod,
				ShippingStatus:  string(orderRow.ShippingStatus),
				TransactionID:   orderRow.TransactionID,
				Username:        orderRow.Username,
				CreatedAt:       models.TimestamptzToTime(orderRow.CreatedAt),
				UpdatedAt:       models.TimestamptzToTime(orderRow.UpdatedAt),
				CustomerName:    orderRow.CustomerName,
				CustomerEmail:   orderRow.CustomerEmail,
				CustomerPhone:   orderRow.CustomerPhone,
				Items:           []models.OrderItemResponse{},
			}
		}

		// Add order item if it exists
		if orderRow.OrderItemID != nil {
			item := models.OrderItemResponse{
				OrderItemID:        *orderRow.OrderItemID,
				ProductVariationID: *orderRow.ProductVariationID,
				Quantity:           *orderRow.Quantity,
				Price:              models.NumericToFloat64(orderRow.ItemPrice),
			}
			orderMap[orderRow.OrderID].Items = append(orderMap[orderRow.OrderID].Items, item)
		}
	}

	// Convert map to slice
	orderResponses := make([]models.OrderResponse, 0, len(orderMap))
	for _, order := range orderMap {
		orderResponses = append(orderResponses, *order)
	}

	response := models.CustomerOrdersResponse{
		Orders: orderResponses,
		Total:  len(orderResponses),
	}

	return api.SuccessResponse(c, fiber.StatusOK, response, "Customer orders fetched successfully")
}
