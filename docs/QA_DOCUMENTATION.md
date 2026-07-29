<h1 align="center">🧪 Manual & API QA Documentation</h1>
<h3 align="center">Go RESTful API Inventory & Stock Management</h3>

<p align="center">
  <a href="../README.md">← Back to Main Repository</a>
</p>

---

## 📌 Executive Summary
Proyek ini adalah dokumentasi pengujian QA menyeluruh (*End-to-End*) untuk sistem **Inventory Management & Stock Control API** berbasis bahasa pemrograman **Go (Golang)** yang berjalan di dalam kontainer **Docker**. 

Pengujian berfokus pada **Product Management** dan **Stock Adjustment**, mencakup validasi batas (*boundary testing*), pengujian negatif (*negative scenarios*), ketahanan server (*resilience/middleware check*), serta sinkronisasi antara logika backend API dan tampilan antarmuka web (UI).

---

## 🎯 Scope of Testing
1. **Product Management (`/products`):**
   - CRUD Produk (Create, Read, Update, Delete) & Validasi SKU Duplikat.
   - Sanitasi Payload JSON & Malformed JSON Handling (`TC-MID-001`).
2. **Stock Management (`/products/{sku}/stock`):**
   - Penambahan & pengurangan kuantitas stok.
   - **Boundary & Negative Testing:** Validasi stok minus/melebihi batas stok tersedia (`TC-STK-003`).
   - Peringatan stok tipis / Low Stock Threshold $\le 5$ (`TC-STK-004`).
3. **Audit Log & UI Dashboard:**
   - Verifikasi riwayat pergerakan stok (`TC-LOG-001`).
   - Visual indicator badge alert pada UI (`TC-STK-005`).

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

---

## 📑 Postman & Spreadsheet Artifacts
- **Postman Collection & Environment:** Tersimpan di folder [`/postman`](../postman)
- **Full Excel Test Sheet:** Tersimpan di [`/docs/QA-Test-Plan-and-Execution-Sheet.xlsx`](./)