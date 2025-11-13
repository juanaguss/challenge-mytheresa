package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/internal/application/catalog"
	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type mockDetailService struct {
	mockService
	product          *product.Product
	discountedPrice  float64
	percentage       int
	variantDiscounts map[string]catalog.VariantDiscount
	err              error
}

func (m *mockDetailService) GetProductByCode(code string) (*product.Product, float64, int, map[string]catalog.VariantDiscount, error) {
	if m.err != nil {
		return nil, 0, 0, nil, m.err
	}
	return m.product, m.discountedPrice, m.percentage, m.variantDiscounts, nil
}

func TestHandleGetByCode_Success(t *testing.T) {
	t.Run("returns product with variants", func(t *testing.T) {
		p := product.Product{
			ID:    1,
			Code:  "PROD001",
			Price: decimal.NewFromFloat(100.00),
			Category: &product.Category{
				Code: "shoes",
			},
			Variants: []product.Variant{
				{
					ID:    1,
					SKU:   "VAR001",
					Name:  "Size 45",
					Price: decimal.NewFromFloat(100.00),
				},
				{
					ID:    2,
					SKU:   "VAR002",
					Name:  "Size 46",
					Price: decimal.NewFromFloat(105.00),
				},
			},
		}

		service := &mockDetailService{product: &p, discountedPrice: 100.00, percentage: 0}
		handler := NewCatalogHandler(service)

		req := httptest.NewRequest("GET", "/catalog/PROD001", nil)
		req.SetPathValue("code", "PROD001")
		w := httptest.NewRecorder()

		handler.HandleGetByCode(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "PROD001")
		assert.Contains(t, w.Body.String(), "VAR001")
		assert.Contains(t, w.Body.String(), "VAR002")
		assert.Contains(t, w.Body.String(), "shoes")
	})

	t.Run("returns product with discount applied", func(t *testing.T) {
		p := product.Product{
			ID:    2,
			Code:  "PROD002",
			Price: decimal.NewFromFloat(100.00),
			Category: &product.Category{
				Code: "boots",
			},
			Variants: []product.Variant{},
		}

		service := &mockDetailService{product: &p, discountedPrice: 70.00, percentage: 30}
		handler := NewCatalogHandler(service)

		req := httptest.NewRequest("GET", "/catalog/PROD002", nil)
		req.SetPathValue("code", "PROD002")
		w := httptest.NewRecorder()

		handler.HandleGetByCode(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"discount":"30%"`)
		assert.Contains(t, w.Body.String(), `"final_price":70`)
	})

	t.Run("returns product without variants", func(t *testing.T) {
		p := product.Product{
			ID:       3,
			Code:     "PROD003",
			Price:    decimal.NewFromFloat(50.00),
			Category: nil,
			Variants: []product.Variant{},
		}

		service := &mockDetailService{product: &p, discountedPrice: 50.00, percentage: 0}
		handler := NewCatalogHandler(service)

		req := httptest.NewRequest("GET", "/catalog/PROD003", nil)
		req.SetPathValue("code", "PROD003")
		w := httptest.NewRecorder()

		handler.HandleGetByCode(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "PROD003")
		assert.Contains(t, w.Body.String(), `"variants":[]`)
	})
}

func TestHandleGetByCode_Error(t *testing.T) {
	t.Run("returns 400 when code is missing", func(t *testing.T) {
		service := &mockDetailService{}
		handler := NewCatalogHandler(service)

		req := httptest.NewRequest("GET", "/catalog/", nil)
		w := httptest.NewRecorder()

		handler.HandleGetByCode(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "product code is required")
	})

	t.Run("returns 404 when product not found", func(t *testing.T) {
		service := &mockDetailService{err: errors.New("not found")}
		handler := NewCatalogHandler(service)

		req := httptest.NewRequest("GET", "/catalog/NONEXISTENT", nil)
		req.SetPathValue("code", "NONEXISTENT")
		w := httptest.NewRecorder()

		handler.HandleGetByCode(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "not found")
	})
}
