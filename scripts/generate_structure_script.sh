#!/bin/bash

# ================================
# Project Management Backend - Structure Generator
# ================================

set -e

PROJECT_NAME="project-management-backend"
CURRENT_DIR=$(pwd)

echo "🚀 Creating Go project structure with Repository Pattern..."
echo "Project: $PROJECT_NAME"
echo "Location: $CURRENT_DIR/$PROJECT_NAME"
echo ""

# Create main project directory
mkdir -p $PROJECT_NAME
cd $PROJECT_NAME

# ================================
# Create Directory Structure
# ================================

echo "📁 Creating directory structure..."

# Main directories
mkdir -p cmd/server
mkdir -p internal/{config,domain,infrastructure,application,interfaces,utils}
mkdir -p internal/domain/{entities,repositories,services}
mkdir -p internal/infrastructure/{database,cache,external}
mkdir -p internal/infrastructure/database/{mysql,migrations}
mkdir -p internal/infrastructure/cache/{redis,interfaces}
mkdir -p internal/infrastructure/external/{email,storage}
mkdir -p internal/application/{services,dto}
mkdir -p internal/interfaces/http/{handlers,middleware,routes,responses}
mkdir -p internal/interfaces/grpc/{handlers,protos}
mkdir -p internal/utils/{validator,jwt,password,database,helpers}
mkdir -p pkg/{logger,errors,constants}
mkdir -p api/{swagger,postman}
mkdir -p scripts
mkdir -p deployments/{docker,kubernetes}
mkdir -p configs
mkdir -p tests/{unit,integration,fixtures}
mkdir -p tests/unit/{services,repositories,handlers}
mkdir -p docs

echo "✅ Directory structure created!"

# ================================
# Create Documentation Files
# ================================

echo "📝 Creating documentation and README files..."

# Root README
cat > README.md << 'EOF'
# Project Management Backend

A robust project management backend built with **Go**, **Gin**, **GORM**, and **Redis** implementing **Repository Pattern** and **Clean Architecture**.

## Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Interface     │    │   Application   │    │     Domain      │
│   (HTTP/gRPC)   │───▶│   (Services)    │───▶│   (Business)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Infrastructure │    │      Utils      │    │      Config     │
│ (DB/Cache/API)  │    │   (Helpers)     │    │   (Settings)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Quick Start

1. **Setup Project**
   ```bash
   go mod init github.com/yourusername/project-management-backend
   go mod tidy
   ```

2. **Configure Environment**
   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

3. **Run with Docker**
   ```bash
   make docker-up
   ```

## Directory Structure

Each directory has its own README.md explaining its purpose and responsibilities.

## Features

- ✅ Clean Architecture with Repository Pattern
- ✅ JWT Authentication & Authorization
- ✅ Redis Caching Layer
- ✅ MySQL Database with GORM
- ✅ RESTful API with Gin
- ✅ Docker Support
- ✅ Comprehensive Testing
- ✅ API Documentation
- ✅ Database Migrations
- ✅ Logging & Monitoring

## Contributing

Please read the individual README files in each directory to understand the codebase structure and conventions.
EOF

# cmd/ directory README
cat > cmd/README.md << 'EOF'
# cmd/ - Application Entry Points

## Purpose
Contains the main applications for this project. Each subdirectory represents a different executable.

## Structure
```
cmd/
└── server/
    └── main.go          # Main HTTP server entry point
```

## Responsibilities
- **Application Bootstrap**: Initialize all dependencies and start services
- **Configuration Loading**: Load environment variables and config files
- **Dependency Injection**: Wire up all components (repositories, services, handlers)
- **Server Startup**: Start HTTP server, gRPC server, or other services
- **Graceful Shutdown**: Handle shutdown signals and cleanup resources

## Guidelines
- Keep main.go files **thin** - delegate to other packages
- Handle **dependency injection** at this level
- Implement **graceful shutdown** for production readiness
- Log **startup information** and errors appropriately
- Each cmd should be a **single responsibility** application

## Example main.go Structure
```go
func main() {
    // 1. Load configuration
    // 2. Initialize logger
    // 3. Setup database connection
    // 4. Initialize Redis
    // 5. Wire up repositories
    // 6. Wire up services  
    // 7. Setup HTTP handlers
    // 8. Start server with graceful shutdown
}
```

## Future Extensions
- `cmd/migrator/` - Database migration tool
- `cmd/seeder/` - Database seeding tool
- `cmd/worker/` - Background job processor
- `cmd/cli/` - Command line interface
EOF

# internal/ directory README
cat > internal/README.md << 'EOF'
# internal/ - Private Application Code

## Purpose
Contains the private application and library code. Code in this directory cannot be imported by external applications.

## Architecture Layers

### 🏗️ Clean Architecture Implementation
```
┌─────────────────────────────────────────────────────────┐
│                    External Interfaces                  │
│                   (HTTP, gRPC, CLI)                    │
├─────────────────────────────────────────────────────────┤
│                   Application Layer                     │
│              (Use Cases, DTOs, Services)               │
├─────────────────────────────────────────────────────────┤
│                     Domain Layer                        │
│           (Entities, Repositories, Services)           │
├─────────────────────────────────────────────────────────┤
│                 Infrastructure Layer                    │
│            (Database, Cache, External APIs)            │
└─────────────────────────────────────────────────────────┘
```

## Directory Structure

### `/domain` - Business Logic Core
- **Entities**: Core business objects
- **Repositories**: Data access interfaces
- **Services**: Business logic interfaces

### `/application` - Application Services
- **Services**: Use case implementations
- **DTOs**: Data Transfer Objects

### `/infrastructure` - External Dependencies
- **Database**: Database implementations
- **Cache**: Caching implementations
- **External**: Third-party integrations

### `/interfaces` - External Communication
- **HTTP**: REST API handlers
- **gRPC**: gRPC service implementations

### `/utils` - Helper Functions
- Common utilities and helpers
- Cross-cutting concerns

### `/config` - Configuration
- Configuration structures
- Environment setup

## Design Principles

1. **Dependency Inversion**: Depend on interfaces, not implementations
2. **Single Responsibility**: Each package has one reason to change
3. **Open/Closed**: Open for extension, closed for modification
4. **Interface Segregation**: Small, focused interfaces
5. **Dependency Injection**: Wire dependencies at startup

## Key Rules
- Domain layer **never** depends on infrastructure
- Infrastructure **implements** domain interfaces
- Application **orchestrates** domain operations
- Interfaces **adapt** external communications
EOF

# internal/domain/ directory README
cat > internal/domain/README.md << 'EOF'
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
EOF

# internal/application/ directory README
cat > internal/application/README.md << 'EOF'
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
EOF

# internal/infrastructure/ directory README
cat > internal/infrastructure/README.md << 'EOF'
# internal/infrastructure/ - Infrastructure Layer

## Purpose
Implements external dependencies and provides concrete implementations of domain interfaces. This layer handles **technical concerns** like databases, caching, external APIs, and file systems.

## Structure
```
infrastructure/
├── database/            # Database implementations
│   ├── mysql/          # MySQL repository implementations
│   └── migrations/     # Database schema migrations
├── cache/              # Caching implementations
│   ├── redis/         # Redis cache implementation
│   └── interfaces/    # Cache interface definitions
└── external/           # External service integrations
    ├── email/         # Email service implementations
    └── storage/       # File storage implementations
```

## Responsibilities

### 🗄️ database/
- **Repository Implementations**: Concrete data access implementations
- **Database Connections**: Connection pooling and management
- **Query Optimization**: Database-specific optimizations
- **Migrations**: Schema versioning and updates
- **Transactions**: Database transaction handling

### ⚡ cache/
- **Cache Implementations**: Redis, in-memory, or other caching solutions
- **Cache Strategies**: TTL, invalidation, warming strategies
- **Serialization**: Data serialization for cache storage
- **Cache Interfaces**: Abstract caching operations

### 🌐 external/
- **Third-party Integrations**: Email, SMS, payment gateways
- **API Clients**: External REST/GraphQL API clients
- **File Storage**: Cloud storage, local filesystem
- **Message Queues**: Pub/sub, message brokers

## Repository Implementation Pattern

```go
// Domain interface (in domain/repositories/)
type UserRepository interface {
    Save(ctx context.Context, user *entities.User) error
    FindByID(ctx context.Context, id int) (*entities.User, error)
    FindByEmail(ctx context.Context, email string) (*entities.User, error)
}

// Infrastructure implementation
type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Save(ctx context.Context, user *entities.User) error {
    return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) FindByID(ctx context.Context, id int) (*entities.User, error) {
    var user entities.User
    err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, domain.ErrUserNotFound
        }
        return nil, err
    }
    return &user, nil
}
```

## Cache Implementation Pattern

```go
// Cache interface
type CacheInterface interface {
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    Get(ctx context.Context, key string, dest interface{}) error
    Delete(ctx context.Context, key string) error
}

// Redis implementation
type redisCache struct {
    client *redis.Client
}

func (c *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return c.client.Set(ctx, key, data, ttl).Err()
}
```

## Database Migrations

### Migration File Structure
```sql
-- 001_create_users_table.sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_users_email (email),
    INDEX idx_users_deleted_at (deleted_at)
);
```

### Migration Management
```go
type Migrator struct {
    db *sql.DB
}

func (m *Migrator) RunMigrations() error {
    migrations := []string{
        "001_create_users_table.sql",
        "002_create_projects_table.sql",
        // ... more migrations
    }
    
    for _, migration := range migrations {
        if err := m.runMigration(migration); err != nil {
            return fmt.Errorf("failed to run migration %s: %w", migration, err)
        }
    }
    return nil
}
```

## External Service Integration

```go
// Email service interface (in domain/)
type EmailService interface {
    SendWelcomeEmail(email string) error
    SendPasswordReset(email, token string) error
}

// SMTP implementation (in infrastructure/)
type smtpEmailService struct {
    host     string
    port     int
    username string
    password string
}

func (s *smtpEmailService) SendWelcomeEmail(email string) error {
    // SMTP implementation details
    msg := s.buildWelcomeMessage(email)
    return s.sendMail(email, msg)
}
```

## Design Principles

### ✅ Good Practices
- Implement domain interfaces
- Handle infrastructure-specific errors
- Use connection pooling
- Implement proper logging
- Handle retries and timeouts
- Use transactions where appropriate

### ❌ Avoid
- Business logic in repositories
- Exposing infrastructure details to domain
- Hard-coded configurations
- Ignoring errors
- Not handling connection failures

## Configuration Management

```go
type DatabaseConfig struct {
    Host            string `yaml:"host"`
    Port            int    `yaml:"port"`
    Username        string `yaml:"username"`
    Password        string `yaml:"password"`
    Database        string `yaml:"database"`
    MaxIdleConns    int    `yaml:"max_idle_conns"`
    MaxOpenConns    int    `yaml:"max_open_conns"`
    ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

func InitDatabase(cfg *DatabaseConfig) (*gorm.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
    
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
    sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
    sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Minute)
    
    return db, nil
}
```

## Testing Infrastructure

```go
// Use test containers for integration tests
func setupTestDB(t *testing.T) *gorm.DB {
    container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: testcontainers.ContainerRequest{
            Image:        "mysql:8.0",
            ExposedPorts: []string{"3306/tcp"},
            Env: map[string]string{
                "MYSQL_ROOT_PASSWORD": "password",
                "MYSQL_DATABASE":      "testdb",
            },
        },
        Started: true,
    })
    require.NoError(t, err)
    
    // Get connection details and create DB connection
    // Return configured test database
}
```

This layer provides **reliable, performant** implementations of your domain contracts while keeping technical details isolated from business logic.
EOF

# internal/interfaces/ directory README
cat > internal/interfaces/README.md << 'EOF'
# internal/interfaces/ - External Interface Layer

## Purpose
Handles external communication and adapts between external protocols and internal application logic. This layer **translates** between the outside world and your application.

## Structure
```
interfaces/
├── http/                    # HTTP REST API
│   ├── handlers/           # HTTP request handlers
│   ├── middleware/         # HTTP middleware (auth, cors, etc.)
│   ├── routes/            # Route definitions
│   └── responses/         # Response formatting helpers
└── grpc/                   # gRPC API (optional)
    ├── handlers/          # gRPC service implementations
    └── protos/           # Protocol buffer definitions
```

## Responsibilities

### 🌐 http/
- **Request Handling**: Process HTTP requests
- **Response Formatting**: Format and return HTTP responses
- **Input Validation**: Validate incoming HTTP data
- **Authentication**: Handle auth tokens and sessions
- **Route Management**: Define API endpoints
- **Middleware**: Cross-cutting concerns (logging, CORS, rate limiting)

### ⚡ grpc/
- **Service Implementation**: Implement gRPC services
- **Protocol Buffers**: Define message contracts
- **Streaming**: Handle streaming RPCs
- **Error Handling**: gRPC-specific error responses

## HTTP Handler Pattern

```go
type UserHandler struct {
    userService services.UserService
    validator   *validator.Validator
}

func NewUserHandler(userService services.UserService, validator *validator.Validator) *UserHandler {
    return &UserHandler{
        userService: userService,
        validator:   validator,
    }
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    // 1. Bind request to DTO
    var req dto.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        responses.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
        return
    }
    
    // 2. Validate input
    if err := h.validator.Validate(&req); err != nil {
        responses.ValidationErrorResponse(c, err)
        return
    }
    
    // 3. Call application service
    user, err := h.userService.CreateUser(c.Request.Context(), &req)
    if err != nil {
        responses.HandleServiceError(c, err)
        return
    }
    
    // 4. Format response
    responses.SuccessResponse(c, http.StatusCreated, "User created successfully", user)
}
```

## Middleware Pattern

### Authentication Middleware
```go
func AuthMiddleware(jwtUtil jwt.JWTUtil) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            responses.ErrorResponse(c, http.StatusUnauthorized, "Missing authorization header", nil)
            c.Abort()
            return
        }
        
        token := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := jwtUtil.ValidateToken(token)
        if err != nil {
            responses.ErrorResponse(c, http.StatusUnauthorized, "Invalid token", nil)
            c.Abort()
            return
        }
        
        // Set user context
        c.Set("user_id", claims.UserID)
        c.Set("user_role", claims.Role)
        c.Next()
    }
}
```

### Logging Middleware
```go
func LoggingMiddleware(logger logger.Logger) gin.HandlerFunc {
    return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
        logger.Info("HTTP Request", map[string]interface{}{
            "method":      param.Method,
            "path":        param.Path,
            "status_code": param.StatusCode,
            "latency":     param.Latency,
            "client_ip":   param.ClientIP,
            "user_agent":  param.Request.UserAgent(),
        })
        return ""
    })
}
```

## Response Helper Pattern

```go
type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   interface{} `json:"error,omitempty"`
    Meta    interface{} `json:"meta,omitempty"`
}

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
    c.JSON(statusCode, Response{
        Success: true,
        Message: message,
        Data:    data,
    })
}

func ErrorResponse(c *gin.Context, statusCode int, message string, err interface{}) {
    c.JSON(statusCode, Response{
        Success: false,
        Message: message,
        Error:   err,
    })
}

// Handle service layer errors
func HandleServiceError(c *gin.Context, err error) {
    switch {
    case errors.Is(err, domain.ErrUserNotFound):
        ErrorResponse(c, http.StatusNotFound, "User not found", nil)
    case errors.Is(err, domain.ErrEmailAlreadyExists):
        ErrorResponse(c, http.StatusConflict, "Email already exists", nil)
    case errors.Is(err, domain.ErrUnauthorized):
        ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
    default:
        ErrorResponse(c, http.StatusInternalServerError, "Internal server error", nil)
    }
}
```

## Route Organization

```go
func SetupRoutes(
    router *gin.Engine,
    userHandler *handlers.UserHandler,
    projectHandler *handlers.ProjectHandler,
    authMiddleware gin.HandlerFunc,
) {
    api := router.Group("/api/v1")
    
    // Public routes
    auth := api.Group("/auth")
    {
        auth.POST("/login", userHandler.Login)
        auth.POST("/register", userHandler.Register)
        auth.POST("/refresh", userHandler.RefreshToken)
    }
    
    // Protected routes
    protected := api.Group("/")
    protected.Use(authMiddleware)
    {
        // User routes
        users := protected.Group("/users")
        {
            users.GET("/profile", userHandler.GetProfile)
            users.PUT("/profile", userHandler.UpdateProfile)
        }
        
        // Project routes
        projects := protected.Group("/projects")
        {
            projects.GET("/", projectHandler.ListProjects)
            projects.POST("/", projectHandler.CreateProject)
            projects.GET("/:id", projectHandler.GetProject)
            projects.PUT("/:id", projectHandler.UpdateProject)
            projects.DELETE("/:id", projectHandler.DeleteProject)
        }
    }
}
```

## Input Validation

```go
type CreateProjectRequest struct {
    Name        string  `json:"name" validate:"required,min=3,max=100"`
    Description string  `json:"description" validate:"max=500"`
    Budget      float64 `json:"budget" validate:"min=0"`
    CategoryID  int     `json:"category_id" validate:"required"`
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
    var req dto.CreateProjectRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        responses.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON", err)
        return
    }
    
    if err := h.validator.Validate(&req); err != nil {
        responses.ValidationErrorResponse(c, err)
        return
    }
    
    // Process request...
}
```

## Error Handling Best Practices

### ✅ Good Error Handling
```go
func (h *UserHandler) GetUser(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        responses.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", nil)
        return
    }
    
    user, err := h.userService.GetUserByID(c.Request.Context(), id)
    if err != nil {
        // Let helper handle service errors
        responses.HandleServiceError(c, err)
        return
    }
    
    responses.SuccessResponse(c, http.StatusOK, "User retrieved", user)
}
```

### ❌ Bad Error Handling
```go
func (h *UserHandler) GetUser(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id")) // ❌ Ignoring error
    
    user, err := h.userService.GetUserByID(c.Request.Context(), id)
    if err != nil {
        c.JSON(500, "Error") // ❌ Not handling different error types
        return
    }
    
    c.JSON(200, user) // ❌ Exposing internal structure
}
```

## Security Considerations

- **Input Validation**: Validate all incoming data
- **Authentication**: Verify user identity
- **Authorization**: Check user permissions
- **Rate Limiting**: Prevent abuse
- **CORS**: Configure cross-origin requests
- **Request Size Limits**: Prevent large payloads
- **Timeout Handling**: Prevent hanging requests

## Testing Handlers

```go
func TestCreateUser(t *testing.T) {
    // Setup
    mockService := mocks.NewUserService()
    handler := handlers.NewUserHandler(mockService, validator.New())
    router := gin.New()
    router.POST("/users", handler.CreateUser)
    
    // Test case
    reqBody := `{"username":"john","email":"john@example.com","password":"password123"}`
    req := httptest.NewRequest("POST", "/users", strings.NewReader(reqBody))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    
    // Execute
    router.ServeHTTP(w, req)
    
    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)
    mockService.AssertExpectations(t)
}
```

This layer ensures your application has **clean, secure, well-documented APIs** that properly handle all edge cases and provide excellent developer experience.
EOF

# internal/utils/ directory README
cat > internal/utils/README.md << 'EOF'
# internal/utils/ - Utility Functions & Helpers

## Purpose
Contains reusable utility functions, helpers, and cross-cutting concerns that are used across multiple layers of the application.

## Structure
```
utils/
├── validator/             # Input validation utilities
├── jwt/                  # JWT token handling
├── password/             # Password hashing utilities
├── database/             # Database helper functions
└── helpers/              # General purpose helpers
```

## Responsibilities

### ✅ validator/
- **Input Validation**: Custom validation rules and functions
- **Data Sanitization**: Clean and normalize input data
- **Validation Messages**: Custom error messages
- **Struct Validation**: Validate struct fields using tags

### 🔐 jwt/
- **Token Generation**: Create JWT tokens
- **Token Validation**: Verify and parse JWT tokens
- **Claims Management**: Handle custom token claims
- **Token Refresh**: Refresh token logic

### 🔒 password/
- **Password Hashing**: Secure password hashing (bcrypt)
- **Password Verification**: Verify passwords against hashes
- **Password Strength**: Check password complexity
- **Salt Generation**: Generate secure salts

### 🗄️ database/
- **Pagination**: Database pagination helpers
- **Transaction Management**: Transaction utilities
- **Query Builders**: Dynamic query construction
- **Connection Helpers**: Database connection utilities

### 🛠️ helpers/
- **String Utilities**: String manipulation functions
- **Time Utilities**: Date/time formatting and parsing
- **File Utilities**: File handling and manipulation
- **Conversion Helpers**: Data type conversions

## Implementation Examples

### Password Utility
```go
package password

import (
    "golang.org/x/crypto/bcrypt"
)

type PasswordUtil interface {
    HashPassword(password string) (string, error)
    VerifyPassword(password, hash string) error
    GenerateRandomPassword(length int) string
}

type bcryptUtil struct {
    cost int
}

func NewBcryptUtil() PasswordUtil {
    return &bcryptUtil{cost: bcrypt.DefaultCost}
}

func (b *bcryptUtil) HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
    return string(bytes), err
}

func (b *bcryptUtil) VerifyPassword(password, hash string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
```

### JWT Utility
```go
package jwt

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type JWTUtil interface {
    GenerateToken(userID int, username, role string) (string, error)
    ValidateToken(tokenString string) (*Claims, error)
    RefreshToken(tokenString string) (string, error)
}

type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

type jwtUtil struct {
    secret     []byte
    expiration time.Duration
}

func NewJWTUtil(secret string, expirationHours int) JWTUtil {
    return &jwtUtil{
        secret:     []byte(secret),
        expiration: time.Duration(expirationHours) * time.Hour,
    }
}

func (j *jwtUtil) GenerateToken(userID int, username, role string) (string, error) {
    claims := &Claims{
        UserID:   userID,
        Username: username,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(j.secret)
}
```

### Validator Utility
```go
package validator

import (
    "github.com/go-playground/validator/v10"
)

type Validator struct {
    validate *validator.Validate
}

func NewValidator() *Validator {
    v := validator.New()
    
    // Register custom validators
    v.RegisterValidation("password_strength", validatePasswordStrength)
    
    return &Validator{validate: v}
}

func (v *Validator) Validate(s interface{}) error {
    return v.validate.Struct(s)
}

func validatePasswordStrength(fl validator.FieldLevel) bool {
    password := fl.Field().String()
    
    // Check minimum length
    if len(password) < 8 {
        return false
    }
    
    // Check for at least one uppercase, lowercase, and digit
    hasUpper := false
    hasLower := false
    hasDigit := false
    
    for _, char := range password {
        switch {
        case 'A' <= char && char <= 'Z':
            hasUpper = true
        case 'a' <= char && char <= 'z':
            hasLower = true
        case '0' <= char && char <= '9':
            hasDigit = true
        }
    }
    
    return hasUpper && hasLower && hasDigit
}
```

### Database Pagination Helper
```go
package database

import (
    "gorm.io/gorm"
)

type PaginationParams struct {
    Page  int `json:"page" validate:"min=1"`
    Limit int `json:"limit" validate:"min=1,max=100"`
}

type PaginationResult struct {
    CurrentPage int   `json:"current_page"`
    LastPage    int   `json:"last_page"`
    PerPage     int   `json:"per_page"`
    Total       int64 `json:"total"`
}

func Paginate(db *gorm.DB, params PaginationParams) (*gorm.DB, *PaginationResult) {
    var total int64
    db.Count(&total)
    
    offset := (params.Page - 1) * params.Limit
    lastPage := int((total + int64(params.Limit) - 1) / int64(params.Limit))
    
    result := &PaginationResult{
        CurrentPage: params.Page,
        LastPage:    lastPage,
        PerPage:     params.Limit,
        Total:       total,
    }
    
    return db.Offset(offset).Limit(params.Limit), result
}
```

### String Helper
```go
package helpers

import (
    "math/rand"
    "strings"
    "time"
)

func ToCamelCase(s string) string {
    words := strings.Fields(s)
    for i := 1; i < len(words); i++ {
        words[i] = strings.Title(words[i])
    }
    return strings.Join(words, "")
}

func ToSnakeCase(s string) string {
    var result strings.Builder
    for i, char := range s {
        if i > 0 && char >= 'A' && char <= 'Z' {
            result.WriteRune('_')
        }
        result.WriteRune(char)
    }
    return strings.ToLower(result.String())
}

func GenerateRandomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    
    rand.Seed(time.Now().UnixNano())
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[rand.Intn(len(charset))]
    }
    return string(b)
}

func Slugify(s string) string {
    s = strings.ToLower(s)
    s = strings.ReplaceAll(s, " ", "-")
    // Remove special characters
    // Add more sophisticated logic as needed
    return s
}
```

## Design Guidelines

### ✅ Good Utilities
- **Pure Functions**: No side effects
- **Reusable**: Used in multiple places
- **Well Tested**: Comprehensive unit tests
- **Single Responsibility**: Each function does one thing
- **Error Handling**: Proper error handling and validation

### ❌ Avoid in Utils
- Business logic (belongs in domain/services)
- Database-specific operations (belongs in infrastructure)
- HTTP-specific code (belongs in interfaces)
- Configuration management (belongs in config)

## Testing Utilities

```go
package validator

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestPasswordStrengthValidation(t *testing.T) {
    v := NewValidator()
    
    type testStruct struct {
        Password string `validate:"password_strength"`
    }
    
    tests := []struct {
        name      string
        password  string
        shouldErr bool
    }{
        {"Valid password", "Password123", false},
        {"Too short", "Pass1", true},
        {"No uppercase", "password123", true},
        {"No lowercase", "PASSWORD123", true},
        {"No digit", "Password", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ts := testStruct{Password: tt.password}
            err := v.Validate(ts)
            
            if tt.shouldErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## Usage Examples

```go
// In service layer
func (s *userService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) error {
    // Use password utility
    hashedPassword, err := s.passwordUtil.HashPassword(req.Password)
    if err != nil {
        return err
    }
    
    // Use helper function
    username := helpers.Slugify(req.Username)
    
    // Create user...
    return nil
}

// In handler layer
func (h *authHandler) Login(c *gin.Context) {
    // Validate input
    if err := h.validator.Validate(&req); err != nil {
        // Handle validation error
        return
    }
    
    // Generate JWT token
    token, err := h.jwtUtil.GenerateToken(user.ID, user.Username, user.Role)
    if err != nil {
        // Handle error
        return
    }
    
    // Return token...
}
```

This layer provides **reliable, reusable utilities** that keep your other layers clean and focused on their primary responsibilities.
EOF

# pkg/ directory README
cat > pkg/README.md << 'EOF'
# pkg/ - Public Packages

## Purpose
Contains library code that can be imported by external applications. These packages should be **stable, well-documented, and backward-compatible**.

## Structure
```
pkg/
├── logger/               # Logging utilities
├── errors/              # Custom error types
└── constants/           # Application constants
```

## Responsibilities

### 📝 logger/
- **Structured Logging**: JSON/structured log output
- **Log Levels**: Support for different log levels (debug, info, warn, error)
- **Context Logging**: Contextual information in logs
- **Multiple Outputs**: Console, file, remote logging services
- **Performance**: High-performance logging suitable for production

### ❌ errors/
- **Custom Error Types**: Application-specific error definitions
- **Error Codes**: Standardized error codes for APIs
- **Error Wrapping**: Enhanced error context and stack traces
- **Error Classification**: Business vs technical errors

### 🔧 constants/
- **Application Constants**: App-wide constants
- **Configuration Keys**: Environment variable keys
- **Cache Keys**: Redis cache key patterns
- **API Constants**: HTTP status codes, headers, etc.

## Logger Implementation

```go
package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

type Logger interface {
    Debug(msg string, fields map[string]interface{})
    Info(msg string, fields map[string]interface{})
    Warn(msg string, fields map[string]interface{})
    Error(msg string, fields map[string]interface{})
    Fatal(msg string, fields map[string]interface{})
    With(fields map[string]interface{}) Logger
}

type zapLogger struct {
    logger *zap.Logger
}

func NewZapLogger(environment string) Logger {
    var config zap.Config
    
    if environment == "production" {
        config = zap.NewProductionConfig()
    } else {
        config = zap.NewDevelopmentConfig()
        config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
    }
    
    logger, _ := config.Build()
    return &zapLogger{logger: logger}
}

func (l *zapLogger) Info(msg string, fields map[string]interface{}) {
    l.logger.Info(msg, l.mapToZapFields(fields)...)
}

func (l *zapLogger) mapToZapFields(fields map[string]interface{}) []zap.Field {
    zapFields := make([]zap.Field, 0, len(fields))
    for key, value := range fields {
        zapFields = append(zapFields, zap.Any(key, value))
    }
    return zapFields
}

// Usage in application
func (s *userService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) error {
    s.logger.Info("Creating user", map[string]interface{}{
        "username": req.Username,
        "email":    req.Email,
    })
    
    // ... service logic
    
    s.logger.Info("User created successfully", map[string]interface{}{
        "user_id": user.ID,
        "username": user.Username,
    })
    
    return nil
}
```

## Custom Errors

```go
package errors

import (
    "fmt"
    "net/http"
)

// Error codes for API responses
const (
    // User errors
    CodeUserNotFound         = "USER_NOT_FOUND"
    CodeEmailAlreadyExists   = "EMAIL_ALREADY_EXISTS"
    CodeUsernameExists       = "USERNAME_ALREADY_EXISTS"
    CodeInvalidCredentials   = "INVALID_CREDENTIALS"
    
    // Project errors
    CodeProjectNotFound      = "PROJECT_NOT_FOUND"
    CodeProjectAccessDenied  = "PROJECT_ACCESS_DENIED"
    
    // Auth errors
    CodeUnauthorized         = "UNAUTHORIZED"
    CodeForbidden           = "FORBIDDEN"
    CodeInvalidToken        = "INVALID_TOKEN"
    
    // General errors
    CodeValidationFailed    = "VALIDATION_FAILED"
    CodeInternalError       = "INTERNAL_ERROR"
)

// AppError represents an application error with context
type AppError struct {
    Code       string                 `json:"code"`
    Message    string                 `json:"message"`
    Details    map[string]interface{} `json:"details,omitempty"`
    StatusCode int                    `json:"-"`
    Err        error                  `json:"-"`
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
    }
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
    return e.Err
}

// Error constructors
func NewUserNotFoundError() *AppError {
    return &AppError{
        Code:       CodeUserNotFound,
        Message:    "User not found",
        StatusCode: http.StatusNotFound,
    }
}

func NewValidationError(details map[string]interface{}) *AppError {
    return &AppError{
        Code:       CodeValidationFailed,
        Message:    "Validation failed",
        Details:    details,
        StatusCode: http.StatusBadRequest,
    }
}

func NewInternalError(err error) *AppError {
    return &AppError{
        Code:       CodeInternalError,
        Message:    "Internal server error",
        StatusCode: http.StatusInternalServerError,
        Err:        err,
    }
}

// Error checking helpers
func IsUserNotFound(err error) bool {
    var appErr *AppError
    return errors.As(err, &appErr) && appErr.Code == CodeUserNotFound
}

func IsValidationError(err error) bool {
    var appErr *AppError
    return errors.As(err, &appErr) && appErr.Code == CodeValidationFailed
}
```

## Constants

```go
package constants

// Application constants
const (
    AppName        = "Project Management API"
    AppVersion     = "1.0.0"
    DefaultTimeout = 30 // seconds
)

// Environment constants
const (
    EnvDevelopment = "development"
    EnvProduction  = "production"
    EnvTest        = "test"
)

// Cache key patterns
const (
    CacheKeyUser        = "user:%d"
    CacheKeyUserProfile = "user_profile:%d"
    CacheKeyProject     = "project:%d"
    CacheKeyUserProjects = "user_projects:%d:page:%d"
    CacheKeyCategories  = "categories"
)

// Cache TTL (in seconds)
const (
    CacheTTLShort  = 300   // 5 minutes
    CacheTTLMedium = 1800  // 30 minutes
    CacheTTLLong   = 3600  // 1 hour
    CacheTTLDay    = 86400 // 24 hours
)

// HTTP headers
const (
    HeaderAuthorization = "Authorization"
    HeaderContentType   = "Content-Type"
    HeaderUserAgent     = "User-Agent"
    HeaderXRequestID    = "X-Request-ID"
)

// Content types
const (
    ContentTypeJSON = "application/json"
    ContentTypeXML  = "application/xml"
    ContentTypeForm = "application/x-www-form-urlencoded"
)

// Database constants
const (
    DefaultPageSize = 20
    MaxPageSize     = 100
    DefaultPage     = 1
)

// Project status constants
const (
    ProjectStatusPlanning  = "planning"
    ProjectStatusActive    = "active"
    ProjectStatusOnHold    = "on_hold"
    ProjectStatusCompleted = "completed"
    ProjectStatusCancelled = "cancelled"
)

// User roles
const (
    RoleAdmin = "admin"
    RoleUser  = "user"
)

// JWT constants
const (
    JWTTokenType = "Bearer"
    JWTClaimUserID = "user_id"
    JWTClaimRole = "role"
    JWTClaimUsername = "username"
)
```

## Usage Examples

### Using Logger
```go
// In service layer
func (s *userService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) error {
    logger := s.logger.With(map[string]interface{}{
        "operation": "CreateUser",
        "username":  req.Username,
    })
    
    logger.Info("Starting user creation")
    
    // ... business logic
    
    if err != nil {
        logger.Error("Failed to create user", map[string]interface{}{
            "error": err.Error(),
        })
        return err
    }
    
    logger.Info("User created successfully", map[string]interface{}{
        "user_id": user.ID,
    })
    
    return nil
}
```

### Using Custom Errors
```go
// In repository layer
func (r *userRepository) FindByID(ctx context.Context, id int) (*entities.User, error) {
    var user entities.User
    err := r.db.WithContext(ctx).First(&user, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, customerrors.NewUserNotFoundError()
        }
        return nil, customerrors.NewInternalError(err)
    }
    return &user, nil
}

// In handler layer
func (h *userHandler) GetUser(c *gin.Context) {
    user, err := h.userService.GetUserByID(ctx, id)
    if err != nil {
        var appErr *customerrors.AppError
        if errors.As(err, &appErr) {
            c.JSON(appErr.StatusCode, gin.H{
                "error": appErr,
            })
            return
        }
        // Handle unknown errors
        c.JSON(500, gin.H{"error": "Internal server error"})
        return
    }
    
    c.JSON(200, user)
}
```

### Using Constants
```go
// In cache service
func (c *cacheService) GetUser(ctx context.Context, userID int) (*entities.User, error) {
    key := fmt.Sprintf(constants.CacheKeyUser, userID)
    
    var user entities.User
    err := c.redis.Get(ctx, key, &user)
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

func (c *cacheService) SetUser(ctx context.Context, user *entities.User) error {
    key := fmt.Sprintf(constants.CacheKeyUser, user.ID)
    return c.redis.Set(ctx, key, user, constants.CacheTTLMedium)
}
```

## Design Guidelines

### ✅ Good Practices for pkg/
- **Stable APIs**: Don't break backward compatibility
- **Comprehensive Documentation**: Include examples and usage
- **Minimal Dependencies**: Reduce external dependencies
- **Interface-Based**: Define interfaces for extensibility
- **Well Tested**: Comprehensive test coverage

### ❌ Avoid in pkg/
- Internal application logic
- Framework-specific code
- Business rules
- Configuration management
- Database models

## Versioning
Since pkg/ contains public APIs, consider semantic versioning:
- **Major version**: Breaking changes
- **Minor version**: New features (backward compatible)
- **Patch version**: Bug fixes

This layer provides **stable, reusable components** that can be shared across multiple projects or used by external consumers.
EOF

# Create other necessary README files for remaining directories

# api/ directory README
cat > api/README.md << 'EOF'
# api/ - API Documentation

## Purpose
Contains API documentation, specifications, and client tools for the project.

## Structure
```
api/
├── swagger/              # OpenAPI/Swagger documentation
└── postman/             # Postman collections
```

## Contents

### swagger/
- OpenAPI 3.0 specifications
- Auto-generated documentation from code annotations
- Interactive API documentation
- Schema definitions

### postman/
- Postman collection files
- Environment configurations
- Pre-request scripts and tests
- API testing scenarios

## Usage

### Generate Swagger Docs
```bash
make swagger
```

### View Documentation
```bash
# Serve swagger UI locally
swagger-ui-serve api/swagger/swagger.yaml
```

### Import Postman Collection
1. Open Postman
2. Import `api/postman/project_management.postman_collection.json`
3. Set up environment variables
4. Start testing endpoints
EOF

# scripts/ directory README
cat > scripts/README.md << 'EOF'
# scripts/ - Build and Deployment Scripts

## Purpose
Contains automation scripts for building, testing, and deploying the application.

## Scripts

### build.sh
- Build the Go application
- Set build information (version, commit, build time)
- Create binaries for different platforms

### migrate.sh  
- Run database migrations
- Wait for database connectivity
- Handle migration rollbacks

### deploy.sh
- Deploy to different environments
- Build and push Docker images
- Update Kubernetes deployments

### seed.sh
- Seed database with initial data
- Create default users and categories
- Set up development data

## Usage
```bash
# Make scripts executable
chmod +x scripts/*.sh

# Run specific script
./scripts/build.sh
./scripts/migrate.sh
```
EOF

# deployments/ directory README
cat > deployments/README.md << 'EOF'
# deployments/ - Deployment Configurations

## Purpose
Contains containerization and orchestration configurations for different deployment environments.

## Structure
```
deployments/
├── docker/              # Docker configurations
└── kubernetes/          # Kubernetes manifests
```

## Docker
- **Dockerfile**: Multi-stage build for Go application
- **docker-compose.yml**: Development environment setup
- **docker-compose.prod.yml**: Production environment setup

## Kubernetes
- **namespace.yaml**: Kubernetes namespace
- **deployment.yaml**: Application deployment
- **service.yaml**: Service definitions
- **configmap.yaml**: Configuration management

## Usage

### Docker
```bash
# Development
docker-compose up -d

# Production
docker-compose -f deployments/docker/docker-compose.prod.yml up -d
```

### Kubernetes
```bash
kubectl apply -f deployments/kubernetes/
```
EOF

# configs/ directory README
cat > configs/README.md << 'EOF'
# configs/ - Configuration Files

## Purpose
Contains configuration files for different environments.

## Files
- **config.yaml**: Default development configuration
- **config.prod.yaml**: Production configuration
- **config.test.yaml**: Test environment configuration

## Environment Variables
Configuration can be overridden using environment variables:
- `DB_HOST`, `DB_PORT`, `DB_USERNAME`, `DB_PASSWORD`
- `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`
- `JWT_SECRET`

## Usage
Application automatically loads the appropriate config based on `APP_ENV` environment variable.
EOF

# tests/ directory README  
cat > tests/README.md << 'EOF'
# tests/ - Test Files

## Purpose
Contains all test files organized by test type and layer.

## Structure
```
tests/
├── unit/                # Unit tests
├── integration/         # Integration tests
└── fixtures/           # Test data fixtures
```

## Test Types

### Unit Tests
- Test individual functions/methods in isolation
- Mock external dependencies
- Fast execution
- High code coverage

### Integration Tests
- Test component interactions
- Use real databases (test containers)
- Test API endpoints end-to-end
- Slower but more comprehensive

### Fixtures
- Sample data for tests
- JSON/YAML test data files
- Reusable test scenarios

## Running Tests
```bash
# All tests
make test

# Unit tests only
go test -short ./...

# Integration tests
go test -run Integration ./...

# With coverage
make test-coverage
```
EOF

# docs/ directory README
cat > docs/README.md << 'EOF'
# docs/ - Project Documentation

## Purpose
Contains comprehensive project documentation.

## Documents
- **README.md**: Project overview and quick start
- **API.md**: Detailed API documentation  
- **DEPLOYMENT.md**: Deployment guides
- **ARCHITECTURE.md**: System architecture documentation
- **CONTRIBUTING.md**: Contribution guidelines

## Usage
Documentation is written in Markdown and should be kept up-to-date with code changes.
EOF

# Create placeholder files with basic content

# Create main.go
cat > cmd/server/main.go << 'EOF'
package main

import (
	"log"
	"os"
)

func main() {
	log.Println("🚀 Starting Project Management Backend...")
	log.Println("📝 TODO: Implement application bootstrap")
	log.Println("💡 See cmd/README.md for implementation guidance")
	
	// TODO: Implement application initialization
	// 1. Load configuration
	// 2. Initialize database
	// 3. Setup Redis
	// 4. Wire up dependencies
	// 5. Start HTTP server
	
	os.Exit(0)
}
EOF

# Create go.mod template
cat > go.mod << 'EOF'
module github.com/yourusername/project-management-backend

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/go-redis/redis/v8 v8.11.5
	github.com/go-playground/validator/v10 v10.15.5
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/spf13/viper v1.17.0
	golang.org/x/crypto v0.14.0
	gorm.io/driver/mysql v1.5.2
	gorm.io/gorm v1.25.5
	go.uber.org/zap v1.26.0
)
EOF

# Create basic config file
cat > configs/config.yaml << 'EOF'
app:
  name: "Project Management Backend"
  version: "1.0.0" 
  environment: "development"

server:
  host: "localhost"
  port: "8080"
  read_timeout: 30
  write_timeout: 30

database:
  driver: "mysql"
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"
  database: "project_management"
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 60

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

jwt:
  secret: "your-super-secret-jwt-key-here"
  expiration: 24
EOF

# Create .env.example
cat > .env.example << 'EOF'
# Application
APP_ENV=development
APP_NAME="Project Management Backend"

# Database
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=password
DB_NAME=project_management

# Redis  
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=your-super-secret-jwt-key-here

# Server
SERVER_PORT=8080
EOF

# Create Makefile
cat > Makefile << 'EOF'
.PHONY: build run test clean deps

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application  
run:
	go run cmd/server/main.go

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Generate swagger docs (requires swag)
swagger:
	swag init -g cmd/server/main.go -o ./api/swagger

# Build and run with docker-compose
docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down

# Database operations
migrate:
	./scripts/migrate.sh

seed:
	./scripts/seed.sh
EOF

# Create .gitignore
cat > .gitignore << 'EOF'
# Binaries
*.exe
*.exe~  
*.dll
*.so
*.dylib

# Test binary
*.test

# Output of the go coverage tool
*.out

# Go workspace file
go.work

# Environment variables
.env
.env.local

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Build artifacts
bin/
dist/

# Logs
*.log

# Coverage
coverage.html

# Dependencies
vendor/

# Air temp files
tmp/
EOF

# Create basic placeholder files for key directories
touch internal/domain/entities/.keep
touch internal/domain/repositories/.keep
touch internal/domain/services/.keep
touch internal/application/services/.keep
touch internal/application/dto/.keep
touch internal/infrastructure/database/mysql/.keep
touch internal/infrastructure/cache/redis/.keep
touch internal/interfaces/http/handlers/.keep
touch internal/interfaces/http/middleware/.keep
touch internal/utils/validator/.keep
touch pkg/logger/.keep
touch pkg/errors/.keep
touch api/swagger/.keep
touch tests/unit/.keep
touch tests/integration/.keep

# Make scripts executable
chmod +x scripts/*.sh 2>/dev/null || true

echo ""
echo "✅ Project structure created successfully!"
echo ""
echo "📁 Directory structure:"
echo "   - $(find . -type d | wc -l) directories created"
echo "   - $(find . -name "README.md" | wc -l) README files with detailed explanations"
echo "   - Essential configuration files and templates"
echo ""
echo "🚀 Next steps:"
echo "1. cd $PROJECT_NAME"
echo "2. go mod tidy"
echo "3. Read the README.md files in each directory"
echo "4. Customize configs/config.yaml"
echo "5. cp .env.example .env && edit .env"
echo "6. Start implementing based on the guidelines in each README.md"
echo ""
echo "📖 Architecture Overview:"
echo "   Domain → Application → Infrastructure → Interface"
echo "   Clean Architecture with Repository Pattern implemented!"
echo ""
echo "🛠️  Available Make commands:"
echo "   make deps     # Install dependencies"
echo "   make build    # Build application"
echo "   make run      # Run application"
echo "   make test     # Run tests"
echo "   make fmt      # Format code"
echo ""
echo "Happy coding! 🎉"
