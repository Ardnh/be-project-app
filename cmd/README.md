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
