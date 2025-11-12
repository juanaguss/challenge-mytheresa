package http

import (
	"errors"
	"net/http"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCatalogHandler(t *testing.T) {
	t.Run("creates handler with service", func(t *testing.T) {
		service := newMockService(nil, nil)
		handler := NewCatalogHandler(service)

		require.NotNil(t, handler)
	})
}

func TestHandleGet_Success(t *testing.T) {
	t.Run("returns products with categories", func(t *testing.T) {
		products := createTestProducts(2)
		service := newMockService(products, nil)
		handler := NewCatalogHandler(service)

		w := makeRequest(handler, "/catalog")

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		response := parseResponse(t, w)
		assert.Len(t, response.Products, 2)
		assert.Equal(t, 2, response.Total)
		assert.Equal(t, "PROD001", response.Products[0].Code)
		assert.Equal(t, "clothing", response.Products[0].Category)
		assert.Equal(t, 65.76, response.Products[0].Price)
	})

	t.Run("returns empty array when no products exist", func(t *testing.T) {
		service := newMockService([]product.Product{}, nil)
		handler := NewCatalogHandler(service)

		w := makeRequest(handler, "/catalog")

		assert.Equal(t, http.StatusOK, w.Code)

		response := parseResponse(t, w)
		assert.Empty(t, response.Products)
		assert.Equal(t, 0, response.Total)
	})

	t.Run("returns empty category when product has none", func(t *testing.T) {
		products := []product.Product{
			newTestProduct(99, "PROD099", 99.01, nil),
		}
		service := newMockService(products, nil)
		handler := NewCatalogHandler(service)

		w := makeRequest(handler, "/catalog")

		assert.Equal(t, http.StatusOK, w.Code)

		response := parseResponse(t, w)
		assert.Len(t, response.Products, 1)
		assert.Equal(t, "", response.Products[0].Category)
	})
}

func TestHandleGet_Error(t *testing.T) {
	t.Run("returns 500 when service fails", func(t *testing.T) {
		service := newMockService(nil, errors.New("database connection failed"))
		handler := NewCatalogHandler(service)

		w := makeRequest(handler, "/catalog")

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "database connection failed")
	})
}
