package product

import "github.com/shopspring/decimal"

// Filter contains information for filtering products.
type Filter struct {
	Category      string
	PriceLessThan *decimal.Decimal
}
