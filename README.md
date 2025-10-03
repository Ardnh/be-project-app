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
- ✅ RESTful API with Fiber
- ✅ Docker Support
- ✅ Comprehensive Testing
- ✅ API Documentation
- ✅ Database Migrations
- ✅ Logging & Monitoring

## Contributing

Please read the individual README files in each directory to understand the codebase structure and conventions.

## Visual Application Flow

┌─────────────────────────────────────────────────────────────┐
│ Client (Postman/Browser) │
└─────────────────────────────────────────────────────────────┘
↓ HTTP Request
┌─────────────────────────────────────────────────────────────┐
│ Fiber Framework (Routing + Middleware) │
│ app.Post("/api/v1/users", userHandler.CreateUserHandler) │
└─────────────────────────────────────────────────────────────┘
↓ c \*fiber.Ctx
┌─────────────────────────────────────────────────────────────┐
│ Handler Layer (interface/http/handler/) │
│ - Parse request │
│ - Validate input │
│ - Call service │
│ - Format response │
└─────────────────────────────────────────────────────────────┘
↓ DTO
┌─────────────────────────────────────────────────────────────┐
│ Service Layer (application/services/) │
│ - Business workflow │
│ - Call repository │
│ - Orchestration │
└─────────────────────────────────────────────────────────────┘
↓ Entity
┌─────────────────────────────────────────────────────────────┐
│ Repository Layer (infrastructure/repository/) │
│ - Database operations │
│ - GORM queries │
└─────────────────────────────────────────────────────────────┘
↓
┌─────────────────────────────────────────────────────────────┐
│ Database (PostgreSQL) │
└─────────────────────────────────────────────────────────────┘
