package handler

import (
	"context"

	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/service"
	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/utils"
	"github.com/fathurzoy/go-grpc-ecommerce-be/pb/cart"
)

type cartHandler struct {
	cart.UnimplementedCartServiceServer

	cartService service.ICartService
}

func (c *cartHandler) AddProductToCart(ctx context.Context, request *cart.AddProductToCartRequest) (*cart.AddProductToCartResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &cart.AddProductToCartResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := c.cartService.AddProductToCart(ctx, request)
	if err != nil {
		return nil, err
	}

	return &cart.AddProductToCartResponse{
		Base: res.Base,
		Id:   res.Id,
	}, nil
}

func (c *cartHandler) ListCart(ctx context.Context, request *cart.ListCartRequest) (*cart.ListCartResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &cart.ListCartResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := c.cartService.ListCart(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *cartHandler) DeleteCart(ctx context.Context, request *cart.DeleteCartRequest) (*cart.DeleteCartResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &cart.DeleteCartResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := c.cartService.DeleteCart(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewCartHandler(cartService service.ICartService) *cartHandler {
	return &cartHandler{
		cartService: cartService,
	}
}
