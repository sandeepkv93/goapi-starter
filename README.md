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
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # Request middleware
│   ├── models/          # Data models
│   ├── routes/          # API route definitions
│   ├── services/        # Business logic
│   └── utils/           # Utility functions
│
├── docker-compose.yml   # Docker orchestration
├── Dockerfile           # Container configuration
└── Taskfile.yml         # Task automation
```

## 🔐 Authentication Endpoints

- `POST /api/auth/signup`: Register a new user
- `POST /api/auth/signin`: User login
- `POST /api/auth/refresh`: Refresh authentication tokens

## 📦 API Endpoints

### Products

- `GET /api/products`: List all products
- `POST /api/products`: Create a new product
- `GET /api/products/{id}`: Get a specific product
- `PUT /api/products/{id}`: Update a product
- `DELETE /api/products/{id}`: Delete a product

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
