package catalog

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleGet_Pagination(t *testing.T) {
	const totalProducts = 50
	testProducts := createTestProducts(totalProducts)

	t.Run("returns default pagination when no query parameters", func(t *testing.T) {
		mockRepo := newMockRepo(testProducts, nil)
		handler := NewCatalogHandler(mockRepo)

		w := makeRequest(handler, http.MethodGet, "/catalog")

		assert.Equal(t, http.StatusOK, w.Code)

		response := parseResponse(t, w)
		assert.Len(t, response.Products, DefaultLimit)
		assert.Equal(t, totalProducts, response.Total)
		assert.Equal(t, "PROD001", response.Products[0].Code)
	})

	t.Run("returns correct page with custom offset and limit", func(t *testing.T) {
		offset := 10
		limit := 5

		mockRepo := newMockRepo(testProducts, nil)
		handler := NewCatalogHandler(mockRepo)

		w := makeRequest(handler, http.MethodGet,
			fmt.Sprintf("/catalog?offset=%d&limit=%d", offset, limit))

		assert.Equal(t, http.StatusOK, w.Code)

		response := parseResponse(t, w)
		assert.Len(t, response.Products, limit)
		assert.Equal(t, totalProducts, response.Total)
		assert.Equal(t, "PROD011", response.Products[0].Code)
		assert.Equal(t, "PROD015", response.Products[len(response.Products)-1].Code)
	})

	t.Run("returns remaining products when offset + limit exceeds total", func(t *testing.T) {
		offset := totalProducts - 3
		limit := 10

		mockRepo := newMockRepo(testProducts, nil)
		handler := NewCatalogHandler(mockRepo)

		w := makeRequest(handler, http.MethodGet,
			fmt.Sprintf("/catalog?offset=%d&limit=%d", offset, limit))

		assert.Equal(t, http.StatusOK, w.Code)

		response := parseResponse(t, w)
		assert.Len(t, response.Products, 3)
		assert.Equal(t, totalProducts, response.Total)
	})

	t.Run("returns empty array if offset exceeds total", func(t *testing.T) {
		offset := totalProducts + 10
		limit := 10

		mockRepo := newMockRepo(testProducts, nil)
		handler := NewCatalogHandler(mockRepo)

		w := makeRequest(handler, http.MethodGet,
			fmt.Sprintf("/catalog?offset=%d&limit=%d", offset, limit))

		assert.Equal(t, http.StatusOK, w.Code)

		response := parseResponse(t, w)
		assert.Empty(t, response.Products)
		assert.Equal(t, totalProducts, response.Total)
	})
}

func TestHandleGet_PaginationValidation(t *testing.T) {
	testProducts := createTestProducts(20)

	tests := []struct {
		name              string
		queryParams       string
		expectedErrorText string
	}{
		{
			name:              "limit exceeds maximum",
			queryParams:       fmt.Sprintf("limit=%d", MaxLimit+1),
			expectedErrorText: "limit",
		},
		{
			name:              "limit below minimum",
			queryParams:       "limit=0",
			expectedErrorText: "limit",
		},
		{
			name:              "negative offset",
			queryParams:       "offset=-1",
			expectedErrorText: "offset",
		},
		{
			name:              "invalid limit format",
			queryParams:       "limit=invalid",
			expectedErrorText: "limit",
		},
		{
			name:              "invalid offset format",
			queryParams:       "offset=invalid",
			expectedErrorText: "offset",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockRepo(testProducts, nil)
			handler := NewCatalogHandler(mockRepo)

			w := makeRequest(handler, http.MethodGet, "/catalog?"+tt.queryParams)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedErrorText)
		})
	}
}
