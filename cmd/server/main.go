package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/mytheresa/go-hiring-challenge/internal/application/catalog"
	"github.com/mytheresa/go-hiring-challenge/internal/application/category"
	"github.com/mytheresa/go-hiring-challenge/internal/domain/discount"
	httpHandler "github.com/mytheresa/go-hiring-challenge/internal/infrastructure/http"
	"github.com/mytheresa/go-hiring-challenge/internal/infrastructure/persistence"
	"github.com/mytheresa/go-hiring-challenge/pkg/database"
)

// buildDiscountEngine constructs the discount engine with the required business rules.
// 30% off boots category, 15% off SKU 000003.
func buildDiscountEngine() *discount.Engine {
	strategies := []discount.Strategy{
		discount.NewCategoryDiscountStrategy("boots", 30),
		discount.NewSKUDiscountStrategy("000003", 15),
	}
	return discount.NewEngine(strategies)
}

func main() {
	_ = godotenv.Load(".env")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}
	log.Printf("Conexion de db %s:%s", host, os.Getenv("POSTGRES_PORT"))

	db, close := database.New(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		host,
		os.Getenv("POSTGRES_PORT"),
	)
	defer close()

	productRepo := persistence.NewProductRepository(db)
	categoryRepo := persistence.NewCategoryRepository(db)

	discountEngine := buildDiscountEngine()
	catalogService := catalog.NewService(productRepo, discountEngine)
	categoryService := category.NewService(categoryRepo)

	catalogHandler := httpHandler.NewCatalogHandler(catalogService)
	categoryHandler := httpHandler.NewCategoryHandler(categoryService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog", catalogHandler.HandleGet)
	mux.HandleFunc("GET /catalog/{code}", catalogHandler.HandleGetByCode)
	mux.HandleFunc("GET /categories", categoryHandler.HandleGet)
	mux.HandleFunc("POST /categories", categoryHandler.HandlePost)

	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", os.Getenv("HTTP_PORT")),
		Handler: mux,
	}

	go func() {
		log.Printf("Starting server on http://%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %s", err)
		}
		log.Println("Server stopped gracefully")
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("Server shutdown error: %s", err)
	}
}
