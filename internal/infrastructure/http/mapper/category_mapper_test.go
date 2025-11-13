package mapper

import (
	"testing"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
	"github.com/stretchr/testify/assert"
)

func TestToCategoryResponse(t *testing.T) {
	t.Run("converts domain category to response DTO", func(t *testing.T) {
		category := product.Category{
			ID:   1,
			Code: "boots",
			Name: "Boots",
		}

		response := ToCategoryResponse(category)

		assert.Equal(t, "boots", response.Code)
		assert.Equal(t, "Boots", response.Name)
	})
}

func TestToCategoryResponses(t *testing.T) {
	t.Run("converts multiple categories", func(t *testing.T) {
		categories := []product.Category{
			{ID: 1, Code: "boots", Name: "Boots"},
			{ID: 2, Code: "sandals", Name: "Sandals"},
			{ID: 3, Code: "sneakers", Name: "Sneakers"},
		}

		responses := ToCategoryResponses(categories)

		assert.Len(t, responses, 3)
		assert.Equal(t, "boots", responses[0].Code)
		assert.Equal(t, "Boots", responses[0].Name)
		assert.Equal(t, "sandals", responses[1].Code)
		assert.Equal(t, "Sandals", responses[1].Name)
		assert.Equal(t, "sneakers", responses[2].Code)
		assert.Equal(t, "Sneakers", responses[2].Name)
	})

	t.Run("returns empty slice for empty input", func(t *testing.T) {
		categories := []product.Category{}

		responses := ToCategoryResponses(categories)

		assert.Empty(t, responses)
		assert.NotNil(t, responses)
	})
}
