package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/entity"
	jwtentity "github.com/fathurzoy/go-grpc-ecommerce-be/internal/entity/jwt"
	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/repository"
	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/utils"
	"github.com/fathurzoy/go-grpc-ecommerce-be/pb/common"
	"github.com/fathurzoy/go-grpc-ecommerce-be/pb/product"
	"github.com/google/uuid"
)

type IProductService interface {
	CreateProduct(ctx context.Context, request *product.CreateProductRequest) (*product.CreateProductResponse, error)
	DetailProduct(ctx context.Context, request *product.DetailProductRequest) (*product.DetailProductResponse, error)
	EditProduct(ctx context.Context, request *product.EditProductRequest) (*product.EditProductResponse, error)
	DeleteProduct(ctx context.Context, request *product.DeleteProductRequest) (*product.DeleteProductResponse, error)
	ListProduct(ctx context.Context, request *product.ListProductRequest) (*product.ListProductResponse, error)
	ListProductAdmin(ctx context.Context, request *product.ListProductAdminRequest) (*product.ListProductAdminResponse, error)
}

type productService struct {
	productRepository repository.IProductRepository
}

func (ps *productService) CreateProduct(ctx context.Context, request *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	//cek apakah user admin?
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if claims.Role != entity.UserRoleAdmin {
		return nil, utils.UnauthenticatedResponse()
	}

	//apakah image ada?
	imagePath := filepath.Join("storage", "product", request.ImageFileName)
	_, err = os.Stat(imagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &product.CreateProductResponse{
				Base: utils.BadRequestResponse("Image not found"),
			}, nil
		}
		return nil, err
	}
	//insert ke database

	// success response
	productEntity := entity.Product{
		Id:            uuid.NewString(),
		Name:          request.Name,
		Description:   request.Description,
		Price:         request.Price,
		ImageFileName: request.ImageFileName,
		CreatedAt:     time.Now(),
		CreatedBy:     &claims.FullName,
	}

	//ngecek email ke database
	err = ps.productRepository.CreateNewProduct(ctx, &productEntity)
	if err != nil {
		return nil, err
	}

	return &product.CreateProductResponse{
		Base: utils.SuccessResponse("Create product success"),
		Id:   productEntity.Id,
	}, nil
}

func (ps *productService) DetailProduct(ctx context.Context, request *product.DetailProductRequest) (*product.DetailProductResponse, error) {
	// query ke db dengan data id
	productEntity, err := ps.productRepository.GetProductById(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	// apabila null, return notfound
	if productEntity == nil {
		return &product.DetailProductResponse{
			Base: utils.NotFoundResponse("Product not found"),
		}, nil
	}

	// kirim response
	return &product.DetailProductResponse{
		Base:        utils.SuccessResponse("Get product success"),
		Id:          productEntity.Id,
		Name:        productEntity.Name,
		Description: productEntity.Description,
		Price:       productEntity.Price,
		ImageUrl:    fmt.Sprintf("%s/product/%s", os.Getenv("STORAGE_SERVICE_URL"), productEntity.ImageFileName),
	}, nil
}

func (ps *productService) EditProduct(ctx context.Context, request *product.EditProductRequest) (*product.EditProductResponse, error) {
	// valdiasi apakah id yang dikirim itu ada di db
	productEntity, err := ps.productRepository.GetProductById(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	if productEntity == nil {
		return &product.EditProductResponse{
			Base: utils.NotFoundResponse("Product not found"),
		}, nil
	}

	// cek apakah user admin
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if claims.Role != entity.UserRoleAdmin {
		return nil, utils.UnauthenticatedResponse()
	}

	// kalau gambarnya berubah, hapus gambar lama
	if productEntity.ImageFileName != request.ImageFileName {
		newImagePath := filepath.Join("storage", "product", request.ImageFileName)
		_, err = os.Stat(newImagePath)
		if err != nil {
			if os.IsNotExist(err) {
				return &product.EditProductResponse{
					Base: utils.BadRequestResponse("Image not found"),
				}, nil
			}
			return nil, err
		}

		oldImagePath := filepath.Join("storage", "product", productEntity.ImageFileName)
		err = os.Remove(oldImagePath)
		if err != nil {
			log.Println(err)
			// return nil, err
		}
	}

	// update ke db
	newProduct := entity.Product{
		Id:            request.Id,
		Name:          request.Name,
		Description:   request.Description,
		Price:         request.Price,
		ImageFileName: request.ImageFileName,
		UpdatedAt:     time.Now(),
		UpdatedBy:     &claims.FullName,
	}

	err = ps.productRepository.UpdateProduct(ctx, &newProduct)
	if err != nil {
		return nil, err
	}

	// kirim response
	return &product.EditProductResponse{
		Base: utils.SuccessResponse("Edit product success"),
		Id:   request.Id,
	}, nil
}

func (ps *productService) DeleteProduct(ctx context.Context, request *product.DeleteProductRequest) (*product.DeleteProductResponse, error) {
	// valdiasi apakah id yang dikirim itu ada di db
	productEntity, err := ps.productRepository.GetProductById(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	if productEntity == nil {
		return &product.DeleteProductResponse{
			Base: utils.NotFoundResponse("Product not found"),
		}, nil
	}

	// cek apakah user admin
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if claims.Role != entity.UserRoleAdmin {
		return nil, utils.UnauthenticatedResponse()
	}

	// delete ke db
	err = ps.productRepository.DeleteProduct(ctx, request.Id, time.Now(), &claims.FullName)
	if err != nil {
		return nil, err
	}

	imagePath := filepath.Join("storage", "product", productEntity.ImageFileName)
	err = os.Remove(imagePath)
	if err != nil {
		log.Println(err)
		// return nil, err
	}

	// kirim response
	return &product.DeleteProductResponse{
		Base: utils.SuccessResponse("Delete product success"),
	}, nil
}

func (ps *productService) ListProduct(ctx context.Context, request *product.ListProductRequest) (*product.ListProductResponse, error) {
	//buat default request pagination
	if request.Pagination == nil {
		request.Pagination = &common.PaginationRequest{
			CurrentPage: 1,
			ItemPerPage: 10,
		}
	}

	products, paginationResponse, err := ps.productRepository.GetProductPagination(ctx, request.Pagination)
	if err != nil {
		return nil, err
	}

	var data []*product.ListProductResponseItem = make([]*product.ListProductResponseItem, 0)
	for _, v := range products {
		data = append(data, &product.ListProductResponseItem{
			Id:          v.Id,
			Name:        v.Name,
			Description: v.Description,
			Price:       v.Price,
			ImageUrl:    fmt.Sprintf("%s/product/%s", os.Getenv("STORAGE_SERVICE_URL"), v.ImageFileName),
		})
	}

	return &product.ListProductResponse{
		Base:       utils.SuccessResponse("Get product success"),
		Data:       data,
		Pagination: paginationResponse,
	}, nil
}

func (ps *productService) ListProductAdmin(ctx context.Context, request *product.ListProductAdminRequest) (*product.ListProductAdminResponse, error) {
	// cek apakah user admin
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if claims.Role != entity.UserRoleAdmin {
		return nil, utils.UnauthenticatedResponse()
	}

	//buat default request pagination
	if request.Pagination == nil {
		request.Pagination = &common.PaginationRequest{
			CurrentPage: 1,
			ItemPerPage: 10,
		}
	}

	products, paginationResponse, err := ps.productRepository.GetProductPagination(ctx, request.Pagination)
	if err != nil {
		return nil, err
	}

	var data []*product.ListProductAdminResponseItem = make([]*product.ListProductAdminResponseItem, 0)
	for _, v := range products {
		data = append(data, &product.ListProductAdminResponseItem{
			Id:          v.Id,
			Name:        v.Name,
			Description: v.Description,
			Price:       v.Price,
			ImageUrl:    fmt.Sprintf("%s/product/%s", os.Getenv("STORAGE_SERVICE_URL"), v.ImageFileName),
		})
	}

	return &product.ListProductAdminResponse{
		Base:       utils.SuccessResponse("Get product success"),
		Data:       data,
		Pagination: paginationResponse,
	}, nil
}

func NewProductService(productRepository repository.IProductRepository) IProductService {
	return &productService{
		productRepository: productRepository,
	}
}
