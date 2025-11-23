# API Key Rotator

[English](README.md) | [中文简体](README_CN.md)

一个企业级的API密钥管理和轮换服务，提供智能的密钥池管理、自动故障转移和负载均衡功能。

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
├── docker-compose.yml                # 轻量级部署
├── docker-compose.enterprise.yml     # 企业级部署
├── Dockerfile                        # 默认构建（轻量级）
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
- **前端**: Vue 3 + JavaScript + Element Plus + Vue Router + Vue I18n
- **数据库**: MySQL 8.0+（企业级）/ SQLite（轻量级）
- **缓存**: Redis 6.0+（企业级）/ 内存缓存（轻量级）
- **容器化**: Docker + Docker Compose
- **架构**: 接口抽象 + 适配器模式

### 🌐 API端点

启动服务后，可以访问以下API：

- **根路径**: `http://localhost:8000/` - 服务状态信息
- **管理API**: `http://localhost:8000/admin/*` - 后台管理接口
  - `GET /admin/app-config` - 获取应用配置
  - `POST /admin/login` - 用户登录
  - `GET/POST/PUT/DELETE /admin/proxy-configs` - 代理配置管理
  - `GET/POST/DELETE /admin/proxy-configs/:id/keys` - API密钥管理
  - `PATCH /admin/keys/:keyID` - 密钥状态管理
- **前端管理界面**: `http://localhost:8000/` - Vue3后台管理界面

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