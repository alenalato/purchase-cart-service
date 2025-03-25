package http

import (
	"context"
	"errors"
	api "github.com/alenalato/purchase-cart-service/internal/api/go"
	"net/http"
)

type PurchaseCartAPIService struct {
}

// NewPurchaseCartAPIService creates a default api service
func NewPurchaseCartAPIService() *PurchaseCartAPIService {
	return &PurchaseCartAPIService{}
}

// ApiV1OrderPost creates a new order
func (s *PurchaseCartAPIService) V1OrderPost(ctx context.Context, createOrderRequest api.CreateOrderRequest) (api.ImplResponse, error) {
	// TODO: Uncomment the next line to return response Response(201, Order{}) or use other options such as http.Ok ...
	// return Response(201, Order{}), nil

	// TODO: Uncomment the next line to return response Response(400, Error{}) or use other options such as http.Ok ...
	// return Response(400, Error{}), nil

	// TODO: Uncomment the next line to return response Response(500, Error{}) or use other options such as http.Ok ...
	// return Response(500, Error{}), nil

	return api.Response(http.StatusNotImplemented, nil), errors.New("ApiV1OrderPost will be implemented")
}
