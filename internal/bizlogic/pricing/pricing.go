package pricing

import (
	"github.com/alenalato/purchase-cart-service/internal/bizlogic"
	"github.com/alenalato/purchase-cart-service/internal/storage"
	"github.com/go-playground/validator/v10"
)

// vatClass is a type for VAT classes
type vatClass = int

var validate = validator.New(validator.WithRequiredStructEnabled())

type Logic struct {
	vatManager     VatManager
	pricingStorage storage.PricingStorage
}

var _ bizlogic.PricingManager = new(Logic)

// NewLogic creates a new pricing logic
func NewLogic(vatManager VatManager, pricingStorage storage.PricingStorage) *Logic {
	return &Logic{
		vatManager:     vatManager,
		pricingStorage: pricingStorage,
	}
}
