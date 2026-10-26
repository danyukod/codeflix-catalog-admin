# codeflix-catalog-admin

A Go-based administrative service for managing the Codeflix video catalog. This project follows Clean Architecture principles and Domain-Driven Design (DDD).

## Project Structure

- `cmd/`: Entry points for the application.
  - `admin/`: The main administrative application.
- `internal/`: Private application and library code.
  - `application/`: Application services and use cases.
  - `domain/`: Business logic and entities (e.g., Category).
  - `infrastructure/`: External implementations (database, external APIs).

## Getting Started

### Prerequisites

- Go 1.25.5 or higher

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/danyukod/codeflix-catalog-admin.git
   cd codeflix-catalog-admin
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

### Running the application

To run the admin service:
```bash
go run cmd/admin/main.go
```

## API documentation

The OpenAPI 3.1.2 specification is in [`docs/openapi.yaml`](docs/openapi.yaml) and is served at `GET /openapi.yaml` when the application is running. Open `http://localhost:8080/docs` for the visual Swagger UI documentation. The page loads Swagger UI 5.33.1 assets from jsDelivr, so the browser needs internet access. The same specification can be opened in GoLand to navigate schemas.

