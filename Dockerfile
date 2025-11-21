# ---- Stage 1: Build Frontend ----
FROM node:22-alpine AS frontend-builder
WORKDIR /app/frontend

# Copy frontend package files and install dependencies
COPY frontend/package.json frontend/pnpm-lock.yaml* ./
RUN npm install -g pnpm && pnpm install

# Copy frontend source and build
COPY frontend/ ./
RUN pnpm build

# ---- Stage 2: Build Backend ----
FROM golang:1.22-alpine AS backend-builder
WORKDIR /app/backend

# Set Go proxy
ENV GOPROXY=https://goproxy.cn,direct

# Copy backend go.mod and go.sum and download dependencies
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source and build
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /api-key-rotator .

# ---- Stage 3: Final Image ----
FROM alpine:3.20
WORKDIR /app

# Install necessary packages
RUN apk add --no-cache ca-certificates tzdata

# Copy backend binary from backend-builder stage
COPY --from=backend-builder /api-key-rotator /app/api-key-rotator

# Copy frontend build artifacts from frontend-builder stage
COPY --from=frontend-builder /app/frontend/dist /app/public

# Copy other necessary files
COPY .env.example /app/.env.example
COPY backend/data /app/data

# Set timezone
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Expose port
EXPOSE 8000

# Set entrypoint
CMD ["/app/api-key-rotator"]
