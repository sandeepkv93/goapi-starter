# 🚀 GoAPI Starter: A Robust Go API Boilerplate

## 📋 Project Overview

GoAPI Starter is a comprehensive, production-ready Go API boilerplate designed to jumpstart your backend development with best practices and modern Go ecosystem tools. 🛠️

### ✨ Features

- 🔐 Authentication System
  - JWT-based authentication
  - Signup and Signin flows
  - Refresh token mechanism
- 🗃️ Database Integration
  - PostgreSQL with GORM ORM
  - Auto-migration support
- 🛡️ Middleware
  - Logging middleware
  - Authentication middleware
- 🧪 Structured Project Layout
  - Clean, modular architecture
  - Separation of concerns

## 🛠️ Technology Stack

- **Language**: Go 1.23
- **Web Framework**: Chi Router
- **ORM**: GORM
- **Database**: PostgreSQL
- **Authentication**: JWT
- **Validation**: go-playground/validator
- **Environment**: godotenv

## 🚦 Prerequisites

- 🐳 Docker
- 🐍 Go 1.23+
- 📦 PostgreSQL

## 🔧 Installation & Setup

### 1. Clone the Repository

```bash
git clone https://github.com/sandeepkv93/goapi-starter.git
cd goapi-starter
```

### 2. Environment Configuration

Copy the example environment file:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

- Set database credentials
- Configure JWT secrets
- Adjust server port

### 3. Running the Application

#### Local Development

```bash
# Install air for hot reloading
task install-tools

# Start development server
task dev
```

#### Docker Deployment

```bash
# Build and run with Docker
task docker-run
```

## 🧪 Running Tests

```bash
# Run all tests
task test
```

## 📂 Project Structure

```
goapi-starter/
│
├── cmd/                 # Application entry points
│   └── api/
│       └── main.go
│
├── internal/            # Core application code
│   ├── config/          # Configuration management
│   ├── database/        # Database connection
│   ├── grafana/         # Grafana configuration
│   ├── handlers/        # HTTP request handlers
│   ├── logger/          # Logger
│   ├── metrics/         # Metrics
│   ├── middleware/      # Request middleware
│   ├── models/          # Data models
│   ├── prometheus/      # Prometheus configuration
│   ├── routes/          # API route definitions
│   ├── services/        # Business logic
│   └── utils/           # Utility functions
│
├── goapi.rest           # Rest client for testing the API
├── docker-compose.yml   # Docker orchestration
├── docker-compose.dev.yml # Docker orchestration for local development
├── Dockerfile           # Container configuration
└── Taskfile.yml         # Task automation
```

## 🔐 Authentication Endpoints

- `POST /api/auth/signup`: Register a new user
- `POST /api/auth/signin`: User login
- `POST /api/auth/refresh`: Refresh authentication tokens

## 📦 API Endpoints

### Products

- `GET /api/dummy-products`: List all the dummy products
- `POST /api/dummy-products`: Create a new dummy product
- `GET /api/dummy-products/{id}`: Get a specific dummy product
- `PUT /api/dummy-products/{id}`: Update a dummy product
- `DELETE /api/dummy-products/{id}`: Delete a dummy product

## 🛡️ Security Features

- Password hashing with bcrypt
- JWT token-based authentication
- Refresh token mechanism
- Input validation
- Middleware-based authentication

## Metrics and Monitoring

The application includes comprehensive metrics for monitoring:

- **HTTP Metrics**: Request counts, durations, and status codes
- **Business Operation Metrics**: Success/failure rates for key operations
- **Error Tracking**: Detailed error tracking with categorization
- **Database Metrics**: Database operation counts and performance

Metrics are exposed via a `/metrics` endpoint in Prometheus format and can be visualized using the included Grafana dashboards.

### Error Tracking

The application includes detailed error tracking that categorizes errors and captures specific error reasons. This helps with:

- Identifying common error patterns
- Debugging specific issues
- Monitoring error trends over time

Error details are available in the Grafana dashboard under the "Error Details" panels.

## 🛠️ Available Tasks

Here are all the available tasks you can run with `task`:

### Development Tasks

| Command           | Description                                        |
| ----------------- | -------------------------------------------------- |
| **default**       | Run the application (alias for `run`)              |
| **ensure-db**     | Ensure the database container is running           |
| **build**         | Build the application binary                       |
| **run**           | Run the application locally                        |
| **dev**           | Run the application with hot reload (requires air) |
| **clean**         | Clean build files                                  |
| **test**          | Run tests                                          |
| **install-tools** | Install development tools                          |

### Docker Tasks

| Command              | Description                                                                                                                                                  |
| -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| **docker-build**     | Build Docker images                                                                                                                                          |
| **docker-run**       | Run all services in Docker                                                                                                                                   |
| **docker-logs**      | Follow Docker logs                                                                                                                                           |
| **docker-stop**      | Stop Docker containers                                                                                                                                       |
| **docker-clean**     | Stop and remove Docker containers and volumes                                                                                                                |
| **docker-clean-run** | Stop and remove Docker containers and volumes and start them again                                                                                           |
| **docker-dev**       | Run all services in Docker for local development except the API service. This allows you to run the API service locally via IDE and attach a debugger to it. |
| **docker-dev-stop**  | Stop all services in Docker for local development except the API service                                                                                     |
| **docker-dev-logs**  | Follow Docker logs for local development except the API service                                                                                              |
| **docker-dev-clean** | Stop and remove Docker containers and volumes for local development except the API service                                                                   |
| **docker-dev-run**   | Run all services in Docker for local development and run the API service locally without IDE                                                                 |

### Monitoring Tasks

| Command             | Description                     |
| ------------------- | ------------------------------- |
| **prometheus-up**   | Start Prometheus container      |
| **grafana-up**      | Start Grafana container         |
| **monitoring-up**   | Start all monitoring containers |
| **monitoring-logs** | Follow monitoring logs          |
