# API Key Rotator

[English](README.md) | [中文简体](README_CN.md)

一个轻量级的API密钥管理和轮换服务，旨在提供一个简单、高效的解决方案。它通过智能的密钥池管理、自动故障转移和负载均衡功能，帮助开发者轻松管理API密钥。项目同时提供企业级的部署选项，以满足不同的使用场景。

## ✨ 主要功能

- 🔑 **API密钥管理**: 支持多个API密钥的集中管理，包括添加、删除、启用/禁用操作
- 🔄 **智能轮换**: 自动在多个密钥间轮换使用，避免单点故障和配额限制
- 🌐 **多API类型支持**: 同时支持通用REST API和LLM模型API（如OpenAI、Claude等）
- 🎯 **代理服务**: 为不同类型的API提供统一的代理接口，简化客户端集成
- 📊 **管理界面**: 基于Vue3的现代化Web管理界面，支持中英双语
- 🔒 **安全认证**: JWT令牌认证和代理密钥验证
- 🏗️ **灵活架构**: 支持SQLite/MySQL数据库和内存/Redis缓存

## 🚀 接口抽象架构 + 优化构建

**本项目采用接口抽象架构配合分离式优化构建** - 根据需求选择合适的构建版本：

### 两种构建方案

| 构建 | 数据库 | 缓存 | 镜像大小 | 适用场景 | QPS支持 |
|------|--------|------|----------|----------|-------------|
| 🟢 **轻量级构建** | SQLite | 内存缓存 | ~50MB | 个人项目、小型应用 | < 5K |
| 🔴 **企业级构建** | MySQL | Redis | ~80MB | 企业应用、大型部署 | > 10K |

### 架构优势

- **接口抽象**: 业务逻辑与基础设施实现通过明确定义的接口进行清晰的分离
- **优化依赖**: 每个构建只包含必要的库文件
- **快速下载**: 更小的镜像便于快速部署
- **易于维护**: 轻量级和企业级功能分离明确
- **适配器模式**: 可插拔的数据库和缓存实现

### 🔧 快速开始

#### 🏗️ Docker构建部署

**轻量级版本（SQLite + 内存缓存）**
```bash
# 构建轻量级镜像
docker build -t api-key-rotator:latest .

# 运行容器
docker run -d \
  --name api-key-rotator \
  -p 8000:8000 \
  -v $(pwd)/data:/app/data \
  -e ADMIN_USERNAME=admin \
  -e ADMIN_PASSWORD=your_admin_password \
  -e JWT_SECRET=your_very_secret_and_random_jwt_key \
  api-key-rotator:latest
```

**企业级版本（MySQL + Redis）**
```bash
# 构建企业级镜像
docker build -f Dockerfile.enterprise -t api-key-rotator:enterprise .

# 运行容器
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

#### 🐳 使用Docker Compose

**轻量级部署**
```bash
docker-compose up -d
```

**企业级部署**
```bash
docker-compose -f docker-compose.enterprise.yml up -d
```

### 📋 环境变量

| 变量名 | 描述 | 默认值 | 示例 |
|---|---|---|---|
| **通用** | | | |
| `BACKEND_PORT` | 后端服务监听的端口。 | `8000` | `8000` |
| `LOG_LEVEL` | 日志级别。 | `info` | `debug` |
| `ADMIN_USERNAME` | 管理员初始用户名。 | `admin` | `admin` |
| `ADMIN_PASSWORD` | 管理员初始密码。 | `your_admin_password` | `mysecretpassword` |
| `JWT_SECRET` | 用于生成JWT令牌的密钥。 | `your_very_secret...` | `a_long_random_string` |
| `GLOBAL_PROXY_KEYS` | 全局代理密钥，用逗号分隔。 | (空) | `key1,key2` |
| `PROXY_TIMEOUT` | 代理请求的超时时间（秒）。 | `30` | `60` |
| `PROXY_PUBLIC_BASE_URL` | 服务的公共访问URL。 | `http://localhost:8000` | `https://your.domain.com` |
| **数据库** | | | |
| `DB_TYPE` | 数据库类型。 | `sqlite` | `mysql` |
| `DATABASE_PATH` | SQLite数据库文件路径。 | `/app/data/rotator.db` | |
| `DB_HOST` | MySQL主机。 | | `localhost` |
| `DB_USER` | MySQL用户名。 | | `dbuser` |
| `DB_PASSWORD` | MySQL密码。 | | `dbpass` |
| `DB_NAME` | MySQL数据库名。 | | `rotator_db` |
| `DB_PORT` | MySQL端口。 | | `3306` |
| `DATABASE_URL` | 数据库连接字符串 (优先)。 | | `mysql://...` |
| **缓存** | | | |
| `CACHE_TYPE` | 缓存类型。 | `memory` | `redis` |
| `REDIS_HOST` | Redis主机。 | | `localhost` |
| `REDIS_PORT` | Redis端口。 | | `6379` |
| `REDIS_PASSWORD` | Redis密码。 | | (空) |
| `REDIS_URL` | Redis连接字符串 (优先)。 | | `redis://...` |

### 🏗️ 项目结构

项目分为两大部分：`backend`（Go语言实现的核心API服务）和 `frontend`（Vue.js实现的管理界面）。每个部分都有其独立的`README.md`文件，其中包含更详细的结构说明。

- `backend/`: 后端服务，负责API代理、密钥管理和认证。
- `frontend/`: 前端应用，提供一个用户友好的Web界面来管理代理配置和密钥。
- `Dockerfile`: 用于构建轻量级版本的Docker镜像。
- `Dockerfile.enterprise`: 用于构建企业级版本的Docker镜像。
- `docker-compose.yml`: 用于快速部署轻量级版本。
- `docker-compose.enterprise.yml`: 用于快速部署企业级版本。

### 🛠️ 技术栈

- **后端**: Go + Gin框架 + GORM ORM
- **前端**: Vue 3 + JavaScript + Element Plus + Vue Router + Vue I18n
- **数据库**: MySQL 8.0+（企业级）/ SQLite（轻量级）
- **缓存**: Redis 6.0+（企业级）/ 内存缓存（轻量级）
- **容器化**: Docker + Docker Compose
- **架构**: 接口抽象 + 适配器模式

### 📖 使用示例

以配置 `OpenRouter` 为例，您可以进行如下设置：

1.  在管理界面创建一个新的代理配置。
2.  **服务标识 (Slug)**: 填入 `openai-openrouter` (可自定义)。
3.  **API 格式**: 选择 `OpenAI Compatible`。
4.  **目标 Base URL**: 填入 `https://openrouter.ai/api/v1`。
5.  添加您的 `OpenRouter` API 密钥到此配置的密钥池中。

配置完成后，您就可以在任何兼容OpenAI的客户端（例如 `Cherry Studio`）中使用了。将客户端的 `Base URL` 或 `API Endpoint` 设置为：

```
${PROXY_PUBLIC_BASE_URL}/llm/openai-openrouter
```

并将 `API 密钥` 字段填写为您在环境变量 `GLOBAL_PROXY_KEYS` 中设置的全局代理密钥。

- `${PROXY_PUBLIC_BASE_URL}` 是您为服务配置的公共访问地址 (例如 `http://localhost:8000`)。
- `/llm/openai-openrouter` 中的 `openai-openrouter` 对应您设置的 **服务标识 (Slug)**。

您也可以直接使用 `curl` 进行测试：
```bash
# 使用 curl 调用代理接口
curl -X POST ${PROXY_PUBLIC_BASE_URL}/llm/openai-openrouter/v1/chat/completions \
-H "Authorization: Bearer ${GLOBAL_PROXY_KEYS}" \
-H "Content-Type: application/json" \
-d '{
  "model": "google/gemini-flash-1.5",
  "messages": [{"role": "user", "content": "Hello!"}],
  "stream": false
}'
```

### 🐳 部署选项

#### 构建选项说明

| 构建类型 | Dockerfile | 数据库 | 缓存 | 镜像大小 | 适用场景 |
|---------|------------|--------|------|----------|----------|
| 轻量级 | `Dockerfile` | SQLite | 内存 | ~50MB | 个人开发、小型部署 |
| 企业级 | `Dockerfile.enterprise` | MySQL | Redis | ~80MB | 生产环境、大型应用 |

#### 生产环境部署建议

**数据持久化配置**
```bash
# 确保数据持久化
-v $(pwd)/data:/app/data          # SQLite数据库文件
-v $(pwd)/logs:/app/logs          # 应用日志
```

**环境变量管理**
```bash
# 创建环境变量文件
cat > .env << EOF
ADMIN_USERNAME=admin
ADMIN_PASSWORD=your_secure_password
JWT_SECRET=your_very_secret_key
GLOBAL_PROXY_KEYS=your_proxy_key_1,proxy_key_2
EOF

# 使用环境变量文件启动
docker run --env-file .env api-key-rotator:latest
```


### 🔒 安全

- 所有代理请求需要 `X-Proxy-Key` 头部认证
- 管理界面需要用户名密码认证
- 生产环境中应妥善保护环境变量
- 数据库密码和API密钥应加密存储

### 📈 性能

- **轻量级**: < 50MB镜像大小，快速启动，资源占用少
- **企业级**: < 80MB镜像大小，高并发，可扩展架构
- **API响应**: 正常负载下大多数操作 < 100ms

### 🤝 贡献

1. Fork 本仓库
2. 创建功能分支
3. 进行更改
4. 如适用，添加测试
5. 提交拉取请求

### 📄 许可证

本项目采用 MIT 许可证 - 详情请参阅 [LICENSE](LICENSE) 文件。
