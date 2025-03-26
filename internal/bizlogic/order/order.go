package order

import (
	"github.com/alenalato/purchase-cart-service/internal/bizlogic"
	"github.com/alenalato/purchase-cart-service/internal/storage"
)

type Logic struct {
	pricingManager bizlogic.PricingManager
	orderStorage   storage.OrderStorage
	converter      modelConverter
}

var _ bizlogic.OrderManager = new(Logic)

func NewLogic(pricingManager bizlogic.PricingManager, orderStorage storage.OrderStorage) *Logic {
	return &Logic{
		pricingManager: pricingManager,
		orderStorage:   orderStorage,
		converter:      newStorageModelConverter(),
	}
}
