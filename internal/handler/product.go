package handler

import (
	"context"

	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/service"
	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/utils"
	"github.com/fathurzoy/go-grpc-ecommerce-be/pb/product"
)

type productHandler struct {
	product.UnimplementedProductSeriviceServer

	productService service.IProductService
}

func (ph *productHandler) CreateProduct(ctx context.Context, request *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &product.CreateProductResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	// Proccess register

	res, err := ph.productService.CreateProduct(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ph *productHandler) DetailProduct(ctx context.Context, request *product.DetailProductRequest) (*product.DetailProductResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &product.DetailProductResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ph.productService.DetailProduct(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ph *productHandler) EditProduct(ctx context.Context, request *product.EditProductRequest) (*product.EditProductResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &product.EditProductResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ph.productService.EditProduct(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewProductHandler(productService service.IProductService) *productHandler {
	return &productHandler{
		productService: productService,
	}
}
