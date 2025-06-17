package repository

import (
	"context"
	"database/sql"

	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/entity"
)

type IProductRepository interface {
	CreateNewProduct(ctx context.Context, product *entity.Product) error
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

func NewProductRepository(db *sql.DB) IProductRepository {
	return &productRepository{
		db: db,
	}
}
