package catalog

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

const (
	// DefaultOffset is the default starting position for pagination.
	DefaultOffset = 0
	// DefaultLimit is the default number of products per page.
	DefaultLimit = 10
	// MinLimit is the minimum allowed products per page.
	MinLimit = 1
	// MaxLimit is the maximum allowed products per page.
	MaxLimit = 100
)

type Response struct {
	Products []Product `json:"products"`
	Total    int       `json:"total"`
}

type Product struct {
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

// Handler handles HTTP requests for the product catalog.
type Handler struct {
	repo models.ProductsRepository
}

// NewCatalogHandler creates a new catalog handler with the given repository.
func NewCatalogHandler(r models.ProductsRepository) *Handler {
	return &Handler{
		repo: r,
	}
}

// HandleGet retrieves all products with their categories. Accepts paginations
func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	offset, limit, err := parsePaginationParams(r)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	products, total, err := h.repo.GetProducts(offset, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := Response{
		Products: mapProductsToDTO(products),
		Total:    int(total),
	}

	api.OKResponse(w, response)
}

// parsePaginationParams extracts and validates offset and limit from query params.
func parsePaginationParams(r *http.Request) (offset, limit int, err error) {
	offset = DefaultOffset
	limit = DefaultLimit

	// Parse offset
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		var parseErr error
		offset, parseErr = strconv.Atoi(offsetStr)
		if parseErr != nil {
			return 0, 0, fmt.Errorf("invalid offset parameter")
		}
		if offset < 0 {
			return 0, 0, fmt.Errorf("offset must be non-negative")
		}
	}

	// Parse limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		var parseErr error
		limit, parseErr = strconv.Atoi(limitStr)
		if parseErr != nil {
			return 0, 0, fmt.Errorf("invalid limit parameter")
		}
		if limit < MinLimit {
			return 0, 0, fmt.Errorf("limit must be at least %d", MinLimit)
		}
		if limit > MaxLimit {
			return 0, 0, fmt.Errorf("limit must not exceed %d", MaxLimit)
		}
	}

	return offset, limit, nil
}

func mapProductsToDTO(products []models.Product) []Product {
	result := make([]Product, len(products))
	for i, p := range products {
		categoryCode := ""
		if p.Category != nil {
			categoryCode = p.Category.Code
		}
		result[i] = Product{
			Code:     p.Code,
			Price:    p.Price.InexactFloat64(),
			Category: categoryCode,
		}
	}
	return result
}
