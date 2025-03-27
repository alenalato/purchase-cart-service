package mongodb

import (
	"context"
	"github.com/alenalato/purchase-cart-service/internal/common"
	"github.com/alenalato/purchase-cart-service/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMongoDB_CreateOrder_success(t *testing.T) {
	orderDetails := storage.OrderDetails{
		TotalPrice: common.AsDecimal(100),
		TotalVat:   common.AsDecimal(20),
		Items: []storage.OrderDetailsItem{
			{
				ProductId: 1,
				Quantity:  2,
				Price:     common.AsDecimal(50),
				Vat:       common.AsDecimal(10),
			},
			{
				ProductId: 2,
				Quantity:  1,
				Price:     common.AsDecimal(40),
				Vat:       common.AsDecimal(20),
			},
		},
	}

	result, err := testMongoStorage.CreateOrder(context.Background(), orderDetails)

	assert.Nil(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Id)
	assert.Equal(t, orderDetails.TotalPrice, result.TotalPrice)
	assert.Equal(t, orderDetails.TotalVat, result.TotalVat)
	assert.Len(t, result.Items, 2)
	assert.Equal(t, orderDetails.Items[0].ProductId, result.Items[0].ProductId)
	assert.Equal(t, orderDetails.Items[0].Quantity, result.Items[0].Quantity)
	assert.Equal(t, orderDetails.Items[0].Price, result.Items[0].Price)
	assert.Equal(t, orderDetails.Items[0].Vat, result.Items[0].Vat)
	assert.Equal(t, orderDetails.Items[1].ProductId, result.Items[1].ProductId)
	assert.Equal(t, orderDetails.Items[1].Quantity, result.Items[1].Quantity)
	assert.Equal(t, orderDetails.Items[1].Price, result.Items[1].Price)
	assert.Equal(t, orderDetails.Items[1].Vat, result.Items[1].Vat)
}
