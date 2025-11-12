package models

import (
	"github.com/shopspring/decimal"
)

// Product represents a product in the catalog.
// It includes a unique code, a price, and belongs to a category.
type Product struct {
	ID         uint            `gorm:"primaryKey"`
	Code       string          `gorm:"uniqueIndex;not null"`
	Price      decimal.Decimal `gorm:"type:decimal(10,2);not null"`
	CategoryID *uint           `gorm:"index"`
	Category   *Category       `gorm:"foreignKey:CategoryID"`
	Variants   []Variant       `gorm:"foreignKey:ProductID"`
}

func (p *Product) TableName() string {
	return "products"
}

// ProductFilters represents the available filters for querying products.
// Category filters by the category code.
// PriceLessThan filters products with price below the given value.
type ProductFilters struct {
	Category      string
	PriceLessThan *decimal.Decimal
}
