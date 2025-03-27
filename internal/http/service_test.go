package http

import (
	"github.com/alenalato/purchase-cart-service/internal/bizlogic"
	"go.uber.org/mock/gomock"
	"reflect"
	"testing"
)

type testSuite struct {
	service            *PurchaseCartAPIService
	mockedOrderManager *bizlogic.MockOrderManager
	mockConverter      *MockmodelConverter
	mockCtrl           *gomock.Controller
}

func getTestSuite(t *testing.T) *testSuite {
	mockCtrl := gomock.NewController(t)

	return &testSuite{
		service:            NewPurchaseCartAPIService(bizlogic.NewMockOrderManager(mockCtrl)),
		mockedOrderManager: bizlogic.NewMockOrderManager(mockCtrl),
		mockConverter:      NewMockmodelConverter(mockCtrl),
		mockCtrl:           mockCtrl,
	}
}

func TestNewPurchaseCartAPIService(t *testing.T) {
	tSuite := getTestSuite(t)
	type args struct {
		orderManager bizlogic.OrderManager
	}
	tests := []struct {
		name string
		args args
		want *PurchaseCartAPIService
	}{
		{
			name: "Test NewPurchaseCartAPIService",
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
