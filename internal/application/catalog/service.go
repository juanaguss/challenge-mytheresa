package catalog

import (
	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
)

// Service defines operations for the catalog business logic.
type Service interface {
	GetProducts(offset, limit int, filters product.Filter) ([]product.Product, int64, error)
}

type service struct {
	repo product.Repository
}

// NewService creates a new catalog service.
func NewService(repo product.Repository) Service {
	return &service{repo: repo}
}

// GetProducts retrieves filtered and paginated products.
func (s *service) GetProducts(offset, limit int, filters product.Filter) ([]product.Product, int64, error) {
	return s.repo.GetFiltered(offset, limit, filters)
}
