package order

import (
	"github.com/alenalato/purchase-cart-service/internal/bizlogic"
	"github.com/alenalato/purchase-cart-service/internal/storage"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type Logic struct {
	pricingManager bizlogic.PricingManager
	orderStorage   storage.OrderStorage
	converter      modelConverter
}

var _ bizlogic.OrderManager = new(Logic)

// NewLogic creates a new order logic
func NewLogic(pricingManager bizlogic.PricingManager, orderStorage storage.OrderStorage) *Logic {
	return &Logic{
		pricingManager: pricingManager,
		orderStorage:   orderStorage,
		converter:      newStorageModelConverter(),
	}
}
