package persistence

import (
	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
	"gorm.io/gorm"
)

// CategoryRepository implements category ops using GORM.
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new GORM category repository.
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// GetAll retrieves all categories.
func (r *CategoryRepository) GetAll() ([]product.Category, error) {
	var models []categoryModel

	err := r.db.Find(&models).Error
	if err != nil {
		return nil, err
	}

	categories := make([]product.Category, len(models))
	for i, m := range models {
		categories[i] = product.Category{
			ID:   m.ID,
			Code: m.Code,
			Name: m.Name,
		}
	}

	return categories, nil
}

// Create creates a new category.
func (r *CategoryRepository) Create(cat product.Category) (*product.Category, error) {
	model := categoryModel{
		Code: cat.Code,
		Name: cat.Name,
	}

	err := r.db.Create(&model).Error
	if err != nil {
		return nil, err
	}

	return &product.Category{
		ID:   model.ID,
		Code: model.Code,
		Name: model.Name,
	}, nil
}
