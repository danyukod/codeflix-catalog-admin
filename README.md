# Codeflix Catalog Admin

A Go HTTP service for administering the Codeflix video catalog. The current API manages categories and is organized using Clean Architecture and Domain-Driven Design principles.

## Requirements

- Go 1.27 or later
- Docker and Docker Compose (for MySQL)

## Getting started

Clone the repository and download its dependencies:

```bash
git clone https://github.com/danyukod/codeflix-catalog-admin.git
cd codeflix-catalog-admin
go mod download
```

Start the MySQL database used by the default configuration:

```bash
docker compose up -d mysql
```

The compose file creates a database named `shop` with the credentials `shop`/`shop` and publishes MySQL on `localhost:3306`. The application reads its settings from [`configs/application.yaml`](configs/application.yaml).

Run the service:

```bash
go run ./cmd/admin
```

The server listens on `http://localhost:8080`. It runs the database migrations automatically at startup.

To stop the database container:

```bash
docker compose down
```

## API

The category endpoints are:

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/categories` | List categories, with pagination and filtering |
| `POST` | `/categories` | Create a category |
| `GET` | `/categories/{id}` | Retrieve a category |
| `PUT` | `/categories/{id}` | Replace a category |
| `PATCH` | `/categories/{id}` | Partially update a category |
| `DELETE` | `/categories/{id}` | Delete a category |

Interactive API documentation is available at [`http://localhost:8080/docs`](http://localhost:8080/docs). The raw OpenAPI 3.1.2 document is available at [`/openapi.yaml`](http://localhost:8080/openapi.yaml) and in [`docs/openapi.yaml`](docs/openapi.yaml). The Swagger UI loads its assets from jsDelivr, so viewing `/docs` requires internet access.

Example requests are in [`requests/categories.http`](requests/categories.http) and can be run from GoLand or another HTTP client that supports `.http` files.

## Tests

Run the complete test suite with:

```bash
go test ./...
```

The integration tests use Testcontainers and require Docker. Unit and repository tests can run without the MySQL container.

## Project structure

- `cmd/admin/`: application entry point
- `configs/`: runtime configuration
- `docs/`: embedded OpenAPI specification and Swagger UI
- `internal/application/`: use cases
- `internal/domain/`: domain entities, validation, and pagination
- `internal/infrastructure/`: MySQL connection, migrations, and repositories
- `internal/interfaces/http/`: HTTP routing and controllers
- `requests/`: example API requests

