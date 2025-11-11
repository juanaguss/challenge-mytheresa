package models

// Category represents a product category in the catalog.
// Categories serve to organize and filter products.
type Category struct {
	ID   uint   `grom:"primary_key"`
	Code string `grom:"uniqueIndex;not null; size:32"`
	Name string `grom:"not null; size:32"`
}

func (c *Category) TableName() string {
	return "categories"
}
