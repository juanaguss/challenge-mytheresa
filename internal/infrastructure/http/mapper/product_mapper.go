package mapper

import (
	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
)

// ProductResponse is a product in the catalog API response.
type ProductResponse struct {
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

// ToProductResponse converts a domain product to a DTO.
func ToProductResponse(p product.Product) ProductResponse {
	categoryCode := ""
	if p.Category != nil {
		categoryCode = p.Category.Code
	}

	return ProductResponse{
		Code:     p.Code,
		Price:    p.Price.InexactFloat64(),
		Category: categoryCode,
	}
}

// ToProductResponses converts a slice of domain products to product response DTOs.
func ToProductResponses(products []product.Product) []ProductResponse {
	responses := make([]ProductResponse, len(products))
	for i, p := range products {
		responses[i] = ToProductResponse(p)
	}
	return responses
}
