# internal/domain/ - Business Domain Layer

## Purpose
The heart of the application containing business logic, entities, and core interfaces. This layer is **framework-agnostic** and represents pure business rules.

## Structure
```
domain/
├── entities/              # Business entities (domain models)
├── repositories/          # Data access interfaces
└── services/             # Business logic interfaces
```

## Responsibilities

### 📦 entities/
- **Domain Models**: Core business objects with behavior
- **Value Objects**: Immutable objects representing concepts
- **Aggregates**: Groups of related entities
- **Business Rules**: Domain-specific validation and logic

**Example**: `User`, `Project`, `Category`, `ProjectTodolist`

### 🔄 repositories/
- **Data Access Interfaces**: Abstract data operations
- **Query Specifications**: Define complex queries
- **Repository Contracts**: Define what data operations are needed

**Example**: `UserRepository`, `ProjectRepository`

### ⚙️ services/
- **Business Logic Interfaces**: Define business operations
- **Domain Services**: Coordinate between entities
- **Business Rules**: Complex business logic that doesn't fit in entities

**Example**: `AuthService`, `ProjectService`

## Design Principles

### ✅ What Goes Here
- Business entities with behavior
- Domain-specific validation rules
- Business logic interfaces
- Domain events and specifications
- Value objects and aggregates

### ❌ What Doesn't Go Here
- Database-specific code
- HTTP request/response models
- Framework dependencies
- External API integrations
- Infrastructure concerns

## Entity Guidelines
```go
// ✅ Good - Business logic in entity
type Project struct {
    ID     int
    Name   string
    Budget decimal.Decimal
    Status ProjectStatus
}

func (p *Project) CanAddExpense(amount decimal.Decimal) error {
    if p.Status != StatusActive {
        return ErrProjectNotActive
    }
    if p.GetRemainingBudget().LessThan(amount) {
        return ErrInsufficientBudget
    }
    return nil
}

// ❌ Bad - Infrastructure concerns in domain
type Project struct {
    ID     int    `gorm:"primary_key"` // Database specific
    Name   string `json:"name"`        // HTTP specific
}
```

## Repository Interface Guidelines
```go
// ✅ Good - Domain-focused interface
type ProjectRepository interface {
    Save(ctx context.Context, project *Project) error
    FindByID(ctx context.Context, id int) (*Project, error)
    FindActiveProjects(ctx context.Context) ([]*Project, error)
}

// ❌ Bad - Implementation-specific
type ProjectRepository interface {
    SaveToDB(project *Project) error    // DB specific
    FindWithSQL(query string) []Project // SQL specific
}
```

## Testing
- Focus on **business logic** testing
- Mock external dependencies
- Test domain rules and validations
- Ensure entities maintain invariants

This layer should be the **most stable** and **well-tested** part of your application.
