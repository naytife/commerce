// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createCategory = `-- name: CreateCategory :one
INSERT INTO categories (slug, title, description, parent_id, shop_id, category_attributes)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING category_id, slug, title, description, parent_id, created_at, updated_at, shop_id, category_attributes
`

type CreateCategoryParams struct {
	Slug               string
	Title              string
	Description        pgtype.Text
	ParentID           pgtype.Int8
	ShopID             int64
	CategoryAttributes []byte
}

func (q *Queries) CreateCategory(ctx context.Context, arg CreateCategoryParams) (Category, error) {
	row := q.db.QueryRow(ctx, createCategory,
		arg.Slug,
		arg.Title,
		arg.Description,
		arg.ParentID,
		arg.ShopID,
		arg.CategoryAttributes,
	)
	var i Category
	err := row.Scan(
		&i.CategoryID,
		&i.Slug,
		&i.Title,
		&i.Description,
		&i.ParentID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ShopID,
		&i.CategoryAttributes,
	)
	return i, err
}

const createCategoryAttribute = `-- name: CreateCategoryAttribute :one
UPDATE categories
SET category_attributes = jsonb_set(
    COALESCE(category_attributes, '{}'), 
    ARRAY[UPPER($1)::text], 
    to_jsonb($2::text)
)
WHERE category_id = $3
RETURNING category_attributes
`

type CreateCategoryAttributeParams struct {
	Title      interface{}
	DataType   string
	CategoryID int64
}

func (q *Queries) CreateCategoryAttribute(ctx context.Context, arg CreateCategoryAttributeParams) ([]byte, error) {
	row := q.db.QueryRow(ctx, createCategoryAttribute, arg.Title, arg.DataType, arg.CategoryID)
	var category_attributes []byte
	err := row.Scan(&category_attributes)
	return category_attributes, err
}

const createProduct = `-- name: CreateProduct :one
INSERT INTO products ( title, description, category_id, shop_id, allowed_attributes, status)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING product_id, title, description, allowed_attributes, created_at, updated_at, status, category_id, shop_id
`

type CreateProductParams struct {
	Title             string
	Description       string
	CategoryID        int64
	ShopID            int64
	AllowedAttributes []byte
	Status            string
}

func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error) {
	row := q.db.QueryRow(ctx, createProduct,
		arg.Title,
		arg.Description,
		arg.CategoryID,
		arg.ShopID,
		arg.AllowedAttributes,
		arg.Status,
	)
	var i Product
	err := row.Scan(
		&i.ProductID,
		&i.Title,
		&i.Description,
		&i.AllowedAttributes,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Status,
		&i.CategoryID,
		&i.ShopID,
	)
	return i, err
}

const createShop = `-- name: CreateShop :one
INSERT INTO shops (owner_id, title, domain, favicon_url,logo_url,email, currency_code, about, status, address,phone_number, seo_description, seo_keywords, seo_title)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING shop_id, owner_id, title, domain, favicon_url, logo_url, email, currency_code, status, about, address, phone_number, seo_description, seo_keywords, seo_title, updated_at, created_at
`

type CreateShopParams struct {
	OwnerID        uuid.UUID
	Title          string
	Domain         string
	FaviconUrl     pgtype.Text
	LogoUrl        pgtype.Text
	Email          string
	CurrencyCode   string
	About          pgtype.Text
	Status         string
	Address        pgtype.Text
	PhoneNumber    pgtype.Text
	SeoDescription pgtype.Text
	SeoKeywords    []string
	SeoTitle       pgtype.Text
}

func (q *Queries) CreateShop(ctx context.Context, arg CreateShopParams) (Shop, error) {
	row := q.db.QueryRow(ctx, createShop,
		arg.OwnerID,
		arg.Title,
		arg.Domain,
		arg.FaviconUrl,
		arg.LogoUrl,
		arg.Email,
		arg.CurrencyCode,
		arg.About,
		arg.Status,
		arg.Address,
		arg.PhoneNumber,
		arg.SeoDescription,
		arg.SeoKeywords,
		arg.SeoTitle,
	)
	var i Shop
	err := row.Scan(
		&i.ShopID,
		&i.OwnerID,
		&i.Title,
		&i.Domain,
		&i.FaviconUrl,
		&i.LogoUrl,
		&i.Email,
		&i.CurrencyCode,
		&i.Status,
		&i.About,
		&i.Address,
		&i.PhoneNumber,
		&i.SeoDescription,
		&i.SeoKeywords,
		&i.SeoTitle,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const deleteCategoryAttribute = `-- name: DeleteCategoryAttribute :one
UPDATE categories
SET category_attributes = category_attributes - $1::text
WHERE category_id = $2
RETURNING category_attributes
`

type DeleteCategoryAttributeParams struct {
	Attribute  string
	CategoryID int64
}

func (q *Queries) DeleteCategoryAttribute(ctx context.Context, arg DeleteCategoryAttributeParams) ([]byte, error) {
	row := q.db.QueryRow(ctx, deleteCategoryAttribute, arg.Attribute, arg.CategoryID)
	var category_attributes []byte
	err := row.Scan(&category_attributes)
	return category_attributes, err
}

const getCategories = `-- name: GetCategories :many
SELECT category_id, slug, title, description, created_at, updated_at
FROM categories
WHERE shop_id = $1 AND category_id > $2
LIMIT $3
`

type GetCategoriesParams struct {
	ShopID int64
	After  int64
	Limit  int32
}

type GetCategoriesRow struct {
	CategoryID  int64
	Slug        string
	Title       string
	Description pgtype.Text
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

func (q *Queries) GetCategories(ctx context.Context, arg GetCategoriesParams) ([]GetCategoriesRow, error) {
	rows, err := q.db.Query(ctx, getCategories, arg.ShopID, arg.After, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCategoriesRow
	for rows.Next() {
		var i GetCategoriesRow
		if err := rows.Scan(
			&i.CategoryID,
			&i.Slug,
			&i.Title,
			&i.Description,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCategory = `-- name: GetCategory :one
SELECT category_id, slug, title, description, created_at, updated_at, parent_id, category_attributes
FROM categories
WHERE shop_id = $1 AND category_id = $2
`

type GetCategoryParams struct {
	ShopID     int64
	CategoryID int64
}

type GetCategoryRow struct {
	CategoryID         int64
	Slug               string
	Title              string
	Description        pgtype.Text
	CreatedAt          pgtype.Timestamptz
	UpdatedAt          pgtype.Timestamptz
	ParentID           pgtype.Int8
	CategoryAttributes []byte
}

func (q *Queries) GetCategory(ctx context.Context, arg GetCategoryParams) (GetCategoryRow, error) {
	row := q.db.QueryRow(ctx, getCategory, arg.ShopID, arg.CategoryID)
	var i GetCategoryRow
	err := row.Scan(
		&i.CategoryID,
		&i.Slug,
		&i.Title,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ParentID,
		&i.CategoryAttributes,
	)
	return i, err
}

const getProducts = `-- name: GetProducts :many
SELECT product_id, title, description, allowed_attributes, created_at, updated_at, status, category_id, shop_id FROM products
`

func (q *Queries) GetProducts(ctx context.Context) ([]Product, error) {
	rows, err := q.db.Query(ctx, getProducts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ProductID,
			&i.Title,
			&i.Description,
			&i.AllowedAttributes,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Status,
			&i.CategoryID,
			&i.ShopID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getShop = `-- name: GetShop :one
SELECT shop_id, owner_id, title, domain, favicon_url, logo_url, email, currency_code, status, about, address, phone_number, seo_description, seo_keywords, seo_title, updated_at, created_at FROM shops
WHERE shop_id = $1
`

func (q *Queries) GetShop(ctx context.Context, shopID int64) (Shop, error) {
	row := q.db.QueryRow(ctx, getShop, shopID)
	var i Shop
	err := row.Scan(
		&i.ShopID,
		&i.OwnerID,
		&i.Title,
		&i.Domain,
		&i.FaviconUrl,
		&i.LogoUrl,
		&i.Email,
		&i.CurrencyCode,
		&i.Status,
		&i.About,
		&i.Address,
		&i.PhoneNumber,
		&i.SeoDescription,
		&i.SeoKeywords,
		&i.SeoTitle,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getShopByDomain = `-- name: GetShopByDomain :one
SELECT shop_id, owner_id, title, domain, favicon_url, logo_url, email, currency_code, status, about, address, phone_number, seo_description, seo_keywords, seo_title, updated_at, created_at FROM shops
WHERE domain = $1
`

func (q *Queries) GetShopByDomain(ctx context.Context, domain string) (Shop, error) {
	row := q.db.QueryRow(ctx, getShopByDomain, domain)
	var i Shop
	err := row.Scan(
		&i.ShopID,
		&i.OwnerID,
		&i.Title,
		&i.Domain,
		&i.FaviconUrl,
		&i.LogoUrl,
		&i.Email,
		&i.CurrencyCode,
		&i.Status,
		&i.About,
		&i.Address,
		&i.PhoneNumber,
		&i.SeoDescription,
		&i.SeoKeywords,
		&i.SeoTitle,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getShopIDByDomain = `-- name: GetShopIDByDomain :one
SELECT shop_id FROM shops
WHERE domain = $1
`

func (q *Queries) GetShopIDByDomain(ctx context.Context, domain string) (int64, error) {
	row := q.db.QueryRow(ctx, getShopIDByDomain, domain)
	var shop_id int64
	err := row.Scan(&shop_id)
	return shop_id, err
}

const getShopsByOwner = `-- name: GetShopsByOwner :many
SELECT shop_id, owner_id, title, domain, favicon_url, logo_url, email, currency_code, status, about, address, phone_number, seo_description, seo_keywords, seo_title, updated_at, created_at FROM shops
WHERE owner_id = $1
`

func (q *Queries) GetShopsByOwner(ctx context.Context, ownerID uuid.UUID) ([]Shop, error) {
	rows, err := q.db.Query(ctx, getShopsByOwner, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Shop
	for rows.Next() {
		var i Shop
		if err := rows.Scan(
			&i.ShopID,
			&i.OwnerID,
			&i.Title,
			&i.Domain,
			&i.FaviconUrl,
			&i.LogoUrl,
			&i.Email,
			&i.CurrencyCode,
			&i.Status,
			&i.About,
			&i.Address,
			&i.PhoneNumber,
			&i.SeoDescription,
			&i.SeoKeywords,
			&i.SeoTitle,
			&i.UpdatedAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUser = `-- name: GetUser :one
SELECT user_id, auth0_sub, email, name, profile_picture_url, created_at, last_login FROM users
WHERE auth0_sub = $1
`

func (q *Queries) GetUser(ctx context.Context, auth0Sub pgtype.Text) (User, error) {
	row := q.db.QueryRow(ctx, getUser, auth0Sub)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Auth0Sub,
		&i.Email,
		&i.Name,
		&i.ProfilePictureUrl,
		&i.CreatedAt,
		&i.LastLogin,
	)
	return i, err
}

const updateCategory = `-- name: UpdateCategory :one
UPDATE categories
SET 
    title = COALESCE($1, title),
    description = COALESCE($2, description),
    parent_id = COALESCE($3, parent_id)
WHERE category_id = $4
RETURNING category_id, slug, title, description, parent_id, created_at, updated_at, shop_id, category_attributes
`

type UpdateCategoryParams struct {
	Title       pgtype.Text
	Description pgtype.Text
	ParentID    pgtype.Int8
	CategoryID  int64
}

func (q *Queries) UpdateCategory(ctx context.Context, arg UpdateCategoryParams) (Category, error) {
	row := q.db.QueryRow(ctx, updateCategory,
		arg.Title,
		arg.Description,
		arg.ParentID,
		arg.CategoryID,
	)
	var i Category
	err := row.Scan(
		&i.CategoryID,
		&i.Slug,
		&i.Title,
		&i.Description,
		&i.ParentID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ShopID,
		&i.CategoryAttributes,
	)
	return i, err
}

const updateShop = `-- name: UpdateShop :one
UPDATE shops
SET 
    title = COALESCE($1, title),
    favicon_url = COALESCE($2, favicon_url),
    currency_code = COALESCE($3, currency_code),
    about = COALESCE($4, about),
    status = COALESCE($5, status),
    phone_number = COALESCE($6, phone_number),
    seo_description = COALESCE($7, seo_description),
    seo_keywords = COALESCE($8, seo_keywords),
    seo_title = COALESCE($9, seo_title),
    address = COALESCE($10, address),
    email = COALESCE($11, email)
WHERE domain = $12
RETURNING shop_id, owner_id, title, domain, favicon_url, logo_url, email, currency_code, status, about, address, phone_number, seo_description, seo_keywords, seo_title, updated_at, created_at
`

type UpdateShopParams struct {
	Title          pgtype.Text
	FaviconUrl     pgtype.Text
	CurrencyCode   pgtype.Text
	About          pgtype.Text
	Status         pgtype.Text
	PhoneNumber    pgtype.Text
	SeoDescription pgtype.Text
	SeoKeywords    []string
	SeoTitle       pgtype.Text
	Address        pgtype.Text
	Email          pgtype.Text
	Domain         string
}

func (q *Queries) UpdateShop(ctx context.Context, arg UpdateShopParams) (Shop, error) {
	row := q.db.QueryRow(ctx, updateShop,
		arg.Title,
		arg.FaviconUrl,
		arg.CurrencyCode,
		arg.About,
		arg.Status,
		arg.PhoneNumber,
		arg.SeoDescription,
		arg.SeoKeywords,
		arg.SeoTitle,
		arg.Address,
		arg.Email,
		arg.Domain,
	)
	var i Shop
	err := row.Scan(
		&i.ShopID,
		&i.OwnerID,
		&i.Title,
		&i.Domain,
		&i.FaviconUrl,
		&i.LogoUrl,
		&i.Email,
		&i.CurrencyCode,
		&i.Status,
		&i.About,
		&i.Address,
		&i.PhoneNumber,
		&i.SeoDescription,
		&i.SeoKeywords,
		&i.SeoTitle,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const upsertUser = `-- name: UpsertUser :one
INSERT INTO users (auth0_sub, email, name, profile_picture_url)
VALUES ($1, $2, $3, $4)
ON CONFLICT (auth0_sub)
DO UPDATE SET email = EXCLUDED.email, name = EXCLUDED.name, profile_picture_url = EXCLUDED.profile_picture_url
RETURNING user_id, auth0_sub, email, name, profile_picture_url
`

type UpsertUserParams struct {
	Auth0Sub          pgtype.Text
	Email             string
	Name              pgtype.Text
	ProfilePictureUrl pgtype.Text
}

type UpsertUserRow struct {
	UserID            uuid.UUID
	Auth0Sub          pgtype.Text
	Email             string
	Name              pgtype.Text
	ProfilePictureUrl pgtype.Text
}

func (q *Queries) UpsertUser(ctx context.Context, arg UpsertUserParams) (UpsertUserRow, error) {
	row := q.db.QueryRow(ctx, upsertUser,
		arg.Auth0Sub,
		arg.Email,
		arg.Name,
		arg.ProfilePictureUrl,
	)
	var i UpsertUserRow
	err := row.Scan(
		&i.UserID,
		&i.Auth0Sub,
		&i.Email,
		&i.Name,
		&i.ProfilePictureUrl,
	)
	return i, err
}
