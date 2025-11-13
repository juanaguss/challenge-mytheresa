package mapper

import (
	"testing"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestToProductResponse(t *testing.T) {
	t.Run("maps product without discount", func(t *testing.T) {
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

		response := ToProductResponse(p, 89.99, 0)

		assert.Equal(t, "PROD001", response.Code)
		assert.Equal(t, 89.99, response.Price)
		assert.Equal(t, "clothing", response.Category)
		assert.Nil(t, response.Discount)
		assert.Nil(t, response.FinalPrice)
	})

	t.Run("maps product with discount", func(t *testing.T) {
		p := product.Product{
			ID:    1,
			Code:  "PROD001",
			Price: decimal.NewFromFloat(100.0),
			Category: &product.Category{
				Code: "boots",
			},
		}

		response := ToProductResponse(p, 70.0, 30)

		assert.Equal(t, "PROD001", response.Code)
		assert.Equal(t, 100.0, response.Price)
		assert.Equal(t, "boots", response.Category)
		assert.NotNil(t, response.Discount)
		assert.Equal(t, "30%", *response.Discount)
		assert.NotNil(t, response.FinalPrice)
		assert.Equal(t, 70.0, *response.FinalPrice)
	})

	t.Run("maps product without category", func(t *testing.T) {
		p := product.Product{
			ID:       2,
			Code:     "PROD002",
			Price:    decimal.NewFromFloat(45.50),
			Category: nil,
		}

		response := ToProductResponse(p, 45.50, 0)

		assert.Equal(t, "PROD002", response.Code)
		assert.Equal(t, 45.50, response.Price)
		assert.Equal(t, "", response.Category)
		assert.Nil(t, response.Discount)
		assert.Nil(t, response.FinalPrice)
	})
}

func TestToProductResponses(t *testing.T) {
	t.Run("maps various products without discounts", func(t *testing.T) {
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
		discountedPrices := []float64{89.99, 129.99}
		percentages := []int{0, 0}

		responses := ToProductResponses(products, discountedPrices, percentages)

		assert.Len(t, responses, 2)

		assert.Equal(t, "PROD001", responses[0].Code)
		assert.Equal(t, 89.99, responses[0].Price)
		assert.Equal(t, "clothing", responses[0].Category)
		assert.Nil(t, responses[0].Discount)
		assert.Nil(t, responses[0].FinalPrice)

		assert.Equal(t, "PROD002", responses[1].Code)
		assert.Equal(t, 129.99, responses[1].Price)
		assert.Equal(t, "shoes", responses[1].Category)
		assert.Nil(t, responses[1].Discount)
		assert.Nil(t, responses[1].FinalPrice)
	})

	t.Run("maps multiple products with different discounts", func(t *testing.T) {
		products := []product.Product{
			{
				ID:    1,
				Code:  "PROD001",
				Price: decimal.NewFromFloat(100.00),
				Category: &product.Category{
					Code: "boots",
				},
			},
			{
				ID:    2,
				Code:  "PROD002",
				Price: decimal.NewFromFloat(50.00),
				Category: &product.Category{
					Code: "shoes",
				},
			},
		}
		discountedPrices := []float64{70.00, 50.00}
		percentages := []int{30, 0}

		responses := ToProductResponses(products, discountedPrices, percentages)

		assert.Len(t, responses, 2)

		// First product
		assert.Equal(t, "PROD001", responses[0].Code)
		assert.Equal(t, 100.00, responses[0].Price)
		assert.NotNil(t, responses[0].Discount)
		assert.Equal(t, "30%", *responses[0].Discount)
		assert.NotNil(t, responses[0].FinalPrice)
		assert.Equal(t, 70.00, *responses[0].FinalPrice)

		// Second product
		assert.Equal(t, "PROD002", responses[1].Code)
		assert.Equal(t, 50.00, responses[1].Price)
		assert.Nil(t, responses[1].Discount)
		assert.Nil(t, responses[1].FinalPrice)
	})

	t.Run("returns empty slice for empty input", func(t *testing.T) {
		products := []product.Product{}
		discountedPrices := []float64{}
		percentages := []int{}

		responses := ToProductResponses(products, discountedPrices, percentages)

		assert.Empty(t, responses)
		assert.NotNil(t, responses)
	})

	t.Run("handles nil category in multiple products", func(t *testing.T) {
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
		discountedPrices := []float64{50.00, 75.00}
		percentages := []int{0, 0}

		responses := ToProductResponses(products, discountedPrices, percentages)

		assert.Len(t, responses, 2)
		assert.Equal(t, "", responses[0].Category)
		assert.Equal(t, "accessories", responses[1].Category)
	})
}
