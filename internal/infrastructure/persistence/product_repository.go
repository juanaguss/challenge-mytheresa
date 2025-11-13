package persistence

import (
	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const (
	relationVariants = "Variants"
	relationCategory = "Category"
)

type productModel struct {
	ID         uint           `gorm:"primaryKey"`
	Code       string         `gorm:"uniqueIndex;not null"`
	Price      string         `gorm:"type:decimal(10,2);not null"`
	CategoryID *uint          `gorm:"index"`
	Category   *categoryModel `gorm:"foreignKey:CategoryID"`
	Variants   []variantModel `gorm:"foreignKey:ProductID"`
}

func (productModel) TableName() string {
	return "products"
}

type categoryModel struct {
	ID   uint   `gorm:"primaryKey"`
	Code string `gorm:"uniqueIndex;not null;size:32"`
	Name string `gorm:"not null;size:256"`
}

func (categoryModel) TableName() string {
	return "categories"
}

type variantModel struct {
	ID        uint    `gorm:"primaryKey"`
	ProductID uint    `gorm:"not null"`
	Name      string  `gorm:"not null"`
	SKU       string  `gorm:"uniqueIndex;not null"`
	Price     *string `gorm:"type:decimal(10,2)"`
}

func (variantModel) TableName() string {
	return "product_variants"
}

// ProductRepository implements product.Repository using GORM.
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new GORM product repository.
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// GetAll retrieves all products with their relations.
func (r *ProductRepository) GetAll() ([]product.Product, error) {
	var models []productModel

	err := r.db.
		Preload(relationVariants).
		Preload(relationCategory).
		Find(&models).Error

	if err != nil {
		return nil, err
	}

	return toDomainProducts(models), nil
}

// GetFiltered retrieves products with pagination and filtering applied.
// Returns the filtered products and the total count of products matching the filters.
func (r *ProductRepository) GetFiltered(offset, limit int, filters product.Filter) ([]product.Product, int64, error) {
	var models []productModel
	var total int64

	query := r.applyFilters(r.db.Model(&productModel{}), filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = r.applyFilters(r.db, filters).
		Preload(relationVariants).
		Preload(relationCategory).
		Offset(offset).
		Limit(limit)

	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	return toDomainProducts(models), total, nil
}

func (r *ProductRepository) applyFilters(query *gorm.DB, filters product.Filter) *gorm.DB {
	if filters.Category != "" {
		query = query.
			Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.code = ?", filters.Category)
	}

	if filters.PriceLessThan != nil {
		query = query.Where("products.price < ?", filters.PriceLessThan)
	}

	return query
}

func toDomainProducts(models []productModel) []product.Product {
	products := make([]product.Product, len(models))
	for i, m := range models {
		products[i] = toDomainProduct(m)
	}
	return products
}

func toDomainProduct(m productModel) product.Product {
	p := product.Product{
		ID:         m.ID,
		Code:       m.Code,
		CategoryID: m.CategoryID,
	}

	p.Price, _ = decimal.NewFromString(m.Price)

	if m.Category != nil {
		p.Category = &product.Category{
			ID:   m.Category.ID,
			Code: m.Category.Code,
			Name: m.Category.Name,
		}
	}

	if len(m.Variants) > 0 {
		p.Variants = make([]product.Variant, len(m.Variants))
		for i, v := range m.Variants {
			var price decimal.Decimal
			if v.Price != nil {
				price, _ = decimal.NewFromString(*v.Price)
			} else {
				price = p.Price
			}
			p.Variants[i] = product.Variant{
				ID:        v.ID,
				ProductID: v.ProductID,
				Name:      v.Name,
				SKU:       v.SKU,
				Price:     price,
			}
		}
	}

	return p
}
