package http

import (
	"github.com/alenalato/purchase-cart-service/internal/bizlogic"
)

// apiFloatPrecision is the precision to use when converting float values to API values
const apiFloatPrecision uint = 2

type PurchaseCartAPIService struct {
	orderManager bizlogic.OrderManager
	converter    modelConverter
}

// NewPurchaseCartAPIService creates a PurchaseCartAPIService
func NewPurchaseCartAPIService(orderManager bizlogic.OrderManager) *PurchaseCartAPIService {
	return &PurchaseCartAPIService{
		orderManager: orderManager,
		converter:    newApiModelConverter(int(apiFloatPrecision)),
	}
}
