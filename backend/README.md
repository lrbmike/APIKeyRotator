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
├── build.sh                   # Linux build script
├── build.bat                  # Windows build script
├── README.md                  # Project documentation
└── internal/                  # Internal packages
    ├── config/                # Configuration management
    │   └── config.go
    ├── database/              # Database connection and migration
    │   └── database.go
    ├── cache/                 # In-memory cache
    │   └── cache.go
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
        ├── gemini_adapter.go
        └── anthropic_adapter.go
```

## Tech Stack

*   **Framework**: [Gin](https://gin-gonic.com/) - A high-performance HTTP web framework
*   **ORM**: [GORM](https://gorm.io/) - The ORM library for Go
*   **Database**: SQLite 3 (using [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3))
*   **Cache**: In-memory cache (self-implemented, thread-safe)
*   **Configuration**: Environment variables + [godotenv](https://github.com/joho/godotenv)
*   **Containerization**: Docker + Docker Compose

## Core Features

*   **Centralized Key Management**: Manage API key pools for all services in a unified web interface
*   **Dynamic Key Rotation**: Atomic rotation based on in-memory cache to effectively distribute API request quotas
*   **Type-Safe Proxies**:
    *   **Generic API Proxy (`/proxy`)**: Provides proxy services for any RESTful API
    *   **LLM API Proxy (`/llm`)**: Offers native streaming support for OpenAI-compatible large model APIs
*   **Highly Extensible Architecture**: Uses an adapter pattern, making it easy to extend support for new types of LLM APIs in the future
*   **Secure Isolation**: All proxy requests are authenticated via global keys, protecting real backend keys from being exposed
*   **Lightweight Deployment**: Single executable file + SQLite database file, no additional services required

## Local Development

### Prerequisites

*   Go 1.21+
*   GCC Compiler (SQLite requires CGO support)
    *   **Windows**: Install [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) or [MinGW-w64](https://www.mingw-w64.org/)
    *   **Linux**: Usually pre-installed. If not, run `sudo apt-get install build-essential` (Ubuntu/Debian)
    *   **Verification**: Run `gcc --version` to confirm installation

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
   
   Main configuration items:
   ```env
   DATABASE_PATH=./data/api_key_rotator.db
   BACKEND_PORT=8000
   GLOBAL_PROXY_KEYS=your_secret_key
   ADMIN_USER=admin
   ADMIN_PASSWORD=your_password
   ```

4. **Build the project**

   **Build Method**
   ```bash
   # Windows (PowerShell)
   $env:CGO_ENABLED=1
   go build -o api-key-rotator.exe .
   
   # Linux/macOS
   CGO_ENABLED=1 go build -o api-key-rotator .
   ```

5. **Run the service**
   ```bash
   # Windows
   .\api-key-rotator.exe
   
   # Linux/macOS
   ./api-key-rotator
   
   # Or run directly (development mode)
   go run main.go
   ```

   The service will start at `http://localhost:8000`

### API Documentation

After starting the service, you can view the APIs as follows:

*   **Root Path**: `http://localhost:8000/` - Welcome message
*   **Management API**: `http://localhost:8000/api/admin/*` - Configuration management interfaces
*   **Generic Proxy**: `http://localhost:8000/proxy/*` - Generic API proxy
*   **LLM Proxy**: `http://localhost:8000/llm/*` - LLM API proxy

## Docker Deployment

We offer multiple ways to deploy using Docker.

### Method 1: Single Image Deployment (Recommended)

This service is part of a larger project that can be deployed as a single Docker image. This is the simplest and recommended method, especially for PaaS platforms like Render.

To do this, use the `Dockerfile` in the project's root directory.

```bash
# In the project root directory
docker build -t api-key-rotator .
```

For detailed instructions, please see the **[Docker Deployment Guide](../DEPLOY_WITH_DOCKER.md)** in the root directory.

### Method 2: Using Docker Compose (For backend-only or multi-container setups)

If you wish to run the backend service separately, you can use `docker-compose`.

**Building the backend image:**
```bash
# Run in the project root directory
docker-compose build backend
```

**Using Docker Compose:**
```bash
# Start all services
docker-compose up -d

# View logs for the backend
docker-compose logs -f backend

# Stop services
docker-compose down

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/handlers

# Run tests and show coverage
go test -cover ./...
```

## FAQ

### Q: Getting "cgo: C compiler not found" during compilation
**A:** You need to install a GCC compiler, refer to the "Prerequisites" section above

### Q: Cannot find gcc on Windows
**A:** 
1. Install TDM-GCC or MinGW-w64
2. Ensure gcc.exe is in your system PATH
3. Restart your command line window

### Q: Database permission errors
**A:** Ensure the `data` directory exists and has write permissions
```bash
mkdir -p data
chmod 755 data
```

## Related Documentation

*   [Quick Start Guide](../QUICKSTART.md) - Detailed quick start steps
*   [Deployment Guide](../DEPLOYMENT.md) - Complete deployment instructions
*   [Technical Decisions](../TECHNICAL_DECISIONS.md) - Technical decision explanations
*   [Migration Guide](../MIGRATION_SQLITE.md) - Guide for migrating from MySQL+Redis

## Performance Characteristics

*   **Lightweight & Efficient**: Single executable file with low memory footprint
*   **Fast Startup**: No need to wait for external services to start
*   **Suitable Scenarios**: Small to medium-scale deployments (< 10000 QPS)
*   **Easy Backup**: Only need to backup the SQLite database file

For higher performance or distributed deployment, consider switching back to MySQL + Redis.
