package mapper

import (
	"testing"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestToProductResponse(t *testing.T) {
	t.Run("maps product with category", func(t *testing.T) {
		p := product.Product{
			ID:    1,
			Code:  "PROD001",
			Price: decimal.NewFromFloat(89.99),
			Category: &product.Category{
				ID:   1,
				Code: "clothing",
				Name: "Clothing",
			},
		}

		response := ToProductResponse(p)

		assert.Equal(t, "PROD001", response.Code)
		assert.Equal(t, 89.99, response.Price)
		assert.Equal(t, "clothing", response.Category)
	})

	t.Run("maps product without category", func(t *testing.T) {
		p := product.Product{
			ID:       2,
			Code:     "PROD002",
			Price:    decimal.NewFromFloat(45.50),
			Category: nil,
		}

		response := ToProductResponse(p)

		assert.Equal(t, "PROD002", response.Code)
		assert.Equal(t, 45.50, response.Price)
		assert.Equal(t, "", response.Category)
	})

	t.Run("handles zero price correctly", func(t *testing.T) {
		p := product.Product{
			ID:    3,
			Code:  "PROD003",
			Price: decimal.Zero,
		}

		response := ToProductResponse(p)

		assert.Equal(t, "PROD003", response.Code)
		assert.Equal(t, 0.0, response.Price)
	})
}

func TestToProductResponses(t *testing.T) {
	t.Run("maps multiple products", func(t *testing.T) {
		products := []product.Product{
			{
				ID:    1,
				Code:  "PROD001",
				Price: decimal.NewFromFloat(89.99),
				Category: &product.Category{
					Code: "clothing",
				},
			},
			{
				ID:    2,
				Code:  "PROD002",
				Price: decimal.NewFromFloat(129.99),
				Category: &product.Category{
					Code: "shoes",
				},
			},
		}

		responses := ToProductResponses(products)

		assert.Len(t, responses, 2)

		assert.Equal(t, "PROD001", responses[0].Code)
		assert.Equal(t, 89.99, responses[0].Price)
		assert.Equal(t, "clothing", responses[0].Category)

		assert.Equal(t, "PROD002", responses[1].Code)
		assert.Equal(t, 129.99, responses[1].Price)
		assert.Equal(t, "shoes", responses[1].Category)
	})

	t.Run("returns empty slice for empty input", func(t *testing.T) {
		products := []product.Product{}

		responses := ToProductResponses(products)

		assert.Empty(t, responses)
		assert.NotNil(t, responses)
	})

	t.Run("handles null category in multiple products", func(t *testing.T) {
		products := []product.Product{
			{
				ID:       1,
				Code:     "PROD001",
				Price:    decimal.NewFromFloat(50.00),
				Category: nil,
			},
			{
				ID:    2,
				Code:  "PROD002",
				Price: decimal.NewFromFloat(75.00),
				Category: &product.Category{
					Code: "accessories",
				},
			},
		}

		responses := ToProductResponses(products)

		assert.Len(t, responses, 2)
		assert.Equal(t, "", responses[0].Category)
		assert.Equal(t, "accessories", responses[1].Category)
	})
}
