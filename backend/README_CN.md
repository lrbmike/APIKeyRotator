# Go后端服务 - API Key Rotator

[English](README.md) | [中文简体](README_CN.md)

本项目是 `API Key Rotator` 的后端服务，基于 **Gin** 框架构建，提供API密钥管理和请求代理功能。

## 架构概览

Go后端采用**接口抽象架构 (Interface Abstraction Architecture)** 和**优化构建 (Optimized Builds)**，确保高度的可维护性和部署灵活性。

*   **接口抽象**: 业务逻辑与基础设施实现通过明确定义的接口进行清晰的分离
*   **基础设施适配器**: 为不同数据库（SQLite、MySQL）和缓存（内存、Redis）提供可插拔的实现
*   **优化构建**: 为轻量级和企业级场景分别构建，减少镜像大小和依赖
*   **模块化设计**: 每个组件都是独立的，可以轻松扩展或替换

## 项目结构

```
backend/
├── main.go                    # 应用入口点
├── go.mod                     # Go模块定义
├── go.sum                     # 依赖版本锁定
├── Dockerfile.lightweight     # 轻量级Docker构建
├── Dockerfile.enterprise      # 企业级Docker构建
├── README.md                  # 项目文档
└── internal/                  # 内部包
    ├── config/                # 配置管理
    │   ├── config.go          # 配置加载
    │   └── factory.go         # 基础设施工厂
    ├── infrastructure/        # 基础设施层 (NEW)
    │   ├── database/
    │   │   ├── interface.go   # 数据库仓库接口
    │   │   ├── sqlite/        # SQLite实现
    │   │   └── mysql/         # MySQL实现
    │   └── cache/
    │       ├── interface.go   # 缓存接口
    │       ├── memory/        # 内存缓存实现
    │       └── redis/         # Redis实现
    ├── adapters/              # LLM适配器 (需要接口更新)
    ├── handlers/              # HTTP处理器
    ├── services/              # 业务服务
    ├── models/                # 数据模型
    ├── dto/                   # 数据传输对象
    ├── logger/                # 日志配置
    ├── middleware/            # 中间件
    ├── router/                # 路由配置
    └── utils/                 # 工具函数
```

## 技术栈

*   **框架**: [Gin](https://gin-gonic.com/) - 高性能HTTP Web框架
*   **ORM**: [GORM](https://gorm.io/) - Go语言ORM库
*   **数据库**: MySQL 8.0+ (企业级) / SQLite (轻量级)
*   **缓存**: Redis 6.0+ (企业级) / 内存缓存 (轻量级)
*   **配置**: 环境变量 + [godotenv](https://github.com/joho/godotenv)
*   **容器化**: Docker + Docker Compose 优化构建
*   **架构**: 接口抽象 + 适配器模式

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

本项目支持通过根目录的Makefile进行优化构建。

### 构建镜像

```bash
# 构建轻量级版本 (SQLite + 内存缓存)
make build-lightweight

# 构建企业级版本 (MySQL + Redis)
make build-enterprise

# 构建所有版本
make build-all
```

### 使用Docker Compose

根据需求运行相应的compose文件：

```bash
# 轻量级部署
docker-compose -f docker-compose.yml up -d

# 企业级部署
docker-compose -f docker-compose.prod.yml up -d
```

### Docker镜像标签

* `api-key-rotator:lightweight` - ~50MB, SQLite + 内存缓存
* `api-key-rotator:enterprise` - ~80MB, MySQL + Redis (包含所有功能)

## 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/handlers

# 运行测试并显示覆盖率
go test -cover ./...
```