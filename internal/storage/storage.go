package storage

import "context"

//go:generate mockgen -destination=storage_mock.go -package=storage github.com/alenalato/purchase-cart-service/internal/storage OrderStorage,PricingStorage

type OrderStorage interface {
	// CreateOrder creates a new order
	CreateOrder(ctx context.Context, orderDetails OrderDetails) (*Order, error)
}

type PricingStorage interface {
	// GetProductPrices returns pricing information for the given product IDs
	GetProductPrices(ctx context.Context, productIds []int) (map[int]ProductPrice, error)
}
