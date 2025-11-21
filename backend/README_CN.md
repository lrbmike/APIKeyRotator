# Go后端服务 - API Key Rotator

[English](README.md) | [中文简体](README_CN.md)

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
├── build.sh                   # Linux编译脚本
├── build.bat                  # Windows编译脚本
├── README.md                  # 项目文档
└── internal/                  # 内部包
    ├── config/                # 配置管理
    │   └── config.go
    ├── database/              # 数据库连接和迁移
    │   └── database.go
    ├── cache/                 # 内存缓存
    │   └── cache.go
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
        ├── gemini_adapter.go
        └── anthropic_adapter.go
```

## 技术栈

*   **框架**: [Gin](https://gin-gonic.com/) - 高性能HTTP Web框架
*   **ORM**: [GORM](https://gorm.io/) - Go语言ORM库
*   **数据库**: SQLite 3 (使用 [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3))
*   **缓存**: 内存缓存 (自实现，线程安全)
*   **配置**: 环境变量 + [godotenv](https://github.com/joho/godotenv)
*   **容器化**: Docker + Docker Compose

## 核心功能

*   **集中化密钥管理**: 在Web界面统一管理所有服务的API密钥池
*   **动态密钥轮询**: 基于内存缓存实现的原子性轮询，有效分摊API请求配额
*   **类型安全的代理**:
    *   **通用API代理 (`/proxy`)**: 为任何RESTful API提供代理服务
    *   **LLM API代理 (`/llm`)**: 为兼容OpenAI格式的大模型API提供原生流式支持
*   **高度可扩展架构**: 采用适配器模式，未来可轻松扩展支持任何新类型的LLM API
*   **安全隔离**: 所有代理请求均通过全局密钥进行认证，保护后端真实密钥不被泄露
*   **轻量部署**: 单一可执行文件 + SQLite 数据库文件，无需额外服务

## 本地开发

### 环境要求

*   Go 1.21+
*   GCC 编译器（SQLite 需要 CGO 支持）
    *   **Windows**: 安装 [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) 或 [MinGW-w64](https://www.mingw-w64.org/)
    *   **Linux**: 通常已预装，如未安装运行 `sudo apt-get install build-essential` (Ubuntu/Debian)
    *   **验证**: 运行 `gcc --version` 确认安装成功

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
   
   主要配置项：
   ```env
   DATABASE_PATH=./data/api_key_rotator.db
   BACKEND_PORT=8000
   GLOBAL_PROXY_KEYS=your_secret_key
   ADMIN_USER=admin
   ADMIN_PASSWORD=your_password
   ```

4. **编译项目**

   **方法 1: 使用编译脚本（推荐）**
   ```bash
   # Windows
   build.bat
   
   # Linux/macOS
   chmod +x build.sh
   ./build.sh
   ```

   **方法 2: 手动编译**
   ```bash
   # Windows (PowerShell)
   $env:CGO_ENABLED=1
   go build -o api-key-rotator.exe .
   
   # Linux/macOS
   CGO_ENABLED=1 go build -o api-key-rotator .
   ```

5. **运行服务**
   ```bash
   # Windows
   .\api-key-rotator.exe
   
   # Linux/macOS
   ./api-key-rotator
   
   # 或直接运行（开发模式）
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

使用 Docker 可以避免 GCC 依赖问题，推荐用于生产环境。

### 构建镜像

```bash
# 在项目根目录运行
docker-compose build backend
```

### 使用Docker Compose

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f backend

# 停止服务
docker-compose down
```

## 数据备份

SQLite 数据库文件位于 `data/api_key_rotator.db`，定期备份此文件即可：

```bash
# Linux
cp data/api_key_rotator.db data/api_key_rotator.db.backup.$(date +%Y%m%d)

# Windows
copy data\api_key_rotator.db data\api_key_rotator.db.backup
```

## 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/handlers

# 运行测试并显示覆盖率
go test -cover ./...
```

## 常见问题

### Q: 编译时提示 "cgo: C compiler not found"
**A:** 需要安装 GCC 编译器，参考上面的"环境要求"部分

### Q: Windows 上找不到 gcc
**A:** 
1. 安装 TDM-GCC 或 MinGW-w64
2. 确保 gcc.exe 在系统 PATH 中
3. 重启命令行窗口

### Q: 数据库权限错误
**A:** 确保 `data` 目录存在且有写权限
```bash
mkdir -p data
chmod 755 data
```

## 相关文档

*   [快速开始指南](../QUICKSTART.md) - 详细的快速开始步骤
*   [部署说明](../DEPLOYMENT.md) - 完整的部署指南
*   [技术选型说明](../TECHNICAL_DECISIONS.md) - 技术决策说明
*   [迁移指南](../MIGRATION_SQLITE.md) - 从 MySQL+Redis 迁移的说明

## 性能特点

*   **轻量高效**: 单一可执行文件，内存占用低
*   **快速启动**: 无需等待外部服务启动
*   **适合场景**: 中小规模部署 (< 10000 QPS)
*   **易于备份**: 只需备份 SQLite 数据库文件

如需更高性能或分布式部署，可考虑切换回 MySQL + Redis。
