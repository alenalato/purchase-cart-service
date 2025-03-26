package http

import (
	"context"
	api "github.com/alenalato/purchase-cart-service/internal/api/go"
)

// V1OrderPost creates an order
func (s *PurchaseCartAPIService) V1OrderPost(
	ctx context.Context,
	createOrderRequest api.CreateOrderRequest,
) (api.ImplResponse, error) {
	order, createErr := s.orderManager.CreateOrder(
		ctx,
		s.converter.fromApiCreateOrderRequestToModel(ctx, createOrderRequest),
	)
	if createErr != nil {
		// TODO: Uncomment the next line to return response Response(400, Error{}) or use other options such as http.Ok ...
		// return Response(400, Error{}), nil
		// TODO: Uncomment the next line to return response Response(500, Error{}) or use other options such as http.Ok ...
		// return Response(500, Error{}), nil
	}
	if order == nil {
		// TODO: Uncomment the next line to return response Response(500, Error{}) or use other options such as http.Ok ...
		// return Response(500, Error{}), nil
	}

	res, convertErr := s.converter.fromModelOrderToApi(ctx, order)
	if convertErr != nil {
		// TODO: Uncomment the next line to return response Response(500, Error{}) or use other options such as http.Ok ...
		// return Response(500, Error{}), nil
	}

	return api.Response(201, res), nil
}
