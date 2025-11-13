package catalog

import (
	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
	"github.com/shopspring/decimal"
)

// ProductRepository defines operations for product persistence.
type ProductRepository interface {
	GetAll() ([]product.Product, error)
	GetFiltered(offset, limit int, filters product.Filter) ([]product.Product, int64, error)
	GetByCode(code string) (*product.Product, error)
}

// DiscountEngine defines operations for discount calculation.
type DiscountEngine interface {
	ApplyDiscount(p product.Product) decimal.Decimal
	GetDiscountPercentage(p product.Product) int
	GetVariantDiscountPercentage(sku string, p product.Product) int
}

// VariantDiscount holds discount information for a variant.
type VariantDiscount struct {
	DiscountedPrice float64
	Percentage      int
}

// Service defines operations for the catalog business logic.
type Service interface {
	GetProducts(offset, limit int, filters product.Filter) ([]product.Product, []float64, []int, int64, error)
	GetProductByCode(code string) (*product.Product, float64, int, map[string]VariantDiscount, error)
}

type service struct {
	repo           ProductRepository
	discountEngine DiscountEngine
}

// NewService creates a new catalog service.
func NewService(repo ProductRepository, discountEngine DiscountEngine) Service {
	return &service{
		repo:           repo,
		discountEngine: discountEngine,
	}
}

// GetProducts retrieves filtered and paginated products with discounts.
// Returns products, discounted prices, discount percentages, and total count.
func (s *service) GetProducts(offset, limit int, filters product.Filter) ([]product.Product, []float64, []int, int64, error) {
	products, total, err := s.repo.GetFiltered(offset, limit, filters)
	if err != nil {
		return nil, nil, nil, 0, err
	}

	discountedPrices := make([]float64, len(products))
	discountPercentages := make([]int, len(products))

	for i, p := range products {
		discountedPrices[i] = s.discountEngine.ApplyDiscount(p).InexactFloat64()
		discountPercentages[i] = s.discountEngine.GetDiscountPercentage(p)
	}

	return products, discountedPrices, discountPercentages, total, nil
}

// GetProductByCode retrieves a product by its code with discount applied.
// Returns product, discounted price, discount percentage, and variant discounts.
func (s *service) GetProductByCode(code string) (*product.Product, float64, int, map[string]VariantDiscount, error) {
	p, err := s.repo.GetByCode(code)
	if err != nil {
		return nil, 0, 0, nil, err
	}

	discountedPrice := s.discountEngine.ApplyDiscount(*p).InexactFloat64()
	discountPercentage := s.discountEngine.GetDiscountPercentage(*p)

	// Calculate discounts for each variant
	variantDiscounts := make(map[string]VariantDiscount)
	for _, v := range p.Variants {
		percentage := s.discountEngine.GetVariantDiscountPercentage(v.SKU, *p)
		discounted := v.Price.InexactFloat64()
		if percentage > 0 {
			discount := v.Price.Mul(decimal.NewFromInt(int64(percentage))).Div(decimal.NewFromInt(100))
			discounted = v.Price.Sub(discount).InexactFloat64()
		}
		variantDiscounts[v.SKU] = VariantDiscount{
			DiscountedPrice: discounted,
			Percentage:      percentage,
		}
	}

	return p, discountedPrice, discountPercentage, variantDiscounts, nil
}
