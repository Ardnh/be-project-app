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
