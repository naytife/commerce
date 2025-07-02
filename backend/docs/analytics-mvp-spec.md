# MVP Analytics Specification (Updated)

## Endpoint & Response Design Principles
- **Base Path:** `/v1/shops/{shop_id}/analytics/...`
- **Response Model:**
  ```json
  {
    "success": true,
    "data": { ... },
    "message": "Optional message"
  }
  ```

---

## 1. Sales Overview

### 1.1. Endpoint: Get Sales Summary
- **URL:** `GET /v1/shops/{shop_id}/analytics/sales-summary?period={period}`
- **Period:** `today`, `week`, `month`, `custom` (with `start_date`, `end_date`)
- **Response:**
  ```json
  {
    "success": true,
    "data": {
      "total_sales": 12345.67,
      "total_orders": 120,
      "average_order_value": 102.88
    },
    "message": "Sales summary fetched successfully."
  }
  ```
- **Handler:** `backend/internal/api/handlers/analytics.handlers.go` → `GetSalesSummary(w http.ResponseWriter, r *http.Request)`
- **SQL:**
  ```sql
  SELECT
    SUM(total_amount) AS total_sales,
    COUNT(*) AS total_orders,
    COALESCE(SUM(total_amount)/NULLIF(COUNT(*),0),0) AS average_order_value
  FROM orders
  WHERE shop_id = $1 AND created_at BETWEEN $2 AND $3 AND status = 'completed';
  ```
- **Dashboard:** `dashboard/src/routes/[shop]/analytics/+page.svelte` (Three cards: Total Sales, Total Orders, Avg. Order Value; Period selector)

---

## 2. Orders Over Time

### 2.1. Endpoint: Get Orders Over Time
- **URL:** `GET /v1/shops/{shop_id}/analytics/orders-over-time?interval={interval}&period={period}`
- **Interval:** `day`, `week`, `month`
- **Response:**
  ```json
  {
    "success": true,
    "data": {
      "labels": ["2024-06-01", "2024-06-02", ...],
      "orders": [5, 8, 3, ...]
    },
    "message": "Orders over time fetched successfully."
  }
  ```
- **Handler:** `backend/internal/api/handlers/analytics.handlers.go` → `GetOrdersOverTime(w http.ResponseWriter, r *http.Request)`
- **SQL:**
  ```sql
  SELECT
    DATE_TRUNC($1, created_at) AS period,
    COUNT(*) AS order_count
  FROM orders
  WHERE shop_id = $2 AND created_at BETWEEN $3 AND $4 AND status = 'completed'
  GROUP BY period
  ORDER BY period ASC;
  ```
- **Dashboard:** `dashboard/src/routes/[shop]/analytics/+page.svelte` (Line or bar chart: X-axis = date, Y-axis = order count)

---

## 3. Top Products

### 3.1. Endpoint: Get Top Products
- **URL:** `GET /v1/shops/{shop_id}/analytics/top-products?period={period}&limit=5`
- **Response:**
  ```json
  {
    "success": true,
    "data": [
      {
        "product_id": "abc123",
        "name": "Product Name",
        "units_sold": 42,
        "revenue": 1234.56
      }
    ],
    "message": "Top products fetched successfully."
  }
  ```
- **Handler:** `backend/internal/api/handlers/analytics.handlers.go` → `GetTopProducts(w http.ResponseWriter, r *http.Request)`
- **SQL:**
  ```sql
  SELECT
    oi.product_id,
    p.name,
    SUM(oi.quantity) AS units_sold,
    SUM(oi.quantity * oi.unit_price) AS revenue
  FROM order_items oi
  JOIN orders o ON oi.order_id = o.id
  JOIN products p ON oi.product_id = p.id
  WHERE o.shop_id = $1 AND o.created_at BETWEEN $2 AND $3 AND o.status = 'completed'
  GROUP BY oi.product_id, p.name
  ORDER BY units_sold DESC
  LIMIT $4;
  ```
- **Dashboard:** `dashboard/src/routes/[shop]/analytics/+page.svelte` (Table: Product Name | Units Sold | Revenue)

---

## 4. Customer Insights

### 4.1. Endpoint: Get Customer Summary
- **URL:** `GET /v1/shops/{shop_id}/analytics/customers-summary?period={period}`
- **Response:**
  ```json
  {
    "success": true,
    "data": {
      "new_customers": 10,
      "returning_customers": 5,
      "top_customers": [
        {
          "customer_id": "xyz789",
          "name": "Jane Doe",
          "orders": 3,
          "total_spent": 500.00
        }
      ]
    },
    "message": "Customer summary fetched successfully."
  }
  ```
- **Handler:** `backend/internal/api/handlers/analytics.handlers.go` → `GetCustomerSummary(w http.ResponseWriter, r *http.Request)`
- **SQL:**
  - **New/Returning:**
    ```sql
    SELECT
      COUNT(DISTINCT CASE WHEN c.created_at BETWEEN $2 AND $3 THEN c.id END) AS new_customers,
      COUNT(DISTINCT CASE WHEN c.created_at < $2 THEN c.id END) AS returning_customers
    FROM customers c
    JOIN orders o ON c.id = o.customer_id
    WHERE o.shop_id = $1 AND o.created_at BETWEEN $2 AND $3 AND o.status = 'completed';
    ```
  - **Top Customers:**
    ```sql
    SELECT
      c.id AS customer_id,
      c.name,
      COUNT(o.id) AS orders,
      SUM(o.total_amount) AS total_spent
    FROM customers c
    JOIN orders o ON c.id = o.customer_id
    WHERE o.shop_id = $1 AND o.created_at BETWEEN $2 AND $3 AND o.status = 'completed'
    GROUP BY c.id, c.name
    ORDER BY total_spent DESC
    LIMIT 5;
    ```
- **Dashboard:** `dashboard/src/routes/[shop]/analytics/+page.svelte` (Cards: New Customers, Returning Customers; Table: Top Customers)

---

## 5. Inventory Snapshots

### 5.1. Endpoint: Get Low Stock Products
- **URL:** `GET /v1/shops/{shop_id}/analytics/low-stock?threshold=5`
- **Response:**
  ```json
  {
    "success": true,
    "data": [
      {
        "product_id": "abc123",
        "name": "Product Name",
        "stock": 3
      }
    ],
    "message": "Low stock products fetched successfully."
  }
  ```
- **Handler:** `backend/internal/api/handlers/analytics.handlers.go` → `GetLowStockProducts(w http.ResponseWriter, r *http.Request)`
- **SQL:**
  ```sql
  SELECT id AS product_id, name, stock
  FROM products
  WHERE shop_id = $1 AND stock <= $2
  ORDER BY stock ASC;
  ```
- **Dashboard:** `dashboard/src/routes/[shop]/analytics/+page.svelte` (Table: Product Name | Stock)

---

# Backend Implementation Notes
- New handler file: `backend/internal/api/handlers/analytics.handlers.go` (create if not exists)
- Register new endpoints in `backend/internal/api/routes/analytics.go` (create if not exists)
- Use SQLC or your preferred query method.

# Dashboard Implementation Notes
- Page: `dashboard/src/routes/[shop]/analytics/+page.svelte` (create if not exists)
- Use `dashboard/src/lib/api.ts` to add fetchers for each endpoint.
- Use cards for KPIs, tables for lists, and simple charts (e.g., Chart.js or Svelte Simple Charts).
- Start with Sales Summary, then Orders Over Time, then Top Products, etc.

# Schema Changes (if needed)
- Ensure `orders` table has: `id`, `shop_id`, `customer_id`, `total_amount`, `created_at`, `status`
- Ensure `order_items` table has: `order_id`, `product_id`, `quantity`, `unit_price`
- Ensure `products` table has: `id`, `shop_id`, `name`, `stock`
- Ensure `customers` table has: `id`, `name`, `created_at`
- Ensure `carts` table has: `id`, `shop_id`, `created_at`

# Next Steps
1. Confirm which endpoints to start with.
2. Create handler and route stubs in backend.
3. Add SQL queries (using SQLC if possible).
4. Add fetchers and UI components in dashboard.
5. Iterate and expand as needed. 