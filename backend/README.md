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
├── README.md                  # Project documentation
└── internal/                  # Internal packages
    ├── config/                # Configuration management
    │   ├── config.go          # Configuration loading
    │   └── factory.go         # Infrastructure factory
    ├── infrastructure/        # Infrastructure layer
    │   ├── database/
    │   │   ├── interface.go   # Database repository interface
    │   │   ├── sqlite/        # SQLite implementation
    │   │   └── mysql/         # MySQL implementation
    │   └── cache/
    │       ├── interface.go   # Cache interface
    │       ├── memory/        # Memory cache implementation
    │       └── redis/         # Redis implementation
    ├── adapters/              # LLM adapters
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

This project supports pure Docker builds.

### Building Images

```bash
# Lightweight version (SQLite + Memory Cache)
docker build -t api-key-rotator .

# Enterprise version (MySQL + Redis)
docker build -f Dockerfile.enterprise -t api-key-rotator:enterprise .
```

### Using Docker Compose

Run the appropriate compose file based on your needs:

#### Quick Deployment (Recommended for Beginners)
If you prefer the simplest approach, you can switch to the `sqlite` branch directly:
```bash
git checkout sqlite
docker-compose up -d
```
The `sqlite` branch is a pure SQLite + memory cache version with simpler configuration, ideal for quick testing.

#### Current Branch Deployment
```bash
# Lightweight version deployment
docker-compose -f docker-compose.yml up -d

# Enterprise version deployment
docker-compose -f docker-compose.enterprise.yml up -d
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
```

## Secondary Development Extension

### PostgreSQL Integration Example

This project uses **Interface Abstraction Architecture** to easily extend support for new database types. Here's a concise example for integrating PostgreSQL:

#### 1. Create PostgreSQL Implementation

**Directory Structure:**
```
internal/infrastructure/database/
├── interface.go          # Existing interface definitions
├── sqlite/               # SQLite implementation
├── mysql/                # MySQL implementation
└── postgres/             # PostgreSQL implementation (new)
    ├── manager.go        # PostgreSQL manager
    └── repository.go     # PostgreSQL repository implementation
```

**Core Code Examples:**

**postgres/manager.go**
```go
package postgres

import (
    "api-key-rotator/backend/internal/infrastructure/database"
)

type Manager struct {
    dsn string
    repo database.Repository
}

func NewPostgresManager(dsn string) *Manager {
    return &Manager{dsn: dsn}
}

func (m *Manager) Initialize() (database.Repository, error) {
    repo, err := NewPostgresRepository(m.dsn)
    if err != nil {
        return nil, err
    }
    m.repo = repo
    return repo, nil
}
```

**postgres/repository.go**
```go
package postgres

import (
    "api-key-rotator/backend/internal/models"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type Repository struct {
    db *gorm.DB
}

func NewPostgresRepository(dsn string) (*Repository, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    return &Repository{db: db}, nil
}

// Implement all database.Repository interface methods
func (r *Repository) GetDB() *gorm.DB { return r.db }
func (r *Repository) CreateProxyConfig(config *models.ProxyConfig) error {
    return r.db.Create(config).Error
}
// ... other methods similar to SQLite implementation
```

#### 2. Update Configuration Factory

Add PostgreSQL support in `internal/config/factory.go`:
```go
// Add PostgreSQL option in CreateDatabaseManager function
if strings.Contains(os.Getenv("DATABASE_URL"), "postgres") {
    return postgres.NewPostgresManager(dsn), nil
}
```

#### 3. Add Dependencies

Add to `go.mod`:
```bash
go get gorm.io/driver/postgres
```

#### 4. Environment Variables Configuration

```bash
# PostgreSQL connection configuration
DATABASE_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
```

#### 5. Build Configuration

Add PostgreSQL build support in Dockerfile, similar to existing MySQL and SQLite builds.

This extension approach maintains the integrity of the existing architecture while providing flexible database support.

### Adding New LLM Adapters (Adapters)

Adapters are responsible for handling direct communication with upstream LLM APIs, including request building, authentication, and key rotation. Core code is located in `backend/internal/adapters`.

#### 1. Define Adapter Structure
Implement the `LLMAdapter` interface:

```go
type LLMAdapter interface {
    ProcessRequest() (*services.TargetRequest, error)
}
```

#### 2. Implement Core Logic
Embed `BaseLLMAdapter` to reuse common logic (like key rotation):

```go
type XAIAdapter struct {
    *BaseLLMAdapter
}

func (a *XAIAdapter) ProcessRequest() (*services.TargetRequest, error) {
    // 1. Verify Proxy Key
    // 2. Rotate Upstream Key
    upstreamKey, err := a.RotateUpstreamKey()
    
    // 3. Build Request (Filtering headers, removing gzip, etc.)
    headers := utils.FilterRequestHeaders(a.c.Request.Header, []string{"authorization", "accept-encoding"})
    headers["Authorization"] = "Bearer " + upstreamKey
    
    // 4. Return TargetRequest object
    return &services.TargetRequest{...}, nil
}
```

### Extending New API Formats (Converters)

Converters are responsible for handling transformations between client-side formats and backend API expected formats. Core code is located in `backend/internal/converters`.

#### 1. Register New Format
Register the new format identifier in `backend/internal/converters/formats/registry.go`:

```go
// Example: Adding a new format identifier
const (
    FormatClaudeNative = "claude_native"
)
```

#### 2. Implement Format Handler
Implement the `FormatHandler` interface in `backend/internal/converters/formats/claude_native/`:

```go
type ClaudeNativeHandler struct{}

// Build Request: Convert universal request format to Claude native format
func (h *ClaudeNativeHandler) BuildRequest(req *types.UniversalRequest) ([]byte, error) {
    // Implement conversion logic...
}

// Parse Response: Convert Claude native response to universal response format
func (h *ClaudeNativeHandler) ParseResponse(body []byte) (*types.UniversalResponse, error) {
    // Implement parsing logic...
}
```

#### 3. Handle Streaming Responses (Optional)
If streaming support is required, implement the `StreamHandler` interface to handle SSE (Server-Sent Events) transformations. This is critical for chat interaction experiences.

```go
type StreamHandler interface {
    // Parse stream chunk: Parse format-specific SSE data line into universal stream chunk
    ParseStreamChunk(chunk []byte) (*UniversalStreamChunk, error)
    
    // Build stream chunk: Build format-specific SSE data line from universal stream chunk
    BuildStreamChunk(chunk *UniversalStreamChunk) ([]byte, error)
    
    // Build start event (e.g., OpenAI's role delta)
    BuildStartEvent(model string, id string) [][]byte
    
    // Build end event (e.g., [DONE])
    BuildEndEvent() [][]byte
}
```

Implementation Example:

```go
func (h *ClaudeStreamHandler) ParseStreamChunk(chunk []byte) (*UniversalStreamChunk, error) {
    // 1. Parse SSE line (e.g., "data: {...}")
    // 2. Extract content delta
    // 3. Construct UniversalStreamChunk
    return &UniversalStreamChunk{
        Content: deltaContent,
        FinishReason: finishReason,
    }, nil
}
```