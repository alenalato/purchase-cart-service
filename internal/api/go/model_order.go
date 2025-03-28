// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * Purchase Cart API
 *
 * Purchase Cart API is an API to create an order with items and get the pricing information for the order
 *
 * API version: 1.0
 */

package api

type Order struct {

	// The ID of the order
	Id string `json:"id,omitempty"`

	// The total price for the order
	TotalPrice float32 `json:"total_price,omitempty"`

	// The total VAT for the order
	TotalVat float32 `json:"total_vat,omitempty"`

	// The items in the order
	Items []OrderItem `json:"items,omitempty"`
}

// AssertOrderRequired checks if the required fields are not zero-ed
func AssertOrderRequired(obj Order) error {
	for _, el := range obj.Items {
		if err := AssertOrderItemRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertOrderConstraints checks if the values respects the defined constraints
func AssertOrderConstraints(obj Order) error {
	for _, el := range obj.Items {
		if err := AssertOrderItemConstraints(el); err != nil {
			return err
		}
	}
	return nil
}
