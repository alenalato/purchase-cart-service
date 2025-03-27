package storage

import "github.com/govalues/decimal"

// ProductPrice represents the price and VAT class of a product
type ProductPrice struct {
	ProductId int             `bson:"product_id,omitempty"`
	Price     decimal.Decimal `bson:"price,omitempty"`
	VatClass  int             `bson:"vat_class,omitempty"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ProductId int             `bson:"product_id,omitempty"`
	Quantity  int             `bson:"quantity,omitempty"`
	Price     decimal.Decimal `bson:"price,omitempty"`
	Vat       decimal.Decimal `bson:"vat,omitempty"`
}

// Order represents an order
type Order struct {
	Id         string
	TotalPrice decimal.Decimal `bson:"total_price,omitempty"`
	TotalVat   decimal.Decimal `bson:"total_vat,omitempty"`
	Items      []OrderItem     `bson:"items,omitempty"`
}

// OrderDetailsItem represents an item in an order to be created
type OrderDetailsItem struct {
	ProductId int             `bson:"product_id,omitempty"`
	Quantity  int             `bson:"quantity,omitempty"`
	Price     decimal.Decimal `bson:"price,omitempty"`
	Vat       decimal.Decimal `bson:"vat,omitempty"`
}

// OrderDetails represents the details of an order to be created
type OrderDetails struct {
	TotalPrice decimal.Decimal    `bson:"total_price,omitempty"`
	TotalVat   decimal.Decimal    `bson:"total_vat,omitempty"`
	Items      []OrderDetailsItem `bson:"items,omitempty"`
}
