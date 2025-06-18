package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/entity"
)

type IProductRepository interface {
	CreateNewProduct(ctx context.Context, product *entity.Product) error
	GetProductById(ctx context.Context, id string) (*entity.Product, error)
	UpdateProduct(ctx context.Context, product *entity.Product) error
	DeleteProduct(ctx context.Context, id string, deletedAt time.Time, deletedBy *string) error
}

type productRepository struct {
	db *sql.DB
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

func NewProductRepository(db *sql.DB) IProductRepository {
	return &productRepository{
		db: db,
	}
}
