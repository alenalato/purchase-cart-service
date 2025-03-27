package order

import (
	"context"
	"github.com/alenalato/purchase-cart-service/internal/bizlogic"
	"github.com/alenalato/purchase-cart-service/internal/storage"
)

//go:generate mockgen -destination=converter_mock.go -package=order github.com/alenalato/purchase-cart-service/internal/bizlogic/order modelConverter

// modelConverter is an interface for converting between storage and bizlogic models
type modelConverter interface {
	// fromModelOrderDetailsToStorage converts a bizlogic.OrderDetails to a storage.OrderDetails
	fromModelOrderDetailsToStorage(ctx context.Context, order bizlogic.OrderDetails) storage.OrderDetails
	// fromStorageOrderToModel converts a storage.Order to a bizlogic.Order
	fromStorageOrderToModel(ctx context.Context, order *storage.Order) *bizlogic.Order
}

type storageModelConverter struct {
}

var _ modelConverter = new(storageModelConverter)

func (c *storageModelConverter) fromModelOrderDetailsToStorage(
	_ context.Context,
	order bizlogic.OrderDetails,
) storage.OrderDetails {
	items := make([]storage.OrderDetailsItem, len(order.Items))
	for i, item := range order.Items {
		items[i] = storage.OrderDetailsItem{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		}
	}

	return storage.OrderDetails{
		Items: items,
	}
}

func (c *storageModelConverter) fromStorageOrderToModel(_ context.Context, order *storage.Order) *bizlogic.Order {
	items := make([]bizlogic.OrderItem, len(order.Items))
	for i, item := range order.Items {
		items[i] = bizlogic.OrderItem{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Vat:       item.Vat,
		}
	}

	return &bizlogic.Order{
		Id:         order.Id,
		TotalPrice: order.TotalPrice,
		TotalVat:   order.TotalVat,
		Items:      items,
	}
}

// newStorageModelConverter creates a new storageModelConverter
func newStorageModelConverter() *storageModelConverter {
	return &storageModelConverter{}
}
