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
