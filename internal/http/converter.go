package http

import (
	"context"
	"errors"
	api "github.com/alenalato/purchase-cart-service/internal/api/go"
	"github.com/alenalato/purchase-cart-service/internal/bizlogic"
	"github.com/alenalato/purchase-cart-service/internal/common"
)

//go:generate mockgen -destination=converter_mock.go -package=http github.com/alenalato/purchase-cart-service/internal/http modelConverter

// modelConverter is an interface for converting between API and bizlogic models
type modelConverter interface {
	// fromApiCreateOrderRequestToModel converts an API CreateOrderRequest to a bizlogic.OrderDetails
	fromModelOrderToApi(ctx context.Context, order *bizlogic.Order) (api.Order, error)
	// fromModelOrderToApi converts a bizlogic.Order to an api.Order
	fromApiCreateOrderRequestToModel(ctx context.Context, req api.CreateOrderRequest) bizlogic.OrderDetails
}

type apiModelConverter struct {
	apiFloatPrecision int
}

var _ modelConverter = new(apiModelConverter)

func (c *apiModelConverter) fromApiCreateOrderRequestToModel(_ context.Context, req api.CreateOrderRequest) bizlogic.OrderDetails {
	items := make([]bizlogic.OrderDetailsItem, len(req.Order.Items))

	for i, item := range req.Order.Items {
		items[i] = bizlogic.OrderDetailsItem{
			ProductId: int(item.ProductId),
			Quantity:  int(item.Quantity),
		}
	}

	return bizlogic.OrderDetails{
		Items: items,
	}
}

func (c *apiModelConverter) fromModelOrderToApi(_ context.Context, order *bizlogic.Order) (api.Order, error) {
	items := make([]api.OrderItem, len(order.Items))

	for i, item := range order.Items {
		itemPrice, ok := item.Price.Round(c.apiFloatPrecision).Float64()
		if !ok {
			return api.Order{}, common.NewError(
				errors.New("could not convert item price to float"),
				common.ErrTypeInternal,
			)
		}
		itemVat, ok := item.Vat.Round(c.apiFloatPrecision).Float64()
		if !ok {
			return api.Order{}, common.NewError(
				errors.New("could not convert item vat to float"),
				common.ErrTypeInternal,
			)
		}

		items[i] = api.OrderItem{
			ProductId: int32(item.ProductId),
			Quantity:  int32(item.Quantity),
			Price:     float32(itemPrice),
			Vat:       float32(itemVat),
		}
	}

	totalPrice, ok := order.TotalPrice.Round(c.apiFloatPrecision).Float64()
	if !ok {
		return api.Order{}, common.NewError(
			errors.New("could not convert order total price to float"),
			common.ErrTypeInternal,
		)
	}
	totalVat, ok := order.TotalVat.Round(c.apiFloatPrecision).Float64()
	if !ok {
		return api.Order{}, common.NewError(
			errors.New("could not convert order total vat to float"),
			common.ErrTypeInternal,
		)
	}

	return api.Order{
		Id:         order.Id,
		TotalPrice: float32(totalPrice),
		TotalVat:   float32(totalVat),
		Items:      items,
	}, nil
}

// newApiModelConverter creates a new apiModelConverter
func newApiModelConverter(apiFloatPrecision int) *apiModelConverter {
	return &apiModelConverter{apiFloatPrecision}
}
