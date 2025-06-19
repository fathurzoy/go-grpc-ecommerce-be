package service

import (
	"context"
	"database/sql"
	"fmt"
	operationsystem "os"
	"runtime/debug"
	"time"

	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/entity"
	jwtentity "github.com/fathurzoy/go-grpc-ecommerce-be/internal/entity/jwt"
	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/repository"
	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/utils"
	"github.com/fathurzoy/go-grpc-ecommerce-be/pb/order"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/invoice"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IOrderService interface {
	CreateOrder(ctx context.Context, request *order.CreateOrderRequest) (*order.CreateOrderResponse, error)
	ListOrderAdmin(ctx context.Context, request *order.ListOrderAdminRequest) (*order.ListOrderAdminResponse, error)
	ListOrder(ctx context.Context, request *order.ListOrderRequest) (*order.ListOrderResponse, error)
}

type orderService struct {
	db                *sql.DB
	orderRepository   repository.IOrderRepository
	productRepository repository.IProductRepository
}

func (os *orderService) CreateOrder(ctx context.Context, request *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := os.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := recover(); e != nil {
			if tx != nil {
				tx.Rollback()
			}
			debug.PrintStack()
			panic(e)
		}
	}()
	defer func() {
		if err != nil && tx != nil {
			tx.Rollback()
		}
	}()

	orderRepo := os.orderRepository.WithTransaction(tx)
	productRepo := os.productRepository.WithTransaction(tx)

	numbering, err := orderRepo.GetNumbering(ctx, "order")
	if err != nil {
		return nil, err
	}

	var productIds = make([]string, len(request.Products))
	for i := range request.Products {
		productIds[i] = request.Products[i].Id
	}

	// cek apakah product ada
	products, err := productRepo.GetProductsByIds(ctx, productIds)
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
		Number:          fmt.Sprintf("ORD-%d%08d", now.Year(), numbering.Number),
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
	invoiceItems := make([]xendit.InvoiceItem, 0)
	for _, p := range request.Products {
		prod := productMap[p.Id]
		if p == nil {
			invoiceItems = append(invoiceItems, xendit.InvoiceItem{
				Name:     prod.Name,
				Quantity: int(p.Quantity),
				Price:    prod.Price,
			})
		}
	}

	xenditInvoice, xenditErr := invoice.CreateWithContext(ctx, &invoice.CreateParams{
		ExternalID: orderEntity.Id,
		Amount:     total,
		Customer: xendit.InvoiceCustomer{
			GivenNames: claims.FullName,
		},
		Currency:           "IDR",
		SuccessRedirectURL: fmt.Sprintf("%s/checkout/%s/success", operationsystem.Getenv("FRONTEND_URL"), orderEntity.Id),
		Items:              invoiceItems,
	})

	if xenditErr != nil {
		err = xenditErr
		return nil, err
	}

	orderEntity.XenditInvoiceId = &xenditInvoice.ID
	orderEntity.XenditInvoiceUrl = &xenditInvoice.InvoiceURL

	log.Info("Order createdx")
	err = orderRepo.CreateOrder(ctx, &orderEntity)
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

		err = orderRepo.CreateOrderItem(ctx, &orderitem)
		if err != nil {
			return nil, err
		}
	}

	numbering.Number++
	err = orderRepo.UpdateNumbering(ctx, numbering)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &order.CreateOrderResponse{
		Id:   orderEntity.Id,
		Base: utils.SuccessResponse("Create order success"),
	}, nil
}

func (os *orderService) ListOrderAdmin(ctx context.Context, request *order.ListOrderAdminRequest) (*order.ListOrderAdminResponse, error) {
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if claims.Role != entity.UserRoleAdmin {
		return nil, utils.UnauthenticatedResponse()
	}

	orders, metadata, err := os.orderRepository.GetListOrderAdminPagination(ctx, request.Pagination)
	if err != nil {
		return nil, err
	}

	items := make([]*order.ListOrderAdminResponseItem, 0)
	for _, orderEntity := range orders {

		products := make([]*order.ListOrderAdminResponseItemProduct, 0)

		for _, oi := range orderEntity.Items {
			products = append(products, &order.ListOrderAdminResponseItemProduct{
				Id:       oi.Id,
				Name:     oi.ProductName,
				Quantity: oi.Quantity,
				Price:    oi.ProductPrice,
			})
		}

		orderStatusCode := orderEntity.OrderStatusCode
		if orderEntity.OrderStatusCode == entity.OrderStatusCodeUnpaid && time.Now().After(*orderEntity.ExpiredAt) {
			orderEntity.OrderStatusCode = entity.OrderStatusCodeExpired
		}

		items = append(items, &order.ListOrderAdminResponseItem{
			Id:         orderEntity.Id,
			Number:     orderEntity.Number,
			Customer:   orderEntity.UserFullName,
			StatusCode: orderStatusCode,
			Total:      orderEntity.Total,
			CreatedAt:  timestamppb.New(orderEntity.CreatedAt),
			Products:   products,
		})
	}

	return &order.ListOrderAdminResponse{
		Base:       utils.SuccessResponse("List order success"),
		Pagination: metadata,
		Items:      items,
	}, nil //.nil, nil
}

func (os *orderService) ListOrder(ctx context.Context, request *order.ListOrderRequest) (*order.ListOrderResponse, error) {
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	orders, metadata, err := os.orderRepository.GetListOrderPagination(ctx, request.Pagination, claims.Subject)
	if err != nil {
		return nil, err
	}

	items := make([]*order.ListOrderResponseItem, 0)
	for _, orderEntity := range orders {

		products := make([]*order.ListOrderResponseItemProduct, 0)

		for _, oi := range orderEntity.Items {
			products = append(products, &order.ListOrderResponseItemProduct{
				Id:       oi.Id,
				Name:     oi.ProductName,
				Quantity: oi.Quantity,
				Price:    oi.ProductPrice,
			})
		}

		orderStatusCode := orderEntity.OrderStatusCode
		if orderEntity.OrderStatusCode == entity.OrderStatusCodeUnpaid && time.Now().After(*orderEntity.ExpiredAt) {
			orderEntity.OrderStatusCode = entity.OrderStatusCodeExpired
		}

		xenditInvoiceUrl := ""
		if orderEntity.XenditInvoiceUrl != nil {
			xenditInvoiceUrl = *orderEntity.XenditInvoiceUrl
		}
		items = append(items, &order.ListOrderResponseItem{
			Id:               orderEntity.Id,
			Number:           orderEntity.Number,
			Customer:         orderEntity.UserFullName,
			StatusCode:       orderStatusCode,
			Total:            orderEntity.Total,
			CreatedAt:        timestamppb.New(orderEntity.CreatedAt),
			Products:         products,
			XenditInvoiceUrl: xenditInvoiceUrl,
		})
	}

	return &order.ListOrderResponse{
		Base:       utils.SuccessResponse("List order success"),
		Pagination: metadata,
		Items:      items,
	}, nil //.nil, nil
}

func NewOrderService(db *sql.DB, orderRepository repository.IOrderRepository, productRepository repository.IProductRepository) IOrderService {
	return &orderService{
		db:                db,
		orderRepository:   orderRepository,
		productRepository: productRepository,
	}
}
