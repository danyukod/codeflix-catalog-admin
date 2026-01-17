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
