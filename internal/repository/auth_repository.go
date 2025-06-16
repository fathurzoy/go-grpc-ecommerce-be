package repository

import (
	"context"
	"database/sql"

	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/entity"
)

type IAuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	InsertUser(ctx context.Context, user *entity.User) error
}

type authRepository struct {
	db *sql.DB
}

// InsertUser implements IAuthRepository.
func (ar *authRepository) InsertUser(ctx context.Context, user *entity.User) error {
	_, err := ar.db.ExecContext(ctx, "INSERT INTO \"user\"  (id, email, password, full_name, role_code, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, is_deleted) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
		user.Id, user.Email, user.Password, user.FullName, user.RoleCode, user.CreatedAt, user.CreatedBy, user.UpdatedAt, user.UpdatedBy, user.DeletedAt, user.DeletedBy, user.IsDeleted)
	if err != nil {
		return err
	}

	return nil

}

func (ar *authRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	row := ar.db.QueryRowContext(ctx, "SELECT id, email, password, full_name, role_code FROM \"user\" WHERE email = $1 and is_deleted is false", email)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var user entity.User
	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.FullName, &user.RoleCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func NewAuthRepository(db *sql.DB) IAuthRepository {
	return &authRepository{
		db: db,
	}
}
