# API Key Rotator

[English](README.md) | [中文简体](README_CN.md)

## 🚀 Optimized Build Architecture

**This project uses separate optimized builds for different deployment scenarios** - choose the right build for your needs:

### Two Build Options

| Build | Database | Cache | Image Size | Use Case | QPS Support |
|-------|----------|-------|------------|----------|-------------|
| 🟢 **Lightweight Build** | SQLite | Memory Cache | ~50MB | Personal Projects, Small Applications | < 5K |
| 🔴 **Enterprise Build** | MySQL + SQLite | Redis + Memory Cache | ~80MB | Business Applications, Large Deployments | > 10K |

### Architecture Benefits

- **Optimized Dependencies**: Each build only includes necessary libraries
- **Faster Downloads**: Smaller images for quick deployment
- **Clean Code**: Interface abstraction separates business logic from infrastructure
- **Easy Maintenance**: Clear separation between lightweight and enterprise features

### 📋 Complete Environment Variables

#### 🔴 Database Configuration (Optional - defaults to SQLite if not set)

```bash
# MySQL Connection String
DATABASE_URL=mysql://user:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True&loc=Local

# OR Individual MySQL Variables
DB_HOST=localhost
DB_USER=appdb
DB_PASSWORD=your_strong_password
DB_NAME=api_key_rotator
DB_PORT=3306

# SQLite Path (only used when SQLite mode is detected)
DATABASE_PATH=/app/data/api_key_rotator.db
```

#### 🟠 Redis Configuration (Optional - defaults to memory cache if not set)

```bash
# Basic Redis Configuration
REDIS_HOST=localhost          # Required to enable Redis
REDIS_PORT=6379               # Optional, defaults to 6379
REDIS_PASSWORD=your_password   # Optional, defaults to empty
REDIS_URL=redis://localhost:6379/0  # Optional, alternative connection string
REDIS_DB=0                    # Optional, defaults to 0
```

#### 🔧 Required Configuration (Must be set)

```bash
# Security Configuration (Required)
ADMIN_PASSWORD=your_admin_password
JWT_SECRET=your_very_long_jwt_secret
GLOBAL_PROXY_KEYS=key1,key2,key3

# Service Configuration (Optional)
BACKEND_PORT=8000
PROXY_PUBLIC_BASE_URL=http://localhost:8000
LOG_LEVEL=info
RESET_DB_TABLES=false
```

### Quick Deployment Examples

**🟢 Lightweight Build (SQLite + Memory Cache)**
```bash
# Using pre-built image
docker run -d \
  -p 8000:8000 \
  -e ADMIN_PASSWORD="your_password" \
  -e JWT_SECRET="your_jwt_secret" \
  -e GLOBAL_PROXY_KEYS="your_proxy_key" \
  -v $(pwd)/data:/app/data \
  api-key-rotator:lightweight

# OR build from source
git clone https://github.com/your-repo/APIKeyRotator.git
cd APIKeyRotator
make build-lightweight
docker run -d \
  -p 8000:8000 \
  -e ADMIN_PASSWORD="your_password" \
  -e JWT_SECRET="your_jwt_secret" \
  -e GLOBAL_PROXY_KEYS="your_proxy_key" \
  -v $(pwd)/data:/app/data \
  api-key-rotator:lightweight
```

**🔴 Enterprise Build (MySQL + Redis)**
```bash
# Using pre-built image
docker run -d \
  -p 8000:8000 \
  -e ADMIN_PASSWORD="your_password" \
  -e JWT_SECRET="your_jwt_secret" \
  -e GLOBAL_PROXY_KEYS="your_proxy_key" \
  -e DB_HOST="mysql-server" \
  -e DB_USER="appdb" \
  -e DB_PASSWORD="your_db_password" \
  -e DB_NAME="api_key_rotator" \
  -e REDIS_HOST="redis-server" \
  api-key-rotator:enterprise

# OR build from source
make build-enterprise
docker run -d \
  -p 8000:8000 \
  -e ADMIN_PASSWORD="your_password" \
  -e JWT_SECRET="your_jwt_secret" \
  -e GLOBAL_PROXY_KEYS="your_proxy_key" \
  -e DB_HOST="mysql-server" \
  -e DB_USER="appdb" \
  -e DB_PASSWORD="your_db_password" \
  -e DB_NAME="api_key_rotator" \
  -e REDIS_HOST="redis-server" \
  api-key-rotator:enterprise
```

📖 **See deployment details below** 👇

---

## 🎯 Interface-Abstraction Architecture

**Clean separation between business logic and infrastructure** using interface abstraction pattern:

- **Core Business Logic**: Independent of specific database or cache implementation
- **Infrastructure Adapters**: Pluggable implementations for SQLite, MySQL, Memory Cache, and Redis
- **Optimized Builds**: Separate builds for different deployment scenarios

## Introduction

**API Key Rotator** is a powerful and flexible API key management and request proxy solution built with Go (Gin). It is designed to centralize the management of all your third-party API keys and provide automatic rotation, load balancing, and secure isolation through a unified proxy endpoint.

Whether you need to provide high availability for traditional RESTful APIs or a unified, SDK-compatible access point for large model APIs like OpenAI, this project offers an elegant and scalable solution.

The project includes a high-performance **Go backend** and a simple, easy-to-use **Vue 3 admin panel**, with "one-click" deployment via Docker Compose.

## ✨ Core Features

*   **🔧 Optimized Builds**: Separate builds for lightweight and enterprise scenarios, reducing image size and dependencies.
*   **🏗️ Interface Abstraction**: Clean architecture separating business logic from infrastructure implementations.
*   **🔑 Centralized Key Management**: Manage API key pools for all services in a unified web interface.
*   **🔄 Dynamic Key Rotation**: Atomic rotation based on cache (supports both memory cache and Redis) to effectively distribute API request quotas.
*   **🚀 Type-Safe Proxies**:
    *   **Generic API Proxy (`/proxy`)**: Provides proxy services for any RESTful API.
    *   **LLM API Proxy (`/llm`)**: Offers native streaming support and an SDK-friendly `base_url` for large model APIs compatible with OpenAI's format. Supported providers include **OpenAI, Gemini, Anthropic**, etc.
*   **🛡️ Secure Isolation**: All proxy requests are authenticated via global keys, with support for multiple keys to protect real backend keys from being exposed.
*   **🐳 Efficient Docker Images**: Optimized multi-stage builds with minimal runtime dependencies.

## 🚀 Quick Start

### Method 1: Using Makefile (Recommended)

The simplest deployment method with optimized builds:

```bash
# Clone the project
git clone https://github.com/your-repo/APIKeyRotator.git
cd APIKeyRotator

# Build lightweight version (SQLite + Memory Cache)
make build-lightweight

# Start the service
docker run -d \
  -p 8000:8000 \
  -e ADMIN_PASSWORD="your_password" \
  -e JWT_SECRET="your_jwt_secret" \
  -e GLOBAL_PROXY_KEYS="your_proxy_key" \
  -v $(pwd)/data:/app/data \
  api-key-rotator:lightweight

# Access the application
open http://localhost:8000
```

### Method 2: Direct Docker Build

For enterprise deployment with external services:

```bash
# Build enterprise version (MySQL + Redis)
make build-enterprise

# Start with external MySQL and Redis
docker run -d \
  -p 8000:8000 \
  -e ADMIN_PASSWORD="your_password" \
  -e JWT_SECRET="your_jwt_secret" \
  -e GLOBAL_PROXY_KEYS="your_proxy_key" \
  -e DB_HOST="mysql-server" \
  -e DB_USER="appdb" \
  -e DB_PASSWORD="your_db_password" \
  -e DB_NAME="api_key_rotator" \
  -e REDIS_HOST="redis-server" \
  api-key-rotator:enterprise
```

### Method 3: Pre-built Images

```bash
# Lightweight version
docker run -d \
  -p 8000:8000 \
  -e ADMIN_PASSWORD="your_password" \
  -e JWT_SECRET="your_jwt_secret" \
  -e GLOBAL_PROXY_KEYS="your_proxy_key" \
  -v $(pwd)/data:/app/data \
  yourusername/api-key-rotator:lightweight

# Enterprise version
docker run -d \
  -p 8000:8000 \
  -e ADMIN_PASSWORD="your_password" \
  -e JWT_SECRET="your_jwt_secret" \
  -e GLOBAL_PROXY_KEYS="your_proxy_key" \
  -e DB_HOST="mysql-server" \
  -e DB_USER="appdb" \
  -e DB_PASSWORD="your_db_password" \
  -e DB_NAME="api_key_rotator" \
  -e REDIS_HOST="redis-server" \
  yourusername/api-key-rotator:enterprise
```

### Method 4: Docker Compose Deployment

#### 1. Prerequisites

Ensure that [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/install/) are installed on your system.

#### 2. Configure the Project

```bash
# Clone the project
git clone https://github.com/your-repo/APIKeyRotator.git
cd APIKeyRotator

# Copy the configuration file template
cp .env.example.en .env
# Or use Chinese version: cp .env.example.cn .env
```

#### 3. Proxy Key Configuration

This project uses the `GLOBAL_PROXY_KEYS` environment variable to configure proxy authentication keys, supporting a single key or multiple keys:

1.  **Single Key**:
    ```bash
    GLOBAL_PROXY_KEYS=your_secret_key
    ```

2.  **Multiple Keys** (Recommended for multi-client scenarios):
    ```bash
    GLOBAL_PROXY_KEYS=key1,key2,key3
    ```

#### 4. Start the Services

**🟢 Lightweight Deployment**
```bash
# Copy English configuration template
cp .env.example.en .env

# Build lightweight version first
make build-lightweight

# Start with Docker Compose
docker-compose -f docker-compose.yml up -d
```

**🔴 Enterprise Deployment**
```bash
# Copy English configuration template
cp .env.example.en .env

# Build enterprise version first
make build-enterprise

# Add database and cache configuration
cat >> .env << EOF
DB_HOST=db
DB_USER=appdb
DB_PASSWORD=your_db_password
DB_NAME=api_key_rotator
REDIS_HOST=redis
REDIS_PASSWORD=your_redis_password
EOF

# Start with external services
docker-compose -f docker-compose.prod.yml up -d
```

**Or use Chinese template**:
```bash
# Copy Chinese configuration template
cp .env.example.cn .env
# ... same as above
```

#### 5. Access URLs

*   **Web Application**: `http://localhost:8000`

## Local Development Without Docker (Optional)

If you prefer to run and debug the source code directly on your local machine without using Docker, you can follow these steps.

### 1. Prerequisites

*   Install [Node.js](https://nodejs.org/) (18+)
*   Install [Go](https://golang.org/) (1.21+)
*   Install and run **MySQL** and **Redis** services locally

### 2. Start the Backend Service

1.  **Enter the Go backend directory**
    ```bash
    cd backend/
    ```

2.  **Install dependencies**
    ```bash
    go mod download
    ```

3.  **Configure environment variables**
    Create a `.env` file in the project root (refer to `.env.example`) and configure the connection information for the database and Redis.

4.  **Start the backend server**
    ```bash
    go run main.go
    ```
    The service will run at `http://127.0.0.1:8000`.

### 3. Start the Frontend Service

1.  **Enter the frontend directory** (in another terminal)
    ```bash
    cd frontend/
    ```

2.  **Install dependencies**
    ```bash
    npm install
    ```

3.  **Start the frontend server**
    ```bash
    npm run dev
    ```
    Vite will automatically handle API proxying. The service will run at `http://localhost:5173`.

Now, you can access the admin panel at `http://localhost:5173`.

## Usage Example

### LLM API Proxy

Using the `openai` Python SDK as an example, combined with an `OpenRouter` model, you can use the proxy service by modifying the `base_url`.

```python
from openai import OpenAI

client = OpenAI(
  # Format: http://<PROXY_PUBLIC_BASE_URL>/llm/<Service Slug>
  base_url="http://PROXY_PUBLIC_BASE_URL/llm/openrouter-api",
  api_key="<GLOBAL_PROXY_KEY>",
)

completion = client.chat.completions.create(
  # Please refer to the specific provider's documentation for model names
  model="openai/gpt-4o",
  messages=[
    {
      "role": "user",
      "content": "What is the meaning of life?"
    }
  ]
)

print(completion.choices[0].message.content)
```

Where `PROXY_PUBLIC_BASE_URL` and `GLOBAL_PROXY_KEY` are the environment variables you configured in your `.env` file.

### Generic API Proxy

The Generic API Proxy can be used for any RESTful API. Here's an example of calling a weather API using Python requests library:

```python
import requests

# Configure proxy parameters
proxy_url = "http://PROXY_PUBLIC_BASE_URL/proxy/weather/current"
proxy_key = "<GLOBAL_PROXY_KEY>"

# Query parameters
params = {
    "query": "London"
    # When proxying requests to the target API, the system polls the real API keys configured in the backend and appends them to the original authorization parameter access_key (which is configured in the backend).
}

# Set headers
headers = {
    "X-Proxy-Key": proxy_key
}

# Make the request
response = requests.get(proxy_url, params=params, headers=headers)

# Handle the response
if response.status_code == 200:
    data = response.json()
    print(f"Weather information: {data}")
else:
    print(f"Request failed with status code: {response.status_code}")
```

In this example:
1. `weather` is the service slug configured in the admin panel
2. `current` is the path of the target API endpoint
3. `PROXY_PUBLIC_BASE_URL` is your proxy service address
4. `<GLOBAL_PROXY_KEY>` is one of the global proxy keys you configured

The proxy will automatically forward the request to the configured target URL, appending the path and query parameters to the target address.

## 📚 Technical Features

*   **🔧 Smart Configuration Detection**: System automatically selects the most suitable database and cache scheme based on environment variables
*   **⚡ High-Performance Architecture**: Supports various performance requirements from lightweight to enterprise-grade
*   **🎯 Zero-Configuration Startup**: Default mode requires no database or cache service configuration
*   **🔄 Seamless Upgrades**: Switch between different deployment modes without code changes
*   **🛡️ Production-Ready**: Includes health checks, logging, error handling, and other production-grade features

## 📖 Related Documentation

If you want to dive deeper into the code, please refer to the following documents:

*   **[Backend Development Guide](./backend/README.md)**
*   **[Frontend Development Guide](./frontend/README.md)**

## 🔧 Deployment Examples

### 🟢 Lightweight Deployment

```bash
# SQLite + Memory Cache - Simple and efficient
docker run -d \
  -p 8000:8000 \
  -e ADMIN_PASSWORD="your_password" \
  -e JWT_SECRET="your_jwt_secret" \
  -e GLOBAL_PROXY_KEYS="your_proxy_key" \
  -v $(pwd)/data:/app/data \
  api-key-rotator
```

### 🔴 Enterprise Deployment

```bash
# MySQL + Redis - High performance and scalable
docker run -d \
  -p 8000:8000 \
  -e ADMIN_PASSWORD="secure_password" \
  -e JWT_SECRET="very_long_jwt_secret" \
  -e GLOBAL_PROXY_KEYS="proxy_key1,proxy_key2" \
  -e DB_HOST="mysql.internal" \
  -e DB_USER="appdb" \
  -e DB_PASSWORD="db_password" \
  -e DB_NAME="api_key_rotator" \
  -e REDIS_HOST="redis.internal" \
  -e REDIS_PORT=6379 \
  -e REDIS_PASSWORD="redis_password" \
  -e LOG_LEVEL=info \
  -v $(pwd)/data:/app/data \
  api-key-rotator
```

### 🐳 Docker Compose Example

```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8000:8000"
    environment:
      - ADMIN_PASSWORD=your_password
      - JWT_SECRET=your_jwt_secret
      - GLOBAL_PROXY_KEYS=your_proxy_key
      # Optional: Add these for enterprise mode
      - DB_HOST=db
      - DB_USER=appdb
      - DB_PASSWORD=your_db_password
      - DB_NAME=api_key_rotator
      - REDIS_HOST=redis
    volumes:
      - ./data:/app/data
    depends_on:
      - db
      - redis

  db:
    image: mysql:8.0
    environment:
      - MYSQL_ROOT_PASSWORD=your_root_password
      - MYSQL_DATABASE=api_key_rotator
      - MYSQL_USER=appdb
      - MYSQL_PASSWORD=your_db_password
    volumes:
      - mysql_data:/var/lib/mysql

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass your_redis_password
    volumes:
      - redis_data:/data

volumes:
  mysql_data:
  redis_data:
```

### 🔧 Build Commands (Makefile)

The project includes a Makefile for simplified build management:

```bash
# Show all available commands
make info

# Build lightweight version
make build-lightweight

# Build enterprise version
make build-enterprise

# Build both versions
make build-all

# Test builds (development)
make test-lightweight
make test-enterprise

# Publish to Docker Hub
make publish-lightweight
make publish-enterprise
make publish-all

# Clean build cache
make clean
```

## ❓ Frequently Asked Questions

### Q: How do I choose the right deployment mode?
**A**:
- **Small Projects/Personal Use**: Use lightweight mode (SQLite + Memory Cache)
- **Business Applications**: Use enterprise mode (MySQL + Redis)

### Q: How can I check which database and cache types are currently being used?
**A**: The application displays log information when starting:
```
Database Type: sqlite
Cache Type: memory
```

### Q: How do I upgrade from lightweight mode to enterprise mode?
**A**: Simply add the corresponding environment variables, and the system will automatically detect and switch:
```bash
# Add MySQL configuration
DB_HOST=mysql-server
DB_USER=appdb
DB_PASSWORD=your_password

# Add Redis configuration
REDIS_HOST=redis-server
```

### Q: How is data migration handled?
**A**: The system automatically creates table structures when starting. To migrate data from SQLite to MySQL:

1. **Backup SQLite data**:
   ```bash
   cp data/api_key_rotator.db backup_$(date +%Y%m%d).db
   ```

2. **Add MySQL environment variables**:
   ```bash
   -e DB_HOST="mysql-server" \
   -e DB_USER="appdb" \
   -e DB_PASSWORD="your_password" \
   -e DB_NAME="api_key_rotator"
   ```

3. **Restart the application** - it will automatically create new tables in MySQL

For data import, you'll need to export from SQLite and import to MySQL manually or use a migration tool.

### Q: Does it support distributed deployment?
**A**: Yes, using MySQL + Redis mode supports fully distributed deployment.
