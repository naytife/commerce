package handlers

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/petrejonn/naytife/internal/api"
	"github.com/petrejonn/naytife/internal/api/models"
	"github.com/petrejonn/naytife/internal/db"
)

// GetOrders fetches orders for a shop
// @Summary      Fetch orders for a shop
// @Description  Fetch orders for a shop with pagination
// @Tags         orders
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        limit query int false "Limit" default(10)
// @Param        offset query int false "Offset" default(0)
// @Success      200  {object}  []models.Order  "Orders fetched successfully"
// @Failure      400  {object}  models.ErrorResponse "Invalid shop ID"
// @Failure      401  {object}  models.ErrorResponse "Unauthorized"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/orders [get]
func (h *Handler) GetOrders(c *fiber.Ctx) error {
	shopID, err := api.ParseIDParameter(c, "shop_id", "Shop")
	if err != nil {
		return err
	}

	// Parse query parameters
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Get orders from database
	orders, err := h.Repository.ListOrders(c.Context(), db.ListOrdersParams{
		ShopID: shopID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return api.SystemErrorResponse(c, err, "Failed to fetch orders")
	}

	// Get total count for pagination
	total, err := h.Repository.CountOrders(c.Context(), shopID)
	if err != nil {
		return api.SystemErrorResponse(c, err, "Failed to count orders")
	}

	// Map database models to API models
	var result []models.Order
	for _, order := range orders {
		// Get order items
		items, err := h.Repository.GetOrderItemsByOrder(c.Context(), db.GetOrderItemsByOrderParams{
			OrderID: order.OrderID,
			ShopID:  shopID,
		})
		if err != nil {
			return api.SystemErrorResponse(c, err, "Failed to fetch order items")
		}

		// Map items
		var orderItems []models.OrderItem
		for _, item := range items {
			orderItems = append(orderItems, models.OrderItem{
				ID:                 item.OrderItemID,
				Quantity:           item.Quantity,
				Price:              item.Price,
				CreatedAt:          item.CreatedAt,
				UpdatedAt:          item.UpdatedAt,
				ProductVariationID: item.ProductVariationID,
				OrderID:            item.OrderID,
				ShopID:             item.ShopID,
			})
		}

		// Add order with items to result
		result = append(result, models.Order{
			ID:              order.OrderID,
			Status:          order.Status,
			CreatedAt:       order.CreatedAt,
			UpdatedAt:       order.UpdatedAt,
			CustomerID:      order.ShopCustomerID.Bytes,
			Amount:          order.Amount,
			Discount:        order.Discount,
			ShippingCost:    order.ShippingCost,
			Tax:             order.Tax,
			ShippingAddress: order.ShippingAddress,
			PaymentMethod:   order.PaymentMethod,
			PaymentStatus:   order.PaymentStatus,
			ShippingMethod:  order.ShippingMethod,
			ShippingStatus:  order.ShippingStatus,
			TransactionID:   order.TransactionID,
			Username:        order.Username,
			ShopID:          order.ShopID,
			Items:           orderItems,
			CustomerName:    order.CustomerName,
			CustomerEmail:   order.CustomerEmail,
			CustomerPhone:   order.CustomerPhone,
		})
	}

	totalPages := (int(total) + limit - 1) / limit
	response := fiber.Map{
		"data": result,
		"pagination": fiber.Map{
			"page":         page,
			"limit":        limit,
			"total":        total,
			"total_pages":  totalPages,
			"has_next":     page < totalPages,
			"has_previous": page > 1,
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetOrder fetches a single order
// @Summary      Fetch a single order
// @Description  Fetch a single order by ID
// @Tags         orders
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        order_id path string true "Order ID"
// @Success      200  {object}  models.SuccessResponse{data=models.Order}  "Order fetched successfully"
// @Failure      400  {object}  models.ErrorResponse "Invalid shop or order ID"
// @Failure      401  {object}  models.ErrorResponse "Unauthorized"
// @Failure      404  {object}  models.ErrorResponse "Order not found"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/orders/{order_id} [get]
func (h *Handler) GetOrder(c *fiber.Ctx) error {
	shopID, err := strconv.ParseInt(c.Params("shop_id"), 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	orderID, err := strconv.ParseInt(c.Params("order_id"), 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid order ID", nil)
	}

	// Get order from database
	order, err := h.Repository.GetOrder(c.Context(), db.GetOrderParams{
		OrderID: orderID,
		ShopID:  shopID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Order not found", nil)
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch order", nil)
	}

	// Get order items
	items, err := h.Repository.GetOrderItemsByOrder(c.Context(), db.GetOrderItemsByOrderParams{
		OrderID: order.OrderID,
		ShopID:  shopID,
	})
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch order items", nil)
	}

	// Map items
	var orderItems []models.OrderItem
	for _, item := range items {
		orderItems = append(orderItems, models.OrderItem{
			ID:                 item.OrderItemID,
			Quantity:           item.Quantity,
			Price:              item.Price,
			CreatedAt:          item.CreatedAt,
			UpdatedAt:          item.UpdatedAt,
			ProductVariationID: item.ProductVariationID,
			OrderID:            item.OrderID,
			ShopID:             item.ShopID,
		})
	}

	// Map order to API model
	result := models.Order{
		ID:              order.OrderID,
		Status:          order.Status,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
		CustomerID:      order.ShopCustomerID.Bytes,
		Amount:          order.Amount,
		Discount:        order.Discount,
		ShippingCost:    order.ShippingCost,
		Tax:             order.Tax,
		ShippingAddress: order.ShippingAddress,
		PaymentMethod:   order.PaymentMethod,
		PaymentStatus:   order.PaymentStatus,
		ShippingMethod:  order.ShippingMethod,
		ShippingStatus:  order.ShippingStatus,
		TransactionID:   order.TransactionID,
		Username:        order.Username,
		ShopID:          order.ShopID,
		Items:           orderItems,
		CustomerName:    order.CustomerName,
		CustomerEmail:   order.CustomerEmail,
		CustomerPhone:   order.CustomerPhone,
	}

	return api.SuccessResponse(c, fiber.StatusOK, result, "Order fetched successfully")
}

// CreateOrder creates a new order from cart items
// @Summary      Create a new order
// @Description  Create a new order with items from cart
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        order body models.CreateOrderRequest true "Order object that needs to be created"
// @Success      201  {object}  models.SuccessResponse{data=models.Order}  "Order created successfully"
// @Failure      400  {object}  models.ErrorResponse "Invalid request body or shop ID"
// @Failure      401  {object}  models.ErrorResponse "Unauthorized"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/orders [post]
func (h *Handler) CreateOrder(c *fiber.Ctx) error {
	shopID, err := strconv.ParseInt(c.Params("shop_id"), 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	// Parse request body
	var orderReq models.CreateOrderRequest
	if err := c.BodyParser(&orderReq); err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	// Validate request body
	validator := &models.XValidator{}
	if errs := validator.Validate(&orderReq); len(errs) > 0 {
		errMsgs := models.FormatValidationErrors(errs)
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: errMsgs,
		}
	}

	// Create order in a transaction to ensure consistency
	var order db.Order
	var orderItems []db.OrderItem

	err = h.Repository.WithTx(c.Context(), func(q *db.Queries) error {
		// Calculate totals
		var totalAmount float64
		for _, item := range orderReq.Items {
			totalAmount += item.Price * float64(item.Quantity)
		}

		// Create the order
		var custID pgtype.UUID
		if orderReq.CustomerID != nil {
			custID.Valid = true
			if parsed, parseErr := uuid.Parse(*orderReq.CustomerID); parseErr == nil {
				custID.Bytes = parsed
			}
		}

		createOrderParams := db.CreateOrderParams{
			Status:          db.OrderStatusType("pending"),
			Amount:          pgtype.Numeric{Int: big.NewInt(int64(totalAmount * 100)), Valid: true}, // Store as cents
			Discount:        pgtype.Numeric{Int: big.NewInt(int64(orderReq.Discount * 100)), Valid: true},
			ShippingCost:    pgtype.Numeric{Int: big.NewInt(int64(orderReq.ShippingCost * 100)), Valid: true},
			Tax:             pgtype.Numeric{Int: big.NewInt(int64(orderReq.Tax * 100)), Valid: true},
			ShippingAddress: orderReq.ShippingAddress,
			PaymentMethod:   db.PaymentMethodType(orderReq.PaymentMethod),
			PaymentStatus:   db.PaymentStatusType("pending"),
			ShippingMethod:  orderReq.ShippingMethod,
			ShippingStatus:  db.ShippingStatusType("pending"),
			TransactionID:   orderReq.TransactionID,
			Username:        orderReq.CustomerName,
			ShopCustomerID:  custID,
			ShopID:          shopID,
			CustomerName:    orderReq.CustomerName,
			CustomerEmail:   orderReq.CustomerEmail,
			CustomerPhone:   orderReq.CustomerPhone,
		}

		createdOrder, err := q.CreateOrder(c.Context(), createOrderParams)
		if err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}
		order = createdOrder

		// Create order items
		for _, item := range orderReq.Items {
			productVariationID, err := strconv.ParseInt(item.ProductVariationID, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid product variation ID: %s", item.ProductVariationID)
			}

			createItemParams := db.CreateOrderItemParams{
				Quantity:           int64(item.Quantity),
				Price:              pgtype.Numeric{Int: big.NewInt(int64(item.Price * 100)), Valid: true},
				ProductVariationID: productVariationID,
				OrderID:            order.OrderID,
				ShopID:             shopID,
			}

			orderItem, err := q.CreateOrderItem(c.Context(), createItemParams)
			if err != nil {
				return fmt.Errorf("failed to create order item: %w", err)
			}
			orderItems = append(orderItems, orderItem)
		}

		return nil
	})

	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create order", nil)
	}

	// Map order items to API models
	var apiOrderItems []models.OrderItem
	for _, item := range orderItems {
		apiOrderItems = append(apiOrderItems, models.OrderItem{
			ID:                 item.OrderItemID,
			Quantity:           item.Quantity,
			Price:              item.Price,
			CreatedAt:          item.CreatedAt,
			UpdatedAt:          item.UpdatedAt,
			ProductVariationID: item.ProductVariationID,
			OrderID:            item.OrderID,
			ShopID:             item.ShopID,
		})
	}

	// Map order to API model
	result := models.Order{
		ID:              order.OrderID,
		Status:          order.Status,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
		CustomerID:      order.ShopCustomerID.Bytes,
		Amount:          order.Amount,
		Discount:        order.Discount,
		ShippingCost:    order.ShippingCost,
		Tax:             order.Tax,
		ShippingAddress: order.ShippingAddress,
		PaymentMethod:   order.PaymentMethod,
		PaymentStatus:   order.PaymentStatus,
		ShippingMethod:  order.ShippingMethod,
		ShippingStatus:  order.ShippingStatus,
		TransactionID:   order.TransactionID,
		Username:        order.Username,
		ShopID:          order.ShopID,
		Items:           apiOrderItems,
		CustomerName:    order.CustomerName,
		CustomerEmail:   order.CustomerEmail,
		CustomerPhone:   order.CustomerPhone,
	}

	return api.SuccessResponse(c, fiber.StatusCreated, result, "Order created successfully")
}

// UpdateOrder updates an order
// @Summary      Update an order
// @Description  Update an order's details
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        order_id path string true "Order ID"
// @Param        order body models.UpdateOrderParams true "Order object that needs to be updated"
// @Success      200  {object}  models.SuccessResponse{data=models.Order}  "Order updated successfully"
// @Failure      400  {object}  models.ErrorResponse "Invalid request body or IDs"
// @Failure      401  {object}  models.ErrorResponse "Unauthorized"
// @Failure      404  {object}  models.ErrorResponse "Order not found"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/orders/{order_id} [put]
func (h *Handler) UpdateOrder(c *fiber.Ctx) error {
	shopID, err := api.ParseIDParameter(c, "shop_id", "Shop")
	if err != nil {
		return err
	}

	orderID, err := api.ParseIDParameter(c, "order_id", "Order")
	if err != nil {
		return err
	}

	// Parse request body
	var orderParams models.UpdateOrderParams
	if err := c.BodyParser(&orderParams); err != nil {
		return api.BusinessLogicErrorResponse(c, "Invalid request body")
	}

	// Validate request body
	if err := api.ValidateRequest(c, &orderParams); err != nil {
		return err
	}

	// Check if order exists
	_, err = h.Repository.GetOrder(c.Context(), db.GetOrderParams{
		OrderID: orderID,
		ShopID:  shopID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.NotFoundErrorResponse(c, "Order")
		}
		return api.SystemErrorResponse(c, err, "Failed to fetch order")
	}

	// Update order
	err = h.Repository.UpdateOrder(c.Context(), db.UpdateOrderParams{
		Status:          orderParams.Status,
		Amount:          orderParams.Amount,
		Discount:        orderParams.Discount,
		ShippingCost:    orderParams.ShippingCost,
		Tax:             orderParams.Tax,
		ShippingAddress: orderParams.ShippingAddress,
		PaymentMethod:   orderParams.PaymentMethod,
		PaymentStatus:   orderParams.PaymentStatus,
		ShippingMethod:  orderParams.ShippingMethod,
		ShippingStatus:  orderParams.ShippingStatus,
		TransactionID:   orderParams.TransactionID,
		Username:        orderParams.Username,
		CustomerName:    orderParams.CustomerName,
		CustomerEmail:   orderParams.CustomerEmail,
		CustomerPhone:   orderParams.CustomerPhone,
		OrderID:         orderID,
		ShopID:          shopID,
	})
	if err != nil {
		return api.SystemErrorResponse(c, err, "Failed to update order")
	}

	// Get updated order
	order, err := h.Repository.GetOrder(c.Context(), db.GetOrderParams{
		OrderID: orderID,
		ShopID:  shopID,
	})
	if err != nil {
		return api.SystemErrorResponse(c, err, "Failed to fetch updated order")
	}

	// Get order items
	items, err := h.Repository.GetOrderItemsByOrder(c.Context(), db.GetOrderItemsByOrderParams{
		OrderID: order.OrderID,
		ShopID:  shopID,
	})
	if err != nil {
		return api.SystemErrorResponse(c, err, "Failed to fetch order items")
	}

	// Map items
	var orderItems []models.OrderItem
	for _, item := range items {
		orderItems = append(orderItems, models.OrderItem{
			ID:                 item.OrderItemID,
			Quantity:           item.Quantity,
			Price:              item.Price,
			CreatedAt:          item.CreatedAt,
			UpdatedAt:          item.UpdatedAt,
			ProductVariationID: item.ProductVariationID,
			OrderID:            item.OrderID,
			ShopID:             item.ShopID,
		})
	}

	// Map order to API model
	result := models.Order{
		ID:              order.OrderID,
		Status:          order.Status,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
		CustomerID:      order.ShopCustomerID.Bytes,
		Amount:          order.Amount,
		Discount:        order.Discount,
		ShippingCost:    order.ShippingCost,
		Tax:             order.Tax,
		ShippingAddress: order.ShippingAddress,
		PaymentMethod:   order.PaymentMethod,
		PaymentStatus:   order.PaymentStatus,
		ShippingMethod:  order.ShippingMethod,
		ShippingStatus:  order.ShippingStatus,
		TransactionID:   order.TransactionID,
		Username:        order.Username,
		ShopID:          order.ShopID,
		Items:           orderItems,
		CustomerName:    order.CustomerName,
		CustomerEmail:   order.CustomerEmail,
		CustomerPhone:   order.CustomerPhone,
	}

	return api.SuccessResponse(c, fiber.StatusOK, result, "Order updated successfully")
}

// UpdateOrderStatus updates an order's status
// @Summary      Update an order's status
// @Description  Update only the status of an order
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        order_id path string true "Order ID"
// @Param        status body models.UpdateOrderStatusParams true "Status object"
// @Success      200  {object}  models.SuccessResponse{data=models.Order}  "Order status updated successfully"
// @Failure      400  {object}  models.ErrorResponse "Invalid request body or IDs"
// @Failure      401  {object}  models.ErrorResponse "Unauthorized"
// @Failure      404  {object}  models.ErrorResponse "Order not found"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/orders/{order_id}/status [patch]
func (h *Handler) UpdateOrderStatus(c *fiber.Ctx) error {
	shopID, err := api.ParseIDParameter(c, "shop_id", "Shop")
	if err != nil {
		return err
	}

	orderID, err := api.ParseIDParameter(c, "order_id", "Order")
	if err != nil {
		return err
	}

	// Parse request body
	var statusParams models.UpdateOrderStatusParams
	if err := c.BodyParser(&statusParams); err != nil {
		return api.BusinessLogicErrorResponse(c, "Invalid request body")
	}

	// Validate request body
	if err := api.ValidateRequest(c, &statusParams); err != nil {
		return err
	}

	// Get current order to preserve other fields
	currentOrder, err := h.Repository.GetOrder(c.Context(), db.GetOrderParams{
		OrderID: orderID,
		ShopID:  shopID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.NotFoundErrorResponse(c, "Order")
		}
		return api.SystemErrorResponse(c, err, "Failed to fetch order")
	}

	// Update only the status
	err = h.Repository.UpdateOrder(c.Context(), db.UpdateOrderParams{
		Status:          statusParams.Status,
		Amount:          currentOrder.Amount,
		Discount:        currentOrder.Discount,
		ShippingCost:    currentOrder.ShippingCost,
		Tax:             currentOrder.Tax,
		ShippingAddress: currentOrder.ShippingAddress,
		PaymentMethod:   currentOrder.PaymentMethod,
		PaymentStatus:   currentOrder.PaymentStatus,
		ShippingMethod:  currentOrder.ShippingMethod,
		ShippingStatus:  currentOrder.ShippingStatus,
		TransactionID:   currentOrder.TransactionID,
		Username:        currentOrder.Username,
		CustomerName:    currentOrder.CustomerName,
		CustomerEmail:   currentOrder.CustomerEmail,
		CustomerPhone:   currentOrder.CustomerPhone,
		OrderID:         orderID,
		ShopID:          shopID,
	})
	if err != nil {
		return api.SystemErrorResponse(c, err, "Failed to update order status")
	}

	// Get updated order
	order, err := h.Repository.GetOrder(c.Context(), db.GetOrderParams{
		OrderID: orderID,
		ShopID:  shopID,
	})
	if err != nil {
		return api.SystemErrorResponse(c, err, "Failed to fetch updated order")
	}

	// Get order items
	items, err := h.Repository.GetOrderItemsByOrder(c.Context(), db.GetOrderItemsByOrderParams{
		OrderID: order.OrderID,
		ShopID:  shopID,
	})
	if err != nil {
		return api.SystemErrorResponse(c, err, "Failed to fetch order items")
	}

	// Map items
	var orderItems []models.OrderItem
	for _, item := range items {
		orderItems = append(orderItems, models.OrderItem{
			ID:                 item.OrderItemID,
			Quantity:           item.Quantity,
			Price:              item.Price,
			CreatedAt:          item.CreatedAt,
			UpdatedAt:          item.UpdatedAt,
			ProductVariationID: item.ProductVariationID,
			OrderID:            item.OrderID,
			ShopID:             item.ShopID,
		})
	}

	// Map order to API model
	result := models.Order{
		ID:              order.OrderID,
		Status:          order.Status,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
		CustomerID:      order.ShopCustomerID.Bytes,
		Amount:          order.Amount,
		Discount:        order.Discount,
		ShippingCost:    order.ShippingCost,
		Tax:             order.Tax,
		ShippingAddress: order.ShippingAddress,
		PaymentMethod:   order.PaymentMethod,
		PaymentStatus:   order.PaymentStatus,
		ShippingMethod:  order.ShippingMethod,
		ShippingStatus:  order.ShippingStatus,
		TransactionID:   order.TransactionID,
		Username:        order.Username,
		ShopID:          order.ShopID,
		Items:           orderItems,
		CustomerName:    order.CustomerName,
		CustomerEmail:   order.CustomerEmail,
		CustomerPhone:   order.CustomerPhone,
	}

	return api.SuccessResponse(c, fiber.StatusOK, result, "Order status updated successfully")
}

// DeleteOrder deletes an order
// @Summary      Delete an order
// @Description  Delete an order and all its items
// @Tags         orders
// @Produce      json
// @Param        shop_id path string true "Shop ID"
// @Param        order_id path string true "Order ID"
// @Success      200  {object}  models.SuccessResponse  "Order deleted successfully"
// @Failure      400  {object}  models.ErrorResponse "Invalid shop or order ID"
// @Failure      401  {object}  models.ErrorResponse "Unauthorized"
// @Failure      404  {object}  models.ErrorResponse "Order not found"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/orders/{order_id} [delete]
func (h *Handler) DeleteOrder(c *fiber.Ctx) error {
	shopID, err := strconv.ParseInt(c.Params("shop_id"), 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid shop ID", nil)
	}

	orderID, err := strconv.ParseInt(c.Params("order_id"), 10, 64)
	if err != nil {
		return api.ErrorResponse(c, fiber.StatusBadRequest, "Invalid order ID", nil)
	}

	// Check if order exists
	_, err = h.Repository.GetOrder(c.Context(), db.GetOrderParams{
		OrderID: orderID,
		ShopID:  shopID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return api.ErrorResponse(c, fiber.StatusNotFound, "Order not found", nil)
		}
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch order", nil)
	}

	// Delete order and its items in a transaction
	err = h.Repository.WithTx(c.Context(), func(q *db.Queries) error {
		// Delete order items first (due to foreign key constraints)
		txErr := q.DeleteOrderItemsByOrder(c.Context(), db.DeleteOrderItemsByOrderParams{
			OrderID: orderID,
			ShopID:  shopID,
		})
		if txErr != nil {
			return txErr
		}

		// Delete the order
		txErr = q.DeleteOrder(c.Context(), db.DeleteOrderParams{
			OrderID: orderID,
			ShopID:  shopID,
		})
		return txErr
	})

	if err != nil {
		return api.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete order", nil)
	}

	return api.SuccessResponse(c, fiber.StatusOK, nil, fmt.Sprintf("Order %d deleted successfully", orderID))
}
