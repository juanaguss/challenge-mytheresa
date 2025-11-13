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

// AppliesTo checks if the product belongs to the target category.
func (s *CategoryDiscountStrategy) AppliesTo(p product.Product) bool {
	return p.Category != nil && p.Category.Code == s.categoryCode
}

// CalculatePercentage returns the configured discount percentage.
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

// AppliesTo checks if the product or any of its variants match the SKU.
func (s *SKUDiscountStrategy) AppliesTo(p product.Product) bool {
	if p.Code == s.sku {
		return true
	}
	// Check if any variant has this SKU
	for _, v := range p.Variants {
		if v.SKU == s.sku {
			return true
		}
	}
	return false
}

// CalculatePercentage returns the configured discount percentage.
func (s *SKUDiscountStrategy) CalculatePercentage(p product.Product) int {
	return s.percentage
}

// AppliesToVariant checks if this discount applies to a specific variant.
func (s *SKUDiscountStrategy) AppliesToVariant(sku string) bool {
	return s.sku == sku
}
