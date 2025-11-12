package catalog

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type Response struct {
	Products []Product `json:"products"`
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

// HandleGet retrieves all products with their categories.
func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	res, err := h.repo.GetAllProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := Response{
		Products: mapProductsToDTO(res),
	}

	api.OKResponse(w, response)
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
