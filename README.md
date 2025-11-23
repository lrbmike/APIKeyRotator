# API Key Rotator

[English](README.md) | [中文简体](README_CN.md)

An enterprise-grade API key management and rotation service, providing intelligent key pool management, automatic failover, and load balancing capabilities.

## ✨ Key Features

- 🔑 **API Key Management**: Centralized management of multiple API keys with add, delete, enable/disable operations
- 🔄 **Smart Rotation**: Automatic rotation between multiple keys to avoid single points of failure and quota limits
- 🌐 **Multi-API Support**: Supports both generic REST APIs and LLM model APIs (like OpenAI, Claude, etc.)
- 🎯 **Proxy Service**: Unified proxy interface for different API types, simplifying client integration
- 📊 **Management Interface**: Modern Vue3-based web management interface with bilingual support (English/Chinese)
- 🔒 **Secure Authentication**: JWT token authentication and proxy key verification
- 🏗️ **Flexible Architecture**: Supports SQLite/MySQL databases and memory/Redis caching

## 🚀 Interface Abstraction Architecture + Optimized Builds

**This project uses interface abstraction architecture with separate optimized builds** - choose the right build for your needs:

### Two Build Options

| Build | Database | Cache | Image Size | Use Case | QPS Support |
|------|--------|------|----------|----------|-------------|
| 🟢 **Lightweight Build** | SQLite | Memory Cache | ~50MB | Personal Projects, Small Applications | < 5K |
| 🔴 **Enterprise Build** | MySQL | Redis | ~80MB | Business Applications, Large Deployments | > 10K |

### Architecture Benefits

- **Interface Abstraction**: Clean separation between business logic and infrastructure implementations
- **Optimized Dependencies**: Each build only includes necessary libraries
- **Faster Downloads**: Smaller images for quick deployment
- **Easy Maintenance**: Clear separation between lightweight and enterprise features
- **Adapter Pattern**: Pluggable database and cache implementations

### 🔧 Quick Start

#### 🏗️ Docker Build & Deployment

**Lightweight Version (SQLite + Memory Cache)**
```bash
# Build lightweight image
docker build -t api-key-rotator:latest .

# Run container
docker run -d \
  --name api-key-rotator \
  -p 8000:8000 \
  -v $(pwd)/data:/app/data \
  -e ADMIN_USERNAME=admin \
  -e ADMIN_PASSWORD=your_admin_password \
  -e JWT_SECRET=your_very_secret_and_random_jwt_key \
  api-key-rotator:latest
```

**Enterprise Version (MySQL + Redis)**
```bash
# Build enterprise image
docker build -f Dockerfile.enterprise -t api-key-rotator:enterprise .

# Run container
docker run -d \
  --name api-key-rotator \
  -p 8000:8000 \
  -e DB_TYPE=mysql \
  -e DB_HOST=your_mysql_host \
  -e DB_USER=your_db_user \
  -e DB_PASSWORD=your_db_password \
  -e DB_NAME=api_key_rotator \
  -e CACHE_TYPE=redis \
  -e REDIS_HOST=your_redis_host \
  -e REDIS_PORT=6379 \
  -e ADMIN_USERNAME=admin \
  -e ADMIN_PASSWORD=your_admin_password \
  -e JWT_SECRET=your_very_secret_and_random_jwt_key \
  api-key-rotator:enterprise
```

#### 🐳 Using Docker Compose

**Lightweight Deployment**
```bash
docker-compose up -d
```

**Enterprise Deployment**
```bash
docker-compose -f docker-compose.enterprise.yml up -d
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

# Or use connection string
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
├── docker-compose.yml                # Lightweight deployment
├── docker-compose.enterprise.yml     # Enterprise deployment
├── Dockerfile                        # Default build (lightweight)
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
- **Frontend**: Vue 3 + JavaScript + Element Plus + Vue Router + Vue I18n
- **Database**: MySQL 8.0+ (Enterprise) / SQLite (Lightweight)
- **Cache**: Redis 6.0+ (Enterprise) / Memory Cache (Lightweight)
- **Containerization**: Docker + Docker Compose
- **Architecture**: Interface Abstraction + Adapter Pattern

### 🌐 API Endpoints

After starting the service, you can access the following APIs:

- **Root Path**: `http://localhost:8000/` - Service status information
- **Admin API**: `http://localhost:8000/admin/*` - Backend management interface
  - `GET /admin/app-config` - Get application configuration
  - `POST /admin/login` - User login
  - `GET/POST/PUT/DELETE /admin/proxy-configs` - Proxy configuration management
  - `GET/POST/DELETE /admin/proxy-configs/:id/keys` - API key management
  - `PATCH /admin/keys/:keyID` - Key status management
- **Frontend Management Interface**: `http://localhost:8000/` - Vue3 admin management interface

### 🐳 Deployment Options

#### Build Options Description

| Build Type | Dockerfile | Database | Cache | Image Size | Use Case |
|-----------|------------|---------|-------|------------|----------|
| Lightweight | `Dockerfile` | SQLite | Memory | ~50MB | Personal development, small deployments |
| Enterprise | `Dockerfile.enterprise` | MySQL | Redis | ~80MB | Production environment, large applications |

#### Production Environment Recommendations

**Data Persistence Configuration**
```bash
# Ensure data persistence
-v $(pwd)/data:/app/data          # SQLite database files
-v $(pwd)/logs:/app/logs          # Application logs
```

**Environment Variable Management**
```bash
# Create environment variable file
cat > .env << EOF
ADMIN_USERNAME=admin
ADMIN_PASSWORD=your_secure_password
JWT_SECRET=your_very_secret_key
GLOBAL_PROXY_KEYS=your_proxy_key_1,proxy_key_2
EOF

# Start with environment file
docker run --env-file .env api-key-rotator:latest
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