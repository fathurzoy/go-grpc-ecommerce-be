package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/handler"
	"github.com/fathurzoy/go-grpc-ecommerce-be/pb/service"
	"github.com/fathurzoy/go-grpc-ecommerce-be/pkg/database"
	"github.com/fathurzoy/go-grpc-ecommerce-be/pkg/grpcmiddleware"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()
	godotenv.Load()

	lis, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Panicf("failed to listen: %v", err)
	}

	database.ConnectDB(ctx, os.Getenv("DB_URI"))
	log.Println("Database connected")

	serviceHandler := handler.NewServiceHandler()

	serv := grpc.NewServer(grpc.ChainUnaryInterceptor(grpcmiddleware.ErrorMiddleware))

	service.RegisterHelloWorldServiceServer(serv, serviceHandler)

	if os.Getenv("ENVIRONMENT") == "dev" {
		reflection.Register(serv)
		log.Println("Reflection server is running on port 50051")
	}

	log.Println("server is running on port 50051")
	if err := serv.Serve(lis); err != nil {
		log.Panicf("failed to serve: %v", err)
	}
}
