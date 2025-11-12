package catalog

import (
	"net/http"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/assert"
)

// setupFilterTestProducts usa specific values to test filters
func setupFilterTestProducts() []models.Product {
	return []models.Product{
		newTestProduct(1, "PROD001", 65.76, categoryClothing),
		newTestProduct(2, "PROD002", 12.54, categoryShoes),
		newTestProduct(3, "PROD003", 23788.75, categoryAccessories),
		newTestProduct(4, "PROD004", 56.00, categoryClothing),
		newTestProduct(5, "PROD005", 9.99, categoryAccessories),
		newTestProduct(6, "PROD006", 1.50, categoryShoes),
		newTestProduct(99, "PROD099", 333.00, nil),
	}
}

func TestHandleGet_Filters(t *testing.T) {
	testProducts := setupFilterTestProducts()

	tests := []struct {
		name           string
		queryParams    string
		expectedCount  int
		expectedTotal  int
		validateResult func(t *testing.T, products []Product)
	}{
		{
			name:          "filters by category clothing",
			queryParams:   "category=clothing",
			expectedCount: 2,
			expectedTotal: 2,
			validateResult: func(t *testing.T, products []Product) {
				assert.Equal(t, "PROD001", products[0].Code)
				assert.Equal(t, "PROD004", products[1].Code)
				for _, p := range products {
					assert.Equal(t, "clothing", p.Category)
				}
			},
		},
		{
			name:          "filters by category shoes",
			queryParams:   "category=shoes",
			expectedCount: 2,
			expectedTotal: 2,
			validateResult: func(t *testing.T, products []Product) {
				assert.Equal(t, "PROD002", products[0].Code)
				assert.Equal(t, "PROD006", products[1].Code)
				for _, p := range products {
					assert.Equal(t, "shoes", p.Category)
				}
			},
		},
		{
			name:          "filters by category accessories",
			queryParams:   "category=accessories",
			expectedCount: 2,
			expectedTotal: 2,
			validateResult: func(t *testing.T, products []Product) {
				assert.Equal(t, "PROD003", products[0].Code)
				assert.Equal(t, "PROD005", products[1].Code)
				for _, p := range products {
					assert.Equal(t, "accessories", p.Category)
				}
			},
		},
		{
			name:          "filters by priceLessThan 10 - only cheapest products",
			queryParams:   "priceLessThan=10",
			expectedCount: 2,
			expectedTotal: 2,
			validateResult: func(t *testing.T, products []Product) {
				// Debe incluir PROD005 (9.99) y PROD006 (1.50)
				codes := []string{products[0].Code, products[1].Code}
				assert.Contains(t, codes, "PROD005")
				assert.Contains(t, codes, "PROD006")
				for _, p := range products {
					assert.Less(t, p.Price, 10.0)
				}
			},
		},
		{
			name:          "filters by priceLessThan 50 - low to medium prices",
			queryParams:   "priceLessThan=50",
			expectedCount: 3,
			expectedTotal: 3,
			validateResult: func(t *testing.T, products []Product) {
				// Debe incluir PROD002 (12.54), PROD005 (9.99), PROD006 (1.50)
				for _, p := range products {
					assert.Less(t, p.Price, 50.0)
				}
			},
		},
		{
			name:          "filters by priceLessThan 100 - excludes expensive products",
			queryParams:   "priceLessThan=100",
			expectedCount: 5,
			expectedTotal: 5,
			validateResult: func(t *testing.T, products []Product) {
				// Debe excluir PROD003 (23788.75) y PROD099 (333.00)
				for _, p := range products {
					assert.Less(t, p.Price, 100.0)
					assert.NotEqual(t, "PROD003", p.Code)
					assert.NotEqual(t, "PROD099", p.Code)
				}
			},
		},
		{
			name:          "filters by category clothing and priceLessThan 60",
			queryParams:   "category=clothing&priceLessThan=60",
			expectedCount: 1,
			expectedTotal: 1,
			validateResult: func(t *testing.T, products []Product) {
				assert.Equal(t, "PROD004", products[0].Code)
				assert.Equal(t, "clothing", products[0].Category)
				assert.Equal(t, 56.00, products[0].Price)
			},
		},
		{
			name:          "filters by category clothing and priceLessThan 70",
			queryParams:   "category=clothing&priceLessThan=70",
			expectedCount: 2,
			expectedTotal: 2,
			validateResult: func(t *testing.T, products []Product) {
				// Debe incluir PROD001 (65.76) y PROD004 (56.00)
				assert.Equal(t, "PROD001", products[0].Code)
				assert.Equal(t, "PROD004", products[1].Code)
				for _, p := range products {
					assert.Equal(t, "clothing", p.Category)
					assert.Less(t, p.Price, 70.0)
				}
			},
		},
		{
			name:          "filters by category shoes and priceLessThan 2",
			queryParams:   "category=shoes&priceLessThan=2",
			expectedCount: 1,
			expectedTotal: 1,
			validateResult: func(t *testing.T, products []Product) {
				assert.Equal(t, "PROD006", products[0].Code)
				assert.Equal(t, 1.50, products[0].Price)
			},
		},
		{
			name:          "filters by category accessories and priceLessThan 100",
			queryParams:   "category=accessories&priceLessThan=100",
			expectedCount: 1,
			expectedTotal: 1,
			validateResult: func(t *testing.T, products []Product) {
				// Solo PROD005 (9.99), excluye PROD003 (23788.75)
				assert.Equal(t, "PROD005", products[0].Code)
				assert.Equal(t, 9.99, products[0].Price)
			},
		},
		{
			name:          "returns empty array when no products match filters",
			queryParams:   "category=clothing&priceLessThan=1",
			expectedCount: 0,
			expectedTotal: 0,
			validateResult: func(t *testing.T, products []Product) {
				assert.Empty(t, products)
			},
		},
		{
			name:          "returns empty array for nonexistent category",
			queryParams:   "category=electronics",
			expectedCount: 0,
			expectedTotal: 0,
			validateResult: func(t *testing.T, products []Product) {
				assert.Empty(t, products)
			},
		},
		{
			name:          "filters work with pagination - first page",
			queryParams:   "category=clothing&offset=0&limit=1",
			expectedCount: 1,
			expectedTotal: 2, // Total clothing: PROD001, PROD004
			validateResult: func(t *testing.T, products []Product) {
				assert.Equal(t, "PROD001", products[0].Code)
				assert.Equal(t, "clothing", products[0].Category)
			},
		},
		{
			name:          "filters work with pagination - second page",
			queryParams:   "category=clothing&offset=1&limit=1",
			expectedCount: 1,
			expectedTotal: 2,
			validateResult: func(t *testing.T, products []Product) {
				assert.Equal(t, "PROD004", products[0].Code)
				assert.Equal(t, "clothing", products[0].Category)
			},
		},
		{
			name:          "filters by priceLessThan with pagination",
			queryParams:   "priceLessThan=100&offset=0&limit=3",
			expectedCount: 3,
			expectedTotal: 5, // Total < 100: PROD001, PROD002, PROD004, PROD005, PROD006
			validateResult: func(t *testing.T, products []Product) {
				for _, p := range products {
					assert.Less(t, p.Price, 100.0)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockRepo(testProducts, nil)
			handler := NewCatalogHandler(mockRepo)

			w := makeRequest(handler, http.MethodGet, "/catalog?"+tt.queryParams)

			assert.Equal(t, http.StatusOK, w.Code)

			response := parseResponse(t, w)
			assert.Len(t, response.Products, tt.expectedCount,
				"Expected %d products but got %d", tt.expectedCount, len(response.Products))
			assert.Equal(t, tt.expectedTotal, response.Total,
				"Expected total %d but got %d", tt.expectedTotal, response.Total)

			if tt.validateResult != nil {
				tt.validateResult(t, response.Products)
			}
		})
	}
}

func TestHandleGet_FilterValidation(t *testing.T) {
	testProducts := setupFilterTestProducts()

	tests := []struct {
		name              string
		queryParams       string
		expectedErrorText string
	}{
		{
			name:              "invalid priceLessThan - non-numeric",
			queryParams:       "priceLessThan=invalid",
			expectedErrorText: "priceLessThan",
		},
		{
			name:              "invalid priceLessThan - negative",
			queryParams:       "priceLessThan=-10",
			expectedErrorText: "priceLessThan",
		},
		{
			name:              "invalid priceLessThan - zero",
			queryParams:       "priceLessThan=0",
			expectedErrorText: "priceLessThan",
		},
		{
			name:              "invalid priceLessThan - special characters",
			queryParams:       "priceLessThan=10.5.5",
			expectedErrorText: "priceLessThan",
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
