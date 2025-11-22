# API Key Rotator

[English](README.md) | [中文简体](README_CN.md)

## 🚀 Interface Abstraction Architecture + Optimized Builds

**This project uses interface abstraction architecture with separate optimized builds** - choose the right build for your needs:

### Two Build Options

| Build | Database | Cache | Image Size | Use Case | QPS Support |
|-------|----------|-------|------------|----------|-------------|
| 🟢 **Lightweight Build** | SQLite | Memory Cache | ~50MB | Personal Projects, Small Applications | < 5K |
| 🔴 **Enterprise Build** | MySQL | Redis | ~80MB | Business Applications, Large Deployments | > 10K |

### Architecture Benefits

- **Interface Abstraction**: Clean separation between business logic and infrastructure implementations
- **Optimized Dependencies**: Each build only includes necessary libraries
- **Faster Downloads**: Smaller images for quick deployment
- **Easy Maintenance**: Clear separation between lightweight and enterprise features
- **Adapter Pattern**: Pluggable database and cache implementations

### 🔧 Quick Start

#### Lightweight Build (Default)
```bash
# Build lightweight version
make build-lightweight

# Run with default SQLite + Memory Cache
docker-compose up -d
```

#### Enterprise Build
```bash
# Build enterprise version
make build-enterprise

# Run with MySQL + Redis
docker-compose -f docker-compose.prod.yml up -d
```

### 📋 Environment Variables

#### Database Configuration
```bash
# SQLite (Lightweight - Default)
DATABASE_PATH=/app/data/api_key_rotator.db

# MySQL (Enterprise)
DB_HOST=localhost
DB_USER=appdb
DB_PASSWORD=your_strong_password
DB_NAME=api_key_rotator
DB_PORT=3306

# OR use connection string
DATABASE_URL=mysql://user:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True&loc=Local
```

#### Cache Configuration
```bash
# Memory Cache (Lightweight - Default)
# No additional configuration needed

# Redis (Enterprise)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your_password
REDIS_URL=redis://localhost:6379/0
```

#### Application Configuration
```bash
# Server
BACKEND_PORT=8000
LOG_LEVEL=info

# Authentication
ADMIN_USERNAME=admin
ADMIN_PASSWORD=your_admin_password_here
JWT_SECRET=your_very_secret_and_random_jwt_key

# Proxy
GLOBAL_PROXY_KEYS=your_secure_global_proxy_key
PROXY_TIMEOUT=30
PROXY_PUBLIC_BASE_URL=http://localhost:8000

# Database reset option
RESET_DB_TABLES=false
```

### 🏗️ Project Structure

```
api-key-rotator/
├── Makefile                          # Build orchestration
├── docker-compose.yml                # Lightweight deployment
├── docker-compose.prod.yml           # Enterprise deployment
├── Dockerfile.lightweight            # Lightweight build
├── Dockerfile.enterprise             # Enterprise build
├── README.md                         # Project documentation
└── backend/                          # Go backend service
    ├── main.go                       # Application entry point
    ├── go.mod                        # Go module definition
    └── internal/                      # Internal packages
        ├── config/                    # Configuration management
        │   ├── config.go              # Configuration loading
        │   └── factory.go             # Infrastructure factory
        ├── infrastructure/            # Infrastructure layer
        │   ├── database/
        │   │   ├── interface.go        # Database repository interface
        │   │   ├── sqlite/             # SQLite implementation
        │   │   └── mysql/              # MySQL implementation
        │   └── cache/
        │       ├── interface.go        # Cache interface
        │       ├── memory/             # Memory cache implementation
        │       └── redis/              # Redis implementation
        ├── handlers/                  # HTTP handlers
        ├── models/                    # Data models
        ├── dto/                       # Data transfer objects
        ├── router/                    # Route configuration
        └── logger/                    # Logger configuration
└── frontend/                         # Vue.js frontend
    ├── src/                          # Source code
    ├── package.json                  # Dependencies
    └── Dockerfile                    # Frontend build
```

### 🛠️ Tech Stack

- **Backend**: Go + Gin Framework + GORM ORM
- **Frontend**: Vue 3 + TypeScript + Element Plus
- **Database**: MySQL 8.0+ (Enterprise) / SQLite (Lightweight)
- **Cache**: Redis 6.0+ (Enterprise) / Memory Cache (Lightweight)
- **Containerization**: Docker + Docker Compose
- **Architecture**: Interface Abstraction + Adapter Pattern

### 🌐 API Endpoints

After starting the service, you can access the following APIs:

- **Root Path**: `http://localhost:8000/` - Welcome message
- **Admin API**: `http://localhost:8000/admin/*` - Configuration management
- **Generic Proxy**: `http://localhost:8000/proxy/*` - Generic API proxy (coming soon)
- **LLM Proxy**: `http://localhost:8000/llm/*` - LLM API proxy (coming soon)

### 📦 Building Images

#### Option 1: Using Default Build (Lightweight)
```bash
# Build lightweight version (default)
docker build -t api-key-rotator .

# Build with custom tag
docker build -t my-api-key-rotator:latest .
```

#### Option 2: Using Makefile (Recommended)
```bash
# Build lightweight version
make build-lightweight

# Build enterprise version
make build-enterprise

# Build both versions
make build-all
```

#### Option 3: Specify Dockerfile Directly
```bash
# Lightweight build (SQLite + Memory Cache)
docker build -f Dockerfile.lightweight -t api-key-rotator:lightweight .

# Enterprise build (MySQL + Redis)
docker build -f Dockerfile.enterprise -t api-key-rotator:enterprise .
```

### 🐳 Docker Deployment

#### Lightweight Deployment
```bash
docker-compose up -d
```

#### Enterprise Deployment
```bash
docker-compose -f docker-compose.prod.yml up -d
```

### 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/handlers

# Run tests with coverage
go test -cover ./...
```

### 🔒 Security

- All proxy requests require `X-Proxy-Key` header authentication
- Admin interface requires username/password authentication
- Environment variables should be properly secured in production
- Database passwords and API keys should be encrypted

### 📈 Performance

- **Lightweight**: < 50MB image size, fast startup, minimal resource usage
- **Enterprise**: < 80MB image size, high concurrency, scalable architecture
- **API Response**: < 100ms for most operations under normal load

### 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

### 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.