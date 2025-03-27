package pricing

import (
	"github.com/alenalato/purchase-cart-service/internal/bizlogic"
	"github.com/alenalato/purchase-cart-service/internal/common"
	"github.com/alenalato/purchase-cart-service/internal/storage"
	"github.com/govalues/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestLogic_GetItemsPrices_success(t *testing.T) {
	tSuite := getTestSuite(t)
	defer tSuite.mockCtrl.Finish()

	orderItems := []bizlogic.OrderDetailsItem{
		{
			ProductId: 1,
			Quantity:  1,
		},
		{
			ProductId: 2,
			Quantity:  2,
		},
		// repeated product to test unique product ids
		{
			ProductId: 2,
			Quantity:  1,
		},
	}
	tSuite.mockedPricingStorage.EXPECT().GetProductPrices(gomock.Any(), []int{1, 2}).
		Return(map[int]storage.ProductPrice{
			1: {
				ProductId: 1,
				Price:     common.AsDecimal(11),
				VatClass:  1,
			},
			2: {
				ProductId: 2,
				Price:     common.AsDecimal(22),
				VatClass:  2,
			},
		}, nil)

	tSuite.mockedVatManager.EXPECT().CalculateVat(gomock.Any(), 1, common.AsDecimal(11), nil).
		Return(common.AsDecimal(1.1), nil)
	tSuite.mockedVatManager.EXPECT().CalculateVat(gomock.Any(), 2, common.AsDecimal(44), nil).
		Return(common.AsDecimal(4.4), nil)
	tSuite.mockedVatManager.EXPECT().CalculateVat(gomock.Any(), 2, common.AsDecimal(22), nil).
		Return(common.AsDecimal(2.2), nil)

	result, err := tSuite.logic.GetItemsPrices(nil, orderItems)

	assert.Nil(t, err)
	assert.Equal(t, []bizlogic.ItemPrice{
		{
			Price: common.AsDecimal(11),
			Vat:   common.AsDecimal(1.1),
		},
		{
			Price: common.AsDecimal(44),
			Vat:   common.AsDecimal(4.4),
		},
		{
			Price: common.AsDecimal(22),
			Vat:   common.AsDecimal(2.2),
		},
	}, result)
}

func TestLogic_GetItemsPrices_validateError(t *testing.T) {
	tSuite := getTestSuite(t)
	defer tSuite.mockCtrl.Finish()

	orderItems := []bizlogic.OrderDetailsItem{
		{
			ProductId: 1,
			Quantity:  0,
		},
		{
			ProductId: 0,
			Quantity:  -1,
		},
	}

	result, err := tSuite.logic.GetItemsPrices(nil, orderItems)

	assert.Error(t, err)
	var commonErr common.Error
	require.ErrorAs(t, err, &commonErr)
	assert.Equal(t, common.ErrTypeInvalidArgument, err.(common.Error).GetType())
	assert.Nil(t, result)
}

func TestLogic_GetItemsPrices_storageError(t *testing.T) {
	tSuite := getTestSuite(t)
	defer tSuite.mockCtrl.Finish()

	orderItems := []bizlogic.OrderDetailsItem{
		{
			ProductId: 1,
			Quantity:  1,
		},
		{
			ProductId: 2,
			Quantity:  2,
		},
	}

	expectedErr := common.NewError(nil, common.ErrTypeNotFound)

	tSuite.mockedPricingStorage.EXPECT().GetProductPrices(gomock.Any(), []int{1, 2}).
		Return(nil, expectedErr)

	result, err := tSuite.logic.GetItemsPrices(nil, orderItems)

	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, expectedErr.GetType(), err.(common.Error).GetType())
	assert.Nil(t, result)
}

func TestLogic_GetItemsPrices_productPriceNotFoundError(t *testing.T) {
	tSuite := getTestSuite(t)
	defer tSuite.mockCtrl.Finish()

	orderItems := []bizlogic.OrderDetailsItem{
		{
			ProductId: 1,
			Quantity:  1,
		},
		{
			ProductId: 2,
			Quantity:  2,
		},
	}

	tSuite.mockedPricingStorage.EXPECT().GetProductPrices(gomock.Any(), []int{1, 2}).
		Return(map[int]storage.ProductPrice{
			1: {
				ProductId: 1,
				Price:     common.AsDecimal(11),
				VatClass:  1,
			},
		}, nil)

	tSuite.mockedVatManager.EXPECT().CalculateVat(gomock.Any(), 1, common.AsDecimal(11), nil).
		Return(common.AsDecimal(1.1), nil)

	result, err := tSuite.logic.GetItemsPrices(nil, orderItems)

	assert.Error(t, err)
	var commonErr common.Error
	require.ErrorAs(t, err, &commonErr)
	assert.Equal(t, common.ErrTypeNotFound, err.(common.Error).GetType())
	assert.Nil(t, result)
}

func TestLogic_GetItemsPrices_vatManagerError(t *testing.T) {
	tSuite := getTestSuite(t)
	defer tSuite.mockCtrl.Finish()

	orderItems := []bizlogic.OrderDetailsItem{
		{
			ProductId: 1,
			Quantity:  1,
		},
		{
			ProductId: 2,
			Quantity:  2,
		},
	}
	tSuite.mockedPricingStorage.EXPECT().GetProductPrices(gomock.Any(), []int{1, 2}).
		Return(map[int]storage.ProductPrice{
			1: {
				ProductId: 1,
				Price:     common.AsDecimal(11),
				VatClass:  1,
			},
			2: {
				ProductId: 2,
				Price:     common.AsDecimal(22),
				VatClass:  2,
			},
		}, nil)

	// first call will be successful
	tSuite.mockedVatManager.EXPECT().CalculateVat(gomock.Any(), 1, common.AsDecimal(11), nil).
		Return(common.AsDecimal(1.1), nil)
	// second call will fail
	expectedErr := common.NewError(nil, common.ErrTypeInternal)
	tSuite.mockedVatManager.EXPECT().CalculateVat(gomock.Any(), 2, common.AsDecimal(44), nil).
		Return(decimal.Zero, expectedErr)

	result, err := tSuite.logic.GetItemsPrices(nil, orderItems)

	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, expectedErr.GetType(), err.(common.Error).GetType())
	assert.Nil(t, result)
}
