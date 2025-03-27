package pricing

import (
	"context"
	"github.com/govalues/decimal"
	"math"
	"reflect"
	"testing"
)

func TestFixedVat_CalculateVat(t *testing.T) {
	type fields struct {
		fixedRate decimal.Decimal
	}
	type args struct {
		in0    context.Context
		in1    vatClass
		amount decimal.Decimal
		in3    interface{}
	}
	asDecimal := func(value float64) decimal.Decimal {
		dv, _ := decimal.NewFromFloat64(value)

		return dv
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    decimal.Decimal
		wantErr bool
	}{
		{
			name: "Test CalculateVat success",
			fields: fields{
				fixedRate: asDecimal(20),
			},
			args: args{
				in0:    context.Background(),
				in1:    1,
				amount: asDecimal(100),
				in3:    nil,
			},
			want:    asDecimal(20),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FixedVat{
				fixedRate: tt.fields.fixedRate,
			}
			got, err := f.CalculateVat(tt.args.in0, tt.args.in1, tt.args.amount, tt.args.in3)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateVat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateVat() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFixedVat(t *testing.T) {
	type args struct {
		rate float64
	}
	asDecimal := func(value float64) decimal.Decimal {
		dv, _ := decimal.NewFromFloat64(value)

		return dv
	}
	tests := []struct {
		name    string
		args    args
		want    *FixedVat
		wantErr bool
	}{
		{
			name: "Test NewFixedVat success",
			args: args{
				rate: 20,
			},
			want: &FixedVat{
				fixedRate: asDecimal(20),
			},
			wantErr: false,
		},
		{
			name: "Test NewFixedVat error",
			args: args{
				rate: math.Inf(1),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFixedVat(tt.args.rate)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFixedVat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFixedVat() got = %v, want %v", got, tt.want)
			}
		})
	}
}
