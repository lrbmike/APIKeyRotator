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

# 复制后端源代码 - 复制全部代码以避免导入错误
COPY backend/internal ./internal
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
    addgroup apikeyrotator && \
    adduser -D -s /bin/sh -G apikeyrotator apikeyrotator

# 创建应用目录
WORKDIR /app

# 复制后端二进制文件
COPY --from=backend-builder /app/backend/api-key-rotator ./

# 复制前端构建文件
COPY --from=frontend-builder /app/frontend/dist ./static

# 创建数据目录并设置权限
RUN mkdir -p /app/data && \
    chown -R apikeyrotator:apikeyrotator /app && \
    chmod -R 755 /app/data

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

# 在切换用户前创建数据库文件并设置权限
RUN touch /app/data/api_key_rotator.db && \
    chown -R apikeyrotator:apikeyrotator /app/data && \
    chmod -R 755 /app/data && \
    chmod 664 /app/data/api_key_rotator.db

# 切换到非特权用户 - 注释掉以避免权限问题
# USER apikeyrotator

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ./api-key-rotator --health-check || exit 1

# 创建启动脚本来确保权限正确
RUN echo '#!/bin/sh' > /app/entrypoint.sh && \
    echo 'echo "Setting up database permissions..."' >> /app/entrypoint.sh && \
    echo 'if [ -f "/app/data/api_key_rotator.db" ]; then' >> /app/entrypoint.sh && \
    echo '  chmod 664 /app/data/api_key_rotator.db' >> /app/entrypoint.sh && \
    echo '  echo "Database file permissions updated"' >> /app/entrypoint.sh && \
    echo 'else' >> /app/entrypoint.sh && \
    echo '  echo "Database file not found, will be created automatically"' >> /app/entrypoint.sh && \
    echo 'fi' >> /app/entrypoint.sh && \
    echo 'exec ./api-key-rotator "$@"' >> /app/entrypoint.sh && \
    chmod +x /app/entrypoint.sh

# 启动应用
CMD ["/app/entrypoint.sh"]