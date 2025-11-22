# Go Backend Service - API Key Rotator

[English](README.md) | [中文简体](README_CN.md)

This project is the backend service for `API Key Rotator`, built with the **Gin** framework, providing API key management and request proxy functionalities.

## Architecture Overview

The Go backend adopts **Interface Abstraction Architecture** with **Optimized Builds**, ensuring high maintainability and deployment flexibility.

*   **Interface Abstraction**: Clean separation between business logic and infrastructure implementations using well-defined interfaces.
*   **Infrastructure Adapters**: Pluggable implementations for different databases (SQLite, MySQL) and caches (Memory, Redis).
*   **Optimized Builds**: Separate builds for lightweight and enterprise scenarios, reducing image size and dependencies.
*   **Modular Design**: Each component is independent and can be easily extended or replaced.

## Project Structure

```
backend/
├── main.go                    # Application entry point
├── go.mod                     # Go module definition
├── go.sum                     # Dependency version lock
├── Dockerfile.lightweight     # Lightweight Docker build
├── Dockerfile.enterprise      # Enterprise Docker build
├── README.md                  # Project documentation
└── internal/                  # Internal packages
    ├── config/                # Configuration management
    │   ├── config.go          # Configuration loading
    │   └── factory.go         # Infrastructure factory
    ├── infrastructure/        # Infrastructure layer (NEW)
    │   ├── database/
    │   │   ├── interface.go   # Database repository interface
    │   │   ├── sqlite/        # SQLite implementation
    │   │   └── mysql/         # MySQL implementation
    │   └── cache/
    │       ├── interface.go   # Cache interface
    │       ├── memory/        # Memory cache implementation
    │       └── redis/         # Redis implementation
    ├── adapters/              # LLM adapters (needs interface update)
    ├── handlers/              # HTTP handlers
    ├── services/              # Business services
    ├── models/                # Data models
    ├── dto/                   # Data Transfer Objects
    ├── logger/                # Logger configuration
    ├── middleware/            # Middleware
    ├── router/                # Route configuration
    └── utils/                 # Utility functions
```

## Tech Stack

*   **Framework**: [Gin](https://gin-gonic.com/) - A high-performance HTTP web framework
*   **ORM**: [GORM](https://gorm.io/) - The ORM library for Go
*   **Database**: MySQL 8.0+ (Enterprise) / SQLite (Lightweight)
*   **Cache**: Redis 6.0+ (Enterprise) / In-Memory (Lightweight)
*   **Configuration**: Environment variables + [godotenv](https://github.com/joho/godotenv)
*   **Containerization**: Docker + Docker Compose with optimized builds
*   **Architecture**: Interface Abstraction + Adapter Pattern

## Core Features

*   **Centralized Key Management**: Manage API key pools for all services in a unified web interface.
*   **Dynamic Key Rotation**: Atomic rotation based on Redis to effectively distribute API request quotas.
*   **Type-Safe Proxies**:
    *   **Generic API Proxy (`/proxy`)**: Provides proxy services for any RESTful API.
    *   **LLM API Proxy (`/llm`)**: Offers native streaming support for OpenAI-compatible large model APIs.
*   **Highly Extensible Architecture**: Uses an adapter pattern, making it easy to extend support for new types of LLM APIs in the future.
*   **Secure Isolation**: All proxy requests are authenticated via global keys, protecting real backend keys from being exposed.

## Local Development

### Prerequisites

*   Go 1.21+
*   MySQL 8.0+
*   Redis 6.0+

### Quick Start

1. **Enter the Go backend directory**
   ```bash
   cd backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment variables**
   
   Create a `.env` file in the project root (refer to `.env.example`):
   ```bash
   cp ../.env.example ../.env
   ```

4. **Run the service**
   ```bash
   go run main.go
   ```

   The service will start at `http://localhost:8000`.

### API Documentation

After starting the service, you can view the APIs as follows:

*   **Root Path**: `http://localhost:8000/` - Welcome message
*   **Management API**: `http://localhost:8000/api/admin/*` - Configuration management interfaces
*   **Generic Proxy**: `http://localhost:8000/proxy/*` - Generic API proxy
*   **LLM Proxy**: `http://localhost:8000/llm/*` - LLM API proxy

## Docker Deployment

This project supports optimized builds via the Makefile in the root directory.

### Building Images

```bash
# Build lightweight version (SQLite + Memory Cache)
make build-lightweight

# Build enterprise version (MySQL + Redis)
make build-enterprise

# Build both versions
make build-all
```

### Using Docker Compose

Run the appropriate compose file based on your needs:

```bash
# Lightweight deployment
docker-compose -f docker-compose.yml up -d

# Enterprise deployment
docker-compose -f docker-compose.prod.yml up -d
```

### Docker Image Tags

* `api-key-rotator:lightweight` - ~50MB, SQLite + Memory Cache
* `api-key-rotator:enterprise` - ~80MB, MySQL + Redis (includes all features)

## Testing

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/handlers

# Run tests and show coverage
go test -cover ./...