package http

import (
	"context"
	"errors"
	api "github.com/alenalato/purchase-cart-service/internal/api/go"
	"github.com/alenalato/purchase-cart-service/internal/bizlogic"
	"github.com/alenalato/purchase-cart-service/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestPurchaseCartAPIService_V1OrderPostV1OrderPost_success(t *testing.T) {
	tSuite := getTestSuite(t)
	defer tSuite.mockCtrl.Finish()

	createOrderRequest := api.CreateOrderRequest{
		Order: api.CreateOrderRequestOrder{
			Items: []api.CreateOrderRequestOrderItemsInner{
				{
					ProductId: 1,
					Quantity:  1,
				},
				{
					ProductId: 2,
					Quantity:  2,
				},
			},
		},
	}
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
	tSuite.mockedConverter.EXPECT().fromApiCreateOrderRequestToModel(gomock.Any(), createOrderRequest).
		Return(orderDetails)

	order := &bizlogic.Order{
		Id:         "1",
		TotalPrice: common.AsDecimal(11),
		TotalVat:   common.AsDecimal(1.1),
		Items: []bizlogic.OrderItem{
			{
				ProductId: 1,
				Quantity:  1,
				Price:     common.AsDecimal(5),
				Vat:       common.AsDecimal(0.5),
			},
			{
				ProductId: 2,
				Quantity:  2,
				Price:     common.AsDecimal(6),
				Vat:       common.AsDecimal(0.6),
			},
		},
	}

	tSuite.mockedOrderManager.EXPECT().CreateOrder(gomock.Any(), orderDetails).
		Return(order, nil)

	apiOrder := api.Order{
		Id:         "1",
		TotalPrice: 11,
		TotalVat:   1.1,
		Items: []api.OrderItem{
			{
				ProductId: 1,
				Quantity:  1,
				Price:     5,
				Vat:       0.5,
			},
			{
				ProductId: 2,
				Quantity:  2,
				Price:     6,
				Vat:       0.6,
			},
		},
	}
	tSuite.mockedConverter.EXPECT().fromModelOrderToApi(gomock.Any(), order).
		Return(apiOrder, nil)

	response, err := tSuite.purchaseCartService.V1OrderPost(context.Background(), createOrderRequest)

	assert.Nil(t, err)
	assert.Equal(t, 201, response.Code)
	assert.Equal(t, apiOrder, response.Body)
}

func TestPurchaseCartAPIService_V1OrderPost_error(t *testing.T) {
	createOrderRequest := api.CreateOrderRequest{
		Order: api.CreateOrderRequestOrder{
			Items: []api.CreateOrderRequestOrderItemsInner{
				{
					ProductId: 1,
					Quantity:  1,
				},
				{
					ProductId: 2,
					Quantity:  2,
				},
			},
		},
	}
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
	order := &bizlogic.Order{
		Id:         "1",
		TotalPrice: common.AsDecimal(11),
		TotalVat:   common.AsDecimal(1.1),
		Items: []bizlogic.OrderItem{
			{
				ProductId: 1,
				Quantity:  1,
				Price:     common.AsDecimal(5),
				Vat:       common.AsDecimal(0.5),
			},
			{
				ProductId: 2,
				Quantity:  2,
				Price:     common.AsDecimal(6),
				Vat:       common.AsDecimal(0.6),
			},
		},
	}

	tests := []struct {
		name        string
		setup       func(tSuite *testSuite)
		want        api.ImplResponse
		wantMessage bool
	}{
		{
			name: "create order error invalid argument",
			setup: func(tSuite *testSuite) {
				tSuite.mockedConverter.EXPECT().fromApiCreateOrderRequestToModel(gomock.Any(), createOrderRequest).
					Return(orderDetails)
				tSuite.mockedOrderManager.EXPECT().CreateOrder(gomock.Any(), orderDetails).
					Return(nil, common.NewError(errors.New("invalid test argument"), common.ErrTypeInvalidArgument))
			},
			want: api.Response(400, api.Error{
				Code: errorCodeInvalidArgument.String(),
			}),
			wantMessage: true,
		},
		{
			name: "create order error not found",
			setup: func(tSuite *testSuite) {
				tSuite.mockedConverter.EXPECT().fromApiCreateOrderRequestToModel(gomock.Any(), createOrderRequest).
					Return(orderDetails)
				tSuite.mockedOrderManager.EXPECT().CreateOrder(gomock.Any(), orderDetails).
					Return(nil, common.NewError(errors.New("test not found"), common.ErrTypeNotFound))
			},
			want: api.Response(422, api.Error{
				Code: errorCodeUnprocessableEntity.String(),
			}),
			wantMessage: true,
		},
		{
			name: "create order generic error",
			setup: func(tSuite *testSuite) {
				tSuite.mockedConverter.EXPECT().fromApiCreateOrderRequestToModel(gomock.Any(), createOrderRequest).
					Return(orderDetails)
				tSuite.mockedOrderManager.EXPECT().CreateOrder(gomock.Any(), orderDetails).
					Return(nil, errors.New("test error"))
			},
			want: api.Response(500, api.Error{
				Code: errorCodeInternal.String(),
			}),
			wantMessage: false,
		},
		{
			name: "unexpected nil order",
			setup: func(tSuite *testSuite) {
				tSuite.mockedConverter.EXPECT().fromApiCreateOrderRequestToModel(gomock.Any(), createOrderRequest).
					Return(orderDetails)
				tSuite.mockedOrderManager.EXPECT().CreateOrder(gomock.Any(), orderDetails).
					Return(nil, nil)
			},
			want: api.Response(500, api.Error{
				Code: errorCodeInternal.String(),
			}),
			wantMessage: false,
		},
		{
			name: "convert order to response error",
			setup: func(tSuite *testSuite) {
				tSuite.mockedConverter.EXPECT().fromApiCreateOrderRequestToModel(gomock.Any(), createOrderRequest).
					Return(orderDetails)
				tSuite.mockedOrderManager.EXPECT().CreateOrder(gomock.Any(), orderDetails).
					Return(order, nil)
				tSuite.mockedConverter.EXPECT().fromModelOrderToApi(gomock.Any(), order).
					Return(api.Order{}, errors.New("test convert error"))
			},
			want: api.Response(500, api.Error{
				Code: errorCodeInternal.String(),
			}),
			wantMessage: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tSuite := getTestSuite(t)
			defer tSuite.mockCtrl.Finish()

			tt.setup(tSuite)

			response, err := tSuite.purchaseCartService.V1OrderPost(context.Background(), createOrderRequest)

			assert.Nil(t, err)
			assert.Equal(t, tt.want.Code, response.Code)
			require.IsType(t, api.Error{}, response.Body)
			assert.Equal(t, tt.want.Body.(api.Error).Code, response.Body.(api.Error).Code)
			if tt.wantMessage {
				assert.NotEmpty(t, response.Body.(api.Error).Message)
			}
		})
	}
}
