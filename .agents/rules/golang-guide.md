---
description: Golang Development Guidelines for Brayat
---

# Golang Guide for Brayat

These rules ensure the Go backend is maintainable, decoupled, and adheres to clean architecture principles.

## 1. Architecture: Repository Pattern

The application must strictly follow the Repository Pattern, divided into four distinct layers:

### A. Handler Layer (Delivery)

- Responsible for HTTP request/response handling, JSON serialization, and parameter extraction.
- Maps HTTP errors from service layer errors.
- **Dependency**: Interfaces from the Service layer.

### B. Service Layer (Business Logic)

- Contains all core business logic and validations.
- **Dependency**: Interfaces from the Repository layer.

### C. Repository Layer (Data Access)

- Handles all database interactions and specific queries.
- Returns domain models.
- **Dependency**: Database connection pool (`*sql.DB`).

### D. Model Layer

- Defines the core business entities/structs (e.g., `Person`, `Session`).
- Defines the interfaces for Services and Repositories.

## 2. Interfaces as Contracts

Every layer (except Handlers) must expose its functionality via an **interface**. This enables mocking and decoupling.

- The Service interfaces and Repository interfaces should ideally be defined in the same package where they are consumed, or centrally in the domain/model package.

## 3. Context Passing

Every single exported function in the Service and Repository layers MUST accept `context.Context` as its first parameter.

- This ensures cancellation, timeout capabilities, and distributed tracing are easily supported.

```go
func (s *personService) GetByID(ctx context.Context, id string) (*Person, error) { ... }
```

## 4. Package Initialization (Merged Dependencies)

Instead of having dedicated `New` functions for every single object, merge related dependencies under a single struct per layer (e.g., `Repositories` for all repo interfaces, `Services` for all service interfaces).

```go
type Repositories struct {
    Person PersonRepository
    Session SessionRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
    return &Repositories{
        Person:  &personRepository{db: db},
        Session: &sessionRepository{db: db},
    }
}

type Services struct {
    Person PersonService
    // Embed other services
}

func NewServices(repos *Repositories) *Services {
    return &Services{
        Person: &personService{repo: repos.Person},
    }
}
```

## 5. Configuration & Environment Variables

- Load all configurations centrally inside a dedicated `config` package.
- Parse environment variables into a single merged `Config` struct.
- Validate the loaded configuration struct using a validator library (e.g., `github.com/go-playground/validator/v10`). The app should fail fast on startup if config is invalid.

## 6. Testing

- Write comprehensive unit tests for each layer (Handler, Service, Repository).
- Ensure both **happy paths** and **edge/error cases** are fully covered.
- Use interface mocks (e.g., via `gomock` or `testify/mock`) to isolate testing for the Handler and Service layers.

## 7. Dependencies & Tools

- **Go Version:** Always use the latest stable Go release (1.23+).
- **Web Framework:** Use `github.com/gin-gonic/gin` for HTTP routing and middleware.
- **ORM:** Use `gorm.io/gorm` along with its SQLite driver (`gorm.io/driver/sqlite`).
- **Logging:** Use `go.uber.org/zap` for highly performant structured JSON logging.
- **Validation:** Use `github.com/go-playground/validator/v10` for validating payloads and config.

## 8. Graceful Shutdown

- The `main` function must implement graceful shutdown mechanics.
- Listen for `SIGINT` and `SIGTERM`.
- Allow active HTTP connections, background jobs, and database connections time to finish and close cleanly before exiting.
