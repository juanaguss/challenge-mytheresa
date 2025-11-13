package category

import "github.com/mytheresa/go-hiring-challenge/internal/domain/product"

// Repository defines ops for category persistence.
type Repository interface {
	GetAll() ([]product.Category, error)
	Create(cat product.Category) (*product.Category, error)
}

// Service defines ops for category business logic.
type Service interface {
	GetCategories() ([]product.Category, error)
	CreateCategory(code, name string) (*product.Category, error)
}

type service struct {
	repo Repository
}

// NewService creates a new category service.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// GetCategories retrieves all categories.
func (s *service) GetCategories() ([]product.Category, error) {
	return s.repo.GetAll()
}

// CreateCategory creates a new category.
func (s *service) CreateCategory(code, name string) (*product.Category, error) {
	cat := product.Category{
		Code: code,
		Name: name,
	}
	return s.repo.Create(cat)
}
