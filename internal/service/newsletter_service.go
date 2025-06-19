package service

import (
	"context"
	"time"

	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/entity"
	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/repository"
	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/utils"
	"github.com/fathurzoy/go-grpc-ecommerce-be/pb/newsletter"
	"github.com/google/uuid"
)

type INewsletterService interface {
	SubscribeNewsletter(ctx context.Context, request *newsletter.SubscribeNewsletterRequest) (*newsletter.SubscribeNewsletterResponse, error)
}

type newsletterService struct {
	newsletterRepository repository.INewsletterRepository
}

func (s *newsletterService) SubscribeNewsletter(ctx context.Context, request *newsletter.SubscribeNewsletterRequest) (*newsletter.SubscribeNewsletterResponse, error) {
	// cek ke db email
	newsletterEntity, err := s.newsletterRepository.GetNewsletterByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	if newsletterEntity != nil {
		return &newsletter.SubscribeNewsletterResponse{
			Base: utils.BadRequestResponse("Email already registered"),
		}, nil
	}

	// insert ke db
	newNewsletterEntity := entity.Newsletter{
		Id:        uuid.NewString(),
		FullName:  request.FullName,
		Email:     request.Email,
		CreatedAt: time.Now(),
		CreatedBy: "Public",
	}
	err = s.newsletterRepository.CreateNewNewsletter(ctx, &newNewsletterEntity)
	if err != nil {
		return nil, err
	}

	return &newsletter.SubscribeNewsletterResponse{
		Base: utils.SuccessResponse("Success"),
	}, nil
}

func NewNewsletterService(newsletterRepository repository.INewsletterRepository) INewsletterService {
	return &newsletterService{
		newsletterRepository: newsletterRepository,
	}
}
