package http

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

type mockService struct {
	products []product.Product
	total    int64
	err      error
}

func (m *mockService) GetProducts(offset, limit int, filters product.Filter) ([]product.Product, []float64, []int, int64, error) {
	if m.err != nil {
		return nil, nil, nil, 0, m.err
	}

	filtered := make([]product.Product, 0)
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
		return []product.Product{}, []float64{}, []int{}, total, nil
	}

	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	result := filtered[offset:end]

	discountedPrices := make([]float64, len(result))
	percentages := make([]int, len(result))
	for i, p := range result {
		price, _ := p.Price.Float64()
		discountedPrices[i] = price
		percentages[i] = 0
	}

	return result, discountedPrices, percentages, total, nil
}

var (
	categoryClothing = &product.Category{
		ID:   1,
		Code: "clothing",
		Name: "Clothing",
	}
	categoryShoes = &product.Category{
		ID:   2,
		Code: "shoes",
		Name: "Shoes",
	}
	categoryAccessories = &product.Category{
		ID:   3,
		Code: "accessories",
		Name: "Accessories",
	}
)

func newTestProduct(id uint, code string, price float64, category *product.Category) product.Product {
	return product.Product{
		ID:       id,
		Code:     code,
		Price:    decimal.NewFromFloat(price),
		Category: category,
	}
}

func newMockService(products []product.Product, err error) *mockService {
	return &mockService{
		products: products,
		total:    int64(len(products)),
		err:      err,
	}
}

func createTestProducts(count int) []product.Product {
	baseProducts := []product.Product{
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

	products := make([]product.Product, count)
	copy(products, baseProducts)

	categories := []*product.Category{categoryClothing, categoryShoes, categoryAccessories}

	for i := 8; i < count; i++ {
		categoryIndex := i % len(categories)
		products[i] = product.Product{
			ID:       uint(i + 1),
			Code:     fmt.Sprintf("PROD%03d", i+1),
			Price:    decimal.NewFromFloat(float64(i+1) * 10.0),
			Category: categories[categoryIndex],
		}
	}

	return products
}

func makeRequest(handler *CatalogHandler, url string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	handler.HandleGet(w, req)
	return w
}

func parseResponse(t *testing.T, w *httptest.ResponseRecorder) catalogResponse {
	var response catalogResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to parse response JSON")
	return response
}
