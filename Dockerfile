# 统一镜像构建 - 支持多种数据库和缓存方案
# 通过环境变量动态切换: SQLite+内存缓存 (默认) 或 MySQL+Redis

FROM golang:1.21-alpine AS builder

# 设置Go模块代理，加速依赖下载
ENV GOPROXY=https://goproxy.cn,direct

# 安装构建依赖 (SQLite需要CGO支持)
RUN apk add --no-cache gcc musl-dev sqlite-dev git

WORKDIR /app

# 复制Go模块文件
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# 复制后端源代码
COPY backend/ ./

# 构建应用 (启用CGO以支持SQLite)
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o ./api-key-rotator .

# 运行时镜像
FROM alpine:3.22.0

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata sqlite

# 创建应用目录和数据目录
RUN mkdir -p /app/data

# 复制可执行文件
COPY --from=builder /app/api-key-rotator /app/api-key-rotator

# 复制前端构建文件 (如果存在)
COPY frontend/dist/ /app/frontend/dist/ 2>/dev/null || true

WORKDIR /app

# 暴露端口
EXPOSE 8000

# 设置时区
ENV TZ=Asia/Shanghai
RUN cp /usr/share/zoneinfo/${TZ} /etc/localtime && echo "${TZ}" > /etc/timezone

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8000/ || exit 1

# 启动命令
CMD ["./api-key-rotator"]