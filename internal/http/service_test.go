package http

import (
	"github.com/alenalato/purchase-cart-service/internal/bizlogic"
	"go.uber.org/mock/gomock"
	"reflect"
	"testing"
)

type testSuite struct {
	purchaseCartService *PurchaseCartAPIService
	mockedOrderManager  *bizlogic.MockOrderManager
	mockedConverter     *MockmodelConverter
	mockCtrl            *gomock.Controller
}

func getTestSuite(t *testing.T) *testSuite {
	mockCtrl := gomock.NewController(t)
	mockedOrderManager := bizlogic.NewMockOrderManager(mockCtrl)
	mockedConverter := NewMockmodelConverter(mockCtrl)
	purchaseCartService := NewPurchaseCartAPIService(mockedOrderManager)
	purchaseCartService.converter = mockedConverter

	return &testSuite{
		purchaseCartService: purchaseCartService,
		mockedOrderManager:  mockedOrderManager,
		mockedConverter:     mockedConverter,
		mockCtrl:            mockCtrl,
	}
}

func TestNewPurchaseCartAPIService(t *testing.T) {
	tSuite := getTestSuite(t)
	defer tSuite.mockCtrl.Finish()

	type args struct {
		orderManager bizlogic.OrderManager
	}
	tests := []struct {
		name string
		args args
		want *PurchaseCartAPIService
	}{
		{
			name: "NewPurchaseCartAPIService",
			args: args{
				orderManager: tSuite.mockedOrderManager,
			},
			want: &PurchaseCartAPIService{
				orderManager: tSuite.mockedOrderManager,
				converter:    newApiModelConverter(int(apiFloatPrecision)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPurchaseCartAPIService(tt.args.orderManager); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPurchaseCartAPIService() = %v, want %v", got, tt.want)
			}
		})
	}
}
