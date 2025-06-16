package service

import (
	"context"
	"time"

	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/entity"
	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/repository"
	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/utils"
	"github.com/fathurzoy/go-grpc-ecommerce-be/pb/auth"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Register(ctx context.Context, request *auth.RegisterRequest) (*auth.RegisterResponse, error)
}

type authService struct {
	authRepository repository.IAuthRepository
}

func (as *authService) Register(ctx context.Context, request *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	if request.Password != request.PasswordConfirmation {
		return &auth.RegisterResponse{
			Base: utils.BadRequestResponse("Password and password confirmation does not match"),
		}, nil
	}

	//ngecek email ke database
	user, err := as.authRepository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	//apabila email sudah terdaftar, kita errorkan karena tidak mau double
	if user != nil {
		return &auth.RegisterResponse{
			Base: utils.BadRequestResponse("Email already registered"),
		}, nil
	}

	//hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
	if err != nil {
		return nil, err
	}

	//buat user baru
	newUser := entity.User{
		Id:        uuid.NewString(),
		FullName:  request.FullName,
		Email:     request.Email,
		Password:  string(hashedPassword),
		RoleCode:  entity.UserRoleCustomer,
		CreatedAt: time.Now(),
		CreatedBy: &request.FullName,
	}

	//apabila belum terdaftar, insert ke database
	err = as.authRepository.InsertUser(ctx, &newUser)
	if err != nil {
		return nil, err
	}

	return &auth.RegisterResponse{
		Base: utils.SuccessResponse("Register success"),
	}, nil
}

func NewAuthService(authRepository repository.IAuthRepository) IAuthService {
	return &authService{
		authRepository: authRepository,
	}
}
