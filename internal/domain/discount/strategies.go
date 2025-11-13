package discount

import "github.com/mytheresa/go-hiring-challenge/internal/domain/product"

// CategoryDiscountStrategy applies discount based on product category.
type CategoryDiscountStrategy struct {
	categoryCode string
	percentage   int
}

// NewCategoryDiscountStrategy creates a discount strategy for a category.
func NewCategoryDiscountStrategy(categoryCode string, percentage int) *CategoryDiscountStrategy {
	return &CategoryDiscountStrategy{
		categoryCode: categoryCode,
		percentage:   percentage,
	}
}

// AppliesTo checks if the product belongs to the category.
func (s *CategoryDiscountStrategy) AppliesTo(p product.Product) bool {
	return p.Category != nil && p.Category.Code == s.categoryCode
}

// CalculatePercentage returns the discount percentage.
func (s *CategoryDiscountStrategy) CalculatePercentage(p product.Product) int {
	return s.percentage
}

// SKUDiscountStrategy applies discount based on product SKU/code.
type SKUDiscountStrategy struct {
	sku        string
	percentage int
}

// NewSKUDiscountStrategy creates a discount strategy for a SKU.
func NewSKUDiscountStrategy(sku string, percentage int) *SKUDiscountStrategy {
	return &SKUDiscountStrategy{
		sku:        sku,
		percentage: percentage,
	}
}

// AppliesTo checks if the product matches the SKU.
func (s *SKUDiscountStrategy) AppliesTo(p product.Product) bool {
	return p.Code == s.sku
}

// CalculatePercentage returns the configured discount percentage.
func (s *SKUDiscountStrategy) CalculatePercentage(p product.Product) int {
	return s.percentage
}
