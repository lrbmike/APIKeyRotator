# Go后端服务 - API Key Rotator

本项目是 `API Key Rotator` 的后端服务，基于 **Gin** 框架构建，提供API密钥管理和请求代理功能。

## 架构概览

Go后端采用**清洁架构 (Clean Architecture)**设计模式，具有高度的规范性和扩展性。

*   **数据模型**: 定义在`internal/models/models.go`，使用一个统一的 `proxy_configs` 表来管理所有类型的代理配置，结构清晰，易于扩展。
*   **业务逻辑**: 核心代理逻辑被抽象到`internal/services/`和`internal/adapters/`中，实现了代码的高度复用和逻辑解耦。
*   **API路由**: 定义在`internal/handlers/`目录下，每个文件负责一个功能模块（管理、通用代理、LLM代理），职责清晰。

## 项目结构

```
backend/
├── main.go                    # 应用入口点
├── go.mod                     # Go模块定义
├── go.sum                     # 依赖版本锁定
├── Dockerfile                 # Docker构建文件
├── README.md                  # 项目文档
└── internal/                  # 内部包
    ├── config/                # 配置管理
    │   └── config.go
    ├── database/              # 数据库连接和迁移
    │   └── database.go
    ├── redis/                 # Redis连接
    │   └── redis.go
    ├── logger/                # 日志配置
    │   └── logger.go
    ├── models/                # 数据模型
    │   └── models.go
    ├── dto/                   # 数据传输对象
    │   └── dto.go
    ├── utils/                 # 工具函数
    │   └── utils.go
    ├── middleware/            # 中间件
    │   └── cors.go
    ├── router/                # 路由配置
    │   └── router.go
    ├── handlers/              # HTTP处理器
    │   ├── management.go      # 管理API
    │   ├── proxy.go          # 通用代理
    │   └── llm_proxy.go      # LLM代理
    ├── services/              # 业务服务
    │   └── proxy_handler.go
    └── adapters/              # LLM适配器
        ├── base_adapter.go
        ├── openai_adapter.go
        └── gemini_adapter.go
```

## 技术栈

*   **框架**: [Gin](https://gin-gonic.com/) - 高性能HTTP Web框架
*   **ORM**: [GORM](https://gorm.io/) - Go语言ORM库
*   **数据库**: MySQL 8.0+
*   **缓存**: Redis 6.0+
*   **配置**: 环境变量 + [godotenv](https://github.com/joho/godotenv)
*   **容器化**: Docker + Docker Compose

## 核心功能

*   **集中化密钥管理**: 在Web界面统一管理所有服务的API密钥池
*   **动态密钥轮询**: 基于Redis实现的原子性轮询，有效分摊API请求配额
*   **类型安全的代理**:
    *   **通用API代理 (`/proxy`)**: 为任何RESTful API提供代理服务
    *   **LLM API代理 (`/llm`)**: 为兼容OpenAI格式的大模型API提供原生流式支持
*   **高度可扩展架构**: 采用适配器模式，未来可轻松扩展支持任何新类型的LLM API
*   **安全隔离**: 所有代理请求均通过全局密钥进行认证，保护后端真实密钥不被泄露

## 本地开发

### 环境要求

*   Go 1.21+
*   MySQL 8.0+
*   Redis 6.0+

### 快速开始

1. **进入Go后端目录**
   ```bash
   cd backend
   ```

2. **安装依赖**
   ```bash
   go mod download
   ```

3. **配置环境变量**
   
   在项目根目录创建 `.env` 文件（参考 `.env.example`）：
   ```bash
   cp ../.env.example ../.env
   ```

4. **运行服务**
   ```bash
   go run main.go
   ```

   服务将在 `http://localhost:8000` 启动

### API文档

启动服务后，可以通过以下方式查看API：

*   **根路径**: `http://localhost:8000/` - 欢迎信息
*   **管理API**: `http://localhost:8000/api/admin/*` - 配置管理接口
*   **通用代理**: `http://localhost:8000/proxy/*` - 通用API代理
*   **LLM代理**: `http://localhost:8000/llm/*` - LLM API代理

## Docker部署

本项目通过根目录的 `docker-compose.yml` 文件进行部署，该文件已将此Go后端作为默认服务。

### 构建镜像

```bash
# 在项目根目录运行
docker-compose build backend
```

### 使用Docker Compose

在项目根目录运行 `docker-compose up` 即可启动所有服务。

## 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/handlers

# 运行测试并显示覆盖率
go test -cover ./...
```