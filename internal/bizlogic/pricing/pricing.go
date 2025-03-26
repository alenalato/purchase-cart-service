package pricing

import (
	"github.com/alenalato/purchase-cart-service/internal/bizlogic"
	"github.com/alenalato/purchase-cart-service/internal/storage"
)

type vatClass = int

type Logic struct {
	vatManager     VatManager
	pricingStorage storage.PricingStorage
}

var _ bizlogic.PricingManager = new(Logic)

func NewLogic(vaManager VatManager, pricingStorage storage.PricingStorage) *Logic {
	return &Logic{
		vatManager:     vaManager,
		pricingStorage: pricingStorage,
	}
}
