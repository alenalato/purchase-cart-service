package order

import (
	"context"
	"errors"
	"github.com/alenalato/purchase-cart-service/internal/bizlogic"
)

func (l *Logic) CreateOrder(ctx context.Context, details bizlogic.OrderDetails) (*bizlogic.Order, error) {
	// TODO request validation

	itemsPrices, pricesErr := l.pricingManager.GetItemsPrices(ctx, details.Items)
	if pricesErr != nil {
		return nil, pricesErr
	}
	if len(itemsPrices) != len(details.Items) {
		return nil, errors.New("prices count mismatch")
	}

	storageOrderDetails := l.converter.fromModelOrderDetailsToStorage(ctx, details)
	if len(storageOrderDetails.Items) != len(details.Items) {
		return nil, errors.New("items count mismatch")
	}

	for i, item := range itemsPrices {
		storageOrderDetails.Items[i].Price = item.Price
		storageOrderDetails.Items[i].Vat = item.Vat

		var priceAddErr error
		storageOrderDetails.TotalPrice, priceAddErr = storageOrderDetails.TotalPrice.Add(item.Price)
		if priceAddErr != nil {
			return nil, priceAddErr
		}
		var vatAddErr error
		storageOrderDetails.TotalVat, vatAddErr = storageOrderDetails.TotalVat.Add(item.Vat)
		if vatAddErr != nil {
			return nil, vatAddErr
		}
	}

	storageOrder, createErr := l.orderStorage.CreateOrder(ctx, storageOrderDetails)
	if createErr != nil {
		return nil, createErr
	}

	order := l.converter.fromStorageOrderToModel(ctx, storageOrder)

	return order, nil
}
