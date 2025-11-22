# 默认构建 - 轻量级版本 (SQLite + 内存缓存)
# 等同于: docker build -f Dockerfile.lightweight -t api-key-rotator:lightweight .
#
# 如需企业级版本，请使用：
# docker build -f Dockerfile.enterprise -t api-key-rotator:enterprise .

# ---- Stage 1: Build Backend ----
FROM golang:1.22 AS backend-builder
WORKDIR /app/backend

# 设置Go模块代理，加速依赖下载
ENV GOPROXY=https://goproxy.cn,direct

# 复制Go模块文件，利用Docker缓存
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# 复制后端源代码 - 只包含轻量级实现
COPY backend/internal/config ./internal/config
COPY backend/internal/infrastructure/database/sqlite ./internal/infrastructure/database/sqlite
COPY backend/internal/infrastructure/cache/memory ./internal/infrastructure/cache/memory
COPY backend/internal/models ./internal/models
COPY backend/internal/dto ./internal/dto
COPY backend/internal/logger ./internal/logger
COPY backend/internal/middleware ./internal/middleware
COPY backend/internal/router ./internal/router
COPY backend/internal/handlers ./internal/handlers
COPY backend/internal/services ./internal/services
COPY backend/internal/utils ./internal/utils
COPY backend/main.go .

# 静态编译应用 (CGO启用，静态链接所有库)
# 这对于在最小alpine镜像中运行CGO二进制文件至关重要
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags="-s -w -extldflags '-static'" \
    -o ./api-key-rotator .

# ---- Stage 2: Build Frontend ----
FROM node:22-alpine AS frontend-builder
WORKDIR /app/frontend

# 安装 pnpm
RUN npm install -g pnpm

# 复制前端依赖文件
COPY frontend/package.json frontend/pnpm-lock.yaml* ./
RUN pnpm install

# 复制前端源代码并构建
COPY frontend/ ./
RUN pnpm run build

# ---- Stage 3: Final Image ----
FROM alpine:3.22.0

# 安装必要的系统包
RUN apk add --no-cache ca-certificates tzdata sqlite && \
    addgroup -g apikeyrotator apikeyrotator && \
    adduser -D -s /bin/sh -G apikeyrotator apikeyrotator

# 创建应用目录
WORKDIR /app

# 复制后端二进制文件
COPY --from=backend-builder /app/backend/api-key-rotator ./

# 复制前端构建文件
COPY --from=frontend-builder /app/frontend/dist ./static

# 创建数据目录
RUN mkdir -p /app/data && \
    chown -R apikeyrotator:apikeyrotator /app /app/data

# 设置权限
RUN chmod +x ./api-key-rotator

# 暴露端口
EXPOSE 8000

# 设置环境变量默认值
ENV BACKEND_PORT=8000 \
    DB_TYPE=sqlite \
    CACHE_TYPE=memory \
    DATABASE_PATH=/app/data/api_key_rotator.db \
    LOG_LEVEL=info

# 切换到非特权用户
USER apikeyrotator

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ./api-key-rotator --health-check || exit 1

# 启动应用
CMD ["./api-key-rotator"]