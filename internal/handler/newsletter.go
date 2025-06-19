package handler

import (
	"context"

	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/service"
	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/utils"
	"github.com/fathurzoy/go-grpc-ecommerce-be/pb/newsletter"
)

type newsletterHandler struct {
	newsletter.UnimplementedNewsletterServiceServer

	newsletterService service.INewsletterService
}

func (h *newsletterHandler) SubscribeNewsletter(ctx context.Context, request *newsletter.SubscribeNewsletterRequest) (*newsletter.SubscribeNewsletterResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &newsletter.SubscribeNewsletterResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	// Proccess
	res, err := h.newsletterService.SubscribeNewsletter(ctx, request)
	if err != nil {
		return nil, err
	}

	return &newsletter.SubscribeNewsletterResponse{
		Base: res.Base,
	}, nil
}

func NewNewsletterHandler(newsletterService service.INewsletterService) newsletter.NewsletterServiceServer {
	return &newsletterHandler{
		newsletterService: newsletterService,
	}
}
