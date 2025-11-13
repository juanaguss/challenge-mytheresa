package discount

import (
	"testing"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCategoryDiscountStrategy(t *testing.T) {
	t.Run("applies to product with matching category", func(t *testing.T) {
		strategy := NewCategoryDiscountStrategy("boots", 30)
		prod := product.Product{
			Code:  "PROD001",
			Price: decimal.NewFromFloat(100.0),
			Category: &product.Category{
				Code: "boots",
				Name: "Boots",
			},
		}

		applies := strategy.AppliesTo(prod)
		percentage := strategy.CalculatePercentage(prod)

		assert.True(t, applies)
		assert.Equal(t, 30, percentage)
	})

	t.Run("does not apply to product with different category", func(t *testing.T) {
		strategy := NewCategoryDiscountStrategy("boots", 30)
		prod := product.Product{
			Code:  "PROD001",
			Price: decimal.NewFromFloat(100.0),
			Category: &product.Category{
				Code: "shoes",
				Name: "Shoes",
			},
		}

		applies := strategy.AppliesTo(prod)

		assert.False(t, applies)
	})

	t.Run("does not apply to product without category", func(t *testing.T) {
		strategy := NewCategoryDiscountStrategy("boots", 30)
		prod := product.Product{
			Code:     "PROD001",
			Price:    decimal.NewFromFloat(100.0),
			Category: nil,
		}

		applies := strategy.AppliesTo(prod)

		assert.False(t, applies)
	})
}

func TestSKUDiscountStrategy(t *testing.T) {
	t.Run("applies to product with equal SKU", func(t *testing.T) {
		strategy := NewSKUDiscountStrategy("000003", 15)
		prod := product.Product{
			Code:  "000003",
			Price: decimal.NewFromFloat(100.0),
		}

		applies := strategy.AppliesTo(prod)
		percentage := strategy.CalculatePercentage(prod)

		assert.True(t, applies)
		assert.Equal(t, 15, percentage)
	})

	t.Run("does not apply to product with different SKU", func(t *testing.T) {
		strategy := NewSKUDiscountStrategy("000003", 15)
		prod := product.Product{
			Code:  "000001",
			Price: decimal.NewFromFloat(100.0),
		}

		applies := strategy.AppliesTo(prod)

		assert.False(t, applies)
	})
}

func TestEngine(t *testing.T) {
	t.Run("applies first matching discount", func(t *testing.T) {
		strategies := []Strategy{
			NewCategoryDiscountStrategy("boots", 30),
			NewSKUDiscountStrategy("000003", 15),
		}
		engine := NewEngine(strategies)

		prod := product.Product{
			Code:  "PROD001",
			Price: decimal.NewFromFloat(100.0),
			Category: &product.Category{
				Code: "boots",
			},
		}

		discountedPrice := engine.ApplyDiscount(prod)
		percentage := engine.GetDiscountPercentage(prod)

		assert.Equal(t, "70", discountedPrice.String())
		assert.Equal(t, 30, percentage)
	})

	t.Run("returns original price when no discount is applicable", func(t *testing.T) {
		strategies := []Strategy{
			NewCategoryDiscountStrategy("boots", 30),
			NewSKUDiscountStrategy("000003", 15),
		}
		engine := NewEngine(strategies)

		prod := product.Product{
			Code:  "PROD001",
			Price: decimal.NewFromFloat(100.0),
			Category: &product.Category{
				Code: "shoes",
			},
		}

		discountedPrice := engine.ApplyDiscount(prod)
		percentage := engine.GetDiscountPercentage(prod)

		assert.Equal(t, "100", discountedPrice.String())
		assert.Equal(t, 0, percentage)
	})

	t.Run("does not do discount stacking", func(t *testing.T) {
		// Product fits both strategies
		strategies := []Strategy{
			NewSKUDiscountStrategy("000003", 15),
			NewCategoryDiscountStrategy("boots", 30),
		}
		engine := NewEngine(strategies)

		prod := product.Product{
			Code:  "000003",
			Price: decimal.NewFromFloat(100.0),
			Category: &product.Category{
				Code: "boots",
			},
		}

		discountedPrice := engine.ApplyDiscount(prod)
		percentage := engine.GetDiscountPercentage(prod)

		assert.Equal(t, "85", discountedPrice.String())
		assert.Equal(t, 15, percentage)
	})

	t.Run("successfully calculates discount for decimal prices", func(t *testing.T) {
		strategies := []Strategy{
			NewCategoryDiscountStrategy("boots", 30),
		}
		engine := NewEngine(strategies)

		prod := product.Product{
			Code:  "PROD001",
			Price: decimal.NewFromFloat(89.99),
			Category: &product.Category{
				Code: "boots",
			},
		}

		discountedPrice := engine.ApplyDiscount(prod)

		// 89.99 - 30% = 62.993 ~= 62.99
		assert.True(t, discountedPrice.LessThan(decimal.NewFromFloat(63.0)))
		assert.True(t, discountedPrice.GreaterThan(decimal.NewFromFloat(62.99)))
	})

	t.Run("works with empty strategies list", func(t *testing.T) {
		engine := NewEngine([]Strategy{})

		prod := product.Product{
			Code:  "PROD001",
			Price: decimal.NewFromFloat(100.0),
		}

		discountedPrice := engine.ApplyDiscount(prod)
		percentage := engine.GetDiscountPercentage(prod)

		assert.Equal(t, "100", discountedPrice.String())
		assert.Equal(t, 0, percentage)
	})
}
