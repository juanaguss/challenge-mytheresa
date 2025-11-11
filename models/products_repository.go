package models

import (
	"gorm.io/gorm"
)

type ProductsRepository interface {
	GetAllProducts() ([]Product, error)
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
	if err := r.db.Preload("Variants").Preload("Category").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
