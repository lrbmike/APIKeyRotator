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
# Use the standard golang image which includes the CGO toolchain.
# This is more efficient than installing build-base on an alpine image.
FROM golang:1.22 AS backend-builder
WORKDIR /app/backend

# Set Go proxy for faster dependency downloads
ENV GOPROXY=https://goproxy.cn,direct

# Copy go.mod and go.sum first to leverage Docker cache
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the backend source code
COPY backend/ ./

# Build the application with CGO enabled for SQLite, output to the current directory
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o ./api-key-rotator .

# ---- Stage 3: Final Image ----
# Use a lightweight alpine image for the final production stage
FROM alpine:3.20
WORKDIR /app

# Install necessary packages for the final image
RUN apk add --no-cache ca-certificates tzdata

# Copy backend binary from the backend-builder stage
COPY --from=backend-builder /app/backend/api-key-rotator /app/api-key-rotator

# Copy frontend build artifacts from the frontend-builder stage
COPY --from=frontend-builder /app/frontend/dist /app/public

# Set timezone
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Expose the application port
EXPOSE 8000

# Set the entrypoint for the container
CMD ["/app/api-key-rotator"]
