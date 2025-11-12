package product

import (
	"github.com/shopspring/decimal"
)

// Product represents a product in the catalog.
type Product struct {
	ID         uint
	Code       string
	Price      decimal.Decimal
	CategoryID *uint
	Category   *Category
	Variants   []Variant
}

// Category represents a product category.
type Category struct {
	ID   uint
	Code string
	Name string
}

// Variant represents a product variant with optional pricing.
type Variant struct {
	ID        uint
	ProductID uint
	Name      string
	SKU       string
	Price     decimal.Decimal
}
