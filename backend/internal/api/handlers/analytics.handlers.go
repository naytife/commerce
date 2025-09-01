package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/petrejonn/naytife/internal/api"
	"github.com/petrejonn/naytife/internal/db"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AnalyticsHandler struct {
	Repository db.Repository
}

func NewAnalyticsHandler(repo db.Repository) *AnalyticsHandler {
	return &AnalyticsHandler{Repository: repo}
}

func toInt64(val interface{}) int64 {
	switch v := val.(type) {
	case int64:
		return v
	case int32:
		return int64(v)
	case float64:
		return int64(v)
	case float32:
		return int64(v)
	case []uint8: // for numeric/decimal from Postgres
		// Try float first, then int
		f, err := strconv.ParseFloat(string(v), 64)
		if err == nil {
			return int64(f)
		}
		i, _ := strconv.ParseInt(string(v), 10, 64)
		return i
	default:
		return 0
	}
}

// GetSalesSummary handles GET /shops/{shop_id}/analytics/sales-summary
// @Summary      Get sales summary
// @Description  Returns total sales, total orders, and average order value for a shop in a given period
// @Tags         analytics
// @Produce      json
// @Param        shop_id path int true "Shop ID"
// @Param        period query string false "Period (today, week, month, custom)" Enums(today, week, month, custom) default(today)
// @Param        start_date query string false "Start date (YYYY-MM-DD, required if period=custom)"
// @Param        end_date query string false "End date (YYYY-MM-DD, required if period=custom)"
// @Success      200 {object} models.SuccessResponse{data=map[string]interface{}} "Sales summary fetched successfully."
// @Failure      400 {object} models.ErrorResponse "Invalid parameters"
// @Failure      500 {object} models.ErrorResponse "Failed to fetch sales summary"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/analytics/sales-summary [get]
func (h *AnalyticsHandler) GetSalesSummary(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id")
	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, 400, "Invalid shop_id", nil)
	}

	period := c.Query("period", "today")
	var start, end, prevStart, prevEnd time.Time
	now := time.Now().UTC()
	switch period {
	case "today":
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		end = start.Add(24 * time.Hour)
		prevStart = start.AddDate(0, 0, -1)
		prevEnd = start
	case "week":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -weekday+1)
		end = start.AddDate(0, 0, 7)
		prevStart = start.AddDate(0, 0, -7)
		prevEnd = start
	case "month":
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		end = start.AddDate(0, 1, 0)
		prevStart = start.AddDate(0, -1, 0)
		prevEnd = start
	case "custom":
		startStr := c.Query("start_date")
		endStr := c.Query("end_date")
		start, err = time.Parse("2006-01-02", startStr)
		if err != nil {
			return api.ErrorResponse(c, 400, "Invalid start_date", nil)
		}
		end, err = time.Parse("2006-01-02", endStr)
		if err != nil {
			return api.ErrorResponse(c, 400, "Invalid end_date", nil)
		}
		end = end.Add(24 * time.Hour)
		days := int(end.Sub(start).Hours() / 24)
		prevStart = start.AddDate(0, 0, -days)
		prevEnd = start
	default:
		return api.ErrorResponse(c, 400, "Invalid period", nil)
	}

	params := db.GetSalesSummaryParams{
		ShopID: shopID,
		Day:    pgtype.Date{Time: start, Valid: true},
		Day_2:  pgtype.Date{Time: end, Valid: true},
	}
	row, err := h.Repository.GetSalesSummary(c.Context(), params)
	if err != nil {
		zap.L().Error("GetSalesSummary: failed to fetch sales summary", zap.Error(err), zap.Int64("shop_id", shopID))
		return api.ErrorResponse(c, 500, "Failed to fetch sales summary", nil)
	}

	prevParams := db.GetSalesSummaryParams{
		ShopID: shopID,
		Day:    pgtype.Date{Time: prevStart, Valid: true},
		Day_2:  pgtype.Date{Time: prevEnd, Valid: true},
	}
	prevRow, err := h.Repository.GetSalesSummary(c.Context(), prevParams)
	if err != nil {
		zap.L().Error("GetSalesSummary: failed to fetch previous sales summary", zap.Error(err), zap.Int64("shop_id", shopID))
		return api.ErrorResponse(c, 500, "Failed to fetch previous sales summary", nil)
	}

	// Customers count
	customers, err := h.Repository.GetCustomersCount(c.Context(), shopID)
	if err != nil {
		zap.L().Error("GetSalesSummary: failed to fetch customers count", zap.Error(err), zap.Int64("shop_id", shopID))
		return api.ErrorResponse(c, 500, "Failed to fetch customers count", nil)
	}
	prevCustomers := customers // For now, as a placeholder (ideally, count for previous period)

	// Percentage change helpers
	percentChange := func(current, prev float64) float64 {
		if prev == 0 {
			if current == 0 {
				return 0
			}
			return 100
		}
		return ((current - prev) / prev) * 100
	}

	// Convert interface{} to float64
	toF := func(val interface{}) float64 {
		switch v := val.(type) {
		case int64:
			return float64(v)
		case int32:
			return float64(v)
		case float64:
			return v
		case float32:
			return float64(v)
		case []uint8:
			f, _ := strconv.ParseFloat(string(v), 64)
			return f
		default:
			return 0
		}
	}

	totalSales := toF(row.TotalSales)
	prevTotalSales := toF(prevRow.TotalSales)
	totalOrders := toF(row.TotalOrders)
	prevTotalOrders := toF(prevRow.TotalOrders)
	// For customers, using current and previous as same for now

	resp := map[string]interface{}{
		"total_sales":          totalSales,
		"total_orders":         totalOrders,
		"average_order_value":  toF(row.AverageOrderValue),
		"sales_change_pct":     percentChange(totalSales, prevTotalSales),
		"orders_change_pct":    percentChange(totalOrders, prevTotalOrders),
		"customers":            customers,
		"customers_change_pct": percentChange(float64(customers), float64(prevCustomers)),
	}
	return api.SuccessResponse(c, 200, resp, "Sales summary fetched successfully.")
}

// GetOrdersOverTime handles GET /shops/{shop_id}/analytics/orders-over-time
// @Summary      Get orders over time
// @Description  Returns order counts grouped by interval (day, week, month) for a shop in a given period
// @Tags         analytics
// @Produce      json
// @Param        shop_id path int true "Shop ID"
// @Param        interval query string false "Interval (day, week, month)" Enums(day, week, month) default(day)
// @Param        period query string false "Period (today, week, month, custom)" Enums(today, week, month, custom) default(today)
// @Param        start_date query string false "Start date (YYYY-MM-DD, required if period=custom)"
// @Param        end_date query string false "End date (YYYY-MM-DD, required if period=custom)"
// @Success      200 {object} models.SuccessResponse{data=map[string]interface{}} "Orders over time fetched successfully."
// @Failure      400 {object} models.ErrorResponse "Invalid parameters"
// @Failure      500 {object} models.ErrorResponse "Failed to fetch orders over time"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/analytics/orders-over-time [get]
func (h *AnalyticsHandler) GetOrdersOverTime(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id")
	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, 400, "Invalid shop_id", nil)
	}

	interval := c.Query("interval", "day")
	switch interval {
	case "day", "week", "month":
		// ok
	default:
		return api.ErrorResponse(c, 400, "Invalid interval", nil)
	}

	period := c.Query("period", "today")
	var start, end time.Time
	now := time.Now().UTC()
	switch period {
	case "today":
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		end = start.Add(24 * time.Hour)
	case "week":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -weekday+1)
		end = start.AddDate(0, 0, 7)
	case "month":
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		end = start.AddDate(0, 1, 0)
	case "custom":
		startStr := c.Query("start_date")
		endStr := c.Query("end_date")
		start, err = time.Parse("2006-01-02", startStr)
		if err != nil {
			return api.ErrorResponse(c, 400, "Invalid start_date", nil)
		}
		end, err = time.Parse("2006-01-02", endStr)
		if err != nil {
			return api.ErrorResponse(c, 400, "Invalid end_date", nil)
		}
		end = end.Add(24 * time.Hour)
	default:
		return api.ErrorResponse(c, 400, "Invalid period", nil)
	}

	params := db.GetOrdersOverTimeParams{
		DateTrunc: interval,
		ShopID:    shopID,
		Day:       pgtype.Date{Time: start, Valid: true},
		Day_2:     pgtype.Date{Time: end, Valid: true},
	}
	rows, err := h.Repository.GetOrdersOverTime(c.Context(), params)
	if err != nil {
		return api.ErrorResponse(c, 500, "Failed to fetch orders over time", nil)
	}

	// Build a map from period string to order count
	orderMap := make(map[string]int64)
	for _, row := range rows {
		var label string
		if row.Period.Valid {
			t := row.Period.Time
			switch interval {
			case "day":
				label = t.Format("2006-01-02")
			case "month":
				label = t.Format("2006-01")
			case "week":
				_, week := t.ISOWeek()
				label = t.Format("2006") + "-W" + fmt.Sprintf("%02d", week)
			}
		}
		orderCount, _ := row.OrderCount.(int64)
		orderMap[label] = orderCount
	}

	labels := []string{}
	orders := []int64{}

	// Generate all intervals between start and end
	cur := start
	switch interval {
	case "day":
		for !cur.After(end.Add(-24 * time.Hour)) {
			label := cur.Format("2006-01-02")
			labels = append(labels, label)
			if v, ok := orderMap[label]; ok {
				orders = append(orders, v)
			} else {
				orders = append(orders, 0)
			}
			cur = cur.Add(24 * time.Hour)
		}
	case "month":
		for cur.Before(end) {
			label := cur.Format("2006-01")
			labels = append(labels, label)
			if v, ok := orderMap[label]; ok {
				orders = append(orders, v)
			} else {
				orders = append(orders, 0)
			}
			cur = cur.AddDate(0, 1, 0)
		}
	case "week":
		for cur.Before(end) {
			year, week := cur.ISOWeek()
			label := fmt.Sprintf("%d-W%02d", year, week)
			labels = append(labels, label)
			if v, ok := orderMap[label]; ok {
				orders = append(orders, v)
			} else {
				orders = append(orders, 0)
			}
			cur = cur.AddDate(0, 0, 7)
		}
	}

	resp := map[string]interface{}{
		"labels": labels,
		"orders": orders,
	}
	return api.SuccessResponse(c, 200, resp, "Orders over time fetched successfully.")
}

// GetTopProducts handles GET /shops/{shop_id}/analytics/top-products
// @Summary      Get top products
// @Description  Returns the top selling products for a shop in a given period
// @Tags         analytics
// @Produce      json
// @Param        shop_id path int true "Shop ID"
// @Param        period query string false "Period (today, week, month, custom)" Enums(today, week, month, custom) default(today)
// @Param        limit query int false "Number of products to return" default(5)
// @Param        start_date query string false "Start date (YYYY-MM-DD, required if period=custom)"
// @Param        end_date query string false "End date (YYYY-MM-DD, required if period=custom)"
// @Success      200 {object} models.SuccessResponse{data=[]map[string]interface{}} "Top products fetched successfully."
// @Failure      400 {object} models.ErrorResponse "Invalid parameters"
// @Failure      500 {object} models.ErrorResponse "Failed to fetch top products"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/analytics/top-products [get]
func (h *AnalyticsHandler) GetTopProducts(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id")
	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, 400, "Invalid shop_id", nil)
	}

	period := c.Query("period", "today")
	limit, err := strconv.Atoi(c.Query("limit", "5"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 5
	}
	var start, end time.Time
	now := time.Now().UTC()
	switch period {
	case "today":
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		end = start.Add(24 * time.Hour)
	case "week":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -weekday+1)
		end = start.AddDate(0, 0, 7)
	case "month":
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		end = start.AddDate(0, 1, 0)
	case "custom":
		startStr := c.Query("start_date")
		endStr := c.Query("end_date")
		start, err = time.Parse("2006-01-02", startStr)
		if err != nil {
			return api.ErrorResponse(c, 400, "Invalid start_date", nil)
		}
		end, err = time.Parse("2006-01-02", endStr)
		if err != nil {
			return api.ErrorResponse(c, 400, "Invalid end_date", nil)
		}
		end = end.Add(24 * time.Hour)
	default:
		return api.ErrorResponse(c, 400, "Invalid period", nil)
	}

	params := db.GetTopProductsParams{
		ShopID:      shopID,
		CreatedAt:   pgtype.Timestamptz{Time: start, Valid: true},
		CreatedAt_2: pgtype.Timestamptz{Time: end, Valid: true},
		Limit:       int32(limit),
	}
	rows, err := h.Repository.GetTopProducts(c.Context(), params)
	if err != nil {
		return api.ErrorResponse(c, 500, "Failed to fetch top products", nil)
	}
	products := make([]map[string]interface{}, 0, len(rows))
	for _, row := range rows {
		unitsSold := toInt64(row.UnitsSold)
		revenue := toInt64(row.Revenue)
		products = append(products, map[string]interface{}{
			"product_variation_id": row.ProductVariationID,
			"product_name":         row.ProductName,
			"units_sold":           unitsSold,
			"revenue":              revenue,
		})
	}
	return api.SuccessResponse(c, 200, products, "Top products fetched successfully.")
}

// GetCustomerSummary handles GET /shops/{shop_id}/analytics/customers-summary
// @Summary      Get customer summary
// @Description  Returns new customers, returning customers, and top customers for a shop in a given period
// @Tags         analytics
// @Produce      json
// @Param        shop_id path int true "Shop ID"
// @Param        period query string false "Period (today, week, month, custom)" Enums(today, week, month, custom) default(today)
// @Param        start_date query string false "Start date (YYYY-MM-DD, required if period=custom)"
// @Param        end_date query string false "End date (YYYY-MM-DD, required if period=custom)"
// @Success      200 {object} models.SuccessResponse{data=map[string]interface{}} "Customer summary fetched successfully."
// @Failure      400 {object} models.ErrorResponse "Invalid parameters"
// @Failure      500 {object} models.ErrorResponse "Failed to fetch customer summary"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/analytics/customers-summary [get]
func (h *AnalyticsHandler) GetCustomerSummary(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id")
	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, 400, "Invalid shop_id", nil)
	}

	period := c.Query("period", "today")
	var start, end time.Time
	now := time.Now().UTC()
	switch period {
	case "today":
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		end = start.Add(24 * time.Hour)
	case "week":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -weekday+1)
		end = start.AddDate(0, 0, 7)
	case "month":
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		end = start.AddDate(0, 1, 0)
	case "custom":
		startStr := c.Query("start_date")
		endStr := c.Query("end_date")
		start, err = time.Parse("2006-01-02", startStr)
		if err != nil {
			return api.ErrorResponse(c, 400, "Invalid start_date", nil)
		}
		end, err = time.Parse("2006-01-02", endStr)
		if err != nil {
			return api.ErrorResponse(c, 400, "Invalid end_date", nil)
		}
		end = end.Add(24 * time.Hour)
	default:
		return api.ErrorResponse(c, 400, "Invalid period", nil)
	}

	params := db.GetCustomerSummaryNewReturningParams{
		ShopID:      shopID,
		CreatedAt:   pgtype.Timestamptz{Time: start, Valid: true},
		CreatedAt_2: pgtype.Timestamptz{Time: end, Valid: true},
	}
	summary, err := h.Repository.GetCustomerSummaryNewReturning(c.Context(), params)
	if err != nil {
		return api.ErrorResponse(c, 500, "Failed to fetch customer summary", nil)
	}
	topParams := db.GetCustomerSummaryTopParams{
		ShopID:      shopID,
		CreatedAt:   pgtype.Timestamptz{Time: start, Valid: true},
		CreatedAt_2: pgtype.Timestamptz{Time: end, Valid: true},
	}
	topRows, err := h.Repository.GetCustomerSummaryTop(c.Context(), topParams)
	if err != nil {
		return api.ErrorResponse(c, 500, "Failed to fetch top customers", nil)
	}
	topCustomers := make([]map[string]interface{}, 0, len(topRows))
	for _, row := range topRows {
		orders := toInt64(row.Orders)
		totalSpent := toInt64(row.TotalSpent)
		name := ""
		if s, ok := row.Name.(string); ok {
			name = s
		} else if b, ok := row.Name.([]uint8); ok {
			name = string(b)
		}
		topCustomers = append(topCustomers, map[string]interface{}{
			"customer_id":   row.CustomerEmail,
			"name":          name,
			"orders":        orders,
			"total_spent":   totalSpent,
			"is_registered": row.IsRegistered,
		})
	}
	resp := map[string]interface{}{
		"new_customers":       toInt64(summary.NewCustomers),
		"returning_customers": toInt64(summary.ReturningCustomers),
		"top_customers":       topCustomers,
	}
	return api.SuccessResponse(c, 200, resp, "Customer summary fetched successfully.")
}

// GetLowStockProducts handles GET /shops/{shop_id}/analytics/low-stock
// @Summary      Get low stock products
// @Description  Returns products with low stock for a shop
// @Tags         analytics
// @Produce      json
// @Param        shop_id path int true "Shop ID"
// @Param        threshold query int false "Stock threshold" default(5)
// @Success      200 {object} models.SuccessResponse{data=[]map[string]interface{}} "Low stock products fetched successfully."
// @Failure      400 {object} models.ErrorResponse "Invalid parameters"
// @Failure      500 {object} models.ErrorResponse "Failed to fetch low stock products"
// @Security     OAuth2AccessCode
// @Router       /shops/{shop_id}/analytics/low-stock [get]
func (h *AnalyticsHandler) GetLowStockProducts(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id")
	shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
	if err != nil {
		return api.ErrorResponse(c, 400, "Invalid shop_id", nil)
	}
	threshold, err := strconv.ParseInt(c.Query("threshold", "5"), 10, 64)
	if err != nil || threshold < 0 {
		threshold = 5
	}
	params := db.GetLowStockProductsParams{
		ShopID:            shopID,
		AvailableQuantity: threshold,
	}
	rows, err := h.Repository.GetLowStockProducts(c.Context(), params)
	if err != nil {
		return api.ErrorResponse(c, 500, "Failed to fetch low stock products", nil)
	}
	variants := make([]map[string]interface{}, 0, len(rows))
	for _, row := range rows {
		variants = append(variants, map[string]interface{}{
			"product_variation_id": row.ProductVariationID,
			"product_name":         row.ProductName,
			"sku":                  row.Sku,
			"description":          row.Description,
			"stock":                row.Stock,
		})
	}
	return api.SuccessResponse(c, 200, variants, "Low stock products fetched successfully.")
}
