package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/entity"
	jwtentity "github.com/fathurzoy/go-grpc-ecommerce-be/internal/entity/jwt"
	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/repository"
	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/utils"
	"github.com/fathurzoy/go-grpc-ecommerce-be/pb/cart"
	"github.com/google/uuid"
)

type ICartService interface {
	AddProductToCart(ctx context.Context, request *cart.AddProductToCartRequest) (*cart.AddProductToCartResponse, error)
	ListCart(ctx context.Context, request *cart.ListCartRequest) (*cart.ListCartResponse, error)
	DeleteCart(ctx context.Context, request *cart.DeleteCartRequest) (*cart.DeleteCartResponse, error)
}

type cartService struct {
	productRepository repository.IProductRepository
	cartRepository    repository.ICartRepository
}

func (cs *cartService) AddProductToCart(ctx context.Context, request *cart.AddProductToCartRequest) (*cart.AddProductToCartResponse, error) {
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// cek terlebih dahulu apakah product id itu ada di db
	productEntity, err := cs.productRepository.GetProductById(ctx, request.ProductId)
	if err != nil {
		return nil, err
	}

	if productEntity == nil {
		return &cart.AddProductToCartResponse{
			Base: utils.NotFoundResponse("Product not found"),
		}, nil
	}
	fmt.Println("x")

	// cek ke db apakah product udah ada di cart user ini
	cartEntity, err := cs.cartRepository.GetCartProductAndUserId(ctx, request.ProductId, claims.Subject)
	if err != nil {
		return nil, err
	}

	// kalau udah ada update db
	if cartEntity != nil {
		now := time.Now()
		cartEntity.Quantity += 1
		cartEntity.UpdatedAt = &now
		cartEntity.UpdatedBy = &claims.Subject

		err = cs.cartRepository.UpdateCart(ctx, cartEntity)
		if err != nil {
			return nil, err
		}

		return &cart.AddProductToCartResponse{
			Base: utils.SuccessResponse("Add product to cart success"),
			Id:   cartEntity.Id,
		}, nil
	}

	// kalau belum ada insert ke db
	newCartEntity := entity.UserCart{
		Id:        uuid.NewString(),
		UserId:    claims.Subject,
		ProductId: request.ProductId,
		Quantity:  1,
		CreatedAt: time.Now(),
		CreatedBy: &claims.Subject,
	}

	err = cs.cartRepository.CreateNewCart(ctx, &newCartEntity)
	if err != nil {
		return nil, err
	}

	return &cart.AddProductToCartResponse{
		Base: utils.SuccessResponse("Add product to cart success"),
		Id:   newCartEntity.Id,
	}, nil
}

func (cs *cartService) ListCart(ctx context.Context, request *cart.ListCartRequest) (*cart.ListCartResponse, error) {
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	carts, err := cs.cartRepository.GetListCart(ctx, claims.Subject)
	if err != nil {
		return nil, err
	}

	var items []*cart.ListCartResponseItem
	for _, cartEntity := range carts {
		item := &cart.ListCartResponseItem{
			CartId:          cartEntity.Id,
			ProductId:       cartEntity.ProductId,
			ProductName:     cartEntity.Product.Name,
			ProductImageUrl: fmt.Sprintf("%s/product/%s", os.Getenv("STORAEGE_SERVICE_URL"), cartEntity.Product.ImageFileName),
			ProductPrice:    cartEntity.Product.Price,
			Quantity:        int32(cartEntity.Quantity),
		}
		items = append(items, item)
	}

	return &cart.ListCartResponse{
		Base:  utils.SuccessResponse("List cart success"),
		Items: items,
	}, nil
}

func (cs *cartService) DeleteCart(ctx context.Context, request *cart.DeleteCartRequest) (*cart.DeleteCartResponse, error) {
	// dapat data user id
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// dapat data cart
	cartEntity, err := cs.cartRepository.GetCartById(ctx, request.CartId)
	if err != nil {
		return nil, err
	}
	if cartEntity == nil {
		return &cart.DeleteCartResponse{
			Base: utils.NotFoundResponse("Cart not found"),
		}, nil
	}

	// cocokan data user id di cart dengan auth
	if cartEntity.UserId != claims.Subject {
		return &cart.DeleteCartResponse{
			Base: utils.BadRequestResponse("Forbidden"),
		}, nil
	}

	err = cs.cartRepository.DeleteCart(ctx, request.CartId)
	if err != nil {
		return nil, err
	}

	return &cart.DeleteCartResponse{
		Base: utils.SuccessResponse("Delete cart success"),
	}, nil
}

func NewCartService(productRepository repository.IProductRepository, cartRepository repository.ICartRepository) ICartService {
	return &cartService{
		productRepository: productRepository,
		cartRepository:    cartRepository,
	}
}
