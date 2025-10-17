# API Key Rotator - 前端管理界面

[English](README.md) | [中文简体](README_CN.md)

这是为 **API Key Rotator** 项目配套开发的前端管理后台。它提供了一个简洁、直观的用户界面，用于完成所有代理服务的配置和密钥管理。

<img width="2169" height="560" alt="frontend_cn" src="https://github.com/user-attachments/assets/0e55e83b-1630-4d60-a754-c899a8d07672" />

## 技术栈

*   **框架**: [Vue 3](https://vuejs.org/) (使用 Composition API 和 `<script setup>`)
*   **构建工具**: [Vite](https://vitejs.dev/)
*   **UI组件库**: [Element Plus](https://element-plus.org/)
*   **路由**: [Vue Router 4](https://router.vuejs.org/)
*   **HTTP请求**: [Axios](https://axios-http.com/)

## 功能列表

*   **用户认证**: 基于`.env`配置的后台登录与登出功能。
*   **仪表盘**: 统一的仪表盘，以表格形式清晰展示所有已配置的代理服务（通用API和LLM API）。
*   **服务配置CRUD**:
    *   **创建**: 通过动态表单创建新的通用API或LLM API服务。
    *   **读取**: 实时从后端获取并展示服务列表。
    *   **更新**: 编辑已存在的服务配置。
    *   **状态切换**: 快速启用或禁用某个服务，并有安全确认提示。
*   **密钥管理 (CRUD)**:
    *   为指定的服务**查看**已配置的密钥列表（脱敏显示）。
    *   **添加**新的API Key。
    *   **更新**密钥的启用/禁用状态。
    *   **删除**指定的API Key，并有安全确认提示。
*   **一键复制**: 快速复制每个代理服务的调用地址到剪贴板。

## 项目结构

前端代码位于`frontend/`目录下，其核心源代码结构如下：

```
frontend/
└── src/
    ├── api/          # 存放所有与后端交互的 Axios API 请求函数。
    ├── components/   # 可复用的UI组件 (如 KeyManager.vue)。
    ├── router/       # Vue Router配置，包括路由表和导航守卫。
    ├── views/        # 页面级组件 (如 Login.vue, Dashboard.vue, Layout.vue)。
    ├── App.vue       # Vue应用的根组件。
    └── main.js       # 应用的入口文件，用于初始化Vue、Element Plus和路由。
```

## 开发与部署

本项目已完全容器化，无论是开发还是生产部署，都通过 Docker 和 Docker Compose 进行管理，无需在本地安装 Node.js 环境。

### 启动开发环境

请参考**项目根目录**下的 `README.md` 或 `backend/README.md` 中关于本地开发的说明。核心步骤如下：

1.  在项目根目录创建并配置好 `.env` 文件。
2.  在项目根目录运行 `docker-compose up --build`。

Docker Compose 会自动完成所有工作，包括：
*   构建前端开发镜像 (`target: development`)。
*   安装所有 npm 依赖。
*   启动 Vite 开发服务器。

启动后，你可以通过 `http://localhost:5173` 访问前端应用。由于源码被挂载到容器中，任何对 `frontend/src` 目录下的文件修改都会触发**热重载**。

### 生产环境部署

生产环境的部署同样由 Docker Compose (`docker-compose.prod.yml`) 管理。

当执行 `docker-compose -f docker-compose.prod.yml up --build` 时，Docker 会执行 `frontend/Dockerfile` 中的生产构建流程：
1.  **构建阶段 (`builder`)**: 在一个临时的容器中，执行 `npm run build`，生成优化的静态文件到 `/app/dist` 目录。
2.  **生产阶段 (`production`)**:
    *   使用一个非常轻量的 `nginx:alpine` 镜像作为最终镜像。
    *   将上一个阶段生成的 `/app/dist` 目录下的所有静态文件，复制到 Nginx 的网站根目录 `/usr/share/nginx/html`。
    *   将项目中的 `nginx.conf` 配置文件复制到容器中，用于处理 API 反向代理和 Vue Router 的 history 模式。

最终，一个只包含 Nginx 和静态文件的、高度优化的镜像被创建并运行，提供了高性能和高安全性的前端服务。

## Dockerfile 解析

`frontend/Dockerfile` 采用了**多阶段构建 (Multi-stage builds)** 策略，这是一个最佳实践，可以同时保证开发时的便利性和生产镜像的轻量化。

*   **`base` & `dependencies` 阶段**: 基础环境和依赖安装。这一层被后续多个阶段共享，可以有效利用 Docker 的层缓存。
*   **`development` 阶段**: 用于本地开发。它直接使用 `dependencies` 阶段的结果，并运行 `npm run dev` 启动 Vite 服务器。`docker-compose.yml` 会选用这个阶段。
*   **`builder` 阶段**: 专门用于执行 `npm run build`，生成生产环境的静态文件。
*   **`production` 阶段**: 最终的生产镜像。它只从 `builder` 阶段拷贝构建产物，完全不包含 Node.js、npm 或源代码，因此镜像体积非常小。`docker-compose.prod.yml` 会选用这个阶段。
