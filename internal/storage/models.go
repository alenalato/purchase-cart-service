package storage

import "github.com/govalues/decimal"

type ProductPrice struct {
	ProductId int             `bson:"product_id,omitempty"`
	Price     decimal.Decimal `bson:"price,omitempty"`
	VatClass  int             `bson:"vat_class,omitempty"`
}

type OrderItem struct {
	ProductId int             `bson:"product_id,omitempty"`
	Quantity  int             `bson:"quantity,omitempty"`
	Price     decimal.Decimal `bson:"price,omitempty"`
	Vat       decimal.Decimal `bson:"vat,omitempty"`
}

type Order struct {
	Id         string
	TotalPrice decimal.Decimal `bson:"total_price,omitempty"`
	TotalVat   decimal.Decimal `bson:"total_vat,omitempty"`
	Items      []OrderItem     `bson:"items,omitempty"`
}

type OrderDetailsItem struct {
	ProductId int             `bson:"product_id,omitempty"`
	Quantity  int             `bson:"quantity,omitempty"`
	Price     decimal.Decimal `bson:"price,omitempty"`
	Vat       decimal.Decimal `bson:"vat,omitempty"`
}

type OrderDetails struct {
	TotalPrice decimal.Decimal    `bson:"total_price,omitempty"`
	TotalVat   decimal.Decimal    `bson:"total_vat,omitempty"`
	Items      []OrderDetailsItem `bson:"items,omitempty"`
}
