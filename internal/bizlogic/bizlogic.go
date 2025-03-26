package bizlogic

import "context"

type PricingManager interface {
	GetItemsPrices(ctx context.Context, Items []ItemDetails) ([]ItemPrice, error)
}

type OrderManager interface {
	CreateOrder(ctx context.Context, details OrderDetails) (*Order, error)
}
