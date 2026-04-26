# 🚀 Short URL Generator

A simple, high-performance URL shortener service built with **Go** and **PostgreSQL**, fully containerized with **Docker**.

## 📖 Description
This application allows users to shorten long URLs into manageable 6-character codes. It features fast redirection, hit counting, and basic statistics for each shortened link. Designed with scalability and efficiency in mind, it leverages Go's concurrency model and PostgreSQL's reliability.

## 🛠 Tech Stack
- **Backend**: Go (Golang) 1.24
- **Database**: PostgreSQL 15
- **Infrastructure**: Docker & Docker Compose
- **Migrations**: golang-migrate
- **Documentation**: OpenAPI 3.0 (JSON)

## ✨ Features
- [x] **URL Shortening**: Generate unique 6-character alphanumeric codes.
- [x] **Fast Redirection**: Instantly redirect short codes to original long URLs.
- [x] **Hit Counter**: Track how many times each short URL has been clicked.
- [x] **Statistics API**: Retrieve metadata and usage stats for any short link.
- [x] **Auto Migrations**: Automatic database schema management with rollback support.
- [x] **Unit Testing**: Comprehensive tests for service and handler layers.
- [x] **OpenAPI 3.0**: Full API documentation in JSON format.
- [x] **Containerized**: Multi-stage Docker builds for optimized images.

## 🚀 Getting Started

### Prerequisites
- [Docker](https://www.docker.com/get-started) & [Docker Compose](https://docs.docker.com/compose/install/)
- (Optional for local) [Go 1.24+](https://go.dev/dl/)
- (Optional for local) [PostgreSQL 15](https://www.postgresql.org/download/)

### Step 1: Clone the Repository
```bash
git clone https://github.com/marifyahya/shorturl-generator-app.git
cd shorturl-generator-app
```

### Step 2: Environment Setup
```bash
cp .env.example .env
```

---

## 🐳 Running with Docker (Recommended)
```bash
docker compose up -d --build
```
The API will be available at `http://localhost:8081`.

---

## 💻 Running Locally (Development)

### 1. Start Database
Ensure you have a PostgreSQL instance running and update your `.env` file with the correct credentials.

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Run Migrations
The application runs migrations automatically on startup. You can also control them via flags:
```bash
# Run migrations UP (default)
go run cmd/api/main.go

# Rollback migrations (1 step)
go run cmd/api/main.go -migrate=down
```

### 4. Start Server
```bash
go run cmd/api/main.go
```
The API will be available at the port defined in `SERVER_PORT` (default: `http://localhost:8080`).

---

## 🛣 API Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/shorten` | Create a short URL. |
| `GET` | `/:short_code` | Redirect to the original long URL. |
| `GET` | `/api/stats/:short_code` | Get statistics (hits, created_at, etc). |

Full specification is available in [docs/openapi.json](docs/openapi.json).

---
Built with ❤️ by M. Arif Yahya
