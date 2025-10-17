# API Key Rotator - Frontend Admin Interface

[English](README.md) | [中文简体](README_CN.md)

This is the frontend admin panel developed for the **API Key Rotator** project. It provides a clean and intuitive user interface for configuring all proxy services and managing keys.

<img width="2160" height="558" alt="api_key_rotator_frontend" src="https://github.com/user-attachments/assets/64d49739-0363-4266-a4dd-ba7162446394" />

## Tech Stack

*   **Framework**: [Vue 3](https://vuejs.org/) (using Composition API and `<script setup>`)
*   **Build Tool**: [Vite](https://vitejs.dev/)
*   **UI Component Library**: [Element Plus](https://element-plus.org/)
*   **Routing**: [Vue Router 4](https://router.vuejs.org/)
*   **HTTP Requests**: [Axios](https://axios-http.com/)

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

## Project Structure

The frontend code is located in the `frontend/` directory. Its core source code structure is as follows:

```
frontend/
└── src/
    ├── api/          # Stores all Axios API request functions for backend interaction.
    ├── components/   # Reusable UI components (e.g., KeyManager.vue).
    ├── router/       # Vue Router configuration, including the routing table and navigation guards.
    ├── views/        # Page-level components (e.g., Login.vue, Dashboard.vue, Layout.vue).
    ├── App.vue       # The root component of the Vue application.
    └── main.js       # The application's entry point for initializing Vue, Element Plus, and the router.
```

## Development and Deployment

This project is fully containerized. Both development and production deployment are managed through Docker and Docker Compose, eliminating the need to install a Node.js environment locally.

### Starting the Development Environment

Please refer to the instructions on local development in the **project root's** `README.md` or `backend/README.md`. The core steps are as follows:

1.  Create and configure the `.env` file in the project root.
2.  Run `docker-compose up --build` in the project root.

Docker Compose will handle everything automatically, including:
*   Building the frontend development image (`target: development`).
*   Installing all npm dependencies.
*   Starting the Vite development server.

After startup, you can access the frontend application at `http://localhost:5173`. Since the source code is mounted into the container, any changes to files in the `frontend/src` directory will trigger **hot reloading**.

### Production Deployment

Production deployment is also managed by Docker Compose (`docker-compose.prod.yml`).

When you run `docker-compose -f docker-compose.prod.yml up --build`, Docker executes the production build process in `frontend/Dockerfile`:
1.  **Build Stage (`builder`)**: In a temporary container, `npm run build` is executed to generate optimized static files in the `/app/dist` directory.
2.  **Production Stage (`production`)**:
    *   A very lightweight `nginx:alpine` image is used as the final image.
    *   All static files from the `/app/dist` directory of the previous stage are copied to the Nginx web root `/usr/share/nginx/html`.
    *   The project's `nginx.conf` file is copied into the container to handle API reverse proxying and Vue Router's history mode.

Ultimately, a highly optimized image containing only Nginx and static files is created and run, providing a high-performance and secure frontend service.

## Dockerfile Explained

The `frontend/Dockerfile` uses a **Multi-stage builds** strategy, which is a best practice that ensures both development convenience and a lightweight production image.

*   **`base` & `dependencies` Stages**: Base environment and dependency installation. This layer is shared by subsequent stages, effectively utilizing Docker's layer caching.
*   **`development` Stage**: Used for local development. It directly uses the result of the `dependencies` stage and runs `npm run dev` to start the Vite server. `docker-compose.yml` selects this stage.
*   **`builder` Stage**: Specifically for running `npm run build` to generate static files for the production environment.
*   **`production` Stage**: The final production image. It only copies the build artifacts from the `builder` stage and contains no Node.js, npm, or source code, resulting in a very small image size. `docker-compose.prod.yml` selects this stage.