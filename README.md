# 🚀 Short URL Generator

A simple, high-performance URL shortener service built with **Go** and **PostgreSQL**, fully containerized with **Docker**.

## 📖 Description
This application allows users to shorten long URLs into manageable 6-character codes. It features fast redirection, hit counting, and basic statistics for each shortened link. Designed with scalability and efficiency in mind, it leverages Go's concurrency model and PostgreSQL's reliability.

## 🛠 Tech Stack
- **Backend**: Go (Golang) 1.21
- **Database**: PostgreSQL 15
- **Infrastructure**: Docker & Docker Compose
- **Documentation**: OpenAPI/Swagger (coming soon)

## ✨ Features
- [x] **URL Shortening**: Generate unique 6-character alphanumeric codes.
- [x] **Fast Redirection**: Instantly redirect short codes to original long URLs.
- [x] **Hit Counter**: Track how many times each short URL has been clicked.
- [x] **Statistics API**: Retrieve metadata and usage stats for any short link.
- [x] **Containerized**: Easy setup with a single command.

## 🚀 How to Run

### Prerequisites
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Step 1: Clone the Repository
```bash
git clone https://github.com/marifyahya/shorturl-generator-app.git
cd shorturl-generator-app
```

### Step 2: Environment Setup
Copy the example environment file and adjust values if necessary (though the defaults work out of the box with Docker).
```bash
cp .env.example .env
```

### Step 3: Spin up with Docker Compose
```bash
docker compose up -d --build
```
This command will:
1. Start a **PostgreSQL** database container.
2. Build and start the **Go API** container.
3. Automatically handle database health checks before starting the API.

### Step 4: Verify
The API will be available at `http://localhost:8081` (default mapping from `docker-compose.yml`).

## 🛣 API Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/shorten` | Create a short URL from a long URL. |
| `GET` | `/:short_code` | Redirect to the original long URL. |
| `GET` | `/api/stats/:short_code` | Get statistics for a specific short URL. |

---
Built with ❤️ by M. Arif Yahya
