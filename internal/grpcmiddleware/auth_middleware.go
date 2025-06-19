package grpcmiddleware

import (
	"context"
	"log"

	jwtentity "github.com/fathurzoy/go-grpc-ecommerce-be/internal/entity/jwt"
	"github.com/fathurzoy/go-grpc-ecommerce-be/internal/utils"
	gocache "github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
)

type authMiddleware struct {
	cacheService *gocache.Cache
}

var publicApis = map[string]bool{
	"/auth.AuthService/Login":                   true,
	"/auth.AuthService/Register":                true,
	"/product.ProductService/DetailProduct":     true,
	"/product.ProductService/ListProduct":       true,
	"/product.ProductSerivice/HighlightProduct": true,
}

func (am *authMiddleware) Middleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	log.Println(info.FullMethod)
	if publicApis[info.FullMethod] {
		return handler(ctx, req)
	}

	// if info.FullMethod == "/auth.AuthService/Login" || info.FullMethod == "/auth.AuthService/Register" || info.FullMethod == "/auth.AuthService/DetailProduct" {
	// 	return handler(ctx, req)
	// }

	//ambil token dari metadata
	tokenStr, err := jwtentity.ParseTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	//cek token dari logout cache
	_, ok := am.cacheService.Get(tokenStr)
	if ok {
		return nil, utils.UnauthenticatedResponse()
	}

	//parse jwtnya hingga jadi entity
	claims, err := jwtentity.GetClaimsFromToken(tokenStr)
	if err != nil {
		return nil, err
	}

	//sematkan entity ke context
	ctx = claims.SetToContext(ctx)

	res, err := handler(ctx, req)

	return res, err
}

func NewAuthMiddleware(cacheService *gocache.Cache) *authMiddleware {
	return &authMiddleware{
		cacheService: cacheService,
	}
}
