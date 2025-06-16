package grpcmiddleware

import (
	"context"
	"log"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "internal server error: %v", r)
		}
	}()

	res, err := handler(ctx, req)
	if err != nil {
		log.Println("Error occurred while handling request:", err)

		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.Unauthenticated {
				return nil, err
			}
		}
		// return nil, status.Error(codes.Internal, "internal server error")

	}
	return res, err
}
