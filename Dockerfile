# 统一镜像构建 - 支持多种数据库和缓存方案
# 通过环境变量动态切换: SQLite+内存缓存 (默认) 或 MySQL+Redis

# ---- Stage 1: Build Backend ----
# 使用标准 golang 镜像，已包含完整的 CGO 工具链，比 alpine 安装构建工具更高效
FROM golang:1.22 AS backend-builder
WORKDIR /app/backend

# 设置Go模块代理，加速依赖下载
ENV GOPROXY=https://goproxy.cn,direct

# 复制Go模块文件，利用Docker缓存
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# 复制后端源代码
COPY backend/ ./

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

# 复制前端依赖文件和源代码
COPY frontend/package.json ./
COPY frontend/ ./

# 安装依赖并构建
RUN pnpm install && pnpm build

# ---- Stage 3: Final Image ----
# 使用轻量级 alpine 镜像作为最终运行环境
FROM alpine:3.20
WORKDIR /app

# 安装运行时依赖 (不包含SQLite运行时，因为已静态编译)
RUN apk add --no-cache ca-certificates tzdata wget

# 创建应用和数据目录
RUN mkdir -p /app/data

# 从构建阶段复制后端二进制文件
COPY --from=backend-builder /app/backend/api-key-rotator /app/api-key-rotator

# 从构建阶段复制前端构建文件
COPY --from=frontend-builder /app/frontend/dist /app/public

# 暴露端口
EXPOSE 8000

# 设置时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8000/ || exit 1

# 启动命令
CMD ["./api-key-rotator"]