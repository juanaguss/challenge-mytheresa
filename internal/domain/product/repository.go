package product

// Repository defines operations for product persistence.
type Repository interface {
	GetAll() ([]Product, error)
	GetFiltered(offset, limit int, filters Filter) ([]Product, int64, error)
}
