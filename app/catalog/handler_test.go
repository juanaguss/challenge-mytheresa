package catalog

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type mockProductsRepository struct {
	products   []models.Product
	totalCount int64
	err        error
}

func (m *mockProductsRepository) GetAllProducts() ([]models.Product, error) {
	return m.products, m.err
}

func (m *mockProductsRepository) GetProducts(offset, limit int) ([]models.Product, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}

	total := m.totalCount
	if total == 0 && len(m.products) > 0 {
		total = int64(len(m.products))
	}

	if offset >= len(m.products) {
		return []models.Product{}, total, nil
	}

	end := offset + limit
	if end > len(m.products) {
		end = len(m.products)
	}

	return m.products[offset:end], total, nil
}

func createTestProducts(count int) []models.Product {
	products := make([]models.Product, count)
	categories := []models.Category{
		{ID: 1, Code: "clothing", Name: "Clothing"},
		{ID: 2, Code: "shoes", Name: "Shoes"},
		{ID: 3, Code: "accessories", Name: "Accessories"},
	}

	for i := 0; i < count; i++ {
		categoryIndex := i % len(categories)
		products[i] = models.Product{
			ID:       uint(i + 1),
			Code:     fmt.Sprintf("PROD%03d", i+1),
			Price:    decimal.NewFromInt(int64((i + 1) * 100)),
			Category: &categories[categoryIndex],
		}
	}
	return products
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

		expectedJSON := `{"products":[{"code":"PROD001","price":270.75,"category":"clothing"},{"code":"PROD002","price":400.99,"category":"shoes"}],"total":2}`
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

		expectedJSON := `{"products":[],"total":0}`
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
		expectedJSON := `{"products":[{"code":"PROD001","price":99.01,"category":""}],"total":1}`
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

func TestHandlerGet_Pagination(t *testing.T) {
	// Arrange for every test
	const totalProducts = 50
	testProducts := createTestProducts(totalProducts)

	t.Run("Returns default pagination when no query parameters", func(t *testing.T) {
		// Arrange
		mockRepo := &mockProductsRepository{
			products:   testProducts,
			totalCount: totalProducts,
			err:        nil,
		}

		handler := NewCatalogHandler(mockRepo)

		req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
		w := httptest.NewRecorder()

		// Act
		handler.HandleGet(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Len(t, response.Products, DefaultLimit, "Should return default limit of products")
		assert.Equal(t, totalProducts, response.Total, "Total should match total products")

		expectedFirstCode := fmt.Sprintf("PROD%03d", DefaultOffset+1)
		assert.Equal(t, expectedFirstCode, response.Products[0].Code, fmt.Sprintf("First product should be %s", expectedFirstCode))
	})

	t.Run("Returns correct page with custom offset and limit", func(t *testing.T) {
		// Arrange
		offset := 10
		limit := 5

		mockRepo := &mockProductsRepository{
			products:   testProducts,
			totalCount: totalProducts,
			err:        nil,
		}

		handler := NewCatalogHandler(mockRepo)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/catalog?offset=%d&limit=%d", offset, limit), nil)
		w := httptest.NewRecorder()

		// Act
		handler.HandleGet(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, limit, len(response.Products), "Should return the requested limit")
		assert.Equal(t, totalProducts, response.Total, "Total should match total products")

		expectedFirstCode := fmt.Sprintf("PROD%03d", offset+1)
		assert.Equal(t, expectedFirstCode, response.Products[0].Code, "First product should match offset position")

		expectedLastCode := fmt.Sprintf("PROD%03d", offset+limit)
		assert.Equal(t, expectedLastCode, response.Products[len(response.Products)-1].Code, "Last product should match offset+limit position")
	})

	t.Run("Returns remaining products in case offset+limit exceeds total", func(t *testing.T) {
		// Arrange
		offset := totalProducts - 3 // last 3 products
		limit := 10                 // ask for 10 but there are 3

		mockRepo := &mockProductsRepository{
			products:   testProducts,
			totalCount: totalProducts,
			err:        nil,
		}

		handler := NewCatalogHandler(mockRepo)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/catalog?offset=%d&limit=%d", offset, limit), nil)
		w := httptest.NewRecorder()

		// Act
		handler.HandleGet(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, 3, len(response.Products),
			"Should return only remaining products")
		assert.Equal(t, totalProducts, response.Total)
	})

	t.Run("Returns empty array if offset exceeds total", func(t *testing.T) {
		// Arrange
		offset := totalProducts + 10
		limit := 10

		mockRepo := &mockProductsRepository{
			products:   testProducts,
			totalCount: totalProducts,
			err:        nil,
		}

		handler := NewCatalogHandler(mockRepo)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/catalog?offset=%d&limit=%d", offset, limit), nil)
		w := httptest.NewRecorder()

		// Act
		handler.HandleGet(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Empty(t, response.Products, "Should return empty array")
		assert.Equal(t, totalProducts, response.Total)
	})

	t.Run("Returns 400 for invalid parameters pagination", func(t *testing.T) {
		testCases := []struct {
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
				queryParams:       fmt.Sprintf("limit=%d", MinLimit-1),
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
			{
				name:              "limit zero",
				queryParams:       "limit=0",
				expectedErrorText: "limit",
			},
			{
				name:              "both parameters invalid",
				queryParams:       "offset=invalid&limit=invalid",
				expectedErrorText: "offset", // Should fail on first validation
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				mockRepo := &mockProductsRepository{
					products:   testProducts,
					totalCount: totalProducts,
				}
				handler := NewCatalogHandler(mockRepo)

				req := httptest.NewRequest(http.MethodGet,
					fmt.Sprintf("/catalog?%s", tc.queryParams), nil)
				w := httptest.NewRecorder()

				handler.HandleGet(w, req)

				assert.Equal(t, http.StatusBadRequest, w.Code, "Should return 400 (Bad Request)")
				assert.Contains(t, w.Body.String(), tc.expectedErrorText, "Error message should have the invalid parameter")
			})
		}
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
