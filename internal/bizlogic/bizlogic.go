package bizlogic

import "context"

type PricingManager interface {
	// GetItemsPrices returns pricing information for the given items
	GetItemsPrices(ctx context.Context, items []ItemDetails) ([]ItemPrice, error)
}

type OrderManager interface {
	// CreateOrder creates a new order
	CreateOrder(ctx context.Context, details OrderDetails) (*Order, error)
}
