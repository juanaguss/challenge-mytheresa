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

type mockDiscountEngine struct {
	discountPercentage int
	discountedPrice    decimal.Decimal
}

func (m *mockDiscountEngine) ApplyDiscount(p product.Product) decimal.Decimal {
	return m.discountedPrice
}

func (m *mockDiscountEngine) GetDiscountPercentage(p product.Product) int {
	return m.discountPercentage
}

func TestService_GetProducts(t *testing.T) {
	t.Run("returns products with discounts from repository", func(t *testing.T) {
		expectedProducts := []product.Product{
			{ID: 1, Code: "PROD001", Price: decimal.NewFromFloat(10.99)},
		}
		repo := &mockRepository{products: expectedProducts, total: 1}
		discountEngine := &mockDiscountEngine{
			discountPercentage: 30,
			discountedPrice:    decimal.NewFromFloat(7.69),
		}
		service := NewService(repo, discountEngine)

		products, discountedPrices, discountPercentages, total, err := service.GetProducts(0, 10, product.Filter{})

		require.NoError(t, err)
		assert.Len(t, products, 1)
		assert.Len(t, discountedPrices, 1)
		assert.Len(t, discountPercentages, 1)
		assert.Equal(t, int64(1), total)
		assert.Equal(t, 7.69, discountedPrices[0])
		assert.Equal(t, 30, discountPercentages[0])
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		repo := &mockRepository{err: errors.New("db error")}
		discountEngine := &mockDiscountEngine{}
		service := NewService(repo, discountEngine)

		_, _, _, _, err := service.GetProducts(0, 10, product.Filter{})

		assert.Error(t, err)
	})
}
