# API Key Rotator - 前端管理界面

[English](README.md) | [中文简体](README_CN.md)

这是为 **API Key Rotator** 项目配套开发的前端管理后台。它提供了一个简洁、直观的用户界面，用于完成所有代理服务的配置和密钥管理。

<img width="2334" height="585" alt="frontend_zh" src="https://github.com/user-attachments/assets/2dc5ee98-7b7a-466d-b455-f2f97bb1328c" />

## 技术栈

*   **框架**: [Vue 3](https://vuejs.org/) (使用 Composition API 和 `<script setup>`)
*   **构建工具**: [Vite](https://vitejs.dev/)
*   **UI组件库**: [Element Plus](https://element-plus.org/)
*   **路由**: [Vue Router 4](https://router.vuejs.org/)
*   **HTTP请求**: [Axios](https://axios-http.com/)
*   **国际化 (i18n)**: [Vue I18n](https://vue-i18n.intlify.dev/)

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
*   **多语言支持**: 支持中英文无缝切换。

## 项目结构

前端代码位于`frontend/`目录下，其核心源代码结构如下：

```
frontend/
└── src/
    ├── api/          # 存放所有与后端交互的 Axios API 请求函数。
    ├── components/   # 可复用的UI组件 (如 KeyManager.vue, LangSwitcher.vue)。
    ├── locales/      # 存放国际化语言包文件 (如 en.json, zh-CN.json)。
    ├── router/       # Vue Router配置，包括路由表和导航守卫。
    ├── views/        # 页面级组件 (如 Login.vue, Dashboard.vue, Layout.vue)。
    ├── App.vue       # Vue应用的根组件。
    ├── i18n.js       # vue-i18n 的配置文件。
    └── main.js       # 应用的入口文件，用于初始化Vue、Element Plus、i18n和路由。
```

## 开发与部署

本项目已完全容器化，并提供多种部署方式。

### 方式一：单一镜像部署（推荐）

此前端项目是整个项目的一部分，可以作为一个单一的 Docker 镜像进行部署，其中 Go 后端直接提供前端文件的服务。这是最简单且推荐的方式。

要使用此方法，请参阅**项目根目录**下的 `Dockerfile` 和说明。详细步骤请参见 **[Docker 部署指南](../DEPLOY_WITH_DOCKER.md)**。

### 方式二：使用 Docker Compose（用于本地开发或独立部署）

此方法适用于支持热重载的本地开发，或者当您希望在不同的容器中独立部署前端和后端时。

#### 启动开发环境

请参考项目根目录 `README.md` 中的说明。核心步骤如下：

1.  在项目根目录创建并配置 `.env` 文件。
2.  在项目根目录运行 `docker-compose up --build`。

Docker Compose 会自动处理开发镜像的构建、依赖安装和 Vite 服务器的启动。您可以在 `http://localhost:5173` 访问应用，任何文件更改都会触发热重载。

#### 生产环境部署

使用 Docker Compose 的生产环境部署在 `docker-compose.prod.yml` 中定义。当您运行 `docker-compose -f docker-compose.prod.yml up --build` 时，它会构建一个优化的、轻量级的 Nginx 镜像来为前端静态文件提供服务。

## Dockerfile 解析

`frontend/Dockerfile` 采用了**多阶段构建 (Multi-stage builds)** 策略，这是一个最佳实践，可以同时保证开发时的便利性和生产镜像的轻量化。

*   **`base` & `dependencies` 阶段**: 基础环境和依赖安装。这一层被后续多个阶段共享，可以有效利用 Docker 的层缓存。
*   **`development` 阶段**: 用于本地开发。它直接使用 `dependencies` 阶段的结果，并运行 `npm run dev` 启动 Vite 服务器。`docker-compose.yml` 会选用这个阶段。
*   **`builder` 阶段**: 专门用于执行 `npm run build`，生成生产环境的静态文件。
*   **`production` 阶段**: 最终的生产镜像。它只从 `builder` 阶段拷贝构建产物，完全不包含 Node.js、npm 或源代码，因此镜像体积非常小。`docker-compose.prod.yml` 会选用这个阶段。
