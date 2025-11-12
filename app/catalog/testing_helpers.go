package catalog

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

// mockProductsRepository simulates testing repository.
type mockProductsRepository struct {
	products   []models.Product
	totalCount int64
	err        error
}

func (m *mockProductsRepository) GetAllProducts() ([]models.Product, error) {
	return m.products, m.err
}

func (m *mockProductsRepository) GetProducts(offset, limit int, filters models.ProductFilters) ([]models.Product, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}

	filtered := make([]models.Product, 0)
	for _, p := range m.products {
		if filters.Category != "" {
			if p.Category == nil || p.Category.Code != filters.Category {
				continue
			}
		}

		if filters.PriceLessThan != nil {
			if p.Price.GreaterThanOrEqual(*filters.PriceLessThan) {
				continue
			}
		}

		filtered = append(filtered, p)
	}

	total := int64(len(filtered))

	if offset >= len(filtered) {
		return []models.Product{}, total, nil
	}

	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[offset:end], total, nil
}

// Predefined categories to use in different tests scenarios.
var (
	categoryClothing = &models.Category{
		ID:   1,
		Code: "clothing",
		Name: "Clothing",
	}
	categoryShoes = &models.Category{
		ID:   2,
		Code: "shoes",
		Name: "Shoes",
	}
	categoryAccessories = &models.Category{
		ID:   3,
		Code: "accessories",
		Name: "Accessories",
	}
)

// newTestProduct creates a test product
func newTestProduct(id uint, code string, price float64, category *models.Category) models.Product {
	return models.Product{
		ID:       id,
		Code:     code,
		Price:    decimal.NewFromFloat(price),
		Category: category,
	}
}

// newMockRepo craetes a repo
func newMockRepo(products []models.Product, err error) *mockProductsRepository {
	return &mockProductsRepository{
		products:   products,
		totalCount: int64(len(products)),
		err:        err,
	}
}

// createTestProducts creates N products for testing.
func createTestProducts(count int) []models.Product {
	baseProducts := []models.Product{
		newTestProduct(1, "PROD001", 65.76, categoryClothing),
		newTestProduct(2, "PROD002", 12.54, categoryShoes),
		newTestProduct(3, "PROD003", 23788.75, categoryAccessories),
		newTestProduct(4, "PROD004", 56.00, categoryClothing),
		newTestProduct(5, "PROD005", 9.99, categoryAccessories),
		newTestProduct(6, "PROD006", 1.50, categoryShoes),
		newTestProduct(7, "PROD007", 667.70, categoryClothing),
		newTestProduct(8, "PROD008", 88.88, categoryAccessories),
	}

	if count <= 8 {
		return baseProducts[:count]
	}

	products := make([]models.Product, count)
	copy(products, baseProducts)

	categories := []*models.Category{categoryClothing, categoryShoes, categoryAccessories}

	for i := 8; i < count; i++ {
		categoryIndex := i % len(categories)
		products[i] = models.Product{
			ID:       uint(i + 1),
			Code:     fmt.Sprintf("PROD%03d", i+1),
			Price:    decimal.NewFromFloat(float64(i+1) * 10.0),
			Category: categories[categoryIndex],
		}
	}

	return products
}

// makeRequest excecutes a requests and returns a ResponseRecorder
func makeRequest(handler *Handler, method, url string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, url, nil)
	w := httptest.NewRecorder()
	handler.HandleGet(w, req)
	return w
}
// parseResponse converts the responseBody to a Response struct.
func parseResponse(t *testing.T, w *httptest.ResponseRecorder) Response {
	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to parse response JSON")
	return response
}