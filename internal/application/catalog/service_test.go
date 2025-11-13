package catalog

import (
	"errors"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRepository struct {
	products []product.Product
	total    int64
	err      error
}

func (m *mockRepository) GetAll() ([]product.Product, error) {
	return m.products, m.err
}

func (m *mockRepository) GetFiltered(offset, limit int, filters product.Filter) ([]product.Product, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	return m.products, m.total, nil
}

func (m *mockRepository) GetByCode(code string) (*product.Product, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, p := range m.products {
		if p.Code == code {
			return &p, nil
		}
	}
	return nil, errors.New("product not found")
}

func TestService_GetProducts(t *testing.T) {
	t.Run("returns products from repository", func(t *testing.T) {
		expectedProducts := []product.Product{
			{ID: 1, Code: "PROD001", Price: decimal.NewFromFloat(10.99)},
		}
		repo := &mockRepository{products: expectedProducts, total: 1}
		service := NewService(repo)

		products, total, err := service.GetProducts(0, 10, product.Filter{})

		require.NoError(t, err)
		assert.Len(t, products, 1)
		assert.Equal(t, int64(1), total)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		repo := &mockRepository{err: errors.New("db error")}
		service := NewService(repo)

		_, _, err := service.GetProducts(0, 10, product.Filter{})

		assert.Error(t, err)
	})
}
