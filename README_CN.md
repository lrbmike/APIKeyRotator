# API Key Rotator

[English](README.md) | [中文简体](README_CN.md)

## 🚀 动态部署方案

**本项目现已支持统一镜像中的动态部署切换**，通过环境变量即可选择不同的部署方案：

### 两种部署方案

| 方案 | 数据库 | 缓存 | 适用场景 | QPS支持 |
|------|--------|------|----------|-------------|
| 🟢 **轻量级部署** | SQLite | 内存缓存 | 个人项目、小型应用 | < 5K |
| 🔴 **企业级部署** | MySQL | Redis | 企业应用、大型部署 | > 10K |

### 智能自动检测

系统根据环境变量自动选择部署方案：
- **检测到MySQL环境变量** (`DB_HOST`, `DB_USER` 等) → 自动使用MySQL
- **检测到Redis环境变量** (`REDIS_HOST`, `REDIS_PORT` 等) → 自动使用Redis
- **未检测到相关变量** → 默认使用SQLite + 内存缓存

### 📋 完整环境变量配置

#### 🔴 数据库配置（可选 - 不设置则默认使用SQLite）

```bash
# MySQL连接字符串
DATABASE_URL=mysql://user:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True&loc=Local

# 或分离的MySQL变量
DB_HOST=localhost
DB_USER=appdb
DB_PASSWORD=your_strong_password
DB_NAME=api_key_rotator
DB_PORT=3306

# SQLite路径（仅在SQLite模式时生效）
DATABASE_PATH=/app/data/api_key_rotator.db
```

#### 🟠 Redis配置（可选 - 不设置则默认使用内存缓存）

```bash
# 基础Redis配置
REDIS_HOST=localhost          # 启用Redis的必需变量
REDIS_PORT=6379               # 可选，默认6379
REDIS_PASSWORD=your_password   # 可选，默认空字符串
REDIS_URL=redis://localhost:6379/0  # 可选，另一种连接字符串
REDIS_DB=0                    # 可选，默认0
```

#### 🔧 必需配置（必须设置）

```bash
# 安全配置（必需）
ADMIN_PASSWORD=your_admin_password
JWT_SECRET=your_very_long_jwt_secret
GLOBAL_PROXY_KEYS=key1,key2,key3

# 服务配置（可选）
BACKEND_PORT=8000
PROXY_PUBLIC_BASE_URL=http://localhost:8000
LOG_LEVEL=info
RESET_DB_TABLES=false
```

### 快速部署示例

**🟢 轻量级部署（推荐用于小型项目）**
```bash
docker run -d \
  -p 8000:8000 \
  -e ADMIN_PASSWORD="your_password" \
  -e JWT_SECRET="your_jwt_secret" \
  -e GLOBAL_PROXY_KEYS="your_proxy_key" \
  -v $(pwd)/data:/app/data \
  api-key-rotator
```

**🔴 企业级部署**
```bash
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
  api-key-rotator
```

📖 **详细部署信息请看下文** 👇

---

## 🎯 单一统一代码库

**所有部署模式现在都在统一代码库中** - 只需通过环境变量控制部署模式，系统将自动配置。

## 项目简介

**API Key Rotator** 是一个基于 Go (Gin) 构建的强大而灵活的API密钥管理与请求代理解决方案。它旨在集中化管理您所有第三方API的密钥，并通过一个统一的代理入口，实现密钥的自动轮询、负载均衡和安全隔离。

无论是为传统的RESTful API提供高可用性，还是为OpenAI等大模型API提供统一的、兼容SDK的访问点，本项目都能提供优雅且可扩展的解决方案。

该项目包含一个高性能的 **Go 后端** 和一个简洁易用的 **Vue 3 管理后台**，并通过 Docker Compose 实现了"一键式"部署。

## ✨ 核心功能

*   **🔧 动态部署切换**: 单一代码库支持多种部署方案，通过环境变量智能选择数据库和缓存类型。
*   **🔑 集中化密钥管理**: 在Web界面统一管理所有服务的API密钥池。
*   **🔄 动态密钥轮询**: 基于缓存实现的原子性轮询，支持内存缓存和Redis，有效分摊API请求配额。
*   **🚀 类型安全的代理**:
    *   **通用API代理 (`/proxy`)**: 为任何RESTful API提供代理服务。
    *   **LLM API代理 (`/llm`)**: 为兼容OpenAI格式的大模型API提供原生流式支持和SDK友好的`base_url`。目前支持的接口格式包括 **OpenAI, Gemini, Anthropic** 等。
*   **🏗️ 高度可扩展架构**: 后端采用适配器模式，未来可轻松扩展支持任何新类型的代理服务。
*   **🛡️ 安全隔离**: 所有代理请求均通过全局密钥进行认证，支持配置多个密钥，保护后端真实密钥不被泄露。
*   **🐳 统一Docker化**: 单一镜像支持所有部署模式，Docker Compose一键部署。

## 🚀 快速开始

### 方式一：轻量级部署（推荐）

最简单的部署方式，只需设置必需的环境变量：

```bash
# 克隆项目
git clone https://github.com/your-repo/APIKeyRotator.git
cd APIKeyRotator

# 构建镜像
docker build -t api-key-rotator .

# 启动服务（SQLite + 内存缓存）
docker run -d \
  -p 8000:8000 \
  -e ADMIN_PASSWORD="your_password" \
  -e JWT_SECRET="your_jwt_secret" \
  -e GLOBAL_PROXY_KEYS="your_proxy_key" \
  -v $(pwd)/data:/app/data \
  api-key-rotator

# 访问应用
open http://localhost:8000
```

### 方式二：Docker Compose 部署

#### 1. 环境准备

确保您的系统中已经安装了 [Docker](https://www.docker.com/) 和 [Docker Compose](https://docs.docker.com/compose/install/)。

#### 2. 配置项目

```bash
# 克隆项目
git clone https://github.com/your-repo/APIKeyRotator.git
cd APIKeyRotator

# 复制配置文件模板
cp .env.example.cn .env
# 或使用英文版本: cp .env.example.en .env
```

#### 3. 代理密钥配置

本项目使用 `GLOBAL_PROXY_KEYS` 环境变量配置代理认证密钥，支持单个密钥或多个密钥：

1.  **单个密钥**:
    ```bash
    GLOBAL_PROXY_KEYS=your_secret_key
    ```

2.  **多个密钥** (推荐用于多客户端场景):
    ```bash
    GLOBAL_PROXY_KEYS=key1,key2,key3
    ```

#### 4. 启动服务

**🟢 轻量级模式 (默认)**
```bash
# 复制中文配置模板
cp .env.example.cn .env

# 根据需要编辑配置
nano .env

docker-compose up --build -d
```

**🔴 企业级模式**
```bash
# 复制中文配置模板
cp .env.example.cn .env

# 添加数据库和缓存配置
cat >> .env << EOF
DB_HOST=db
DB_USER=appdb
DB_PASSWORD=your_db_password
DB_NAME=api_key_rotator
REDIS_HOST=redis
REDIS_PASSWORD=your_redis_password
EOF

docker-compose -f docker-compose.prod.yml up --build -d
```

**或使用英文模板**:
```bash
# 复制英文配置模板
cp .env.example.en .env
# ... 同上操作
```

#### 5. 访问地址

**开发模式** (使用 Vite 和热重载):
*   **前端开发服务器**: `http://localhost:5173`
*   **后端 API 根路径**: `http://localhost:8000/`

**运行模式** (独立服务):
*   **Web 应用 (前端 + 后端 API)**: `http://localhost:8000`

## 非 Docker 本地开发 (可选)

如果你希望在不使用 Docker 的情况下，在本地直接运行和调试源代码，可以遵循以下步骤。

### 1. 环境准备

*   安装 [Node.js](https://nodejs.org/) (18+)
*   安装 [Go](https://golang.org/) (1.21+)
*   在本地安装并运行 **MySQL** 和 **Redis** 服务

### 2. 启动后端服务

1.  **进入Go后端目录**
    ```bash
    cd backend/
    ```

2.  **安装依赖**
    ```bash
    go mod download
    ```

3.  **配置环境变量**
    在项目根目录创建 `.env` 文件（参考 `.env.example`），并配置数据库和 Redis 的连接信息。

4.  **启动后端服务器**
    ```bash
    go run main.go
    ```
    服务将在 `http://127.0.0.1:8000` 上运行。

### 3. 启动前端服务

1.  **进入前端目录** (在另一个终端中)
    ```bash
    cd frontend/
    ```

2.  **安装依赖**
    ```bash
    npm install
    ```

3.  **启动前端服务器**
    ```bash
    npm run dev
    ```
    Vite 会自动处理 API 代理。服务将在 `http://localhost:5173` 上运行。

现在，你可以通过 `http://localhost:5173` 访问管理后台。

## 使用示例

### LLM API 代理

以 `openai` Python SDK 为例，结合使用 `OpenRouter` 模型，你可以通过修改 `base_url` 来使用本项目的代理服务。

```python
from openai import OpenAI

client = OpenAI(
  # 格式为 http://<PROXY_PUBLIC_BASE_URL>/llm/<服务标识 (Slug)>
  base_url="http://PROXY_PUBLIC_BASE_URL/llm/openrouter-api",
  api_key="<GLOBAL_PROXY_KEY>",
)

completion = client.chat.completions.create(
  # 模型名称请参考具体提供商的文档
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

其中 `PROXY_PUBLIC_BASE_URL` 和 `GLOBAL_PROXY_KEY` 是您在 `.env` 文件中配置的环境变量。

### 通用 API 代理

通用 API 代理可用于任何 RESTful API。以下是一个使用 Python requests 库调用天气 API 的示例：

```python
import requests

# 配置代理参数
proxy_url = "http://PROXY_PUBLIC_BASE_URL/proxy/weather/current"
proxy_key = "<GLOBAL_PROXY_KEY>"

# 查询参数
params = {
    "query": "London"
    # 在代理请求转发至目标 API 时，系统会轮询后台配置的真实 API 密钥，并将其拼接到原始授权参数 access_key（该参数由后台配置）中。
）中
}

# 设置请求头
headers = {
    "X-Proxy-Key": proxy_key
}

# 发起请求
response = requests.get(proxy_url, params=params, headers=headers)

# 处理响应
if response.status_code == 200:
    data = response.json()
    print(f"天气信息: {data}")
else:
    print(f"请求失败，状态码: {response.status_code}")
```

在这个示例中：
1. `weather` 是您在管理界面配置的服务标识 (Slug)
2. `current` 是目标API的路径
3. `PROXY_PUBLIC_BASE_URL` 是您的代理服务地址
4. `<GLOBAL_PROXY_KEY>` 是您配置的全局代理密钥之一

代理会自动将请求转发到配置的目标URL，并将路径和查询参数附加到目标地址上。

## 📚 技术特点

*   **🔧 智能配置检测**: 系统根据环境变量自动选择最适合的数据库和缓存方案
*   **⚡ 高性能架构**: 支持从轻量级到企业级的各种性能需求
*   **🎯 零配置启动**: 默认模式下无需任何数据库或缓存服务配置
*   **🔄 无缝升级**: 可在不同部署模式间无缝切换，无需修改代码
*   **🛡️ 生产就绪**: 包含健康检查、日志记录、错误处理等生产级特性

## 📖 相关文档

如果您希望深入代码功能，请参考以下文档：

*   **[后端开发指南](./backend/README.md)**
*   **[前端开发指南](./frontend/README.md)**

## 🔧 部署示例

### 🟢 轻量级部署

```bash
# SQLite + 内存缓存 - 简单高效
docker run -d \
  -p 8000:8000 \
  -e ADMIN_PASSWORD="your_password" \
  -e JWT_SECRET="your_jwt_secret" \
  -e GLOBAL_PROXY_KEYS="your_proxy_key" \
  -v $(pwd)/data:/app/data \
  api-key-rotator
```

### 🔴 企业级部署

```bash
# MySQL + Redis - 高性能和可扩展
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

### 🐳 Docker Compose 示例

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
      # 可选：企业级模式添加这些
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

## ❓ 常见问题

### Q: 如何选择合适的部署模式？
**A**:
- **小型项目/个人使用**: 使用轻量级模式 (SQLite + 内存缓存)
- **企业应用**: 使用企业级模式 (MySQL + Redis)

### Q: 如何查看当前使用的数据库和缓存类型？
**A**: 启动应用时会显示日志信息：
```
Database Type: sqlite
Cache Type: memory
```

### Q: 如何从轻量级模式升级到企业级模式？
**A**: 只需添加相应的环境变量即可，系统会自动检测并切换：
```bash
# 添加MySQL配置
DB_HOST=mysql-server
DB_USER=appdb
DB_PASSWORD=your_password

# 添加Redis配置
REDIS_HOST=redis-server
```

### Q: 数据迁移如何处理？
**A**: 系统启动时会自动创建表结构。要从SQLite迁移到MySQL：

1. **备份SQLite数据**:
   ```bash
   cp data/api_key_rotator.db backup_$(date +%Y%m%d).db
   ```

2. **添加MySQL环境变量**:
   ```bash
   -e DB_HOST="mysql-server" \
   -e DB_USER="appdb" \
   -e DB_PASSWORD="your_password" \
   -e DB_NAME="api_key_rotator"
   ```

3. **重启应用程序** - 系统会在MySQL中自动创建新表

数据导入需要您手动从SQLite导出并导入到MySQL，或使用迁移工具。

### Q: 支持分布式部署吗？
**A**: 是的，使用MySQL + Redis模式支持完全的分布式部署。
