package pricing

import (
	"context"
	"errors"
	"fmt"
	"github.com/alenalato/purchase-cart-service/internal/bizlogic"
	"github.com/alenalato/purchase-cart-service/internal/common"
	"github.com/alenalato/purchase-cart-service/internal/logger"
	"github.com/govalues/decimal"
	"github.com/thoas/go-funk"
)

func (l *Logic) GetItemsPrices(ctx context.Context, items []bizlogic.OrderDetailsItem) ([]bizlogic.ItemPrice, error) {
	// validate input
	validateErrs := make([]error, 0)
	for _, item := range items {
		validateErr := validate.Struct(item)
		if validateErr != nil {
			validateErrs = append(validateErrs, validateErr)
		}
	}
	if len(validateErrs) > 0 {
		logger.Log.Errorf("validation error: %v", validateErrs)

		return nil, common.NewError(errors.Join(validateErrs...), common.ErrTypeInvalidArgument)
	}

	// get unique product ids from items
	productIds := make([]int, 0, len(items))
	for _, item := range items {
		productIds = append(productIds, item.ProductId)
	}
	productIds = funk.Uniq(productIds).([]int)

	// get items prices from storage
	prices, pricesErr := l.pricingStorage.GetProductPrices(ctx, productIds)
	if pricesErr != nil {
		return nil, pricesErr
	}

	itemsPrices := make([]bizlogic.ItemPrice, len(items))
	for i, item := range items {
		productPrice, ok := prices[item.ProductId]
		if !ok {
			err := fmt.Errorf("product price not found for product id: %d", item.ProductId)
			logger.Log.Error(err)

			return nil, common.NewError(err, common.ErrTypeNotFound)
		}

		decimalQuantity, decErr := decimal.NewFromInt64(int64(item.Quantity), 0, 0)
		if decErr != nil {
			logger.Log.Errorf("error converting quantity to decimal: %v", decErr)

			return nil, common.NewError(decErr, common.ErrTypeInternal)
		}

		// item price = product price * quantity
		itemPrice, mulErr := productPrice.Price.Mul(decimalQuantity)
		if mulErr != nil {
			logger.Log.Errorf("error calculating item price: %v", mulErr)

			return nil, common.NewError(mulErr, common.ErrTypeInternal)
		}

		// calculate item vat
		itemVat, vatErr := l.vatManager.CalculateVat(ctx, productPrice.VatClass, itemPrice, nil)
		if vatErr != nil {
			return nil, vatErr
		}

		itemsPrices[i] = bizlogic.ItemPrice{
			Price: itemPrice,
			Vat:   itemVat,
		}
	}

	return itemsPrices, nil
}
