package order

import (
	"context"
	"github.com/alenalato/purchase-cart-service/internal/bizlogic"
	"github.com/alenalato/purchase-cart-service/internal/common"
	"github.com/alenalato/purchase-cart-service/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestLogic_CreateOrder_success(t *testing.T) {
	tSuite := getTestSuite(t)
	defer tSuite.mockCtrl.Finish()

	orderDetails := bizlogic.OrderDetails{
		Items: []bizlogic.OrderDetailsItem{
			{
				ProductId: 1,
				Quantity:  1,
			},
			{
				ProductId: 2,
				Quantity:  2,
			},
		},
	}
	tSuite.mockedPricingManager.EXPECT().GetItemsPrices(gomock.Any(), orderDetails.Items).
		Return([]bizlogic.ItemPrice{
			{
				Price: common.AsDecimal(10),
				Vat:   common.AsDecimal(2),
			},
			{
				Price: common.AsDecimal(20),
				Vat:   common.AsDecimal(4),
			},
		}, nil)

	storageOrderDetails := storage.OrderDetails{
		TotalPrice: common.AsDecimal(0),
		TotalVat:   common.AsDecimal(0),
		Items: []storage.OrderDetailsItem{
			{
				ProductId: 1,
				Quantity:  1,
				Price:     common.AsDecimal(10),
				Vat:       common.AsDecimal(2),
			},
			{
				ProductId: 2,
				Quantity:  2,
				Price:     common.AsDecimal(20),
				Vat:       common.AsDecimal(4),
			},
		},
	}

	tSuite.mockedConverter.EXPECT().fromModelOrderDetailsToStorage(gomock.Any(), orderDetails).
		Return(storageOrderDetails)

	storageOrderDetailsWithTotals := storageOrderDetails
	storageOrderDetailsWithTotals.TotalPrice = common.AsDecimal(30)
	storageOrderDetailsWithTotals.TotalVat = common.AsDecimal(6)

	storageOrder := storage.Order{
		Id:         "1",
		TotalPrice: common.AsDecimal(30),
		TotalVat:   common.AsDecimal(6),
		Items: []storage.OrderItem{
			{
				ProductId: 1,
				Quantity:  1,
				Price:     common.AsDecimal(10),
				Vat:       common.AsDecimal(2),
			},
			{
				ProductId: 2,
				Quantity:  2,
				Price:     common.AsDecimal(20),
				Vat:       common.AsDecimal(4),
			},
		},
	}

	tSuite.mockedOrderStorage.EXPECT().CreateOrder(gomock.Any(), storageOrderDetailsWithTotals).
		Return(&storageOrder, nil)

	order := bizlogic.Order{
		Id:         "1",
		TotalPrice: common.AsDecimal(50),
		TotalVat:   common.AsDecimal(10),
		Items: []bizlogic.OrderItem{
			{
				ProductId: 1,
				Quantity:  1,
				Price:     common.AsDecimal(10),
				Vat:       common.AsDecimal(2),
			},
			{
				ProductId: 2,
				Quantity:  2,
				Price:     common.AsDecimal(20),
				Vat:       common.AsDecimal(4),
			},
		},
	}

	tSuite.mockedConverter.EXPECT().fromStorageOrderToModel(gomock.Any(), &storageOrder).
		Return(&order)

	result, err := tSuite.logic.CreateOrder(context.Background(), orderDetails)

	assert.Nil(t, err)
	require.NotNil(t, result)
	assert.Equal(t, order, *result)
}

func TestLogic_CreateOrder_validateError(t *testing.T) {
	tSuite := getTestSuite(t)
	defer tSuite.mockCtrl.Finish()

	orderDetails := bizlogic.OrderDetails{
		Items: []bizlogic.OrderDetailsItem{
			{
				ProductId: 0,
				Quantity:  1,
			},
			{
				ProductId: -1,
				Quantity:  2,
			},
		},
	}

	result, err := tSuite.logic.CreateOrder(context.Background(), orderDetails)

	assert.Error(t, err)
	var commonErr common.Error
	require.ErrorAs(t, err, &commonErr)
	assert.Equal(t, common.ErrTypeInvalidArgument, err.(common.Error).GetType())
	assert.Nil(t, result)
}

func TestLogic_GetItemsPrices_pricingManagerError(t *testing.T) {
	tSuite := getTestSuite(t)
	defer tSuite.mockCtrl.Finish()

	orderDetails := bizlogic.OrderDetails{
		Items: []bizlogic.OrderDetailsItem{
			{
				ProductId: 1,
				Quantity:  1,
			},
			{
				ProductId: 2,
				Quantity:  2,
			},
		},
	}

	expectedErr := common.NewError(nil, common.ErrTypeNotFound)

	tSuite.mockedPricingManager.EXPECT().GetItemsPrices(gomock.Any(), orderDetails.Items).
		Return(nil, expectedErr)

	result, err := tSuite.logic.CreateOrder(context.Background(), orderDetails)

	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, expectedErr.GetType(), err.(common.Error).GetType())
	assert.Nil(t, result)
}

func TestLogic_GetItemsPrices_itemsCountMismatch(t *testing.T) {
	tSuite := getTestSuite(t)
	defer tSuite.mockCtrl.Finish()

	orderDetails := bizlogic.OrderDetails{
		Items: []bizlogic.OrderDetailsItem{
			{
				ProductId: 1,
				Quantity:  1,
			},
			{
				ProductId: 2,
				Quantity:  2,
			},
		},
	}
	tSuite.mockedPricingManager.EXPECT().GetItemsPrices(gomock.Any(), orderDetails.Items).
		Return([]bizlogic.ItemPrice{
			{
				Price: common.AsDecimal(10),
				Vat:   common.AsDecimal(2),
			},
			// missing price for the second item to trigger the error
		}, nil)

	storageOrderDetails := storage.OrderDetails{
		TotalPrice: common.AsDecimal(0),
		TotalVat:   common.AsDecimal(0),
		Items: []storage.OrderDetailsItem{
			{
				ProductId: 1,
				Quantity:  1,
				Price:     common.AsDecimal(10),
				Vat:       common.AsDecimal(2),
			},
			{
				ProductId: 2,
				Quantity:  2,
				Price:     common.AsDecimal(20),
				Vat:       common.AsDecimal(4),
			},
		},
	}

	tSuite.mockedConverter.EXPECT().fromModelOrderDetailsToStorage(gomock.Any(), orderDetails).
		Return(storageOrderDetails)

	result, err := tSuite.logic.CreateOrder(context.Background(), orderDetails)

	assert.Error(t, err)
	var commonErr common.Error
	require.ErrorAs(t, err, &commonErr)
	assert.Equal(t, common.ErrTypeInternal, err.(common.Error).GetType())
	assert.Nil(t, result)
}

func TestLogic_GetItemsPrices_storageError(t *testing.T) {
	tSuite := getTestSuite(t)
	defer tSuite.mockCtrl.Finish()

	orderDetails := bizlogic.OrderDetails{
		Items: []bizlogic.OrderDetailsItem{
			{
				ProductId: 1,
				Quantity:  1,
			},
			{
				ProductId: 2,
				Quantity:  2,
			},
		},
	}

	tSuite.mockedPricingManager.EXPECT().GetItemsPrices(gomock.Any(), orderDetails.Items).
		Return([]bizlogic.ItemPrice{
			{
				Price: common.AsDecimal(10),
				Vat:   common.AsDecimal(2),
			},
			{
				Price: common.AsDecimal(20),
				Vat:   common.AsDecimal(4),
			},
		}, nil)

	storageOrderDetails := storage.OrderDetails{
		TotalPrice: common.AsDecimal(0),
		TotalVat:   common.AsDecimal(0),
		Items: []storage.OrderDetailsItem{
			{
				ProductId: 1,
				Quantity:  1,
				Price:     common.AsDecimal(10),
				Vat:       common.AsDecimal(2),
			},
			{
				ProductId: 2,
				Quantity:  2,
				Price:     common.AsDecimal(20),
				Vat:       common.AsDecimal(4),
			},
		},
	}

	tSuite.mockedConverter.EXPECT().fromModelOrderDetailsToStorage(gomock.Any(), orderDetails).
		Return(storageOrderDetails)

	storageOrderDetailsWithTotals := storageOrderDetails
	storageOrderDetailsWithTotals.TotalPrice = common.AsDecimal(30)
	storageOrderDetailsWithTotals.TotalVat = common.AsDecimal(6)

	expectedErr := common.NewError(nil, common.ErrTypeInternal)

	tSuite.mockedOrderStorage.EXPECT().CreateOrder(gomock.Any(), storageOrderDetailsWithTotals).
		Return(nil, expectedErr)

	result, err := tSuite.logic.CreateOrder(context.Background(), orderDetails)

	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, expectedErr.GetType(), err.(common.Error).GetType())
	assert.Nil(t, result)
}
