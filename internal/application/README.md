# internal/application/ - Application Layer

## Purpose

Contains the application's use cases and orchestrates the flow between the domain layer and external interfaces. This layer **implements business workflows**.

## Structure

```
application/
├── services/             # Use case implementations
└── dto/                 # Data Transfer Objects
```

## Responsibilities

### 🔧 services/

- **Use Case Implementation**: Execute business workflows
- **Transaction Management**: Coordinate database transactions
- **Business Process Orchestration**: Coordinate multiple domain services
- **Data Validation**: Validate input before domain operations
- **External Service Integration**: Call external APIs as part of workflows

**Key Characteristics**:

- Implements interfaces defined in `domain/services/`
- Depends on domain repositories and services
- Contains **no business rules** (delegates to domain)
- Focuses on **workflow coordination**

### 📄 dto/

- **Data Transfer Objects**: Represent data contracts with external layers
- **Request/Response Models**: HTTP API contracts
- **Mapping Functions**: Convert between DTOs and domain entities
- **Validation Rules**: Input validation annotations

## Service Implementation Pattern

```go
type userService struct {
    userRepo     repositories.UserRepository
    emailService external.EmailService
    cache        cache.CacheInterface
    logger       logger.Logger
}

func (s *userService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
    // 1. Validate input (application concern)
    if err := s.validateCreateUserRequest(req); err != nil {
        return nil, err
    }

    // 2. Check business rules (delegate to domain)
    if exists, _ := s.userRepo.ExistsByEmail(ctx, req.Email); exists {
        return nil, domain.ErrEmailAlreadyExists
    }

    // 3. Create domain entity
    user := &entities.User{
        Username: req.Username,
        Email:    req.Email,
    }

    // 4. Domain operation
    if err := user.SetPassword(req.Password); err != nil {
        return nil, err
    }

    // 5. Persist (infrastructure)
    if err := s.userRepo.Save(ctx, user); err != nil {
        return nil, err
    }

    // 6. Side effects (external services)
    go s.emailService.SendWelcomeEmail(user.Email)

    // 7. Cache update
    s.cache.Delete(ctx, fmt.Sprintf("user_list_*"))

    // 8. Return response DTO
    return dto.ToUserResponse(user), nil
}
```

## DTO Guidelines

### ✅ Good DTO Design

```go
// Request DTOs - Input validation
type CreateProjectRequest struct {
    Name        string  `json:"name" validate:"required,min=3,max=100"`
    Description string  `json:"description" validate:"max=500"`
    Budget      float64 `json:"budget" validate:"min=0"`
    CategoryID  int     `json:"category_id" validate:"required"`
}

// Response DTOs - Output formatting
type ProjectResponse struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Budget      float64   `json:"budget"`
    CreatedAt   time.Time `json:"created_at"`
}

// Conversion functions
func ToProjectResponse(project *entities.Project) *ProjectResponse {
    return &ProjectResponse{
        ID:        project.ID,
        Name:      project.Name,
        Budget:    project.Budget,
        CreatedAt: project.CreatedAt,
    }
}
```

### ❌ Bad DTO Design

```go
// Don't expose domain entities directly
type CreateProjectRequest struct {
    Project entities.Project `json:"project"` // ❌ Tight coupling
}

// Don't include infrastructure concerns
type ProjectResponse struct {
    entities.Project          // ❌ Exposing internal structure
    Password        string    // ❌ Sensitive data
}
```

## Service Guidelines

### ✅ Application Service Responsibilities

- Validate input from external sources
- Coordinate multiple domain operations
- Manage transactions across repositories
- Handle external service integrations
- Cache management
- Logging and monitoring
- DTO conversions

### ❌ What Not To Include

- Business rules (belongs in domain)
- Database queries (belongs in infrastructure)
- HTTP status codes (belongs in interface)
- Framework-specific code

## Transaction Management

```go
func (s *projectService) CreateProjectWithTasks(ctx context.Context, req *dto.CreateProjectWithTasksRequest) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Use transaction-aware repositories
        projectRepo := s.projectRepo.WithTx(tx)
        taskRepo := s.taskRepo.WithTx(tx)

        // Create project
        project := dto.ToProject(req.Project)
        if err := projectRepo.Save(ctx, project); err != nil {
            return err
        }

        // Create tasks
        for _, taskReq := range req.Tasks {
            task := dto.ToTask(taskReq)
            task.ProjectID = project.ID
            if err := taskRepo.Save(ctx, task); err != nil {
                return err
            }
        }

        return nil
    })
}
```

## Testing Strategy

- **Unit Tests**: Mock repositories and external services
- **Integration Tests**: Test with real databases
- **Contract Tests**: Ensure DTOs match API contracts

This layer ensures **clean separation** between your business logic and external concerns.
