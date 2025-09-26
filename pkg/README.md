# pkg/ - Public Packages

## Purpose

Contains library code that can be imported by external applications. These packages should be **stable, well-documented, and backward-compatible**.

## Structure

```
pkg/
├── logger/              # Logging utilities
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
