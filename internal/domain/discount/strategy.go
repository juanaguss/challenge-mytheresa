package discount

import (
	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
	"github.com/shopspring/decimal"
)

// Strategy defines the interface for discount calculation strategies.
// Each strategy determines if a discount applies and calculates the percentage.
type Strategy interface {
	// AppliesTo checks if this discount strategy applies to the given product.
	AppliesTo(p product.Product) bool

	// CalculatePercentage returns the discount %.
	CalculatePercentage(p product.Product) int
}

// Engine orchestrates multiple discount strategies.
// It applies the first matching strategy (discounts are not stackable).
type Engine struct {
	strategies []Strategy
}

// NewEngine creates a discount engine with the given strategies.
// Strategies are evaluated in order, first match wins.
func NewEngine(strategies []Strategy) *Engine {
	return &Engine{strategies: strategies}
}

// ApplyDiscount calculates the discounted price for a product.
// Returns the original price if no discount applies.
func (e *Engine) ApplyDiscount(p product.Product) decimal.Decimal {
	for _, strategy := range e.strategies {
		if strategy.AppliesTo(p) {
			percentage := strategy.CalculatePercentage(p)
			discount := p.Price.Mul(decimal.NewFromInt(int64(percentage))).Div(decimal.NewFromInt(100))
			return p.Price.Sub(discount)
		}
	}
	return p.Price
}

// GetDiscountPercentage returns the discount percentage for a product.
// Returns 0 if no discount applies.
func (e *Engine) GetDiscountPercentage(p product.Product) int {
	for _, strategy := range e.strategies {
		if strategy.AppliesTo(p) {
			return strategy.CalculatePercentage(p)
		}
	}
	return 0
}
