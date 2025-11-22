# API Key Rotator

[English](README.md) | [中文简体](README_CN.md)

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

#### 轻量级构建（默认）
```bash
# 构建轻量级版本
make build-lightweight

# 运行默认的 SQLite + 内存缓存
docker-compose up -d
```

#### 企业级构建
```bash
# 构建企业级版本
make build-enterprise

# 运行 MySQL + Redis
docker-compose -f docker-compose.prod.yml up -d
```

### 📋 环境变量

#### 数据库配置
```bash
# SQLite（轻量级 - 默认）
DATABASE_PATH=/app/data/api_key_rotator.db

# MySQL（企业级）
DB_HOST=localhost
DB_USER=appdb
DB_PASSWORD=your_strong_password
DB_NAME=api_key_rotator
DB_PORT=3306

# 或使用连接字符串
DATABASE_URL=mysql://user:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True&loc=Local
```

#### 缓存配置
```bash
# 内存缓存（轻量级 - 默认）
# 无需额外配置

# Redis（企业级）
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your_password
REDIS_URL=redis://localhost:6379/0
```

#### 应用配置
```bash
# 服务器
BACKEND_PORT=8000
LOG_LEVEL=info

# 认证
ADMIN_USERNAME=admin
ADMIN_PASSWORD=your_admin_password_here
JWT_SECRET=your_very_secret_and_random_jwt_key

# 代理
GLOBAL_PROXY_KEYS=your_secure_global_proxy_key
PROXY_TIMEOUT=30
PROXY_PUBLIC_BASE_URL=http://localhost:8000

# 数据库重置选项
RESET_DB_TABLES=false
```

### 🏗️ 项目结构

```
api-key-rotator/
├── Makefile                          # 构建编排
├── docker-compose.yml                # 轻量级部署
├── docker-compose.prod.yml           # 企业级部署
├── Dockerfile.lightweight            # 轻量级构建
├── Dockerfile.enterprise             # 企业级构建
├── README.md                         # 项目文档
└── backend/                          # Go后端服务
    ├── main.go                       # 应用入口点
    ├── go.mod                        # Go模块定义
    └── internal/                      # 内部包
        ├── config/                    # 配置管理
        │   ├── config.go              # 配置加载
        │   └── factory.go             # 基础设施工厂
        ├── infrastructure/            # 基础设施层
        │   ├── database/
        │   │   ├── interface.go        # 数据库仓库接口
        │   │   ├── sqlite/             # SQLite实现
        │   │   └── mysql/              # MySQL实现
        │   └── cache/
        │       ├── interface.go        # 缓存接口
        │       ├── memory/             # 内存缓存实现
        │       └── redis/              # Redis实现
        ├── handlers/                  # HTTP处理器
        ├── models/                    # 数据模型
        ├── dto/                       # 数据传输对象
        ├── router/                    # 路由配置
        └── logger/                    # 日志配置
└── frontend/                         # Vue.js前端
    ├── src/                          # 源代码
    ├── package.json                  # 依赖
    └── Dockerfile                    # 前端构建
```

### 🛠️ 技术栈

- **后端**: Go + Gin框架 + GORM ORM
- **前端**: Vue 3 + TypeScript + Element Plus
- **数据库**: MySQL 8.0+（企业级）/ SQLite（轻量级）
- **缓存**: Redis 6.0+（企业级）/ 内存缓存（轻量级）
- **容器化**: Docker + Docker Compose
- **架构**: 接口抽象 + 适配器模式

### 🌐 API端点

启动服务后，可以访问以下API：

- **根路径**: `http://localhost:8000/` - 欢迎信息
- **管理API**: `http://localhost:8000/admin/*` - 配置管理
- **通用代理**: `http://localhost:8000/proxy/*` - 通用API代理（即将推出）
- **LLM代理**: `http://localhost:8000/llm/*` - LLM API代理（即将推出）

### 📦 构建镜像

#### 选项 1：使用默认构建（轻量级）
```bash
# 构建轻量级版本（默认）
docker build -t api-key-rotator .

# 使用自定义标签构建
docker build -t my-api-key-rotator:latest .
```

#### 选项 2：使用 Makefile（推荐）
```bash
# 构建轻量级版本
make build-lightweight

# 构建企业级版本
make build-enterprise

# 构建所有版本
make build-all
```

#### 选项 3：直接指定 Dockerfile
```bash
# 轻量级构建（SQLite + 内存缓存）
docker build -f Dockerfile.lightweight -t api-key-rotator:lightweight .

# 企业级构建（MySQL + Redis）
docker build -f Dockerfile.enterprise -t api-key-rotator:enterprise .
```

### 🐳 Docker部署

#### 轻量级部署
```bash
docker-compose up -d
```

#### 企业级部署
```bash
docker-compose -f docker-compose.prod.yml up -d
```

### 🧪 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/handlers

# 运行测试并显示覆盖率
go test -cover ./...
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