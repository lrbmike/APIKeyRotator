# API Key Rotator - Frontend Admin Interface

[English](README.md) | [中文简体](README_CN.md)

This is the frontend admin panel developed for the **API Key Rotator** project. It provides a clean and intuitive user interface for configuring all proxy services and managing keys.

<img width="2334" height="591" alt="frontend_en" src="https://github.com/user-attachments/assets/dcfd9580-60ea-475a-abd6-ccc07944a229" />

## Tech Stack

*   **Framework**: [Vue 3](https://vuejs.org/) (using Composition API and `<script setup>`)
*   **Build Tool**: [Vite](https://vitejs.dev/)
*   **UI Component Library**: [Element Plus](https://element-plus.org/)
*   **Routing**: [Vue Router 4](https://router.vuejs.org/)
*   **HTTP Requests**: [Axios](https://axios-http.com/)
*   **Internationalization (i18n)**: [Vue I18n](https://vue-i18n.intlify.dev/)

## Feature List

*   **User Authentication**: Login and logout functionality based on the `.env` configuration.
*   **Dashboard**: A unified dashboard that clearly displays all configured proxy services (Generic API and LLM API) in a table format.
*   **Service Configuration CRUD**:
    *   **Create**: Create new Generic API or LLM API services through a dynamic form.
    *   **Read**: Fetch and display the list of services from the backend in real-time.
    *   **Update**: Edit existing service configurations.
    *   **Status Toggle**: Quickly enable or disable a service with a security confirmation prompt.
*   **Key Management (CRUD)**:
    *   **View** the list of configured keys for a specific service (with sensitive information redacted).
    *   **Add** a new API Key.
    *   **Update** the enabled/disabled status of a key.
    *   **Delete** a specific API Key with a security confirmation prompt.
*   **One-Click Copy**: Quickly copy the invocation address of each proxy service to the clipboard.
*   **Multi-language Support**: Supports seamless switching between Chinese and English.

## Project Structure

The frontend code is located in the `frontend/` directory. Its core source code structure is as follows:

```
frontend/
└── src/
    ├── api/          # Stores all Axios API request functions for backend interaction.
    ├── components/   # Reusable UI components (e.g., KeyManager.vue, LangSwitcher.vue).
    ├── locales/      # Stores language files for internationalization (e.g., en.json, zh-CN.json).
    ├── router/       # Vue Router configuration, including the routing table and navigation guards.
    ├── views/        # Page-level components (e.g., Login.vue, Dashboard.vue, Layout.vue).
    ├── App.vue       # The root component of the Vue application.
    ├── i18n.js       # Configuration file for vue-i18n.
    └── main.js       # The application's entry point for initializing Vue, Element Plus, i18n, and the router.
```

## Development and Deployment

This project is fully containerized and offers multiple deployment methods.

### Method 1: Single Image Deployment (Recommended)

This frontend is part of a larger project that can be deployed as a single Docker image, where the Go backend serves the frontend files directly. This is the simplest and recommended method.

To use this method, please refer to the `Dockerfile` and instructions in the **project's root directory**. For detailed steps, see the **[Docker Deployment Guide](../DEPLOY_WITH_DOCKER.md)**.

### Method 2: Using Docker Compose (For local development or separate deployments)

This method is suitable for local development with hot-reloading or if you wish to deploy the frontend and backend in separate containers.

#### Starting the Development Environment

Please refer to the instructions in the project root's `README.md`. The core steps are:

1.  Create and configure the `.env` file in the project root.
2.  Run `docker-compose up --build` in the project root.

Docker Compose will automatically handle building the development image, installing dependencies, and starting the Vite server. You can access the application at `http://localhost:5173`, and any file changes will trigger hot reloading.

#####PduoDey

The production deployment with Docker Compose is defined in `docker-compose.prod.yml`. When you run `docker-compose -f docker-compose.prod.yml up --build`, it builds an optimized, lightweight Nginx image to serve the frontend static files.

## Dockerfile Explained

The `frontend/Dockerfile` uses a **Multi-stage builds** strategy, which is a best practice that ensures both development convenience and a lightweight production image.

*   **`base` & `dependencies` Stages**: Base environment and dependency installation. This layer is shared by subsequent stages, effectively utilizing Docker's layer caching.
*   **`development` Stage**: Used for local development. It directly uses the result of the `dependencies` stage and runs `npm run dev` to start the Vite server. `docker-compose.yml` selects this stage.
*   **`builder` Stage**: Specifically for running `npm run build` to generate static files for the production environment.
*   **`production` Stage**: The final production image. It only copies the build artifacts from the `builder` stage and contains no Node.js, npm, or source code, resulting in a very small image size. `docker-compose.prod.yml` selects this stage.
