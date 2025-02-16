package db

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/sirupsen/logrus"
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

type Repository interface {
	WithTx(ctx context.Context, txFn func(*Queries) error) error
	SetShopIDInSession(ctx context.Context, shopID int64) error
	// USER
	UpsertUser(ctx context.Context, arg UpsertUserParams) (User, error)
	GetUser(ctx context.Context, email *string) (User, error)
	GetUserById(ctx context.Context, userID uuid.UUID) (User, error)
	GetUserBySub(ctx context.Context, sub *string) (User, error)
	// SHOP
	CreateShop(ctx context.Context, shopArg CreateShopParams) (Shop, error)
	GetShop(ctx context.Context, shopID int64) (Shop, error)
	DeleteShop(ctx context.Context, shopID int64) error
	UpdateShop(ctx context.Context, arg UpdateShopParams) (Shop, error)
	GetShopsByOwner(ctx context.Context, ownerID uuid.UUID) ([]Shop, error)
	GetShopByDomain(ctx context.Context, defaultDomain string) (Shop, error)
	GetShopIDByDomain(ctx context.Context, domain string) (int64, error)
	GetShopImages(ctx context.Context, shopID int64) (ShopImage, error)
	// PRODUCT-TYPE
	CreateProductType(ctx context.Context, arg CreateProductTypeParams) (ProductType, error)
	DeleteProductType(ctx context.Context, arg DeleteProductTypeParams) (ProductType, error)
	GetProductType(ctx context.Context, arg GetProductTypeParams) (ProductType, error)
	GetProductTypes(ctx context.Context, shopID int64) ([]ProductType, error)
	UpdateProductType(ctx context.Context, arg UpdateProductTypeParams) (ProductType, error)
	// ATTRIBUTE
	CreateAttribute(ctx context.Context, arg CreateAttributeParams) (Attribute, error)
	DeleteAttribute(ctx context.Context, arg DeleteAttributeParams) (Attribute, error)
	GetAttribute(ctx context.Context, arg GetAttributeParams) (Attribute, error)
	GetAttributes(ctx context.Context, arg GetAttributesParams) ([]Attribute, error)
	UpdateAttribute(ctx context.Context, arg UpdateAttributeParams) (Attribute, error)
	GetProductsAttributes(ctx context.Context, arg GetProductsAttributesParams) ([]Attribute, error)
	GetVariationsAttributes(ctx context.Context, arg GetVariationsAttributesParams) ([]Attribute, error)
	// ATTRIBUTE-OPTION
	CreateAttributeOption(ctx context.Context, arg CreateAttributeOptionParams) (AttributeOption, error)
	DeleteAttributeOption(ctx context.Context, arg DeleteAttributeOptionParams) (AttributeOption, error)
	GetAttributeOption(ctx context.Context, arg GetAttributeOptionParams) (AttributeOption, error)
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
	DeleteProduct(ctx context.Context, arg DeleteProductParams) error
	UpdateProduct(ctx context.Context, arg UpdateProductParams) error
	GetProductsByType(ctx context.Context, arg GetProductsByTypeParams) ([]GetProductsByTypeRow, error)
	GetProductsByCategory(ctx context.Context, arg GetProductsByCategoryParams) ([]GetProductsByCategoryRow, error)
	// GetProductAllowedAttributes(ctx context.Context, productID int64) ([]byte, error)
	// CreateProductAllowedAttribute(ctx context.Context, arg CreateProductAllowedAttributeParams) ([]byte, error)
	// DeleteProductAllowedAttribute(ctx context.Context, arg DeleteProductAllowedAttributeParams) ([]byte, error)
	UpsertProductVariants(ctx context.Context, arg []UpsertProductVariantsParams) *UpsertProductVariantsBatchResults
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repoSvc{
		Queries: New(db),
		db:      db,
	}
}

func InitDB(dataSourceName string) (*pgxpool.Pool, error) {
	// Initialize a new logger (using Logrus)
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel) // Set log level to Debug for detailed SQL logging

	// Wrap the Logrus logger in a TraceLog object
	traceLogger := &tracelog.TraceLog{
		Logger: tracelog.LoggerFunc(func(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
			switch level {
			case tracelog.LogLevelError:
				logger.WithFields(logrus.Fields(data)).Error(msg)
			case tracelog.LogLevelWarn:
				logger.WithFields(logrus.Fields(data)).Warn(msg)
			case tracelog.LogLevelInfo:
				logger.WithFields(logrus.Fields(data)).Info(msg)
			case tracelog.LogLevelDebug:
				logger.WithFields(logrus.Fields(data)).Debug(msg)
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

	// Set pool settings
	config.MaxConns = 1
	config.MinConns = 1
	config.MaxConnIdleTime = 5 * time.Minute

	// Create the connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
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
