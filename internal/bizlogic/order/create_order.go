package order

import (
	"context"
	"errors"
	"github.com/alenalato/purchase-cart-service/internal/bizlogic"
	"github.com/alenalato/purchase-cart-service/internal/common"
	"github.com/alenalato/purchase-cart-service/internal/logger"
)

func (l *Logic) CreateOrder(ctx context.Context, details bizlogic.OrderDetails) (*bizlogic.Order, error) {
	// validate input
	validateErr := validate.Struct(details)
	if validateErr != nil {
		logger.Log.Errorf("validation error: %v", validateErr)

		return nil, common.NewError(validateErr, common.ErrTypeInvalidArgument)
	}

	// get prices for all items
	itemsPrices, pricesErr := l.pricingManager.GetItemsPrices(ctx, details.Items)
	if pricesErr != nil {
		return nil, pricesErr
	}
	if len(itemsPrices) != len(details.Items) {
		err := errors.New("prices count mismatch")
		logger.Log.Error(err)

		return nil, common.NewError(err, common.ErrTypeInternal)
	}

	// prepare storage order details
	storageOrderDetails := l.converter.fromModelOrderDetailsToStorage(ctx, details)
	if len(storageOrderDetails.Items) != len(details.Items) {
		err := errors.New("items count mismatch")
		logger.Log.Error(err)

		return nil, common.NewError(err, common.ErrTypeInternal)
	}

	// calculate total price and total vat while setting prices to items
	for i, item := range itemsPrices {
		storageOrderDetails.Items[i].Price = item.Price
		storageOrderDetails.Items[i].Vat = item.Vat

		var priceAddErr error
		storageOrderDetails.TotalPrice, priceAddErr = storageOrderDetails.TotalPrice.Add(item.Price)
		if priceAddErr != nil {
			logger.Log.Errorf("error adding item price: %v", priceAddErr)

			return nil, common.NewError(priceAddErr, common.ErrTypeInternal)
		}

		var vatAddErr error
		storageOrderDetails.TotalVat, vatAddErr = storageOrderDetails.TotalVat.Add(item.Vat)
		if vatAddErr != nil {
			logger.Log.Errorf("error adding item vat: %v", vatAddErr)

			return nil, common.NewError(vatAddErr, common.ErrTypeInternal)
		}
	}

	// create order in storage
	storageOrder, createErr := l.orderStorage.CreateOrder(ctx, storageOrderDetails)
	if createErr != nil {
		return nil, createErr
	}

	// convert storage order to model order
	order := l.converter.fromStorageOrderToModel(ctx, storageOrder)

	return order, nil
}
