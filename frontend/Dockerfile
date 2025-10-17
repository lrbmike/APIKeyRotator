# ---- Stage 1: Base ----
# 定义一个包含 Node.js 的基础镜像，用于共享
FROM node:22-alpine AS base

WORKDIR /app

# ---- Stage 2: Dependencies ----
# 专门用于安装依赖，以便在源码未改变时利用 Docker 缓存
FROM base AS dependencies

# 复制package文件
COPY package.json ./

# 安装依赖
RUN npm install

# ---- Stage 3: Development ----
# 开发环境阶段：基于已安装的依赖，准备运行开发服务器
FROM dependencies AS development

# 复制源代码
COPY . .

EXPOSE 5173

# 设置默认的开发启动命令
CMD ["npm", "run", "dev"]

# ---- Stage 4: Builder ----
# 构建阶段：用于生成生产环境的静态文件
FROM dependencies AS builder
COPY . .
RUN npm run build

# ---- Stage 5: Production ----
# 生产阶段：最终的、轻量的 Nginx 镜像
FROM nginx:1.25-alpine AS production
# 从 builder 阶段复制构建好的静态文件到 Nginx 的服务目录
COPY --from=builder /app/dist /usr/share/nginx/html
# 复制自定义的 Nginx 配置文件
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]