package order

import (
	"context"
	"github.com/alenalato/purchase-cart-service/internal/bizlogic"
	"github.com/alenalato/purchase-cart-service/internal/common"
	"github.com/alenalato/purchase-cart-service/internal/storage"
	"reflect"
	"testing"
)

func Test_newStorageModelConverter(t *testing.T) {
	tests := []struct {
		name string
		want *storageModelConverter
	}{
		{
			name: "newStorageModelConverter",
			want: &storageModelConverter{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newStorageModelConverter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newStorageModelConverter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_storageModelConverter_fromModelOrderDetailsToStorage(t *testing.T) {
	type args struct {
		in0   context.Context
		order bizlogic.OrderDetails
	}
	tests := []struct {
		name string
		args args
		want storage.OrderDetails
	}{
		{
			name: "fromModelOrderDetailsToStorage",
			args: args{
				in0: context.Background(),
				order: bizlogic.OrderDetails{
					Items: []bizlogic.OrderDetailsItem{
						{
							ProductId: 1,
							Quantity:  2,
						},
						{
							ProductId: 3,
							Quantity:  4,
						},
					},
				},
			},
			want: storage.OrderDetails{
				Items: []storage.OrderDetailsItem{
					{
						ProductId: 1,
						Quantity:  2,
					},
					{
						ProductId: 3,
						Quantity:  4,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &storageModelConverter{}
			if got := c.fromModelOrderDetailsToStorage(tt.args.in0, tt.args.order); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromModelOrderDetailsToStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_storageModelConverter_fromStorageOrderToModel(t *testing.T) {
	type args struct {
		in0   context.Context
		order *storage.Order
	}
	tests := []struct {
		name string
		args args
		want *bizlogic.Order
	}{
		{
			name: "fromStorageOrderToModel",
			args: args{
				in0: context.Background(),
				order: &storage.Order{
					Id:         "1",
					TotalPrice: common.AsDecimal(2),
					TotalVat:   common.AsDecimal(3),
					Items: []storage.OrderItem{
						{
							ProductId: 4,
							Quantity:  5,
							Price:     common.AsDecimal(6),
							Vat:       common.AsDecimal(7),
						},
						{
							ProductId: 8,
							Quantity:  9,
							Price:     common.AsDecimal(10),
							Vat:       common.AsDecimal(11),
						},
					},
				},
			},
			want: &bizlogic.Order{
				Id:         "1",
				TotalPrice: common.AsDecimal(2),
				TotalVat:   common.AsDecimal(3),
				Items: []bizlogic.OrderItem{
					{
						ProductId: 4,
						Quantity:  5,
						Price:     common.AsDecimal(6),
						Vat:       common.AsDecimal(7),
					},
					{
						ProductId: 8,
						Quantity:  9,
						Price:     common.AsDecimal(10),
						Vat:       common.AsDecimal(11),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &storageModelConverter{}
			if got := c.fromStorageOrderToModel(tt.args.in0, tt.args.order); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromStorageOrderToModel() = %v, want %v", got, tt.want)
			}
		})
	}
}
