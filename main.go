package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/grpcmiddleware"
	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/handler"
	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/repository"
	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/service"
	"github.com/fathurzoy/go-grpc-ecommerce-be/pb/auth"
	"github.com/fathurzoy/go-grpc-ecommerce-be/pb/product"

	// "github.com/fathurzoy/go-grpc-ecommerce-be/pb/service"
	"github.com/fathurzoy/go-grpc-ecommerce-be/pkg/database"
	"github.com/joho/godotenv"
	gocache "github.com/patrickmn/go-cache"
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

	db := database.ConnectDB(ctx, os.Getenv("DB_URI"))
	log.Println("Database connected")

	cacheService := gocache.New(time.Hour*24, time.Hour)

	authMiddleware := grpcmiddleware.NewAuthMiddleware(cacheService)

	// serviceHandler := handler.NewServiceHandler()

	authRepository := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepository, cacheService)
	authHandler := handler.NewAuthHandler(authService)

	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	serv := grpc.NewServer(grpc.ChainUnaryInterceptor(grpcmiddleware.ErrorMiddleware, authMiddleware.Middleware))

	auth.RegisterAuthServiceServer(serv, authHandler)
	product.RegisterProductSeriviceServer(serv, productHandler)

	// service.RegisterHelloWorldServiceServer(serv, serviceHandler)

	if os.Getenv("ENVIRONMENT") == "dev" {
		reflection.Register(serv)
		log.Println("Reflection server is running on port 50051")
	}

	log.Println("server is running on port 50051")
	if err := serv.Serve(lis); err != nil {
		log.Panicf("failed to serve: %v", err)
	}
}
