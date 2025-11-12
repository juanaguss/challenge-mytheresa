package models

import (
	"gorm.io/gorm"
)

const (
	// RelationVariants is the name of the Variants relationship for GORM.
	RelationVariants = "Variants"
	// RelationCategory is the name of the Category relationship for GORM.
	RelationCategory = "Category"
)

type ProductsRepository interface {
	GetAllProducts() ([]Product, error)
	GetProducts(offset, limit int) ([]Product, int64, error)
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

func (r *gormProductsRepository) GetProducts(offset, limit int) ([]Product, int64, error) {
	var products []Product
	var total int64

	if err := r.db.Model(&Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload(RelationVariants).Preload(RelationCategory).
		Offset(offset).Limit(limit).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
