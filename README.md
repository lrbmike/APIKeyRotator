# API Key Rotator

[English](README.md) | [中文简体](README_CN.md)

A lightweight API key management and rotation service designed to provide a simple and efficient solution. It helps developers easily manage API keys through intelligent key pool management, automatic failover, and load balancing. The project also offers enterprise-level deployment options to suit different use cases.

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

| Variable | Description | Default | Example |
|---|---|---|---|
| **General** | | | |
| `BACKEND_PORT` | Port for the backend service. | `8000` | `8000` |
| `LOG_LEVEL` | Logging level. | `info` | `debug` |
| `ADMIN_USERNAME` | Initial admin username. | `admin` | `admin` |
| `ADMIN_PASSWORD` | Initial admin password. | `your_admin_password` | `mysecretpassword` |
| `JWT_SECRET` | Secret key for JWT tokens. | `your_very_secret...` | `a_long_random_string` |
| `GLOBAL_PROXY_KEYS` | Global proxy keys, comma-separated. | (empty) | `key1,key2` |
| `PROXY_TIMEOUT` | Proxy request timeout in seconds. | `30` | `60` |
| `PROXY_PUBLIC_BASE_URL` | Public access URL for the service. | `http://localhost:8000` | `https://your.domain.com` |
| **Database** | | | |
| `DB_TYPE` | Database type. | `sqlite` | `mysql` |
| `DATABASE_PATH` | Path for SQLite database file. | `/app/data/rotator.db` | |
| `DB_HOST` | MySQL host. | | `localhost` |
| `DB_USER` | MySQL username. | | `dbuser` |
| `DB_PASSWORD` | MySQL password. | | `dbpass` |
| `DB_NAME` | MySQL database name. | | `rotator_db` |
| `DB_PORT` | MySQL port. | | `3306` |
| `DATABASE_URL` | Database connection string (priority). | | `mysql://...` |
| **Cache** | | | |
| `CACHE_TYPE` | Cache type. | `memory` | `redis` |
| `REDIS_HOST` | Redis host. | | `localhost` |
| `REDIS_PORT` | Redis port. | | `6379` |
| `REDIS_PASSWORD` | Redis password. | | (empty) |
| `REDIS_URL` | Redis connection string (priority). | | `redis://...` |

### 🏗️ Project Structure

The project is divided into two main parts: `backend` (a core API service written in Go) and `frontend` (a management interface built with Vue.js). Each part has its own `README.md` file with a more detailed structure description.

- `backend/`: The backend service responsible for API proxying, key management, and authentication.
- `frontend/`: The frontend application that provides a user-friendly web interface for managing proxy configurations and keys.
- `Dockerfile`: Used to build the Docker image for the lightweight version.
- `Dockerfile.enterprise`: Used to build the Docker image for the enterprise version.
- `docker-compose.yml`: For quick deployment of the lightweight version.
- `docker-compose.enterprise.yml`: For quick deployment of the enterprise version.

### 🛠️ Tech Stack

- **Backend**: Go + Gin Framework + GORM ORM
- **Frontend**: Vue 3 + JavaScript + Element Plus + Vue Router + Vue I18n
- **Database**: MySQL 8.0+ (Enterprise) / SQLite (Lightweight)
- **Cache**: Redis 6.0+ (Enterprise) / Memory Cache (Lightweight)
- **Containerization**: Docker + Docker Compose
- **Architecture**: Interface Abstraction + Adapter Pattern

### 📖 Usage Example

Let's take `OpenRouter` as an example. You can set it up as follows:

1.  Create a new proxy configuration in the management interface.
2.  **Service Slug**: Enter `openai-openrouter` (customizable).
3.  **API Format**: Select `OpenAI Compatible`.
4.  **Target Base URL**: Enter `https://openrouter.ai/api/v1`.
5.  Add your `OpenRouter` API keys to the key pool for this configuration.

Once configured, you can use it in any OpenAI-compatible client (e.g., `Cherry Studio`). Set the client's `Base URL` or `API Endpoint` to:

```
${PROXY_PUBLIC_BASE_URL}/llm/openai-openrouter
```

And fill the `API Key` field with the global proxy key you set in the `GLOBAL_PROXY_KEYS` environment variable.

- `${PROXY_PUBLIC_BASE_URL}` is the public access address you configure for the service (e.g., `http://localhost:8000`).
- The `openai-openrouter` in `/llm/openai-openrouter` corresponds to the **Service Slug** you set.

You can also test it directly with `curl`:
```bash
# Call the proxy endpoint using curl
curl -X POST ${PROXY_PUBLIC_BASE_URL}/llm/openai-openrouter/v1/chat/completions \
-H "Authorization: Bearer ${GLOBAL_PROXY_KEYS}" \
-H "Content-Type: application/json" \
-d '{
  "model": "google/gemini-flash-1.5",
  "messages": [{"role": "user", "content": "Hello!"}],
  "stream": false
}'
```

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
