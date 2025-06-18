package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/entity"
)

type ICartRepository interface {
	GetCartProductAndUserId(ctx context.Context, productId, userId string) (*entity.UserCart, error)
	CreateNewCart(ctx context.Context, cartEntity *entity.UserCart) error
	UpdateCart(ctx context.Context, cartEntity *entity.UserCart) error
	GetListCart(ctx context.Context, userId string) ([]*entity.UserCart, error)
	GetCartById(ctx context.Context, cartId string) (*entity.UserCart, error)
	DeleteCart(ctx context.Context, cartId string) error
}

type cartRepository struct {
	db *sql.DB
}

func (cr *cartRepository) GetCartProductAndUserId(ctx context.Context, productId, userId string) (*entity.UserCart, error) {

	// TODO: implement
	row := cr.db.QueryRowContext(
		ctx,
		"SELECT id, user_id, product_id, quantity, created_at, created_by, updated_at, updated_by FROM user_cart WHERE user_id = $1 AND product_id = $2",
		userId,
		productId,
	)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var cartEntity entity.UserCart
	err := row.Scan(
		&cartEntity.Id,
		&cartEntity.UserId,
		&cartEntity.ProductId,
		&cartEntity.Quantity,
		&cartEntity.CreatedAt,
		&cartEntity.CreatedBy,
		&cartEntity.UpdatedAt,
		&cartEntity.UpdatedBy,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &cartEntity, nil
}

func (cr *cartRepository) CreateNewCart(ctx context.Context, cartEntity *entity.UserCart) error {
	_, err := cr.db.ExecContext(ctx, "INSERT INTO user_cart (id, user_id, product_id, quantity, created_at, created_by) VALUES ($1, $2, $3, $4, $5, $6)",
		cartEntity.Id,
		cartEntity.UserId,
		cartEntity.ProductId,
		cartEntity.Quantity,
		cartEntity.CreatedAt,
		cartEntity.CreatedBy,
	)
	if err != nil {
		return err
	}

	return nil
}

func (cr *cartRepository) UpdateCart(ctx context.Context, cartEntity *entity.UserCart) error {
	_, err := cr.db.ExecContext(ctx, `
		UPDATE user_cart
		SET product_id = $1,
		    user_id = $2,
		    quantity = $3,
		    updated_at = $4,
		    updated_by = $5
		WHERE id = $6`,
		cartEntity.ProductId,
		cartEntity.UserId,
		cartEntity.Quantity,
		cartEntity.UpdatedAt,
		cartEntity.UpdatedBy,
		cartEntity.Id,
	)
	return err
}

func (cr *cartRepository) GetListCart(ctx context.Context, userId string) ([]*entity.UserCart, error) {
	rows, err := cr.db.QueryContext(
		ctx,
		`SELECT uc.id, uc.user_id, uc.product_id, uc.quantity, uc.created_at, uc.created_by, uc.updated_at, uc.updated_by,
		        p.id, p.name, p.description, p.price, p.image_file_name
		   FROM user_cart uc
		   JOIN product p ON uc.product_id = p.id
		  WHERE uc.user_id = $1`,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carts []*entity.UserCart
	for rows.Next() {
		var cart entity.UserCart
		cart.Product = &entity.Product{}

		err := rows.Scan(
			&cart.Id,
			&cart.UserId,
			&cart.ProductId,
			&cart.Quantity,
			&cart.CreatedAt,
			&cart.CreatedBy,
			&cart.UpdatedAt,
			&cart.UpdatedBy,
			&cart.Product.Id,
			&cart.Product.Name,
			&cart.Product.Description,
			&cart.Product.Price,
			&cart.Product.ImageFileName,
		)
		if err != nil {
			return nil, err
		}

		carts = append(carts, &cart)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return carts, nil
}

func (cr *cartRepository) GetCartById(ctx context.Context, cartId string) (*entity.UserCart, error) {
	// TODO: implement
	row := cr.db.QueryRowContext(
		ctx,
		"SELECT id, user_id, product_id, quantity, created_at, created_by, updated_at, updated_by FROM user_cart WHERE id = $1",
		cartId,
	)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var cartEntity entity.UserCart
	err := row.Scan(
		&cartEntity.Id,
		&cartEntity.UserId,
		&cartEntity.ProductId,
		&cartEntity.Quantity,
		&cartEntity.CreatedAt,
		&cartEntity.CreatedBy,
		&cartEntity.UpdatedAt,
		&cartEntity.UpdatedBy,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &cartEntity, nil
}

func (cr *cartRepository) DeleteCart(ctx context.Context, cartId string) error {
	_, err := cr.db.ExecContext(
		ctx,
		"DELETE FROM user_cart WHERE id = $1",
		cartId,
	)
	if err != nil {
		return err
	}

	return err
}

func NewCartRepository(db *sql.DB) ICartRepository {
	return &cartRepository{
		db: db,
	}
}
