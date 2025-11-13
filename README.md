# Go Hiring Challenge

Backend API for product catalog management with discount system, built with Go and Clean Architecture principles.

## Overview

This project implements a REST API for managing products, variants, and categories with an escalable discount system. The implementation focuses on maintainability, testability, and following Go best idiomatic practices.

## Key Features

- Product catalog with pagination and filtering
- Dynamic discount system using a Strategy Pattern
- Product variants with price inheritance
- Category management (CRUD operations)
- Postgres database with GORM
- Clean Architecture with proper layer separation
- Proper(-ish) test coverage, with unit and integration tests

## Tech Stack

- Go 1.22+
- PostgreSQL 16
- GORM (ORM)
- testcontainers-go (integration testing)
- Docker & Docker Compose

## Project Structure

The project follows Clean Architecture principles with clear separation of concerns:

```
cmd/
  server/         - Main application entry point
  seed/           - Database seeding utility

internal/
  domain/         - Business entities and core logic
    product/      - Product entities and repository interfaces
    discount/     - Discount strategies (Strategy Pattern)
  
  application/    - Use cases and business rules
    catalog/      - Product catalog service
    category/     - Category management service
  
  infrastructure/ - External concerns (frameworks, databases, HTTP)
    http/         - HTTP handlers and DTOs
    persistence/  - Database repositories (GORM)

pkg/
  database/       - Database connection utilities

sql/              - Migration scripts
```

## Getting Started

### Prerequisites

- Go 1.22 or higher
- Docker Desktop
- Make

### Quick Start with Docker (Recommended)

The fastest way to run the entire application:

```bash
# Build and start everything (Postgres + API)
make docker-build
make docker-run

# Seed the database
make docker-seed

# View logs
make docker-logs

# Test the API
curl http://localhost:8484/catalog
```

The API will be available at `http://localhost:8484`

To stop everything:
```bash
make docker-down
```

### Local Development Setup

If you prefer to run the application locally:

1. Clone the repository
```bash
git clone https://github.com/juanaguss/challenge-mytheresa.git
cd go-hiring-challenge-1.2.0
```

2. Install dependencies
```bash
make tidy
```

3. Start the database (Docker)
```bash
make docker-up
```

4. Seed the database
```bash
make seed
```

5. Run the server
```bash
make run
```

The API will be available at `http://localhost:8484`

## Available Commands

### Development
```bash
make tidy              # Install and vendor dependencies
make run               # Start the API server locally
make seed              # Seed database (destructive operation)
make fmt               # Format code
make lint              # Run linters
make check             # Run fmt + lint + test
make clean             # Remove generated files
```

### Docker
```bash
make docker-build      # Build application Docker image
make docker-run        # Start full stack (DB + API)
make docker-seed       # Seed database in Docker
make docker-logs       # Show container logs
make docker-down       # Stop all containers
make docker-up         # Start Postgres only (for local dev)
```

### Testing
```bash
make test              # Run unit tests
make test-integration  # Run integration tests (requires Docker)
make test-all          # Run all tests (unit + integration)
make test-coverage     # Generate HTML coverage report
make test-race         # Run tests with race detector
```
## API Endpoints

### Products

- `GET /catalog` - List products with pagination and filters
    - Query params: `offset`, `limit`, `category`, `priceLessThan`

- `GET /catalog/{code}` - Get product details with variants

### Categories

- `GET /categories` - List all categories
- `POST /categories` - Create a new category

## Architecture Decisions

### Clean Architecture

The project implements Clean Architecture with dependency inversion. Domain layer has no dependencies, application layer depends only on domain, and infrastructure depends on both.

### Strategy Pattern for Discounts

Discounts are implemented using the Strategy Pattern, making the system extensible without modifying existing code. New discount types can be added by implementing the Strategy interface.
This could be improved in further iterations to avoid manual loading of discount types.

### Repository Pattern

Data access is abstracted through repository interfaces declared in the application layer. This allows easy testing with mocks and potential database changes without affecting business logic.

### Testability

The project includes both unit tests and integration tests. Integration tests use testcontainers to run against real PostgreSQL instances, ensuring database interactions work correctly.

## Testing

Run unit tests (no Docker required):
```bash
make test
```

Run integration tests:
```bash
make test-integration
```

Run all tests with coverage:
```bash
make test-all
```

Current coverage: 53.9% (focusing on required and some additional paths)

## Environment Variables

The application uses a `.env` file for configuration:

```
HTTP_PORT=8484
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_DB=challenge
POSTGRES_PORT=5432
POSTGRES_SQL_DIR=./sql
```

## Business Rules

### Discount System

- Products in the "boots" category receive 30% discount
- Product with SKU "000003" receives 15% discount
- Discounts are not cumulative (first matching strategy wins)
- Original price is always shown alongside discounted price

### Product Variants

- Variants can have their own price or inherit from parent product
- All variants of a product are returned in the detail endpoint
- Price inheritance is handled at the repository level

## Dev Notes

- Uses `decimal.Decimal` for monetary calculations to avoid floating-point precision issues
- Build tags separate integration tests from unit tests
- Graceful shutdown handling for the HTTP server
