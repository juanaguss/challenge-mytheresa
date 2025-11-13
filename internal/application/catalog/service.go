package catalog

import (
	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
)

// ProductRepository defines operations for product persistence.
type ProductRepository interface {
	GetAll() ([]product.Product, error)
	GetFiltered(offset, limit int, filters product.Filter) ([]product.Product, int64, error)
	GetByCode(code string) (*product.Product, error)
}

// Service defines operations for the catalog business logic.
type Service interface {
	GetProducts(offset, limit int, filters product.Filter) ([]product.Product, int64, error)
}

type service struct {
	repo ProductRepository
}

// NewService creates a new catalog service.
func NewService(repo ProductRepository) Service {
	return &service{repo: repo}
}

// GetProducts retrieves filtered and paginated products.
func (s *service) GetProducts(offset, limit int, filters product.Filter) ([]product.Product, int64, error) {
	return s.repo.GetFiltered(offset, limit, filters)
}
