package pricing

import (
	"context"
	"errors"
	"github.com/alenalato/purchase-cart-service/internal/bizlogic"
	"github.com/govalues/decimal"
)

func (l *Logic) GetItemsPrices(ctx context.Context, Items []bizlogic.ItemDetails) ([]bizlogic.ItemPrice, error) {
	// TODO validate request

	productIds := make([]int, 0, len(Items))
	for _, item := range Items {
		productIds = append(productIds, item.ProductId)
	}

	// Get items prices from storage
	prices, pricesErr := l.pricingStorage.GetProductPrices(ctx, productIds)
	if pricesErr != nil {
		return nil, pricesErr
	}

	itemsPrices := make([]bizlogic.ItemPrice, len(Items))
	for i, item := range Items {
		productPrice, ok := prices[item.ProductId]
		if !ok {
			return nil, errors.New("product price not found")
		}
		decimalQuantity, decErr := decimal.NewFromInt64(int64(item.Quantity), 0, 0)
		if decErr != nil {
			return nil, decErr
		}
		// item price = product price * quantity
		itemPrice, mulErr := productPrice.Price.Mul(decimalQuantity)
		if mulErr != nil {
			return nil, mulErr
		}
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
