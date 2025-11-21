# API Key Rotator

[English](README.md) | [中文简体](README_CN.md)

## 分支说明

本项目提供两个主要分支，满足不同的部署需求：

- **`main` 分支**：基于 **MySQL + Redis** 的标准架构
  - 适合高并发、分布式部署场景
  - 需要外部 MySQL 和 Redis 服务
  - 支持横向扩展和集群部署
  
- **`sqlite` 分支**：轻量级单文件部署版本
  - 使用 **SQLite + 内存缓存**
  - 单一可执行文件，无需额外服务
  - 适合中小规模部署（< 10000 QPS）
  - 易于备份和迁移

**当前分支**: `sqlite` - 轻量级版本

切换分支：
```bash
# 切换到标准架构
git checkout main

# 切换到轻量级版本
git checkout sqlite
```

## 项目简介

**API Key Rotator** 是一个基于 Go (Gin) 构建的强大而灵活的API密钥管理与请求代理解决方案。它旨在集中化管理您所有第三方API的密钥，并通过一个统一的代理入口，实现密钥的自动轮询、负载均衡和安全隔离。

无论是为传统的RESTful API提供高可用性，还是为OpenAI等大模型API提供统一的、兼容SDK的访问点，本项目都能提供优雅且可扩展的解决方案。

该项目包含一个高性能的 **Go 后端** 和一个简洁易用的 **Vue 3 管理后台**，并通过 Docker Compose 实现了"一键式"部署。

## 核心功能

*   **集中化密钥管理**: 在Web界面统一管理所有服务的API密钥池。
*   **动态密钥轮询**: 基于内存缓存实现的原子性轮询，有效分摊API请求配额。
*   **类型安全的代理**:
    *   **通用API代理 (`/proxy`)**: 为任何RESTful API提供代理服务。
    *   **LLM API代理 (`/llm`)**: 为兼容OpenAI格式的大模型API提供原生流式支持和SDK友好的`base_url`。目前支持的接口格式包括 **OpenAI, Gemini, Anthropic** 等。
*   **高度可扩展架构**: 后端采用适配器模式，未来可轻松扩展支持任何新类型的代理服务。
*   **安全隔离**: 所有代理请求均通过全局密钥进行认证，支持配置多个密钥，保护后端真实密钥不被泄露。
*   **轻量部署**: 使用 SQLite 数据库和内存缓存，单一可执行文件即可运行，无需额外服务。
*   **Docker化部署**: 提供完整的 Docker Compose 配置，一键启动所有服务。

## 快速开始

本项目支持多种容器化方式，您可以根据需要选择最适合您的方式。

### 方式一：单一 Docker 镜像部署（推荐用于 Render 等 PaaS 平台）

我们在项目根目录提供了一个 `Dockerfile`，它能将前端和后端打包成一个单一的镜像。这是最简单的部署方式，强烈推荐在云平台上使用。

1.  **构建镜像**:
    ```bash
    docker build -t api-key-rotator .
    ```

2.  **运行容器**:
    ```bash
    docker run -d -p 8000:8000 --name api-key-rotator-app -v $(pwd)/backend/data:/app/data api-key-rotator
    ```

运行后，您可以通过 `http://localhost:8000` 访问应用。

更详细的说明，包括如何在 Render 上进行部署，请参阅我们的 **[Docker 部署指南](./DEPLOY_WITH_DOCKER.md)**。

### 方式二：使用 Docker Compose（适用于本地开发和多容器环境）

这种方式适合本地开发，因为它支持热重载。

#### 1. 环境准备

确保您的系统中已经安装了 [Docker](https://www.docker.com/) 和 [Docker Compose](https://docs.docker.com/compose/install/)。

### 2. 配置项目

在运行应用之前，您必须配置好环境变量。

1.  **从模板创建 `.env` 文件**:
    ```bash
    cp .env.example .env
    ```

2.  **编辑 `.env` 文件**: 您 **必须** 为必填项 (`ADMIN_PASSWORD`, `JWT_SECRET`, `GLOBAL_PROXY_KEYS`) 提供值。可选项可以保留不变，以使用其默认值。

主要配置项：
```env
DATABASE_PATH=./data/api_key_rotator.db
BACKEND_PORT=8000
GLOBAL_PROXY_KEYS=your_secret_key
ADMIN_USER=admin
ADMIN_PASSWORD=your_password
PROXY_PUBLIC_BASE_URL=http://localhost:8000
```

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
*   安装 **GCC 编译器** (SQLite 需要 CGO 支持)
    *   **Windows**: 安装 [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) 或 [MinGW-w64](https://www.mingw-w64.org/)
    *   **Linux**: 通常已预装，如未安装运行 `sudo apt-get install build-essential` (Ubuntu/Debian)
    *   **验证**: 运行 `gcc --version` 确认安装成功

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
    在项目根目录创建 `.env` 文件（参考 `.env.example`），并配置必要的环境变量。

4.  **编译项目**
    ```bash
    # Windows (PowerShell)
    $env:CGO_ENABLED=1
    go build -o api-key-rotator.exe .
    
    # Linux/macOS
    CGO_ENABLED=1 go build -o api-key-rotator .
    
    # 或使用编译脚本
    # Windows: build.bat
    # Linux: ./build.sh
    ```

5.  **启动后端服务器**
    ```bash
    # Windows
    .\api-key-rotator.exe
    
    # Linux/macOS
    ./api-key-rotator
    
    # 或直接运行（开发模式）
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
    # 或使用 pnpm
    pnpm install
    ```

3.  **启动前端服务器**
    ```bash
    npm run dev
    # 或
    pnpm dev
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
    # 在代理请求转发至目标 API 时，系统会轮询后台配置的真实 API 密钥，
    # 并将其拼接到原始授权参数（该参数由后台配置）中
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

## 技术特点

*   **轻量高效**: 使用 SQLite 和内存缓存，单一可执行文件即可运行
*   **快速启动**: 无需等待外部数据库和缓存服务启动
*   **适合场景**: 中小规模部署（单机 < 10000 QPS）
*   **易于备份**: 只需备份 SQLite 数据库文件 (`data/api_key_rotator.db`)
*   **Docker 友好**: 完全容器化，生产环境零配置

## 开发指南

如果您希望深入代码功能，请参考以下文档：

*   **[后端开发指南](./backend/README.md)** - Go 后端详细说明
*   **[前端开发指南](./frontend/README.md)** - Vue 3 前端详细说明
*   **[快速开始指南](./QUICKSTART.md)** - 详细的快速开始步骤
*   **[技术选型说明](./TECHNICAL_DECISIONS.md)** - 技术决策说明

## 部署说明

### 环境准备

在开始之前，请确保您的本地环境中已经安装了 [Docker](https://docs.docker.com/get-docker/)。您可以从 Docker 官网下载并安装适合您操作系统的版本。

安装完成后，您可以通过以下命令验证 Docker 是否安装成功：

```bash
docker --version
```

### 构建 Docker 镜像

我们推荐使用项目根目录下的 `Dockerfile` 来构建一个包含前端和后端的统一镜像。这个镜像会由 Go 后端来提供前端静态文件的服务。

在项目的根目录下，打开终端并运行以下命令：

```bash
docker build -t api-key-rotator .
```

这个命令会执行以下操作：
- `-t api-key-rotator`：为您的镜像指定一个名称（tag），这里我们将其命名为 `api-key-rotator`。
- `.`：告诉 Docker 在当前目录下查找 `Dockerfile`。

构建过程可能需要几分钟时间，因为它需要下载依赖、编译代码。

### 运行 Docker 容器

您可以通过在命令行中直接传递所需的环境变量来运行容器。这种方法使得基本配置更加明确。

```bash
docker run -d -p 8000:8000 --name api-key-rotator-app \
  -e ADMIN_PASSWORD="your_strong_password" \
  -e JWT_SECRET="your_very_secret_and_random_jwt_key" \
  -e GLOBAL_PROXY_KEYS="your_secure_global_proxy_key" \
  -v $(pwd)/backend/data:/app/data \
  api-key-rotator
```

命令解释:
- `-d`：以后台模式运行容器。
- `-p 8000:8000`：将容器的 `8000` 端口映射到您主机的 `8000` 端口。
- `--name api-key-rotator-app`：为您的容器指定一个名称。
- `-e VAR="value"`：在容器内设置一个环境变量。您 **必须** 为 `ADMIN_PASSWORD`, `JWT_SECRET`, 和 `GLOBAL_PROXY_KEYS` 提供值。
- `-v $(pwd)/backend/data:/app/data`：将本地的 `backend/data` 目录挂载到容器的 `/app/data` 路径。**这对数据持久化至关重要。** 如果没有这个卷挂载，您保存的任何数据（如新的API密钥或配置）都将在容器被移除时永久丢失。
- `api-key-rotator`：指定要运行的镜像。

容器启动后，您就可以在浏览器中打开 `http://localhost:8000` 来访问您的应用了。

### 访问地址说明

当您通过 `http://localhost:8000` 访问时：

- **前端应用**：直接在浏览器中打开 `http://localhost:8000`，您将看到的是 Vue.js 构建的前端用户界面。
- **后端 API**：所有的后端 API 服务都通过 `/api` 前缀来访问。例如：
  - 登录接口：`http://localhost:8000/api/admin/login`
  - 获取配置接口：`http://localhost:8000/api/admin/proxy-configs`
  - 代理服务接口：`http://localhost:8000/api/proxy/...`

Go 后端同时承担了 Web 服务器和 API 服务器的角色。

### 在 Render 上部署

[Render](https://render.com/) 提供了对 Docker 部署的良好支持。您可以按照以下步骤在 Render 上部署您的项目：

1.  **将代码推送到 GitHub/GitLab**：确保您所有的代码，包括我们新创建的 `Dockerfile`，都已提交并推送到代码仓库。

2.  **在 Render 创建新服务**：
    - 登录 Render，点击 "New" -> "Web Service"。
    - 连接您的 GitHub 或 GitLab 账号，并选择您的项目仓库。

3.  **配置服务**：
    - **Environment**：选择 `Docker`。
    - **Name**：为您的服务取一个名字。
    - **Root Directory**：如果您的 `Dockerfile` 不在根目录，需要指定路径。在我们的项目中，它就在根目录，所以留空即可。
    - **Port**：在 "Advanced" 设置中，确保 "Port" 设置为 `8000`，这与我们在 `Dockerfile` 中暴露的端口一致。
    - **Persistent Storage** (可选，但推荐)：为了持久化您的 SQLite 数据库，您可以添加一个 "Disk"：
        - **Mount Path**：设置为 `/app/data`。
        - **Size**：根据您的需要选择磁盘大小。

4.  **添加环境变量**：如果您的应用需要环境变量（例如，在 `.env` 文件中定义的变量），您需要在 Render 的 "Environment" 标签页中进行配置。

5.  **部署**：点击 "Create Web Service"，Render 将会自动从您的仓库拉取代码，使用 `Dockerfile` 构建镜像，并部署您的应用。

部署完成后，Render 会为您提供一个公开的 URL，您可以通过这个 URL 访问您的应用。

### (可选) 单独构建前后端

如果您希望单独构建和运行前端或后端，您依然可以使用 `frontend/Dockerfile` 和 `backend/Dockerfile`。

- **构建前端**：
  ```bash
  docker build -t frontend-app -f frontend/Dockerfile .
  ```
- **构建后端**：
  ```bash
  docker build -t backend-app -f backend/Dockerfile .
  ```

这种方式更适合在开发或需要将前后端分离部署的场景。

## 常见问题

### Q: 编译时提示 "cgo: C compiler not found"
**A:** 需要安装 GCC 编译器，Windows 用户安装 TDM-GCC，Linux 用户安装 build-essential

### Q: 如何备份数据？
**A:** 只需定期备份 `backend/data/api_key_rotator.db` 文件即可

### Q: 支持分布式部署吗？
**A:** 当前版本使用内存缓存，适合单机部署。如需分布式部署，可考虑切换回 Redis

### Q: Docker 部署是否需要 GCC？
**A:** 不需要，Docker 镜像内已包含所有必要的编译环境

## License

本项目采用 MIT 许可证。详情请参阅 [LICENSE](LICENSE) 文件。
