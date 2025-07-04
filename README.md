# 🏦 Internal Transfer Service

> A high-precision cryptocurrency-style money transfer service built with Go and PostgreSQL

[![Go Version](https://img.shields.io/badge/Go-1.20+-blue.svg)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-green.svg)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://www.docker.com/)

This service allows creation of user accounts and supports internal balance transfers between them with **5-decimal precision**. It is built using Go and uses PostgreSQL as the backend with database-level transaction locking.

## 🎯 Features

- ✅ **High Precision**: 5-decimal place accuracy for crypto-like transfers
- ✅ **ACID Transactions**: Database-level locking with `SELECT FOR UPDATE`
- ✅ **RESTful API**: Clean HTTP endpoints for account and transfer operations
- ✅ **Comprehensive Logging**: Structured logging with Zap
- ✅ **Docker Ready**: Easy setup with Docker Compose
- ✅ **Test Coverage**: Unit and integration tests included

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Layer    │    │  Business Logic │    │  Data Access    │
│  (Handlers)     │◄──►│   (Services)    │◄──►│ (Repositories)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Middleware    │    │     Models      │    │   PostgreSQL    │
│ (Logging, DB)   │    │ (Account, Tx)   │    │   Database      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 📁 Project Structure

```
transfer-service/
├── 📂 api/handler/     # HTTP request handlers
├── 📂 cmd/            # Application entry point
├── 📂 model/          # Domain models (Account, Transaction)
├── 📂 repository/     # Database access layer
├── 📂 service/        # Business logic layer
├── 📂 middleware/     # Cross-cutting concerns (logging, DB)
├── 📂 tests/          # Unit and integration tests
├── 📄 docker-compose.yml  # Database setup
├── 📄 schema.sql      # Database schema
└── 📄 start.sh        # Startup script
```

## 🚀 Quick Start

### Prerequisites

- **Go** 1.20 or higher
- **Docker** & **Docker Compose**
- **Git**

### Setup & Run

1. **Clone and navigate to the project**
   ```bash
   git clone <repository-url>
   cd transfer-service
   ```

2. **Configure environment**
   ```bash
   cp .env-sample .env
   # Edit .env with your database credentials
   ```

3. **Start the service**
   ```bash
   ./start.sh
   ```

4. **Run tests**
   ```bash
   ./tests/run_tests.sh
   ```

The service will be available at `http://localhost:8080`

## 📋 Assumptions

- **🔢 Decimal Precision**: All monetary operations use 5 decimal places or lower (e.g., `100.12345`, `100.1`, `100`)
- **🔒 Database Locks**: We use database-level locks (`SELECT FOR UPDATE`) to ensure transaction consistency, not application-level mutexes or caching

## 🔌 API Endpoints

### Create Account
```http
POST /accounts
Content-Type: application/json

{
  "account_id": 123,
  "balance": "100.12345"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Account created successfully",
  "data": {
    "account_id": 123,
    "balance": 100.12345
  }
}
```

### Get Account Balance
```http
GET /accounts/{id}
```

**Response:**
```json
{
  "success": true,
  "message": "Account retrieved successfully",
  "data": {
    "account_id": 123,
    "balance": 100.12345
  }
}
```

### Transfer Money
```http
POST /transactions
Content-Type: application/json

{
  "source_account_id": 123,
  "destination_account_id": 456,
  "amount": "50.12345"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Transfer completed successfully",
  "data": {
    "message": "Transfer completed successfully",
    "transaction": {
      "id": 1,
      "source_account_id": 123,
      "destination_account_id": 456,
      "amount": 50.12345,
      "created_at": "2025-07-05T03:45:43.347Z"
    }
  }
}
```

## ⚠️ Things to Note

1. **🐳 Database Setup**: We are creating the database in Docker. The `docker-compose.yml` only contains PostgreSQL database setup and not the application itself.

2. **⚙️ Environment Configuration**: Before running the application, make sure to rename `.env-sample` to `.env` and configure your environment variables.

## 🧪 Testing

Run the comprehensive test suite:

```bash
# Run all tests
./tests/run_tests.sh

# Run specific test categories
go test ./tests/service -v -run "TestCreateAccount"
go test ./tests/service -v -run "TestTransfer"
```

## 📚 Documentation

- **[Postman Collection](https://documenter.getpostman.com/view/4623773/2sB34bMjS4)** - Interactive API documentation
