package bizlogic

import (
	"github.com/govalues/decimal"
)

// OrderDetailsItem represents the details of an item in an order to be created
type OrderDetailsItem struct {
	ProductId int `validate:"required,gt=0"`
	Quantity  int `validate:"required,gt=0"`
}

// OrderDetails represents the details of an order to be created
type OrderDetails struct {
	Items []OrderDetailsItem `validate:"required,dive"`
}

// ItemPrice represents the price and VAT of an item
type ItemPrice struct {
	Price decimal.Decimal
	Vat   decimal.Decimal
}

// OrderItem represents an item in an order
type OrderItem struct {
	ProductId int
	Quantity  int
	Price     decimal.Decimal
	Vat       decimal.Decimal
}

// Order represents an order
type Order struct {
	Id         string
	TotalPrice decimal.Decimal
	TotalVat   decimal.Decimal
	Items      []OrderItem
}
