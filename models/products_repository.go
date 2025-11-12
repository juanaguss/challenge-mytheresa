package models

import (
	"gorm.io/gorm"
)

const (
	RelationVariants = "Variants"
	RelationCategory = "Category"
)

type ProductsRepository interface {
	GetAllProducts() ([]Product, error)
	GetProducts(offset, limit int, filters ProductFilters) ([]Product, int64, error)
}

type gormProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) ProductsRepository {
	return &gormProductsRepository{
		db: db,
	}
}

func (r *gormProductsRepository) GetAllProducts() ([]Product, error) {
	var products []Product
	if err := r.db.Preload(RelationVariants).Preload(RelationCategory).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *gormProductsRepository) GetProducts(offset, limit int, filters ProductFilters) ([]Product, int64, error) {
	var products []Product
	var total int64

	query := r.applyFilters(r.db.Model(&Product{}), filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = r.applyFilters(r.db, filters).
		Preload(RelationVariants).
		Preload(RelationCategory).
		Offset(offset).
		Limit(limit)

	if err := query.Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// applyFilters implements filters to a GROM query
func (r *gormProductsRepository) applyFilters(query *gorm.DB, filters ProductFilters) *gorm.DB {
	if filters.Category != "" {
		query = query.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.code = ?", filters.Category)
	}

	if filters.PriceLessThan != nil {
		query = query.Where("products.price < ?", filters.PriceLessThan)
	}

	return query
}