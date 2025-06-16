package utils

import "github.com/fathurzoy/go-grpc-ecommerce-be/pb/common"

func SuccessResponse(message string) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode: 200,
		Message:    message,
	}
}

func ValidationErrorResponse(validationErrors []*common.ValidationError) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode:      400,
		Message:         "Validation Error",
		IsError:         true,
		ValidationError: validationErrors,
	}
}
