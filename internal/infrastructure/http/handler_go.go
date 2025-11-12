package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/internal/application/catalog"
	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
	"github.com/mytheresa/go-hiring-challenge/internal/infrastructure/http/mapper"
	"github.com/shopspring/decimal"
)

const (
	defaultOffset = 0
	defaultLimit  = 10
	minLimit      = 1
	maxLimit      = 100
)

type catalogResponse struct {
	Products []mapper.ProductResponse `json:"products"`
	Total    int                      `json:"total"`
}

// CatalogHandler handles HTTP requests for the product catalog.
type CatalogHandler struct {
	service catalog.Service
}

// NewCatalogHandler creates a new catalog HTTP handler.
func NewCatalogHandler(service catalog.Service) *CatalogHandler {
	return &CatalogHandler{service: service}
}

// HandleGet handles GET /catalog requests.
// Supports optional query parameters: offset, limit, category, and priceLessThan.
func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	offset, limit, err := parsePaginationParams(r)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	filters, err := parseFilterParams(r)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	products, total, err := h.service.GetProducts(offset, limit, filters)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := catalogResponse{
		Products: mapper.ToProductResponses(products),
		Total:    int(total),
	}

	okResponse(w, response)
}

func parsePaginationParams(r *http.Request) (offset, limit int, err error) {
	offset = defaultOffset
	limit = defaultLimit

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid offset parameter")
		}
		if offset < 0 {
			return 0, 0, fmt.Errorf("offset must be non-negative")
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid limit parameter")
		}
		if limit < minLimit {
			return 0, 0, fmt.Errorf("limit must be at least %d", minLimit)
		}
		if limit > maxLimit {
			return 0, 0, fmt.Errorf("limit must not exceed %d", maxLimit)
		}
	}

	return offset, limit, nil
}

func parseFilterParams(r *http.Request) (product.Filter, error) {
	var filters product.Filter

	if category := r.URL.Query().Get("category"); category != "" {
		filters.Category = category
	}

	if priceStr := r.URL.Query().Get("priceLessThan"); priceStr != "" {
		price, err := decimal.NewFromString(priceStr)
		if err != nil {
			return filters, fmt.Errorf("invalid priceLessThan parameter")
		}
		if price.LessThanOrEqual(decimal.Zero) {
			return filters, fmt.Errorf("priceLessThan must be greater than 0")
		}
		filters.PriceLessThan = &price
	}

	return filters, nil
}
