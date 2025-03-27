package storage

import "context"

type OrderStorage interface {
	// CreateOrder creates a new order
	CreateOrder(ctx context.Context, orderDetails OrderDetails) (*Order, error)
}

type PricingStorage interface {
	// GetProductPrices returns pricing information for the given product IDs
	GetProductPrices(ctx context.Context, productIds []int) (map[int]ProductPrice, error)
}
