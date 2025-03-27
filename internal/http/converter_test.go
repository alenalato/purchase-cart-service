package http

import (
	"context"
	api "github.com/alenalato/purchase-cart-service/internal/api/go"
	"github.com/alenalato/purchase-cart-service/internal/bizlogic"
	"github.com/govalues/decimal"
	"reflect"
	"testing"
)

func Test_apiModelConverter_fromApiCreateOrderRequestToModel(t *testing.T) {
	type fields struct {
		apiFloatPrecision int
	}
	type args struct {
		in0 context.Context
		req api.CreateOrderRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bizlogic.OrderDetails
	}{
		{
			name: "Test fromApiCreateOrderRequestToModel",
			fields: fields{
				apiFloatPrecision: 2,
			},
			args: args{
				in0: context.Background(),
				req: api.CreateOrderRequest{
					Order: api.CreateOrderRequestOrder{
						Items: []api.CreateOrderRequestOrderItemsInner{
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
			},
			want: bizlogic.OrderDetails{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &apiModelConverter{
				apiFloatPrecision: tt.fields.apiFloatPrecision,
			}
			if got := c.fromApiCreateOrderRequestToModel(tt.args.in0, tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromApiCreateOrderRequestToModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_apiModelConverter_fromModelOrderToApi(t *testing.T) {
	type fields struct {
		apiFloatPrecision int
	}
	type args struct {
		in0   context.Context
		order *bizlogic.Order
	}
	asDecimal := func(value float64) decimal.Decimal {
		dv, _ := decimal.NewFromFloat64(value)

		return dv
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    api.Order
		wantErr bool
	}{
		{
			name: "Test fromModelOrderToApi success",
			fields: fields{
				apiFloatPrecision: 2,
			},
			args: args{
				in0: context.Background(),
				order: &bizlogic.Order{
					Id:         "1",
					TotalPrice: asDecimal(2.222),
					TotalVat:   asDecimal(1.111),
					Items: []bizlogic.OrderItem{
						{
							ProductId: 1,
							Quantity:  2,
							Price:     asDecimal(2.5555),
							Vat:       asDecimal(1.1111),
						},
						{
							ProductId: 3,
							Quantity:  4,
							Price:     asDecimal(4.4444),
							Vat:       asDecimal(2.99999),
						},
					},
				},
			},
			want: api.Order{
				Id:         "1",
				TotalPrice: 2.22,
				TotalVat:   1.11,
				Items: []api.OrderItem{
					{
						ProductId: 1,
						Quantity:  2,
						Price:     2.56,
						Vat:       1.11,
					},
					{
						ProductId: 3,
						Quantity:  4,
						Price:     4.44,
						Vat:       3.00,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &apiModelConverter{
				apiFloatPrecision: tt.fields.apiFloatPrecision,
			}
			got, err := c.fromModelOrderToApi(tt.args.in0, tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("fromModelOrderToApi() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromModelOrderToApi() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newApiModelConverter(t *testing.T) {
	type args struct {
		apiFloatPrecision int
	}
	tests := []struct {
		name string
		args args
		want *apiModelConverter
	}{
		{
			name: "Test newApiModelConverter",
			args: args{
				apiFloatPrecision: 5,
			},
			want: &apiModelConverter{
				apiFloatPrecision: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newApiModelConverter(tt.args.apiFloatPrecision); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newApiModelConverter() = %v, want %v", got, tt.want)
			}
		})
	}
}
