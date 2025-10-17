# Go Backend Service - API Key Rotator

[English](README.md) | [中文简体](README_CN.md)

This project is the backend service for `API Key Rotator`, built with the **Gin** framework, providing API key management and request proxy functionalities.

## Architecture Overview

The Go backend adopts the **Clean Architecture** design pattern, ensuring high standardization and extensibility.

*   **Data Models**: Defined in `internal/models/models.go`, using a unified `proxy_configs` table to manage all types of proxy configurations. The structure is clear and easy to extend.
*   **Business Logic**: The core proxy logic is abstracted into `internal/services/` and `internal/adapters/`, achieving high code reuse and logical decoupling.
*   **API Routing**: Defined in the `internal/handlers/` directory, where each file is responsible for a functional module (management, generic proxy, LLM proxy), ensuring clear responsibilities.

## Project Structure

```
backend/
├── main.go                    # Application entry point
├── go.mod                     # Go module definition
├── go.sum                     # Dependency version lock
├── Dockerfile                 # Docker build file
├── README.md                  # Project documentation
└── internal/                  # Internal packages
    ├── config/                # Configuration management
    │   └── config.go
    ├── database/              # Database connection and migration
    │   └── database.go
    ├── redis/                 # Redis connection
    │   └── redis.go
    ├── logger/                # Logger configuration
    │   └── logger.go
    ├── models/                # Data models
    │   └── models.go
    ├── dto/                   # Data Transfer Objects
    │   └── dto.go
    ├── utils/                 # Utility functions
    │   └── utils.go
    ├── middleware/            # Middleware
    │   └── cors.go
    ├── router/                # Route configuration
    │   └── router.go
    ├── handlers/              # HTTP handlers
    │   ├── management.go      # Management API
    │   ├── proxy.go          # Generic proxy
    │   └── llm_proxy.go      # LLM proxy
    ├── services/              # Business services
    │   └── proxy_handler.go
    └── adapters/              # LLM adapters
        ├── base_adapter.go
        ├── openai_adapter.go
        └── gemini_adapter.go
```

## Tech Stack

*   **Framework**: [Gin](https://gin-gonic.com/) - A high-performance HTTP web framework
*   **ORM**: [GORM](https://gorm.io/) - The ORM library for Go
*   **Database**: MySQL 8.0+
*   **Cache**: Redis 6.0+
*   **Configuration**: Environment variables + [godotenv](https://github.com/joho/godotenv)
*   **Containerization**: Docker + Docker Compose

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

This project is deployed via the `docker-compose.yml` file in the root directory, which includes this Go backend as a default service.

### Building the Image

```bash
# Run in the project root directory
docker-compose build backend
```

### Using Docker Compose

Run `docker-compose up` in the project root to start all services.

## Testing

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/handlers

# Run tests and show coverage
go test -cover ./...