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
├── README.md                  # 项目文档
└── internal/                  # 内部包
    ├── config/                # 配置管理
    │   ├── config.go          # 配置加载
    │   └── factory.go         # 基础设施工厂
    ├── infrastructure/        # 基础设施层
    │   ├── database/
    │   │   ├── interface.go   # 数据库仓库接口
    │   │   ├── sqlite/        # SQLite实现
    │   │   └── mysql/         # MySQL实现
    │   └── cache/
    │       ├── interface.go   # 缓存接口
    │       ├── memory/        # 内存缓存实现
    │       └── redis/         # Redis实现
    ├── adapters/              # LLM适配器
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

    需要配置数据库文件路径 `DATABASE_PATH` ，并保存目录存在，如：`DATABASE_PATH=./data/api_key_rotator.db`

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

本项目支持纯Docker构建。

### 构建镜像

```bash
# 轻量级版本 (SQLite + 内存缓存)
docker build -t api-key-rotator .

# 企业级版本 (MySQL + Redis)
docker build -f Dockerfile.enterprise -t api-key-rotator:enterprise .
```

### 使用Docker Compose

根据需求运行相应的compose文件：

#### 快速部署（推荐新手）
如果您想使用最简单的方式，可以直接切换到 `sqlite` 分支：
```bash
git checkout sqlite
docker-compose up -d
```
`sqlite` 分支是纯SQLite + 内存缓存版本，配置更简单，适合快速体验。

#### 当前分支部署
```bash
# 轻量级版本部署
docker-compose -f docker-compose.yml up -d

# 企业级版本部署
docker-compose -f docker-compose.enterprise.yml up -d
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

## 二次开发扩展

### 接入PostgreSQL示例

本项目采用**接口抽象架构**，可以轻松扩展支持新的数据库类型。以下是一个接入PostgreSQL的简明示例：

#### 1. 创建PostgreSQL实现

**目录结构：**
```
internal/infrastructure/database/
├── interface.go          # 现有接口定义
├── sqlite/               # SQLite实现
├── mysql/                # MySQL实现
└── postgres/             # PostgreSQL实现（新增）
    ├── manager.go        # PostgreSQL管理器
    └── repository.go     # PostgreSQL仓库实现
```

**核心代码示例：**

**postgres/manager.go**
```go
package postgres

import (
    "api-key-rotator/backend/internal/infrastructure/database"
)

type Manager struct {
    dsn string
    repo database.Repository
}

func NewPostgresManager(dsn string) *Manager {
    return &Manager{dsn: dsn}
}

func (m *Manager) Initialize() (database.Repository, error) {
    repo, err := NewPostgresRepository(m.dsn)
    if err != nil {
        return nil, err
    }
    m.repo = repo
    return repo, nil
}
```

**postgres/repository.go**
```go
package postgres

import (
    "api-key-rotator/backend/internal/models"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type Repository struct {
    db *gorm.DB
}

func NewPostgresRepository(dsn string) (*Repository, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    return &Repository{db: db}, nil
}

// 实现database.Repository接口的所有方法
func (r *Repository) GetDB() *gorm.DB { return r.db }
func (r *Repository) CreateProxyConfig(config *models.ProxyConfig) error {
    return r.db.Create(config).Error
}
// ... 其他方法类似SQLite实现
```

#### 2. 更新配置工厂

在 `internal/config/factory.go` 中添加PostgreSQL支持：
```go
// 在CreateDatabaseManager函数中添加PostgreSQL选项
if strings.Contains(os.Getenv("DATABASE_URL"), "postgres") {
    return postgres.NewPostgresManager(dsn), nil
}
```

#### 3. 添加依赖

在 `go.mod` 中添加：
```bash
go get gorm.io/driver/postgres
```

#### 4. 环境变量配置

```bash
# PostgreSQL连接配置
DATABASE_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
```

#### 5. 构建配置

在Dockerfile中添加PostgreSQL构建支持，类似于现有的MySQL和SQLite构建。

这样的扩展方式保持了现有架构的完整性，同时提供了灵活的数据库支持。

### 接入新的LLM适配器 (Adapters)

适配器负责处理与上游 LLM API 的直接通信，包括请求构建、认证和密钥轮询。核心代码位于 `backend/internal/adapters`。

#### 1. 定义适配器结构
实现 `LLMAdapter` 接口：

```go
type LLMAdapter interface {
    ProcessRequest() (*services.TargetRequest, error)
}
```

#### 2. 实现核心逻辑
继承 `BaseLLMAdapter` 以复用通用逻辑（如密钥轮询）：

```go
type XAIAdapter struct {
    *BaseLLMAdapter
}

func (a *XAIAdapter) ProcessRequest() (*services.TargetRequest, error) {
    // 1. 验证代理密钥 (Proxy Key)
    // 2. 轮询上游密钥 (Upstream Key)
    upstreamKey, err := a.RotateUpstreamKey()
    
    // 3. 构建请求 (Header处理，移除 Gzip 等)
    headers := utils.FilterRequestHeaders(a.c.Request.Header, []string{"authorization", "accept-encoding"})
    headers["Authorization"] = "Bearer " + upstreamKey
    
    // 4. 返回目标请求对象
    return &services.TargetRequest{...}, nil
}
```

### 扩展新的API格式 (Converters)

转换器负责处理客户端格式与后端 API 期望格式之间的相互转换。核心代码位于 `backend/internal/converters`。

#### 1. 注册新格式
在 `backend/internal/converters/formats/registry.go` 中注册新的格式标识：

```go
// 示例：添加新的格式标识
const (
    FormatClaudeNative = "claude_native"
)
```

#### 2. 实现格式处理器
在 `backend/internal/converters/formats/claude_native/` 下实现 `FormatHandler` 接口：

```go
type ClaudeNativeHandler struct{}

// 转换请求：将通用请求格式转换为 Claude 原生格式
func (h *ClaudeNativeHandler) BuildRequest(req *types.UniversalRequest) ([]byte, error) {
    // 实现转换逻辑...
}

// 转换响应：将 Claude 原生响应转换为通用响应格式
func (h *ClaudeNativeHandler) ParseResponse(body []byte) (*types.UniversalResponse, error) {
    // 实现解析逻辑...
}
```

#### 3. 处理流式响应 (可选)
如果需要支持流式传输，还需实现 `StreamHandler` 接口以处理 SSE (Server-Sent Events) 转换。这对于聊天交互体验至关重要。

```go
type StreamHandler interface {
    // 解析流分块：将特定格式的 SSE 数据行解析为通用流分块
    ParseStreamChunk(chunk []byte) (*UniversalStreamChunk, error)
    
    // 构建流分块：将通用流分块构建为特定格式的 SSE 数据行
    BuildStreamChunk(chunk *UniversalStreamChunk) ([]byte, error)
    
    // 构建起始事件 (如 OpenAI 的 role delta)
    BuildStartEvent(model string, id string) [][]byte
    
    // 构建结束事件 (如 [DONE])
    BuildEndEvent() [][]byte
}
```

实现示例：

```go
func (h *ClaudeStreamHandler) ParseStreamChunk(chunk []byte) (*UniversalStreamChunk, error) {
    // 1. 解析 SSE 行 (e.g., "data: {...}")
    // 2. 提取 content delta
    // 3. 构造 UniversalStreamChunk
    return &UniversalStreamChunk{
        Content: deltaContent,
        FinishReason: finishReason,
    }, nil
}
```