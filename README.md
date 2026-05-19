# Fleetify - Maintenance Report System

Project ini adalah aplikasi manajemen maintenance report untuk armada kendaraan (fleet maintenance). Aplikasi ini digunakan untuk melaporkan kerusakan kendaraan, menyetujui perbaikan oleh manager, dan menyelesaikan pengerjaan pemeliharaan.

Sistem ini menggunakan Role-Based Access Control (RBAC) dengan 2 role utama:
1. **SA (Service Advisor)**: Membuat laporan dan menyelesaikan perbaikan.
2. **APPROVAL (Manager)**: Menyetujui laporan pemeliharaan.

---

## Tech Stack

### Backend
* Go v1.23
* Fiber Framework
* GORM ORM
* MySQL 8.0

### Frontend
* HTML & CSS (Bootstrap 5)
* Vanilla JavaScript (Tanpa Framework)

### DevOps & Env
* Docker & Docker Compose
* GoDotEnv

---

## Catatan Keamanan dan Integritas Data

* **Validasi RBAC menggunakan Middleware**: Hak akses peran (SA vs APPROVAL) divalidasi langsung di backend menggunakan middleware custom `RequireRole` dengan membaca header `X-User-ID`.
* **Manipulasi DOM Aman**: Frontend menghindari penggunaan `innerHTML` dan menggunakan `document.createElement` serta `textContent` untuk mencegah celah keamanan Cross-Site Scripting (XSS).
* **Transaksi Database**: Pembuatan report baru dibungkus di dalam transaksi database GORM (`database.DB.Transaction`). Jika terjadi error pada salah satu item perbaikan, semua perubahan akan dibatalkan (rollback).
* **Snapshot Harga**: Harga barang saat pembuatan laporan langsung disimpan ke dalam tabel `report_items` sebagai snapshot. Ini bertujuan agar total biaya laporan lama tidak berubah jika harga barang di tabel master diubah di kemudian hari.

---

## Struktur Folder

```text
fleetify-test/
├── backend/
│   ├── cmd/
│   │   └── main.go              # Main entrypoint backend
│   ├── database/                # Koneksi database
│   ├── dto/                     # DTO untuk request payload
│   ├── handlers/                # HTTP handlers (Controller & Query)
│   ├── middleware/              # Auth middleware (RequireRole)
│   ├── models/                  # Struct model GORM (Schema database)
│   ├── routes/                  # API routing
│   ├── seeders/                 # Seeder otomatis untuk user dan master data
│   ├── services/                # Logika tambahan (Webhook service)
│   ├── .env.example             # Contoh file env
│   └── Dockerfile               # Dockerfile backend
├── frontend/
│   ├── app.js                   # Logic fetch dan manipulasi DOM
│   ├── index.html               # Tampilan web utama
│   └── style.css                # Custom CSS tambahan
├── docker-compose.yml           # Docker Compose MySQL dan Backend
└── README.md                    # Dokumentasi project
```

---

## Cara Menjalankan Project

### Prasyarat
Pastikan sudah menginstal Docker Desktop di komputer Anda.

### Menggunakan Docker Compose (Direkomendasikan)

1. Buka terminal di root direktori project ini.
2. Jalankan perintah berikut:
   ```bash
   docker compose up --build
   ```
3. Backend akan berjalan di `http://localhost:8080` dan MySQL berjalan di port `3307` (port internal container `3306`).
4. Buka file `frontend/index.html` langsung di browser untuk menggunakan aplikasi.

### Menjalankan secara Lokal (Tanpa Docker)

1. Pastikan MySQL lokal sudah aktif dan buat database baru bernama `fleetify`.
2. Masuk ke folder `backend` dan salin `.env.example` menjadi `.env`:
   ```bash
   cp .env.example .env
   ```
3. Sesuaikan isi file `.env` dengan konfigurasi MySQL lokal Anda.
4. Jalankan aplikasi backend:
   ```bash
   go run cmd/main.go
   ```
5. Buka file `frontend/index.html` di browser.

---

## Data User Seeder

Sistem menggunakan header request `X-User-ID` untuk membedakan hak akses. Data user awal yang otomatis terisi saat seeder berjalan:

| User ID | Username | Role | Hak Akses Utama |
| :--- | :--- | :--- | :--- |
| **1** | `sandy_sa` | `SA` | Create Report, Complete Report, View Reports |
| **2** | `manager_approval` | `APPROVAL` | Approve Report, View Reports |

---

## Daftar API Endpoint

Semua request API menggunakan prefix `/api`.

### 1. Master Data
* `GET /api/users` - Mengambil data user
* `GET /api/vehicles` - Mengambil data kendaraan
* `GET /api/master-items` - Mengambil data barang master (sparepart & jasa)

### 2. Maintenance Reports
* `GET /api/reports` - Mengambil riwayat laporan
* `POST /api/reports` - Membuat laporan pemeliharaan baru
  * **Role**: `SA` (Header `X-User-ID: 1`)
  * **Payload**:
    ```json
    {
      "vehicle_id": 1,
      "odometer": 125000,
      "complaint": "Rem bunyi berdecit saat pengereman",
      "initial_photo": "https://url-foto/rem_aus.jpg",
      "items": [
        { "item_id": 1, "quantity": 1 },
        { "item_id": 4, "quantity": 1 }
      ]
    }
    ```
* `PATCH /api/reports/:id/approve` - Menyetujui laporan
  * **Role**: `APPROVAL` (Header `X-User-ID: 2`)
  * **Status flow**: Mengubah status dari `PENDING_APPROVAL` ke `APPROVED`.
* `PATCH /api/reports/:id/complete` - Menyelesaikan pengerjaan laporan
  * **Role**: `SA` (Header `X-User-ID: 1`)
  * **Status flow**: Mengubah status dari `APPROVED` ke `COMPLETED`.
  * **Payload**:
    ```json
    {
      "proof_photo": "https://url-foto/rem_selesai.jpg"
    }
    ```

---

## Notifikasi Webhook Async

Setiap ada perubahan status laporan (saat Approve atau Complete), backend akan mengirim notifikasi webhook secara async di background menggunakan goroutine ke URL yang diatur pada env `WEBHOOK_URL`.

**Payload Webhook:**
```json
{
  "report_id": 1,
  "status": "APPROVED",
  "message": "Maintenance report status updated"
}
```

---

## Ekspor CSV

Terdapat tombol ekspor data di halaman web untuk mengunduh semua riwayat laporan dalam bentuk file CSV. Proses pengolahan data CSV dilakukan langsung di sisi client (native JavaScript) dengan menggunakan API `Blob` dan `URL.createObjectURL`. Output CSV berisi kolom ID, License Plate Kendaraan, Status, Deskripsi Keluhan, dan Odometer.