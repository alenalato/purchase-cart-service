package mongodb

import (
	"context"
	"github.com/alenalato/purchase-cart-service/internal/common"
	"github.com/alenalato/purchase-cart-service/internal/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMongoDB_GetProductPrices_success(t *testing.T) {
	// 1 is repeated twice to test that the function can handle duplicates
	productIds := []int{1, 1, 2, 3}

	result, err := testMongoStorage.GetProductPrices(context.Background(), productIds)

	assert.Nil(t, err)
	assert.Equal(t, map[int]storage.ProductPrice{
		1: {
			ProductId: 1,
			Price:     common.AsDecimal(2),
			VatClass:  1,
		},
		2: {
			ProductId: 2,
			Price:     common.AsDecimal(1.5),
			VatClass:  1,
		},
		3: {
			ProductId: 3,
			Price:     common.AsDecimal(3),
			VatClass:  1,
		},
	}, result)
}

func TestMongoDB_GetProductPrices_emptyProductIds(t *testing.T) {
	productIds := make([]int, 0)

	result, err := testMongoStorage.GetProductPrices(context.Background(), productIds)

	assert.Nil(t, err)
	assert.Empty(t, result)
}
