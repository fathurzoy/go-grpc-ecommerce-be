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
	"github.com/fathurzoy/go-grpc-ecommerce-be/pb/cart"
	"github.com/fathurzoy/go-grpc-ecommerce-be/pb/newsletter"
	"github.com/fathurzoy/go-grpc-ecommerce-be/pb/order"
	"github.com/fathurzoy/go-grpc-ecommerce-be/pb/product"
	"github.com/xendit/xendit-go"

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

	xendit.Opt.SecretKey = os.Getenv("XENDIT_SECRET_KEY")

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

	cartRepository := repository.NewCartRepository(db)
	cartService := service.NewCartService(productRepository, cartRepository)
	cartHandler := handler.NewCartHandler(cartService)

	orderRepository := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(db, orderRepository, productRepository)
	orderHandler := handler.NewOrderHandler(orderService)

	newsletterRepository := repository.NewNewsletterRepository(db)
	newsletterService := service.NewNewsletterService(newsletterRepository)
	newsletterHandler := handler.NewNewsletterHandler(newsletterService)

	serv := grpc.NewServer(grpc.ChainUnaryInterceptor(grpcmiddleware.ErrorMiddleware, authMiddleware.Middleware))

	auth.RegisterAuthServiceServer(serv, authHandler)
	product.RegisterProductSeriviceServer(serv, productHandler)
	cart.RegisterCartServiceServer(serv, cartHandler)
	order.RegisterOrderServiceServer(serv, orderHandler)
	newsletter.RegisterNewsletterServiceServer(serv, newsletterHandler)

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
