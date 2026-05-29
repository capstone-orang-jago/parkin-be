# Parkin Backend (parkin-be)

API backend untuk aplikasi Parkin (capstone project).

## Ringkasan

- Bahasa: Go
- Framework: Gin, Cobra (CLI)
- Database: PostgreSQL

Proyek ini menyediakan API backend lengkap termasuk migrasi, seeder, dan integrasi WhatsApp (whatsmeow) untuk fitur notifikasi/OTP.

## Fitur

- Autentikasi (JWT)
- Manajemen pengguna dan transaksi parkir
- Migrasi dan seeder berbasis file SQL
- Integrasi WhatsApp (QR pertama kali untuk autentikasi device)

## Prasyarat

- Go versi 1.25 atau lebih baru
- PostgreSQL

## Instalasi

1. Clone repository:

```bash
git clone https://github.com/capstone-orang-jago/parkin-be
cd parkin-be
```

2. Ambil dependency Go (ambil dan rapihkan modul):

```bash
go mod download
```

3. Buat file `.env` di root proyek (lihat contoh di bawah).

## Contoh `.env`

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=parkin_db
JWT_SECRET=your_jwt_secret
DEV_ORIGIN=http://localhost:3000
PROD_ORIGIN=https://your-production-url.com
QR_SECRET_KEY=some_random_key
MIDTRANS_SERVER_KEY=midtrans_key_here
```

Sesuaikan nilai sesuai environment Anda.

## Setup Database

1. Jalankan migrasi "up" untuk membuat tabel:

```bash
go run main.go migrate:up
```

2. Membuat migration baru:

```bash
go run main.go migrate:create add_example_table
```

3. Revert (rollback) semua migration:

```bash
go run main.go migrate:down
```

## Seeder

1. Membuat file seeder baru (SQL):

```bash
go run main.go seeder:create initial_users
```

2. Menjalankan semua seeder:

```bash
go run main.go seeder:all
```

3. Menjalankan seeder tertentu:

```bash
go run main.go seeder:run initial_users
```

## Menjalankan Server

Jalankan server untuk development:

```bash
go run main.go serve
```

Server akan berjalan di `:8080` secara default. Saat pertama kali menjalankan integrasi WhatsApp, terminal akan menampilkan QR code untuk discan (lihat log terminal).

## Dokumentasi API

Terdapat file Swagger di `docs/swagger.yaml` dan `docs/swagger.json`. Anda dapat membuka file tersebut dengan Swagger UI (lokal atau eksternal) untuk eksplorasi endpoint. Anda juga dapat mengakses dokumentasi saat server berjalan di:

http://localhost:8080/swagger/index.html

## Catatan

- Pastikan variabel environment sudah benar sebelum menjalankan migrasi atau server.
- Untuk fitur WhatsApp, ikuti instruksi log pada terminal saat pertama kali menghubungkan device.
