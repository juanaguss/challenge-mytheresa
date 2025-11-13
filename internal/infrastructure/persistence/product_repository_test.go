//go:build integration
// +build integration

package persistence

import (
	"context"
	"testing"
	"time"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/product"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	pgdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupTestDB creates a postgres testcontainer for IT
func setupTestDB(t *testing.T) *gorm.DB {
	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second)),
	)
	require.NoError(t, err, "Failed to start PostgreSQL container")

	t.Cleanup(func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	})

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err, "Failed to get connection string")

	db, err := gorm.Open(pgdriver.Open(connStr), &gorm.Config{})
	require.NoError(t, err, "Failed to connect to PostgreSQL container")

	err = db.AutoMigrate(&productModel{}, &categoryModel{}, &variantModel{})
	require.NoError(t, err, "Failed to migrate database schema")

	return db
}

func seedTestData(t *testing.T, db *gorm.DB) {
	categories := []categoryModel{
		{ID: 1, Code: "clothing", Name: "Clothing"},
		{ID: 2, Code: "shoes", Name: "Shoes"},
		{ID: 3, Code: "accessories", Name: "Accessories"},
	}
	for _, cat := range categories {
		require.NoError(t, db.Create(&cat).Error)
	}

	products := []productModel{
		{
			ID:         1,
			Code:       "PROD001",
			Price:      "89.99",
			CategoryID: uintPtr(1),
		},
		{
			ID:         2,
			Code:       "PROD002",
			Price:      "129.99",
			CategoryID: uintPtr(2),
		},
		{
			ID:         3,
			Code:       "PROD003",
			Price:      "49.99",
			CategoryID: uintPtr(3),
		},
		{
			ID:         4,
			Code:       "PROD004",
			Price:      "199.99",
			CategoryID: uintPtr(1),
		},
		{
			ID:    5,
			Code:  "PROD005",
			Price: "15.50",
			// empty category
		},
	}
	for _, prod := range products {
		require.NoError(t, db.Create(&prod).Error)
	}

	priceSmall := "89.99"
	variants := []variantModel{
		{
			ID:        1,
			ProductID: 1,
			Name:      "Small",
			SKU:       "PROD001-S",
			Price:     &priceSmall,
		},
		{
			ID:        2,
			ProductID: 1,
			Name:      "Large",
			SKU:       "PROD001-L",
			Price:     nil,
		},
	}
	for _, v := range variants {
		require.NoError(t, db.Create(&v).Error)
	}
}

func uintPtr(u uint) *uint {
	return &u
}

func TestProductRepository_GetAll(t *testing.T) {
	t.Run("returns all products with category and variants", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestData(t, db)
		repo := NewProductRepository(db)

		products, err := repo.GetAll()

		require.NoError(t, err)
		assert.Len(t, products, 5)

		prod1 := findProductByCode(products, "PROD001")
		require.NotNil(t, prod1)
		assert.Equal(t, "PROD001", prod1.Code)
		assert.Equal(t, "89.99", prod1.Price.String())
		require.NotNil(t, prod1.Category)
		assert.Equal(t, "clothing", prod1.Category.Code)
		assert.Len(t, prod1.Variants, 2)

		prod5 := findProductByCode(products, "PROD005")
		require.NotNil(t, prod5)
		assert.Nil(t, prod5.Category)
	})

	t.Run("returns empty slice when no products exist", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewProductRepository(db)

		products, err := repo.GetAll()

		require.NoError(t, err)
		assert.Empty(t, products)
		assert.NotNil(t, products)
	})
}

func TestProductRepository_GetFiltered(t *testing.T) {
	t.Run("returns filtered products by category", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestData(t, db)
		repo := NewProductRepository(db)

		filters := product.Filter{
			Category: "clothing",
		}

		products, total, err := repo.GetFiltered(0, 10, filters)

		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, products, 2)

		for _, p := range products {
			require.NotNil(t, p.Category)
			assert.Equal(t, "clothing", p.Category.Code)
		}
	})

	t.Run("returns filtered products by price less than", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestData(t, db)
		repo := NewProductRepository(db)

		maxPrice := decimal.NewFromFloat(100.0)
		filters := product.Filter{
			PriceLessThan: &maxPrice,
		}

		products, total, err := repo.GetFiltered(0, 10, filters)

		require.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, products, 3)

		for _, p := range products {
			assert.True(t, p.Price.LessThan(maxPrice),
				"Product %s price %s should be less than %s", p.Code, p.Price, maxPrice)
		}
	})

	t.Run("returns filtered products by category and price", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestData(t, db)
		repo := NewProductRepository(db)

		maxPrice := decimal.NewFromFloat(150.0)
		filters := product.Filter{
			Category:      "clothing",
			PriceLessThan: &maxPrice,
		}

		products, total, err := repo.GetFiltered(0, 10, filters)

		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, products, 1)
		assert.Equal(t, "PROD001", products[0].Code)
	})

	t.Run("applies pagination correctly", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestData(t, db)
		repo := NewProductRepository(db)

		//First page
		products1, total1, err := repo.GetFiltered(0, 2, product.Filter{})

		require.NoError(t, err)
		assert.Equal(t, int64(5), total1)
		assert.Len(t, products1, 2)

		//Second page
		products2, total2, err := repo.GetFiltered(2, 2, product.Filter{})

		require.NoError(t, err)
		assert.Equal(t, int64(5), total2)
		assert.Len(t, products2, 2)

		assert.NotEqual(t, products1[0].Code, products2[0].Code)
	})

	t.Run("returns empty when no products match filter", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestData(t, db)
		repo := NewProductRepository(db)

		filters := product.Filter{
			Category: "nonexistent",
		}

		products, total, err := repo.GetFiltered(0, 10, filters)

		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, products)
	})

	t.Run("preloads category and variants", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestData(t, db)
		repo := NewProductRepository(db)

		products, _, err := repo.GetFiltered(0, 1, product.Filter{Category: "clothing"})

		require.NoError(t, err)
		require.Len(t, products, 1)

		prod := products[0]
		assert.NotNil(t, prod.Category, "Category should be preloaded")
		assert.Equal(t, "clothing", prod.Category.Code)
		assert.NotEmpty(t, prod.Variants, "Variants should be preloaded")
	})
}

func TestProductRepository_GetByCode(t *testing.T) {
	t.Run("returns product by code with relations", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestData(t, db)
		repo := NewProductRepository(db)

		prod, err := repo.GetByCode("PROD001")

		require.NoError(t, err)
		require.NotNil(t, prod)
		assert.Equal(t, "PROD001", prod.Code)
		assert.Equal(t, "89.99", prod.Price.String())

		require.NotNil(t, prod.Category)
		assert.Equal(t, "clothing", prod.Category.Code)

		assert.Len(t, prod.Variants, 2)
		assert.Equal(t, "PROD001-S", prod.Variants[0].SKU)
	})

	t.Run("returns error when product not found", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestData(t, db)
		repo := NewProductRepository(db)

		prod, err := repo.GetByCode("NONEXISTENT")

		assert.Error(t, err)
		assert.Nil(t, prod)
	})
}

func TestProductRepository_DomainMapping(t *testing.T) {
	t.Run("correctly maps GORM model to domain entity", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestData(t, db)
		repo := NewProductRepository(db)

		products, err := repo.GetAll()

		require.NoError(t, err)
		prod := findProductByCode(products, "PROD001")
		require.NotNil(t, prod)

		assert.IsType(t, product.Product{}, *prod)
		assert.IsType(t, decimal.Decimal{}, prod.Price)
		assert.IsType(t, &product.Category{}, prod.Category)
		assert.IsType(t, []product.Variant{}, prod.Variants)
	})

	t.Run("variant without price inherits from product", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestData(t, db)
		repo := NewProductRepository(db)

		prod, err := repo.GetByCode("PROD001")

		require.NoError(t, err)

		var variantWithoutPrice *product.Variant
		for i, v := range prod.Variants {
			if v.SKU == "PROD001-L" {
				variantWithoutPrice = &prod.Variants[i]
				break
			}
		}

		require.NotNil(t, variantWithoutPrice)
		assert.Equal(t, prod.Price, variantWithoutPrice.Price,
			"Variant without price should inherit price from parent product")
	})
}

func findProductByCode(products []product.Product, code string) *product.Product {
	for i, p := range products {
		if p.Code == code {
			return &products[i]
		}
	}
	return nil
}
