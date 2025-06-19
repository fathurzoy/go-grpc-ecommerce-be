package service

import (
	"context"
	"errors"
	"time"

	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/dto"
	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/entity"
	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/repository"
)

type IWebhookService interface {
	ReceiveInvoice(ctx context.Context, request *dto.XenditInvoiceRequest) error
}

type webhookService struct {
	orderRepository repository.IOrderRepository
}

func (w *webhookService) ReceiveInvoice(ctx context.Context, request *dto.XenditInvoiceRequest) error {
	// find order di db
	orderEntity, err := w.orderRepository.GetOrderById(ctx, request.ExternalID)
	if err != nil {
		return err
	}
	if orderEntity == nil {
		return errors.New("order not found")
	}

	// ganti / update entity
	now := time.Now()
	updatedBy := "System"
	orderEntity.OrderStatusCode = entity.OrderStatusCodePaid
	orderEntity.UpdatedBy = &updatedBy
	orderEntity.XenditPaidAt = &now
	orderEntity.XenditPaymentMethod = &request.PaymentMethod
	orderEntity.XenditPaymentChannel = &request.PaymentChannel

	// update ke db
	err = w.orderRepository.UpdateOrder(ctx, orderEntity)
	if err != nil {
		return err
	}

	return nil
}

func NewWebhookService(orderRepository repository.IOrderRepository) IWebhookService {
	return &webhookService{
		orderRepository: orderRepository,
	}
}
