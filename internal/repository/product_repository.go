package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/entity"
	"github.com/fathurzoy/go-grpc-ecommerce-be/pb/common"
	"github.com/fathurzoy/go-grpc-ecommerce-be/pkg/database"
)

type IProductRepository interface {
	WithTransaction(tx *sql.Tx) IProductRepository
	CreateNewProduct(ctx context.Context, product *entity.Product) error
	GetProductById(ctx context.Context, id string) (*entity.Product, error)
	GetProductsByIds(ctx context.Context, ids []string) ([]*entity.Product, error)
	UpdateProduct(ctx context.Context, product *entity.Product) error
	DeleteProduct(ctx context.Context, id string, deletedAt time.Time, deletedBy *string) error
	GetProductPagination(ctx context.Context, pagination *common.PaginationRequest) ([]*entity.Product, *common.PaginationResponse, error)
	GetProductAdminPagination(ctx context.Context, pagination *common.PaginationRequest) ([]*entity.Product, *common.PaginationResponse, error)
	GetProductHighlight(ctx context.Context) ([]*entity.Product, error)
}

type productRepository struct {
	db database.DatabaseQuery
}

func (repo *productRepository) WithTransaction(tx *sql.Tx) IProductRepository {
	return &productRepository{
		db: tx,
	}
}

func (repo *productRepository) CreateNewProduct(ctx context.Context, product *entity.Product) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO product (id, name, description, price, image_file_name, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, is_deleted) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)", product.Id, product.Name, product.Description, product.Price, product.ImageFileName, product.CreatedAt, product.CreatedBy, product.UpdatedAt, product.UpdatedBy, product.DeletedAt, product.DeletedBy, product.IsDeleted)
	if err != nil {
		return err
	}

	return nil
}

func (repo *productRepository) GetProductById(ctx context.Context, id string) (*entity.Product, error) {
	var productEntity entity.Product

	row := repo.db.QueryRowContext(ctx, "SELECT id, name, description, price, image_file_name FROM product WHERE id = $1", id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(
		&productEntity.Id,
		&productEntity.Name,
		&productEntity.Description,
		&productEntity.Price,
		&productEntity.ImageFileName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &productEntity, nil
}

func (repo *productRepository) GetProductsByIds(ctx context.Context, ids []string) ([]*entity.Product, error) {
	queryIds := make([]string, len(ids))
	for i, id := range ids {
		queryIds[i] = fmt.Sprintf("'%s'", id)
	}

	rows, err := repo.db.QueryContext(
		ctx,
		fmt.Sprintf("SELECT id, name, price, image_file_name FROM product WHERE id in (%s) and is_deleted = false", strings.Join(queryIds, ",")),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*entity.Product
	for rows.Next() {
		var product entity.Product
		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Price,
			&product.ImageFileName,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}

func (repo *productRepository) UpdateProduct(ctx context.Context, product *entity.Product) error {
	_, err := repo.db.ExecContext(
		ctx,
		"UPDATE product SET name = $1, description = $2, price = $3, image_file_name = $4, updated_at = $5, updated_by = $6 WHERE id = $7",
		product.Name,
		product.Description,
		product.Price,
		product.ImageFileName,
		product.UpdatedAt,
		product.UpdatedBy,
		product.Id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *productRepository) DeleteProduct(ctx context.Context, id string, deletedAt time.Time, deletedBy *string) error {
	_, err := repo.db.ExecContext(ctx, "UPDATE product SET is_deleted = $1, deleted_at = $2, deleted_by = $3 WHERE id = $4", true, deletedAt, deletedBy, id)
	if err != nil {
		return err
	}

	return nil
}

func (repo *productRepository) GetProductPagination(ctx context.Context, pagination *common.PaginationRequest) ([]*entity.Product, *common.PaginationResponse, error) {
	row := repo.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM product WHERE is_deleted = $1", false)
	if row.Err() != nil {
		return nil, nil, row.Err()
	}

	var totalCount int
	err := row.Scan(&totalCount)
	if err != nil {
		return nil, nil, err
	}

	offset := (pagination.CurrentPage - 1) * pagination.ItemPerPage

	// Hitung total halaman
	totalPages := (totalCount + int(pagination.ItemPerPage) - 1) / int(pagination.ItemPerPage)

	rows, err := repo.db.QueryContext(ctx, `
		SELECT id, name, description, price, image_file_name 
		FROM product 
		WHERE is_deleted = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`, false, pagination.ItemPerPage, offset)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var products []*entity.Product
	for rows.Next() {
		var productEntity entity.Product
		if err := rows.Scan(
			&productEntity.Id,
			&productEntity.Name,
			&productEntity.Description,
			&productEntity.Price,
			&productEntity.ImageFileName,
		); err != nil {
			return nil, nil, err
		}
		products = append(products, &productEntity)
	}

	paginationResponse := &common.PaginationResponse{
		CurrentPage:    pagination.CurrentPage,
		ItemPerPage:    pagination.ItemPerPage,
		TotalPageCount: int32(totalPages),
		TotalItemCount: int32(totalCount),
	}

	return products, paginationResponse, nil
}

func (repo *productRepository) GetProductAdminPagination(ctx context.Context, pagination *common.PaginationRequest) ([]*entity.Product, *common.PaginationResponse, error) {
	row := repo.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM product")
	if row.Err() != nil {
		return nil, nil, row.Err()
	}

	var totalCount int
	err := row.Scan(&totalCount)
	if err != nil {
		return nil, nil, err
	}

	offset := (pagination.CurrentPage - 1) * pagination.ItemPerPage

	// Hitung total halaman
	totalPages := (totalCount + int(pagination.ItemPerPage) - 1) / int(pagination.ItemPerPage)

	allowedSorts := map[string]string{
		"name":        "name",
		"description": "description",
		"price":       "price",
	}

	orderQuery := " ORDER BY created_at DESC"
	if pagination.Sort != nil && allowedSorts[pagination.Sort.Field] != "" {
		direction := "ASC"
		if pagination.Sort.Direction == "desc" {
			direction = "DESC"
		}
		orderQuery = fmt.Sprintf(" ORDER BY %s %s", pagination.Sort.Field, direction)
	}

	baseQuery := fmt.Sprintf("SELECT id, name, description, price, image_file_name FROM product WHERE is_deleted = false %s LIMIT $1 OFFSET $2", orderQuery)
	rows, err := repo.db.QueryContext(ctx, baseQuery, pagination.ItemPerPage, offset)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var products []*entity.Product
	for rows.Next() {
		var productEntity entity.Product
		if err := rows.Scan(
			&productEntity.Id,
			&productEntity.Name,
			&productEntity.Description,
			&productEntity.Price,
			&productEntity.ImageFileName,
		); err != nil {
			return nil, nil, err
		}
		products = append(products, &productEntity)
	}

	paginationResponse := &common.PaginationResponse{
		CurrentPage:    pagination.CurrentPage,
		ItemPerPage:    pagination.ItemPerPage,
		TotalPageCount: int32(totalPages),
		TotalItemCount: int32(totalCount),
	}

	return products, paginationResponse, nil
}

func (repo *productRepository) GetProductHighlight(ctx context.Context) ([]*entity.Product, error) {
	query := `
		SELECT id, name, description, price, image_file_name
		FROM product
		WHERE id IN (
			SELECT p.id
			FROM product p
			JOIN order_item oi ON oi.product_id = p.id
			WHERE p.is_deleted = false AND oi.is_deleted = false
			GROUP BY p.id
			ORDER BY COUNT(*) DESC
			LIMIT 5
		)
	`

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*entity.Product
	for rows.Next() {
		var product entity.Product
		if err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.ImageFileName,
		); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func NewProductRepository(db database.DatabaseQuery) IProductRepository {
	return &productRepository{
		db: db,
	}
}
