package storage

import "context"

type OrderStorage interface {
	CreateOrder(ctx context.Context, orderDetails OrderDetails) (*Order, error)
}

type PricingStorage interface {
	GetProductPrices(ctx context.Context, productIds []int) (map[int]ProductPrice, error)
}
