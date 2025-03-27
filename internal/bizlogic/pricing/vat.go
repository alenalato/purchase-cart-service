package pricing

import (
	"context"
	"github.com/alenalato/purchase-cart-service/internal/common"
	"github.com/govalues/decimal"
)

type VatManager interface {
	// CalculateVat calculates the VAT amount for a given amount
	CalculateVat(ctx context.Context, class vatClass, amount decimal.Decimal, vatContext interface{}) (decimal.Decimal, error)
}

// FixedVat is a VatManager that applies a fixed VAT rate
type FixedVat struct {
	fixedRate decimal.Decimal
}

var _ VatManager = new(FixedVat)

func (f *FixedVat) CalculateVat(_ context.Context, _ vatClass, amount decimal.Decimal, _ interface{}) (decimal.Decimal, error) {
	hundred, _ := decimal.New(100, 0)
	divided, divideErr := amount.Quo(hundred)
	if divideErr != nil {
		return decimal.Zero, common.NewError(divideErr, common.ErrTypeInternal)
	}
	multiplied, multiplyErr := divided.Mul(f.fixedRate)
	if multiplyErr != nil {
		return decimal.Zero, common.NewError(multiplyErr, common.ErrTypeInternal)
	}

	return multiplied, nil
}

// NewFixedVat creates a new FixedVat
func NewFixedVat(rate float64) (*FixedVat, error) {
	decimalRate, decimalErr := decimal.NewFromFloat64(rate)
	if decimalErr != nil {
		return nil, common.NewError(decimalErr, common.ErrTypeInternal)
	}

	return &FixedVat{
		fixedRate: decimalRate,
	}, nil
}
