package handler

import (
	"context"
	"fmt"

	"github.com/fathurzoy/go-grpc-ecommerce-be/pb/service"
)

// type IServiceHandler interface {
// 	HelloWorld(ctx context.Context, request *service.HelloWolrdRequest) (*service.HelloWorldResponse, error)
// }

type serviceHandler struct {
	service.UnimplementedHelloWorldServiceServer
}

func (sh *serviceHandler) HelloWorld(ctx context.Context, request *service.HelloWolrdRequest) (*service.HelloWorldResponse, error) {
	// panic(errors.New("pointer nil"))
	return &service.HelloWorldResponse{
		Message: fmt.Sprintf("Hello %s", request.Name),
	}, nil
}

func NewServiceHandler() *serviceHandler {
	return &serviceHandler{}
}
