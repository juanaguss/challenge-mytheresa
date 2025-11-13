package mapper

import "github.com/mytheresa/go-hiring-challenge/internal/domain/product"

// CategoryResponse is a category in the API.
type CategoryResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// CreateCategoryRequest represents the request body for creating a category.
type CreateCategoryRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// ToCategoryResponse converts a domain category to a response DTO.
func ToCategoryResponse(cat product.Category) CategoryResponse {
	return CategoryResponse{
		Code: cat.Code,
		Name: cat.Name,
	}
}

// ToCategoryResponses converts a slice of domain categories to response DTOs.
func ToCategoryResponses(categories []product.Category) []CategoryResponse {
	responses := make([]CategoryResponse, len(categories))
	for i, cat := range categories {
		responses[i] = ToCategoryResponse(cat)
	}
	return responses
}
