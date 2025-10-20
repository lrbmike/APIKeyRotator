# API Key Rotator

[English](README.md) | [中文简体](README_CN.md)

## 项目简介

**API Key Rotator** 是一个基于 Go (Gin) 构建的强大而灵活的API密钥管理与请求代理解决方案。它旨在集中化管理您所有第三方API的密钥，并通过一个统一的代理入口，实现密钥的自动轮询、负载均衡和安全隔离。

无论是为传统的RESTful API提供高可用性，还是为OpenAI等大模型API提供统一的、兼容SDK的访问点，本项目都能提供优雅且可扩展的解决方案。

该项目包含一个高性能的 **Go 后端** 和一个简洁易用的 **Vue 3 管理后台**，并通过 Docker Compose 实现了"一键式"部署。

## 核心功能

*   **集中化密钥管理**: 在Web界面统一管理所有服务的API密钥池。
*   **动态密钥轮询**: 基于Redis实现的原子性轮询，有效分摊API请求配额。
*   **类型安全的代理**:
    *   **通用API代理 (`/proxy`)**: 为任何RESTful API提供代理服务。
    *   **LLM API代理 (`/llm`)**: 为兼容OpenAI格式的大模型API提供原生流式支持和SDK友好的`base_url`。目前支持的接口格式包括 **OpenAI, Gemini, Anthropic** 等。
*   **高度可扩展架构**: 后端采用适配器模式，未来可轻松扩展支持任何新类型的代理服务。
*   **安全隔离**: 所有代理请求均通过全局密钥进行认证，支持配置多个密钥，保护后端真实密钥不被泄露。
*   **Docker化部署**: 提供完整的 Docker Compose 配置，一键启动后端、前端、数据库和 Redis。

## 快速开始

本项目已完全容器化，推荐使用 Docker Compose 进行一键部署和开发。

### 1. 环境准备

确保您的系统中已经安装了 [Docker](https://www.docker.com/) 和 [Docker Compose](https://docs.docker.com/compose/install/)。

### 2. 配置项目

克隆本项目后，在项目根目录下，从 `.env.example` 模板创建一个 `.env` 文件。

```bash
# 复制配置文件模板
cp .env.example .env
```

然后，根据你的需要编辑 `.env` 文件，至少需要设置数据库密码和管理员密码等敏感信息。

#### 代理密钥配置

本项目使用 `GLOBAL_PROXY_KEYS` 环境变量配置代理认证密钥，支持单个密钥或多个密钥：

1.  **单个密钥**:
    ```bash
    GLOBAL_PROXY_KEYS=your_secret_key
    ```

2.  **多个密钥** (推荐用于多客户端场景):
    ```bash
    GLOBAL_PROXY_KEYS=key1,key2,key3
    ```

多个密钥功能允许您为不同的客户端或服务分配不同的认证密钥，提高安全性和管理灵活性。

### 3. 启动服务

我们提供了标准的 Docker Compose 配置，支持开发和生产环境。

**开发环境**
```bash
# 使用开发环境配置启动
docker-compose up --build -d
```

**生产环境**
```bash
# 使用生产环境配置启动
docker-compose -f docker-compose.prod.yml up --build -d
```

#### 访问地址

**开发环境** (使用 Vite 和热重载):
*   **前端开发服务器**: `http://localhost:5173`
*   **后端 API 根路径**: `http://localhost:8000/`

**生产环境** (使用 Nginx):
*   **Web 应用 (前端 + 后端 API)**: `http://localhost` (或 `http://localhost:80`，取决于你的 `.env` 配置)

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

## 开发指南

如果您希望深入代码功能，请参考以下文档：

*   **[后端开发指南](./backend/README.md)**
*   **[前端开发指南](./frontend/README.md)**