package bizlogic

import (
	"github.com/govalues/decimal"
)

type ItemDetails struct {
	ProductId int
	Quantity  int
}

type OrderDetails struct {
	Items []ItemDetails
}

type ItemPrice struct {
	Price decimal.Decimal
	Vat   decimal.Decimal
}

type OrderItem struct {
	ProductId int
	Quantity  int
	Price     decimal.Decimal
	Vat       decimal.Decimal
}

type Order struct {
	Id         string
	TotalPrice decimal.Decimal
	TotalVat   decimal.Decimal
	Items      []OrderItem
}
