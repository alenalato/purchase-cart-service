package bizlogic

import "context"

//go:generate mockgen -destination=bizlogic_mock.go -package=bizlogic github.com/alenalato/purchase-cart-service/internal/bizlogic PricingManager,OrderManager

type PricingManager interface {
	// GetItemsPrices returns pricing information for the given items
	GetItemsPrices(ctx context.Context, items []OrderDetailsItem) ([]ItemPrice, error)
}

type OrderManager interface {
	// CreateOrder creates a new order
	CreateOrder(ctx context.Context, details OrderDetails) (*Order, error)
}
