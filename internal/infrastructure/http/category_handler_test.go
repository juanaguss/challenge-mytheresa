package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
	"github.com/mytheresa/go-hiring-challenge/internal/infrastructure/http/mapper"
	"github.com/stretchr/testify/assert"
)

type mockCategoryService struct {
	categories      []product.Category
	createdCategory *product.Category
	err             error
}

func (m *mockCategoryService) GetCategories() ([]product.Category, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.categories, nil
}

func (m *mockCategoryService) CreateCategory(code, name string) (*product.Category, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.createdCategory, nil
}

func TestCategoryHandler_HandleGet(t *testing.T) {
	t.Run("returns all categories", func(t *testing.T) {
		categories := []product.Category{
			{ID: 1, Code: "clothing", Name: "Clothing"},
			{ID: 2, Code: "shoes", Name: "Shoes"},
			{ID: 3, Code: "accessories", Name: "Accessories"},
		}
		service := &mockCategoryService{categories: categories}
		handler := NewCategoryHandler(service)

		req := httptest.NewRequest("GET", "/categories", nil)
		w := httptest.NewRecorder()

		handler.HandleGet(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response categoriesResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Categories, 3)
		assert.Equal(t, "clothing", response.Categories[0].Code)
		assert.Equal(t, "Clothing", response.Categories[0].Name)
	})

	t.Run("returns empty array when no categories exist", func(t *testing.T) {
		service := &mockCategoryService{categories: []product.Category{}}
		handler := NewCategoryHandler(service)

		req := httptest.NewRequest("GET", "/categories", nil)
		w := httptest.NewRecorder()

		handler.HandleGet(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response categoriesResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Empty(t, response.Categories)
	})

	t.Run("returns 500 when service fails", func(t *testing.T) {
		service := &mockCategoryService{err: errors.New("database error")}
		handler := NewCategoryHandler(service)

		req := httptest.NewRequest("GET", "/categories", nil)
		w := httptest.NewRecorder()

		handler.HandleGet(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "database error")
	})
}

func TestCategoryHandler_HandlePost(t *testing.T) {
	t.Run("creates category successfully", func(t *testing.T) {
		createdCategory := &product.Category{
			ID:   4,
			Code: "kids",
			Name: "kids",
		}
		service := &mockCategoryService{createdCategory: createdCategory}
		handler := NewCategoryHandler(service)

		reqBody := mapper.CreateCategoryRequest{
			Code: "kids",
			Name: "kids",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.HandlePost(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response mapper.CategoryResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "kids", response.Code)
		assert.Equal(t, "kids", response.Name)
	})

	t.Run("returns 400 when code is non existent", func(t *testing.T) {
		service := &mockCategoryService{}
		handler := NewCategoryHandler(service)

		reqBody := mapper.CreateCategoryRequest{
			Code: "",
			Name: "Electronics",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.HandlePost(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "code is required")
	})

	t.Run("returns 400 when name is non existent", func(t *testing.T) {
		service := &mockCategoryService{}
		handler := NewCategoryHandler(service)

		reqBody := mapper.CreateCategoryRequest{
			Code: "electronics",
			Name: "",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.HandlePost(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "name is required")
	})

	t.Run("returns 400 when request body is invalid", func(t *testing.T) {
		service := &mockCategoryService{}
		handler := NewCategoryHandler(service)

		req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.HandlePost(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid request body")
	})

	t.Run("returns 500 when service fails", func(t *testing.T) {
		service := &mockCategoryService{err: errors.New("database error")}
		handler := NewCategoryHandler(service)

		reqBody := mapper.CreateCategoryRequest{
			Code: "electronics",
			Name: "Electronics",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.HandlePost(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "database error")
	})
}
