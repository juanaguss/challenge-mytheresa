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

// VariantResponse represents a product variant in the response.
type VariantResponse struct {
	Code       string   `json:"code"`
	Price      float64  `json:"price"`
	Discount   *string  `json:"discount,omitempty"`
	FinalPrice *float64 `json:"final_price,omitempty"`
}

// ProductDetailResponse represents product information with variants.
type ProductDetailResponse struct {
	Code       string            `json:"code"`
	Price      float64           `json:"price"`
	Category   string            `json:"category"`
	Discount   *string           `json:"discount,omitempty"`
	FinalPrice *float64          `json:"final_price,omitempty"`
	Variants   []VariantResponse `json:"variants"`
}

// VariantDiscountInfo holds discount information for a variant.
type VariantDiscountInfo struct {
	DiscountedPrice float64
	Percentage      int
}

// ToProductDetailResponse converts a domain product with variants to DTO.
func ToProductDetailResponse(p product.Product, discountedPrice float64, discountPercentage int, variantDiscounts map[string]VariantDiscountInfo) ProductDetailResponse {
	categoryCode := ""
	if p.Category != nil {
		categoryCode = p.Category.Code
	}

	variants := make([]VariantResponse, len(p.Variants))
	for i, v := range p.Variants {
		variant := VariantResponse{
			Code:  v.SKU,
			Price: v.Price.InexactFloat64(),
		}

		// Apply variant-specific discount if available
		if discountInfo, ok := variantDiscounts[v.SKU]; ok && discountInfo.Percentage > 0 {
			discountStr := fmt.Sprintf("%d%%", discountInfo.Percentage)
			variant.Discount = &discountStr
			variant.FinalPrice = &discountInfo.DiscountedPrice
		}

		variants[i] = variant
	}

	response := ProductDetailResponse{
		Code:     p.Code,
		Price:    p.Price.InexactFloat64(),
		Category: categoryCode,
		Variants: variants,
	}

	if discountPercentage > 0 {
		discountStr := fmt.Sprintf("%d%%", discountPercentage)
		response.Discount = &discountStr
		response.FinalPrice = &discountedPrice
	}

	return response
}
