package mapper

import (
	"fmt"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
)

// ProductResponse is a product in the catalog API response.
type ProductResponse struct {
	Code       string   `json:"code"`
	Price      float64  `json:"price"`
	Category   string   `json:"category"`
	Discount   *string  `json:"discount,omitempty"`
	FinalPrice *float64 `json:"final_price,omitempty"`
}

// ToProductResponse converts a domain product to a DTO.
// If discountPercentage > 0, includes discount info and final price.
func ToProductResponse(p product.Product, discountedPrice float64, discountPercentage int) ProductResponse {
	categoryCode := ""
	if p.Category != nil {
		categoryCode = p.Category.Code
	}

	response := ProductResponse{
		Code:     p.Code,
		Price:    p.Price.InexactFloat64(),
		Category: categoryCode,
	}

	if discountPercentage > 0 {
		discountStr := fmt.Sprintf("%d%%", discountPercentage)
		response.Discount = &discountStr
		response.FinalPrice = &discountedPrice
	}

	return response
}

// ToProductResponses converts a slice of domain products to product response DTOs.
// Takes parallel slices of discounted prices and percentages.
func ToProductResponses(products []product.Product, discountedPrices []float64, discountPercentages []int) []ProductResponse {
	responses := make([]ProductResponse, len(products))
	for i, p := range products {
		responses[i] = ToProductResponse(p, discountedPrices[i], discountPercentages[i])
	}
	return responses
}
