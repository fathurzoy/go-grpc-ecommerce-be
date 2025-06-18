package service

import (
	"context"
	"fmt"
	"time"

	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/entity"
	jwtentity "github.com/fathurzoy/go-grpc-ecommerce-be/internal/entity/jwt"
	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/repository"
	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/utils"
	"github.com/fathurzoy/go-grpc-ecommerce-be/pb/order"
	"github.com/google/uuid"
)

type IOrderService interface {
	CreateOrder(ctx context.Context, request *order.CreateOrderRequest) (*order.CreateOrderResponse, error)
}

type orderService struct {
	orderRepository   repository.IOrderRepository
	productRepository repository.IProductRepository
}

func (os *orderService) CreateOrder(ctx context.Context, request *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	numbering, err := os.orderRepository.GetNumbering(ctx, "order")
	if err != nil {
		return nil, err
	}

	var productIds = make([]string, len(request.Products))
	for i := range request.Products {
		productIds[i] = request.Products[i].Id
	}

	// cek apakah product ada
	products, err := os.productRepository.GetProductsByIds(ctx, productIds)
	if err != nil {
		return nil, err
	}

	productMap := make(map[string]*entity.Product)
	for i := range products {
		productMap[products[i].Id] = products[i]
	}

	var total float64 = 0
	for _, p := range request.Products {
		if productMap[p.Id] == nil {
			return &order.CreateOrderResponse{
				Base: utils.NotFoundResponse(fmt.Sprintf("Product %s not found", p.Id)),
			}, nil
		}
		total += productMap[p.Id].Price * float64(p.Quantity)
	}

	// simpan order ke database
	now := time.Now()
	expiredAt := now.Add(time.Hour * 24)
	orderEntity := entity.Order{
		Id: uuid.NewString(),
		//  ORD-(YEAR)(NUMBER-7) ORD-20251111111
		Number:          fmt.Sprintf("ORD-%d%08d", now.Year(), numbering),
		UserId:          claims.Subject,
		OrderStatusCode: entity.OrderStatusCodeUnpaid,
		UserFullName:    request.FullName,
		Address:         request.Address,
		PhoneNumber:     request.PhoneNumber,
		Notes:           request.Notes,
		Total:           total,
		ExpiredAt:       &expiredAt,
		CreatedAt:       now,
		CreatedBy:       claims.FullName,
	}

	err = os.orderRepository.CreateOrder(ctx, &orderEntity)
	if err != nil {
		return nil, err
	}

	// iterasi semua data product di request
	// setiap iterasinya, simpan order_item ke database
	for _, p := range request.Products {
		var orderitem = entity.OrderItem{
			Id:                   uuid.NewString(),
			ProductId:            p.Id,
			ProductName:          productMap[p.Id].Name,
			ProductImageFileName: productMap[p.Id].ImageFileName,
			ProductPrice:         productMap[p.Id].Price,
			Quantity:             p.Quantity,
			OrderId:              orderEntity.Id,
			CreatedAt:            now,
			CreatedBy:            claims.FullName,
		}

		err = os.orderRepository.CreateOrderItem(ctx, &orderitem)
		if err != nil {
			return nil, err
		}
	}

	numbering.Number++
	err = os.orderRepository.UpdateNumbering(ctx, numbering)
	if err != nil {
		return nil, err
	}

	return &order.CreateOrderResponse{
		Id:   orderEntity.Id,
		Base: utils.SuccessResponse("Create order success"),
	}, nil
}

func NewOrderService(orderRepository repository.IOrderRepository, productRepository repository.IProductRepository) IOrderService {
	return &orderService{
		orderRepository:   orderRepository,
		productRepository: productRepository,
	}
}
