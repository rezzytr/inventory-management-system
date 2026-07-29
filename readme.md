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

# 🧪 Manual & API QA Portfolio: Go RESTful API Inventory & Stock Management

![QA Status](https://img.shields.io/badge/QA%20Status-Completed-brightgreen)
![Test Framework](https://img.shields.io/badge/Testing-Manual%20%26%20API-blue)
![Backend](https://img.shields.io/badge/Backend-Go%20%2F%20Docker-00ADD8)
![Tools](https://img.shields.io/badge/Tools-Postman%20%7C%20Web%20UI-orange)

## 📌 Executive Summary
Proyek ini adalah dokumentasi pengujian QA menyeluruh (*End-to-End*) untuk sistem **Inventory Management & Stock Control API** berbasis bahasa pemrograman **Go (Golang)** yang berjalan di dalam kontainer **Docker**. 

Pengujian berfokus pada **Product Management** dan **Stock Adjustment**, mencakup validasi batas (*boundary testing*), pengujian negatif (*negative scenarios*), ketahanan server (*resilience/middleware check*), serta sinkronisasi antara logika backend API dan tampilan antarmuka web (UI).

---

## 🎯 Scope of Testing (Cakupan Pengujian)
1. **Product Management (`/products`):**
   - CRUD Produk (Create, Read, Update, Delete) & Validasi SKU Duplikat.
   - Sanitasi Payload JSON & Malformed JSON Handling (`TC-MID-001`).
2. **Stock Management (`/products/{sku}/stock`):**
   - Penambahan & pengurangan kuantitas stok.
   - **Boundary & Negative Testing:** Validasi stok minus/melebihi batas stok tersedia (`TC-STK-003`).
   - Validasi ketersediaan SKU (`TC-STK-005`).
   - Peringatan stok tipis / Low Stock Threshold $\le 5$ (`TC-STK-004`).
3. **Audit Log & UI Dashboard:**
   - Verifikasi riwayat pergerakan stok (`TC-LOG-001`).
   - Visual indicator badge peringatan pada UI (`TC-STK-005`).

---

## 🛠️ Tech Stack & Tools
- **Backend API:** Go (Golang) REST API
- **Containerization:** Docker & Docker Compose
- **Testing Tools:** Postman Desktop App, Web UI Dashboard
- **Documentation:** Google Sheets / MS Excel (Horizontal QA Reporting Standard)

---

## 🔍 Key QA Findings & Bug Highlights

> ### 🛑 Critical Bug Discovery (`TC-STK-003`)
> * **Isu:** Saat menguji reduksi stok menggunakan payload acuan awal (`"action": "reduce"`), backend Go menerima input bernilai positif tanpa melakukan validasi kondisi `current_stock < reduce_amount`. Hal ini memicu kesalahan *integer underflow/data corruption* yang menyebabkan nilai stok di database melonjak menjadi angka tidak valid (**235**).
> * **Resolusi Penemuan:** Pengujian dikonfirmasi kembali via Web UI dan Postman menggunakan input nilai negatif (`-n`). Hasil akhir mengonfirmasi bahwa backend secara sukses menolak nilai reduksi yang melampaui sisa stok dengan pesan alert: *"Gagal update stok: produk tidak ditemukan atau jumlah stok tidak mencukupi"*.
> * **Status:** **PASS** *(Boundary validation terverifikasi bekerja dengan benar)*.

---

## 📊 Summary of Test Cases

| Test Case ID | Modul / Feature | Scenario | Type | Result |
| :--- | :--- | :--- | :--- | :---: |
| **TC-PROD-001 - 006** | Product | Create, Read, Update, Delete & Validation | Pos/Neg | **PASS** |
| **TC-STK-001** | Stock | Restock / Add Stock Quantity | Positive | **PASS** |
| **TC-STK-002** | Stock | Reduce Stock Quantity (Valid Amount) | Positive | **PASS** |
| **TC-STK-003** | Stock | Reduce Stock Exceeding Current Quantity | Negative | **PASS** |
| **TC-STK-004** | Stock | Check Low Stock List API ($\le 5$) | Positive | **PASS** |
| **TC-STK-005** | Stock | Low Stock Visual Alert Badge in Web UI | Positive | **PASS** |
| **TC-LOG-001** | Audit Log | Verify Stock Movement History & Timestamp | Positive | **PASS** |
| **TC-MID-001** | Resilience | Panicking Route & Malformed JSON Recovery Check | Negative | **PASS** |

