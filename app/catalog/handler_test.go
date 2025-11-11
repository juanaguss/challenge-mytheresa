package catalog

import (
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

// tests para hacer de handleGet
// happy path: vuelven productos - vuelve un array vacio
// error path: vuelve 500
// crear handler con el repo

func testHandleGet_Success(t *testing.T) {
	t.Run("Successfully returns products", func(t *testing.T) {
		//arrange
		mockRepo := &mockProductsRepository{
			products: []models.Product{
				{
					ID:    1,
					Code:  "PROD001",
					Price: decimal.NewFromFloat(270.75),
				},
				{
					ID:    2,
					Code:  "PROD002",
					Price: decimal.NewFromFloat(400.99),
				},
			},
			err: nil,
		}

		handler := NewCatalogHandler(mockRepo)

		//act
		req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
		w := httptest.NewRecorder()

		handler.HandleGet(w, req)

		//assert
		assert.Equal(t, http.StatusOK, w.Code)

		expectedJSON := `{"products":[{"code":"PROD001","price":270.75},{"code":"PROD002","price":400.99}]}`
		assert.JSONEq(t, expectedJSON, w.Body.String())
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
		assert.Equal(t, http.StatusOK, w.Code)

		expectedJSON := `{"products":[]}`
		assert.JSONEq(t, expectedJSON, w.Body.String())

	})
}

func testHandleGet_Error(t *testing.T) {
	t.Run("Returns error 500 when products repository falis", func(t *testing.T) {
		// Arrange
		mockRepo := &mockProductsRepository{}

		// Act
		handler := NewCatalogHandler(mockRepo)

		//Assert
		assert.NotNil(t, handler)
		assert.Equal(t, mockRepo, handler.repo)
	})
}
