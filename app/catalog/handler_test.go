package catalog

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type mockProductsRepository struct {
	products []models.Product
	err      error
}

func (m *mockProductsRepository) GetAllProducts() ([]models.Product, error) {
	return m.products, m.err
}

func TestHandleGet_Success(t *testing.T) {
	t.Run("Successfully returns products with their category", func(t *testing.T) {
		// Arrange
		mockRepo := &mockProductsRepository{
			products: []models.Product{
				{
					ID:    1,
					Code:  "PROD001",
					Price: decimal.NewFromFloat(270.75),
					Category: &models.Category{
						ID:   1,
						Code: "clothing",
						Name: "Clothing",
					},
				},
				{
					ID:    2,
					Code:  "PROD002",
					Price: decimal.NewFromFloat(400.99),
					Category: &models.Category{
						ID:   2,
						Code: "shoes",
						Name: "Shoes",
					},
				},
			},
			err: nil,
		}

		handler := NewCatalogHandler(mockRepo)

		// Act
		req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
		w := httptest.NewRecorder()

		handler.HandleGet(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code, "Response code does not match expected (Status 200 OK)")
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"), "Expected Content-Type application/json")

		expectedJSON := `{"products":[{"code":"PROD001","price":270.75,"category":"clothing"},{"code":"PROD002","price":400.99,"category":"shoes"}]}`
		assert.JSONEq(t, expectedJSON, w.Body.String(), "Response body does not match expected JSON")
	})

	t.Run("Returns empty array when no products exist", func(t *testing.T) {
		// Arrange
		mockRepo := &mockProductsRepository{
			products: []models.Product{},
			err:      nil,
		}

		handler := NewCatalogHandler(mockRepo)

		// Act
		req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
		w := httptest.NewRecorder()

		handler.HandleGet(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code, "Response code does not match expected (Status 200 OK)")

		expectedJSON := `{"products":[]}`
		assert.JSONEq(t, expectedJSON, w.Body.String(), "Response body does not match expected empty JSON")

	})

	t.Run("Returns empty category when product has none", func(t *testing.T) {
		// Arrange
		mockRepo := &mockProductsRepository{
			products: []models.Product{
				{
					ID:       1,
					Code:     "PROD001",
					Price:    decimal.NewFromFloat(99.01),
					Category: nil,
				},
			},
			err: nil,
		}

		handler := NewCatalogHandler(mockRepo)

		// Act
		req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
		w := httptest.NewRecorder()

		handler.HandleGet(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code, "Expected status 200 OK")
		expectedJSON := `{"products":[{"code":"PROD001","price":99.01,"category":""}]}`
		assert.JSONEq(t, expectedJSON, w.Body.String(), "Expected empty string for category when null")
	})
}

func TestHandleGet_Error(t *testing.T) {
	t.Run("Return error 500 when repository fails", func(t *testing.T) {
		// Arrange
		mockRepo := &mockProductsRepository{
			products: nil,
			err:      errors.New("database connection failed"),
		}

		handler := NewCatalogHandler(mockRepo)

		// Act
		req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
		w := httptest.NewRecorder()

		handler.HandleGet(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected status 500 error")
		assert.Contains(t, w.Body.String(), "database connection failed", "Error message should be in response")
	})
}

func TestNewCatalogHandler(t *testing.T) {
	t.Run("Creates handler with repo", func(t *testing.T) {
		// Arrange
		mockRepo := &mockProductsRepository{}

		// Act
		handler := NewCatalogHandler(mockRepo)

		// Assert
		assert.NotNil(t, handler, "Handler should not be nil")
		assert.Equal(t, mockRepo, handler.repo, "Handler should have a repo")
	})
}
