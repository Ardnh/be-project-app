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
