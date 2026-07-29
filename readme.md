# 📦 Inventory Management System (Go Clean Architecture)

A production-ready, highly robust Inventory Management System written in **Go (Golang)** using **Clean Architecture** principles. Includes atomic stock operations, negative stock protection, audit logging, a functional SPA Web Dashboard, and full Docker containerization.

---

## 🚀 Features

* **Clean Architecture**: Strict separation of concerns (Domain Models, Repositories, Services, Handlers, Middleware).
* **Data Safety**:
  * Atomic stock adjustments (prevents race conditions).
  * Negative stock protection (rejects invalid reductions).
  * Automated audit logs for every stock movement.
* **RESTful API & Web UI**:
  * Full CRUD for products.
  * Real-time Interactive SPA Web Dashboard built with HTML5 & Tailwind CSS.
* **Resiliency**: Recovery middleware against panic, request logger middleware, and unit test coverage with mocks.
* **Containerized**: Production-ready multi-stage Docker build with persistence via named volumes.

---

## 🛠️ Tech Stack

* **Language**: Go 1.22+
* **Database**: SQLite (embedded) / PostgreSQL ready
* **Frontend**: HTML5, Vanilla JavaScript, Tailwind CSS (CDN)
* **Configuration**: `godotenv`
* **Containerization**: Docker & Docker Compose

---

## 🏗️ Architecture Overview

```text
INVENTORY-MANAGEMENT-SYSTEM/
├── cmd/
│   └── api/
│       └── main.go           # Application entrypoint & dependency injection
├── internal/
│   ├── handler/              # HTTP Delivery Layer (Request/Response mapping)
│   ├── middleware/           # HTTP Middlewares (Logger, Recovery)
│   ├── model/                # Core Domain Models & Business Entities
│   ├── repository/           # Data Access Layer (SQLite, Postgres, Mock)
│   └── service/              # Business Logic Layer & Unit Tests
├── web/
│   └── index.html            # Single Page Application (SPA) Web Dashboard
├── .env.example              # Template environment variables
├── Dockerfile                # Multi-stage Docker build
└── docker-compose.yml        # Orchestration setup
