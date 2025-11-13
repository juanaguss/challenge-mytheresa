package http

import (
	"encoding/json"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/internal/application/category"
	"github.com/mytheresa/go-hiring-challenge/internal/infrastructure/http/mapper"
)

type categoriesResponse struct {
	Categories []mapper.CategoryResponse `json:"categories"`
}

// CategoryHandler handles HTTP requests for categories.
type CategoryHandler struct {
	service category.Service
}

// NewCategoryHandler creates a new category HTTP handler.
func NewCategoryHandler(service category.Service) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// HandleGet handles GET /categories requests.
func (h *CategoryHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetCategories()
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := categoriesResponse{
		Categories: mapper.ToCategoryResponses(categories),
	}

	okResponse(w, response)
}

// HandlePost handles POST /categories requests.
func (h *CategoryHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	var req mapper.CreateCategoryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Code == "" {
		errorResponse(w, http.StatusBadRequest, "code is required")
		return
	}

	if req.Name == "" {
		errorResponse(w, http.StatusBadRequest, "name is required")
		return
	}

	cat, err := h.service.CreateCategory(req.Code, req.Name)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := mapper.ToCategoryResponse(*cat)

	w.WriteHeader(http.StatusCreated)
	okResponse(w, response)
}
