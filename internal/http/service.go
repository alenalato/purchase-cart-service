package http

import (
	"github.com/alenalato/purchase-cart-service/internal/bizlogic"
)

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
