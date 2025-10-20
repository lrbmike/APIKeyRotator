# API Key Rotator

[English](README.md) | [中文简体](README_CN.md)

## Introduction

**API Key Rotator** is a powerful and flexible API key management and request proxy solution built with Go (Gin). It is designed to centralize the management of all your third-party API keys and provide automatic rotation, load balancing, and secure isolation through a unified proxy endpoint.

Whether you need to provide high availability for traditional RESTful APIs or a unified, SDK-compatible access point for large model APIs like OpenAI, this project offers an elegant and scalable solution.

The project includes a high-performance **Go backend** and a simple, easy-to-use **Vue 3 admin panel**, with "one-click" deployment via Docker Compose.

## Core Features

*   **Centralized Key Management**: Manage API key pools for all services in a unified web interface.
*   **Dynamic Key Rotation**: Atomic rotation based on Redis to effectively distribute API request quotas.
*   **Type-Safe Proxies**:
    *   **Generic API Proxy (`/proxy`)**: Provides proxy services for any RESTful API.
    *   **LLM API Proxy (`/llm`)**: Offers native streaming support and an SDK-friendly `base_url` for large model APIs compatible with OpenAI's format. Supported providers include **OpenAI, Gemini, Anthropic**, etc.
*   **Highly Extensible Architecture**: The backend uses an adapter pattern, making it easy to extend support for new types of proxy services in the future.
*   **Secure Isolation**: All proxy requests are authenticated via global keys, with support for multiple keys to protect real backend keys from being exposed.
*   **Dockerized Deployment**: Provides a complete Docker Compose configuration for one-click startup of the backend, frontend, database, and Redis.

## Quick Start

This project is fully containerized, and it is recommended to use Docker Compose for one-click deployment and development.

### 1. Prerequisites

Ensure that [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/install/) are installed on your system.

### 2. Configure the Project

After cloning the project, create a `.env` file from the `.env.example` template in the project root directory.

```bash
# Copy the configuration file template
cp .env.example .env
```

Then, edit the `.env` file according to your needs, at least setting sensitive information such as the database password and administrator password.

#### Proxy Key Configuration

This project uses the `GLOBAL_PROXY_KEYS` environment variable to configure proxy authentication keys, supporting a single key or multiple keys:

1.  **Single Key**:
    ```bash
    GLOBAL_PROXY_KEYS=your_secret_key
    ```

2.  **Multiple Keys** (Recommended for multi-client scenarios):
    ```bash
    GLOBAL_PROXY_KEYS=key1,key2,key3
    ```

The multiple keys feature allows you to assign different authentication keys to different clients or services, improving security and management flexibility.

### 3. Start the Services

We provide standard Docker Compose configurations for development and production environments.

**Development Environment**
```bash
# Start with the development environment configuration
docker-compose up --build -d
```

**Production Environment**
```bash
# Start with the production environment configuration
docker-compose -f docker-compose.prod.yml up --build -d
```

#### Access URLs

**Development Environment** (with Vite and Hot Reload):
*   **Frontend Dev Server**: `http://localhost:5173`
*   **Backend API Root**: `http://localhost:8000/`

**Production Environment** (with Nginx):
*   **Web Application (Frontend + Backend API)**: `http://localhost` (or `http://localhost:80`, depending on your `.env` configuration)

## Local Development Without Docker (Optional)

If you prefer to run and debug the source code directly on your local machine without using Docker, you can follow these steps.

### 1. Prerequisites

*   Install [Node.js](https://nodejs.org/) (18+)
*   Install [Go](https://golang.org/) (1.21+)
*   Install and run **MySQL** and **Redis** services locally

### 2. Start the Backend Service

1.  **Enter the Go backend directory**
    ```bash
    cd backend/
    ```

2.  **Install dependencies**
    ```bash
    go mod download
    ```

3.  **Configure environment variables**
    Create a `.env` file in the project root (refer to `.env.example`) and configure the connection information for the database and Redis.

4.  **Start the backend server**
    ```bash
    go run main.go
    ```
    The service will run at `http://127.0.0.1:8000`.

### 3. Start the Frontend Service

1.  **Enter the frontend directory** (in another terminal)
    ```bash
    cd frontend/
    ```

2.  **Install dependencies**
    ```bash
    npm install
    ```

3.  **Start the frontend server**
    ```bash
    npm run dev
    ```
    Vite will automatically handle API proxying. The service will run at `http://localhost:5173`.

Now, you can access the admin panel at `http://localhost:5173`.

## Usage Example

### LLM API Proxy

Using the `openai` Python SDK as an example, combined with an `OpenRouter` model, you can use the proxy service by modifying the `base_url`.

```python
from openai import OpenAI

client = OpenAI(
  # Format: http://<PROXY_PUBLIC_BASE_URL>/llm/<Service Slug>
  base_url="http://PROXY_PUBLIC_BASE_URL/llm/openrouter-api",
  api_key="<GLOBAL_PROXY_KEY>",
)

completion = client.chat.completions.create(
  # Please refer to the specific provider's documentation for model names
  model="openai/gpt-4o",
  messages=[
    {
      "role": "user",
      "content": "What is the meaning of life?"
    }
  ]
)

print(completion.choices[0].message.content)
```

Where `PROXY_PUBLIC_BASE_URL` and `GLOBAL_PROXY_KEY` are the environment variables you configured in your `.env` file.

## Development Guide

If you want to dive deeper into the code, please refer to the following documents:

*   **[Backend Development Guide](./backend/README.md)**
*   **[Frontend Development Guide](./frontend/README.md)**