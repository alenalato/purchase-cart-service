package pricing

import (
	"context"
	"github.com/govalues/decimal"
)

type VatManager interface {
	CalculateVat(ctx context.Context, class vatClass, amount decimal.Decimal, vatContext interface{}) (decimal.Decimal, error)
}

type FixedVat struct {
	fixedRate decimal.Decimal
}

func (f *FixedVat) CalculateVat(_ context.Context, _ vatClass, amount decimal.Decimal, _ interface{}) (decimal.Decimal, error) {
	hundred, _ := decimal.New(100, 0)
	divided, divideErr := amount.Quo(hundred)
	if divideErr != nil {
		return decimal.Zero, divideErr
	}

	return divided.Mul(f.fixedRate)
}

func NewFixedVat(rate float64) (*FixedVat, error) {
	decimalRate, decimalErr := decimal.NewFromFloat64(rate)
	if decimalErr != nil {
		return nil, decimalErr
	}

	return &FixedVat{
		fixedRate: decimalRate,
	}, nil
}
