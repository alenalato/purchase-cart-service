package pricing

import (
	"github.com/alenalato/purchase-cart-service/internal/storage"
	"go.uber.org/mock/gomock"
	"reflect"
	"testing"
)

type testSuite struct {
	logic                *Logic
	mockedVatManager     *MockVatManager
	mockedPricingStorage *storage.MockPricingStorage
	mockCtrl             *gomock.Controller
}

func getTestSuite(t *testing.T) *testSuite {
	mockCtrl := gomock.NewController(t)

	return &testSuite{
		logic:                NewLogic(NewMockVatManager(mockCtrl), storage.NewMockPricingStorage(mockCtrl)),
		mockedVatManager:     NewMockVatManager(mockCtrl),
		mockedPricingStorage: storage.NewMockPricingStorage(mockCtrl),
		mockCtrl:             mockCtrl,
	}
}

func TestNewLogic(t *testing.T) {
	tSuite := getTestSuite(t)
	type args struct {
		vatManager     VatManager
		pricingStorage storage.PricingStorage
	}
	tests := []struct {
		name string
		args args
		want *Logic
	}{
		{
			name: "Test NewLogic",
			args: args{
				vatManager:     tSuite.mockedVatManager,
				pricingStorage: tSuite.mockedPricingStorage,
			},
			want: &Logic{
				vatManager:     tSuite.mockedVatManager,
				pricingStorage: tSuite.mockedPricingStorage,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLogic(tt.args.vatManager, tt.args.pricingStorage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLogic() = %v, want %v", got, tt.want)
			}
		})
	}
}
