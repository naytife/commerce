package db

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"go.uber.org/zap"
)

type repoSvc struct {
	*Queries
	db *pgxpool.Pool
}

func (r *repoSvc) WithTx(ctx context.Context, txFn func(*Queries) error) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	q := New(tx)
	err = txFn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx failed: %w, unable to rollback: %w", err, rbErr)
		}
	} else {
		err = tx.Commit(ctx)
	}
	return err
}

// PgConn returns the underlying pgxpool.Pool for direct database access
func (r *repoSvc) PgConn() *pgxpool.Pool {
	return r.db
}

type Repository interface {
	WithTx(ctx context.Context, txFn func(*Queries) error) error
	SetShopIDInSession(ctx context.Context, shopID int64) error
	PgConn() *pgxpool.Pool
	// USER
	UpsertUser(ctx context.Context, arg UpsertUserParams) (User, error)
	GetUser(ctx context.Context, email *string) (User, error)
	GetUserById(ctx context.Context, userID uuid.UUID) (User, error)
	GetUserBySub(ctx context.Context, sub *string) (User, error)
	GetUserBySubWithShops(ctx context.Context, sub *string) (GetUserBySubWithShopsRow, error)
	UpsertCustomer(ctx context.Context, arg UpsertCustomerParams) (ShopCustomer, error)
	GetCustomerByEmail(ctx context.Context, arg GetCustomerByEmailParams) (ShopCustomer, error)
	// Customer Management
	GetCustomers(ctx context.Context, arg GetCustomersParams) ([]ShopCustomer, error)
	GetCustomersCount(ctx context.Context, shopID int64) (int64, error)
	SearchCustomers(ctx context.Context, arg SearchCustomersParams) ([]ShopCustomer, error)
	GetCustomerById(ctx context.Context, arg GetCustomerByIdParams) (ShopCustomer, error)
	UpdateCustomer(ctx context.Context, arg UpdateCustomerParams) (ShopCustomer, error)
	DeleteCustomer(ctx context.Context, arg DeleteCustomerParams) error
	GetCustomerOrders(ctx context.Context, arg GetCustomerOrdersParams) ([]GetCustomerOrdersRow, error)
	// Inventory Management
	GetLowStockVariants(ctx context.Context, arg GetLowStockVariantsParams) ([]GetLowStockVariantsRow, error)
	GetProductVariation(ctx context.Context, arg GetProductVariationParams) (ProductVariation, error)
	UpdateVariantStock(ctx context.Context, arg UpdateVariantStockParams) (ProductVariation, error)
	DeductVariantStock(ctx context.Context, arg DeductVariantStockParams) (ProductVariation, error)
	AddVariantStock(ctx context.Context, arg AddVariantStockParams) (ProductVariation, error)
	GetInventoryReport(ctx context.Context, arg GetInventoryReportParams) ([]GetInventoryReportRow, error)
	GetStockMovements(ctx context.Context, arg GetStockMovementsParams) ([]GetStockMovementsRow, error)
	CreateStockMovement(ctx context.Context, arg CreateStockMovementParams) (StockMovement, error)
	// SHOP
	CreateShop(ctx context.Context, shopArg CreateShopParams) (Shop, error)
	GetShop(ctx context.Context, shopID int64) (Shop, error)
	DeleteShop(ctx context.Context, shopID int64) error
	UpdateShop(ctx context.Context, arg UpdateShopParams) (Shop, error)
	GetShopsByOwner(ctx context.Context, ownerID uuid.UUID) ([]Shop, error)
	GetShopBySubDomain(ctx context.Context, defaultSubDomain string) (Shop, error)
	GetShopIDBySubDomain(ctx context.Context, subDomain string) (int64, error)
	GetShopImages(ctx context.Context, shopID int64) (ShopImage, error)
	// PRODUCT-TYPE
	CreateProductType(ctx context.Context, arg CreateProductTypeParams) (ProductType, error)
	DeleteProductType(ctx context.Context, arg DeleteProductTypeParams) (ProductType, error)
	GetProductType(ctx context.Context, arg GetProductTypeParams) (ProductType, error)
	GetProductTypes(ctx context.Context, shopID int64) ([]ProductType, error)
	UpdateProductType(ctx context.Context, arg UpdateProductTypeParams) (ProductType, error)
	// ATTRIBUTE
	CreateAttribute(ctx context.Context, arg CreateAttributeParams) (Attribute, error)
	DeleteAttribute(ctx context.Context, arg DeleteAttributeParams) error
	GetAttribute(ctx context.Context, arg GetAttributeParams) (GetAttributeRow, error)
	GetAttributes(ctx context.Context, arg GetAttributesParams) ([]GetAttributesRow, error)
	UpdateAttribute(ctx context.Context, arg UpdateAttributeParams) (Attribute, error)
	GetProductsAttributes(ctx context.Context, arg GetProductsAttributesParams) ([]GetProductsAttributesRow, error)
	GetVariationsAttributes(ctx context.Context, arg GetVariationsAttributesParams) ([]GetVariationsAttributesRow, error)
	// ATTRIBUTE-OPTION
	BatchUpsertAttributeOption(ctx context.Context, arg []BatchUpsertAttributeOptionParams) *BatchUpsertAttributeOptionBatchResults
	BatchDeleteAttributeOptions(ctx context.Context, arg []BatchDeleteAttributeOptionsParams) *BatchDeleteAttributeOptionsBatchResults
	GetAttributeOptions(ctx context.Context, arg GetAttributeOptionsParams) ([]AttributeOption, error)
	UpdateAttributeOption(ctx context.Context, arg UpdateAttributeOptionParams) (AttributeOption, error)
	// ATTRIBUTE-VALUE
	GetProductAttributeValues(ctx context.Context, arg GetProductAttributeValuesParams) ([]GetProductAttributeValuesRow, error)
	// CATEGORY
	CreateCategory(ctx context.Context, arg CreateCategoryParams) (Category, error)
	GetCategory(ctx context.Context, arg GetCategoryParams) (GetCategoryRow, error)
	UpdateCategory(ctx context.Context, arg UpdateCategoryParams) (Category, error)
	GetCategoryChildren(ctx context.Context, arg GetCategoryChildrenParams) ([]GetCategoryChildrenRow, error)
	GetCategories(ctx context.Context, arg GetCategoriesParams) ([]GetCategoriesRow, error)
	// CreateCategoryAttribute(ctx context.Context, arg CreateCategoryAttributeParams) ([]byte, error)
	// DeleteCategoryAttribute(ctx context.Context, arg DeleteCategoryAttributeParams) ([]byte, error)
	// GetCategoryAttributes(ctx context.Context, categoryID int64) ([]byte, error)
	// PRODUCT
	CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error)
	GetProducts(ctx context.Context, arg GetProductsParams) ([]GetProductsRow, error)
	GetProduct(ctx context.Context, arg GetProductParams) (GetProductRow, error)
	GetProductById(ctx context.Context, arg GetProductByIdParams) (Product, error)
	DeleteProduct(ctx context.Context, arg DeleteProductParams) error
	UpdateProduct(ctx context.Context, arg UpdateProductParams) error
	GetProductsByType(ctx context.Context, arg GetProductsByTypeParams) ([]GetProductsByTypeRow, error)
	GetProductsByCategory(ctx context.Context, arg GetProductsByCategoryParams) ([]GetProductsByCategoryRow, error)
	// GetProductAllowedAttributes(ctx context.Context, productID int64) ([]byte, error)
	// CreateProductAllowedAttribute(ctx context.Context, arg CreateProductAllowedAttributeParams) ([]byte, error)
	// DeleteProductAllowedAttribute(ctx context.Context, arg DeleteProductAllowedAttributeParams) ([]byte, error)
	UpsertProductVariants(ctx context.Context, arg []UpsertProductVariantsParams) *UpsertProductVariantsBatchResults
	// PRODUCT IMAGES
	CreateProductImage(ctx context.Context, arg CreateProductImageParams) (ProductImage, error)
	GetProductImages(ctx context.Context, arg GetProductImagesParams) ([]ProductImage, error)
	DeleteProductImage(ctx context.Context, arg DeleteProductImageParams) error
	DeleteAllProductImages(ctx context.Context, arg DeleteAllProductImagesParams) error
	// ORDER
	CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error)
	GetOrder(ctx context.Context, arg GetOrderParams) (Order, error)
	ListOrders(ctx context.Context, arg ListOrdersParams) ([]Order, error)
	UpdateOrder(ctx context.Context, arg UpdateOrderParams) error
	DeleteOrder(ctx context.Context, arg DeleteOrderParams) error
	// ORDER ITEMS
	CreateOrderItem(ctx context.Context, arg CreateOrderItemParams) (OrderItem, error)
	GetOrderItemsByOrder(ctx context.Context, arg GetOrderItemsByOrderParams) ([]OrderItem, error)
	UpdateOrderItem(ctx context.Context, arg UpdateOrderItemParams) error
	DeleteOrderItem(ctx context.Context, arg DeleteOrderItemParams) error
	DeleteOrderItemsByOrder(ctx context.Context, arg DeleteOrderItemsByOrderParams) error
	CountOrders(ctx context.Context, shopID int64) (int64, error)
	// PAYMENT STATUS MANAGEMENT
	UpdateOrderPaymentStatus(ctx context.Context, arg UpdateOrderPaymentStatusParams) (Order, error)
	GetOrderByTransactionID(ctx context.Context, arg GetOrderByTransactionIDParams) (Order, error)
	UpdateOrderStatusByTransactionID(ctx context.Context, arg UpdateOrderStatusByTransactionIDParams) (Order, error)
	// SHOP PAYMENT METHODS
	GetShopPaymentMethods(ctx context.Context, shopID int64) ([]ShopPaymentMethod, error)
	GetShopPaymentMethod(ctx context.Context, arg GetShopPaymentMethodParams) (ShopPaymentMethod, error)
	UpsertShopPaymentMethod(ctx context.Context, arg UpsertShopPaymentMethodParams) (ShopPaymentMethod, error)
	UpdateShopPaymentMethodStatus(ctx context.Context, arg UpdateShopPaymentMethodStatusParams) (ShopPaymentMethod, error)
	DeleteShopPaymentMethod(ctx context.Context, arg DeleteShopPaymentMethodParams) error
	// DEPLOYMENT TRACKING
	CreateDeployment(ctx context.Context, arg CreateDeploymentParams) (ShopDeployment, error)
	UpdateDeploymentStatus(ctx context.Context, arg UpdateDeploymentStatusParams) error
	CompleteDeployment(ctx context.Context, arg CompleteDeploymentParams) error
	GetDeploymentByID(ctx context.Context, deploymentID int64) (ShopDeployment, error)
	GetLatestDeploymentByShop(ctx context.Context, shopID int64) (ShopDeployment, error)
	GetDeploymentsByShop(ctx context.Context, arg GetDeploymentsByShopParams) ([]ShopDeployment, error)
	GetShopCurrentTemplate(ctx context.Context, shopID int64) (GetShopCurrentTemplateRow, error)
	IsShopDeployed(ctx context.Context, shopID int64) (bool, error)
	UpdateShopLastDeployment(ctx context.Context, arg UpdateShopLastDeploymentParams) error
	// DATA UPDATE TRACKING
	CreateDataUpdate(ctx context.Context, arg CreateDataUpdateParams) (ShopDataUpdate, error)
	UpdateDataUpdateStatus(ctx context.Context, arg UpdateDataUpdateStatusParams) error
	CompleteDataUpdate(ctx context.Context, arg CompleteDataUpdateParams) error
	GetDataUpdateByID(ctx context.Context, updateID int64) (ShopDataUpdate, error)
	GetLatestDataUpdateByShop(ctx context.Context, shopID int64) (ShopDataUpdate, error)
	UpdateShopLastDataUpdate(ctx context.Context, shopID int64) error
	// DEPLOYMENT URL TRACKING
	CreateDeploymentURL(ctx context.Context, arg CreateDeploymentURLParams) (ShopDeploymentUrl, error)
	GetDeploymentURLs(ctx context.Context, deploymentID int64) ([]ShopDeploymentUrl, error)
	GetSalesSummary(ctx context.Context, arg GetSalesSummaryParams) (GetSalesSummaryRow, error)
	GetOrdersOverTime(ctx context.Context, arg GetOrdersOverTimeParams) ([]GetOrdersOverTimeRow, error)
	GetTopProducts(ctx context.Context, arg GetTopProductsParams) ([]GetTopProductsRow, error)
	GetCustomerSummaryNewReturning(ctx context.Context, arg GetCustomerSummaryNewReturningParams) (GetCustomerSummaryNewReturningRow, error)
	GetCustomerSummaryTop(ctx context.Context, arg GetCustomerSummaryTopParams) ([]GetCustomerSummaryTopRow, error)
	GetLowStockProducts(ctx context.Context, arg GetLowStockProductsParams) ([]GetLowStockProductsRow, error)
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repoSvc{
		Queries: New(db),
		db:      db,
	}
}

func InitDB(dataSourceName string) (*pgxpool.Pool, error) {
	// Wrap the Zap logger in a TraceLog object using the global zap logger
	traceLogger := &tracelog.TraceLog{
		Logger: tracelog.LoggerFunc(func(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
			// Convert map[string]interface{} into Zap fields
			fields := make([]zap.Field, 0, len(data))
			for k, v := range data {
				fields = append(fields, zap.Any(k, v))
			}

			switch level {
			case tracelog.LogLevelError:
				zap.L().Error(msg, fields...)
			case tracelog.LogLevelWarn:
				zap.L().Warn(msg, fields...)
			case tracelog.LogLevelInfo:
				zap.L().Info(msg, fields...)
			case tracelog.LogLevelDebug:
				zap.L().Debug(msg, fields...)
			}
		}),
		LogLevel: tracelog.LogLevelDebug, // Set log level to Debug
	}

	// Parse the pool config from the data source name
	config, err := pgxpool.ParseConfig(dataSourceName)
	if err != nil {
		return nil, err
	}

	// Attach the tracer to the config
	config.ConnConfig.Tracer = traceLogger

	// Set pool settings - optimized for e-commerce workload
	config.MaxConns = 25                       // Increase for better concurrency
	config.MinConns = 5                        // Keep more connections warm
	config.MaxConnLifetime = 30 * time.Minute  // Prevent stale connections
	config.MaxConnIdleTime = 10 * time.Minute  // Longer idle time for efficiency
	config.HealthCheckPeriod = 1 * time.Minute // Regular health checks

	// Create the connection pool
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func (r *repoSvc) SetShopIDInSession(ctx context.Context, shopID int64) error {
	// Assuming you want to execute a SQL command to set shop_id for the current session
	query := fmt.Sprintf("SET commerce.current_shop_id = %d", shopID)
	_, err := r.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to set shop_id in session: %w", err)
	}
	return nil
}

func (r *repoSvc) CreateShop(ctx context.Context, shopArg CreateShopParams) (Shop, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	shop := Shop{}
	err := r.WithTx(ctx, func(q *Queries) error {
		var err error
		shop, err = q.CreateShop(ctx, shopArg)
		return err
	})
	return shop, err
}

func (r *repoSvc) UpdateShop(ctx context.Context, arg UpdateShopParams) (Shop, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	shop := Shop{}
	err := r.WithTx(ctx, func(q *Queries) error {
		var err error
		shop, err = q.UpdateShop(ctx, arg)
		return err
	})
	return shop, err
}

func (r *repoSvc) CreateCategory(ctx context.Context, arg CreateCategoryParams) (Category, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	category := Category{}
	err := r.WithTx(ctx, func(q *Queries) error {
		var err error
		category, err = q.CreateCategory(ctx, arg)
		return err
	})
	return category, err
}

func (r *repoSvc) GetSalesSummary(ctx context.Context, arg GetSalesSummaryParams) (GetSalesSummaryRow, error) {
	return r.Queries.GetSalesSummary(ctx, arg)
}

func (r *repoSvc) GetOrdersOverTime(ctx context.Context, arg GetOrdersOverTimeParams) ([]GetOrdersOverTimeRow, error) {
	return r.Queries.GetOrdersOverTime(ctx, arg)
}

func (r *repoSvc) GetTopProducts(ctx context.Context, arg GetTopProductsParams) ([]GetTopProductsRow, error) {
	return r.Queries.GetTopProducts(ctx, arg)
}

func (r *repoSvc) GetCustomerSummaryNewReturning(ctx context.Context, arg GetCustomerSummaryNewReturningParams) (GetCustomerSummaryNewReturningRow, error) {
	return r.Queries.GetCustomerSummaryNewReturning(ctx, arg)
}

func (r *repoSvc) GetCustomerSummaryTop(ctx context.Context, arg GetCustomerSummaryTopParams) ([]GetCustomerSummaryTopRow, error) {
	return r.Queries.GetCustomerSummaryTop(ctx, arg)
}

func (r *repoSvc) GetLowStockProducts(ctx context.Context, arg GetLowStockProductsParams) ([]GetLowStockProductsRow, error) {
	return r.Queries.GetLowStockProducts(ctx, arg)
}
