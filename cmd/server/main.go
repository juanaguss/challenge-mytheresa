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
	"github.com/mytheresa/go-hiring-challenge/internal/domain/discount"
	httpHandler "github.com/mytheresa/go-hiring-challenge/internal/infrastructure/http"
	"github.com/mytheresa/go-hiring-challenge/internal/infrastructure/persistence"
	"github.com/mytheresa/go-hiring-challenge/pkg/database"
)

// buildDiscountEngine constructs the discount engine with business rules.
func buildDiscountEngine() *discount.Engine {
	strategies := []discount.Strategy{
		discount.NewCategoryDiscountStrategy("boots", 30), // 30% off boots
		discount.NewSKUDiscountStrategy("000003", 15),     // 15% off on this SKU
	}
	return discount.NewEngine(strategies)
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, close := database.New(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)
	defer close()

	productRepo := persistence.NewProductRepository(db)
	discountEngine := buildDiscountEngine()
	catalogService := catalog.NewService(productRepo, discountEngine)
	catalogHandler := httpHandler.NewCatalogHandler(catalogService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog", catalogHandler.HandleGet)

	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", os.Getenv("HTTP_PORT")),
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
