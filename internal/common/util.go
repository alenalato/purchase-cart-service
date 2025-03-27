package common

import "github.com/govalues/decimal"

// AsDecimal converts a float64 to a decimal.Decimal suppressing any errors
func AsDecimal(value float64) decimal.Decimal {
	dv, _ := decimal.NewFromFloat64(value)

	return dv
}
