package order

import (
	"github.com/alenalato/purchase-cart-service/internal/bizlogic"
	"github.com/alenalato/purchase-cart-service/internal/storage"
	"go.uber.org/mock/gomock"
	"reflect"
	"testing"
)

type testSuite struct {
	logic                *Logic
	mockedPricingManager *bizlogic.MockPricingManager
	mockedOrderStorage   *storage.MockOrderStorage
	mockedConverter      *MockmodelConverter
	mockCtrl             *gomock.Controller
}

func getTestSuite(t *testing.T) *testSuite {
	mockCtrl := gomock.NewController(t)
	mockedPricingManager := bizlogic.NewMockPricingManager(mockCtrl)
	mockedOrderStorage := storage.NewMockOrderStorage(mockCtrl)
	mockedConverter := NewMockmodelConverter(mockCtrl)
	logic := NewLogic(mockedPricingManager, mockedOrderStorage)
	logic.converter = mockedConverter

	return &testSuite{
		logic:                logic,
		mockedPricingManager: mockedPricingManager,
		mockedOrderStorage:   mockedOrderStorage,
		mockedConverter:      mockedConverter,
		mockCtrl:             mockCtrl,
	}
}

func TestNewLogic(t *testing.T) {
	tSuite := getTestSuite(t)
	defer tSuite.mockCtrl.Finish()

	type args struct {
		pricingManager bizlogic.PricingManager
		orderStorage   storage.OrderStorage
	}
	tests := []struct {
		name string
		args args
		want *Logic
	}{
		{
			name: "NewLogic",
			args: args{
				pricingManager: tSuite.mockedPricingManager,
				orderStorage:   tSuite.mockedOrderStorage,
			},
			want: &Logic{
				pricingManager: tSuite.mockedPricingManager,
				orderStorage:   tSuite.mockedOrderStorage,
				converter:      newStorageModelConverter(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLogic(tt.args.pricingManager, tt.args.orderStorage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLogic() = %v, want %v", got, tt.want)
			}
		})
	}
}
