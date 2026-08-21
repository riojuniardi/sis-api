# SIS API

REST API untuk **Sistem Informasi Sarana dan Prasarana (SIS)**. API ini digunakan untuk mengelola data pengguna, kategori, lokasi, inventaris barang, peminjaman barang, dan riwayat audit kondisi barang.

## Daftar Isi

- [Fitur](#fitur)
- [Teknologi](#teknologi)
- [Persyaratan](#persyaratan)
- [Menjalankan Project](#menjalankan-project)
- [Konfigurasi Environment](#konfigurasi-environment)
- [Autentikasi dan Hak Akses](#autentikasi-dan-hak-akses)
- [Format Data dan Enum](#format-data-dan-enum)
- [Daftar Endpoint](#daftar-endpoint)
- [Contoh Penggunaan](#contoh-penggunaan)
- [Struktur Project](#struktur-project)
- [Catatan Implementasi](#catatan-implementasi)

## Fitur

- Registrasi dan login pengguna menggunakan JWT.
- Pembatasan akses berdasarkan role `admin`.
- Manajemen pengguna, kategori, lokasi, dan item inventaris.
- Pengajuan, persetujuan, penolakan, dan pengembalian peminjaman.
- Pencatatan serta penelusuran riwayat kondisi item.
- Auto migration tabel database melalui GORM.
- Seeder role awal: `admin`, `guru`, dan `siswa`.

## Teknologi

- Go `1.26.4`
- Gin `1.12.0`
- GORM
- PostgreSQL
- JSON Web Token (JWT)
- `godotenv` untuk memuat konfigurasi dari file `.env`

## Persyaratan

Pastikan perangkat sudah memiliki:

- Go versi `1.26.4` atau yang kompatibel dengan `go.mod`
- PostgreSQL yang aktif
- Git, bila project diambil dari repository

## Menjalankan Project

1. Masuk ke folder project:

   ```bash
   cd sis-api
   ```

2. Buat file `.env` di root project. Lihat contoh pada bagian [Konfigurasi Environment](#konfigurasi-environment).

3. Unduh dependency:

   ```bash
   go mod download
   ```

4. Jalankan API:

   ```bash
   go run .
   ```

   Server tersedia di `http://localhost:8080`.

Saat aplikasi mulai, sistem akan:

- Membaca file `.env`.
- Membuka koneksi ke PostgreSQL.
- Menjalankan `AutoMigrate` untuk seluruh model.
- Membuat role awal jika belum tersedia.
- Menjalankan server Gin pada port `8080`.

## Konfigurasi Environment

Buat `.env` dengan isi minimal berikut:

```env
DATABASE_URL=postgres://postgres:password@localhost:5432/sis_db?sslmode=disable
JWT_SECRET=ganti_dengan_secret_yang_panjang_dan_aman
```

Keterangan:

| Variabel | Wajib | Kegunaan |
| --- | --- | --- |
| `DATABASE_URL` | Ya | DSN koneksi PostgreSQL yang digunakan GORM. |
| `JWT_SECRET` | Ya | Secret untuk menandatangani dan memvalidasi token JWT. |

Jangan memasukkan `.env` ke repository. File tersebut sudah diabaikan oleh `.gitignore`.

## Autentikasi dan Hak Akses

Endpoint yang membutuhkan autentikasi harus menerima header berikut:

```http
Authorization: Bearer <JWT_TOKEN>
```

Token diterbitkan oleh `POST /api/auth/login` dan berlaku selama **7 hari**.

### Role

Role yang dibuat otomatis oleh seeder:

- `admin`: dapat mengakses seluruh endpoint admin.
- `guru`: pengguna biasa.
- `siswa`: pengguna biasa.

Endpoint admin berada di bawah prefix `/api/admin` dan hanya dapat diakses oleh token dengan claim role `admin`. Endpoint yang tidak ditandai `Admin` pada tabel endpoint tidak menggunakan middleware admin.

### Error autentikasi umum

- `401 Unauthorized`: token tidak ada atau token tidak valid.
- `403 Forbidden`: token valid, tetapi role bukan `admin`.
- `400 Bad Request`: body JSON atau format data tidak valid.
- `404 Not Found`: data dengan ID yang diminta tidak ditemukan.
- `409 Conflict`: kode unik sudah digunakan.

## Format Data dan Enum

### Format tanggal

- Item: `purchase_date` menggunakan `YYYY-MM-DD`.
- Peminjaman: `borrow_date`, `due_date`, dan `return_date` menerima `YYYY-MM-DD` atau `YYYY-MM-DD HH:mm:ss`.

### Kondisi item

Nilai `condition` dan `condition_before`/`condition_after` yang didukung:

- `baik`
- `rusak_ringan`
- `rusak_berat`

### Status item

- `tersedia`
- `dipinjam`
- `perbaikan`
- `afkir`

### Status peminjaman

- `pending`
- `disetujui`
- `ditolak`
- `dipinjam`
- `dikembalikan`
- `terlambat`

### Alasan audit kondisi

- `pemeriksaan_rutin`
- `pengembalian_pinjaman`
- `laporan_kerusakan`
- `selesai_perbaikan`

## Daftar Endpoint

Base URL: `http://localhost:8080/api`

Keterangan akses:

- `Publik`: tidak memerlukan token.
- `Auth`: memerlukan JWT valid.
- `Admin`: memerlukan JWT valid dengan role `admin`.

### Autentikasi

| Method | Path | Akses | Kegunaan |
| --- | --- | --- | --- |
| `POST` | `/auth/register` | Publik | Mendaftarkan pengguna baru. |
| `POST` | `/auth/login` | Publik | Login dan mendapatkan JWT. |

### Pengguna

| Method | Path | Akses | Kegunaan |
| --- | --- | --- | --- |
| `GET` | `/users/me` | Auth | Mengambil profil pengguna dari token. |
| `GET` | `/admin/users` | Admin | Mengambil seluruh pengguna. |
| `GET` | `/admin/users/:id` | Admin | Mengambil pengguna berdasarkan ID. |

### Kategori

| Method | Path | Akses | Kegunaan |
| --- | --- | --- | --- |
| `GET` | `/categories` | Publik | Mengambil seluruh kategori. |
| `GET` | `/categories/:id` | Publik | Mengambil kategori berdasarkan ID. |
| `POST` | `/admin/categories` | Admin | Membuat kategori. |
| `PUT` | `/admin/categories/:id` | Admin | Memperbarui kategori. |
| `DELETE` | `/admin/categories/:id` | Admin | Menghapus kategori. |

### Lokasi

| Method | Path | Akses | Kegunaan |
| --- | --- | --- | --- |
| `GET` | `/locations` | Publik | Mengambil seluruh lokasi. |
| `GET` | `/locations/:id` | Publik | Mengambil lokasi berdasarkan ID. |
| `POST` | `/admin/locations` | Admin | Membuat lokasi. |
| `PUT` | `/admin/locations/:id` | Admin | Memperbarui lokasi. |
| `DELETE` | `/admin/locations/:id` | Admin | Menghapus lokasi. |

### Item inventaris

| Method | Path | Akses | Kegunaan |
| --- | --- | --- | --- |
| `GET` | `/items` | Publik | Mengambil seluruh item beserta kategori dan lokasi. |
| `GET` | `/items/:id` | Publik | Mengambil item berdasarkan ID. |
| `POST` | `/admin/items` | Admin | Membuat item inventaris. |
| `PUT` | `/admin/items/:id` | Admin | Memperbarui item. |
| `DELETE` | `/admin/items/:id` | Admin | Menghapus item. |

### Peminjaman

| Method | Path | Akses | Kegunaan |
| --- | --- | --- | --- |
| `GET` | `/borrowings` | Publik | Mengambil seluruh data peminjaman. |
| `GET` | `/borrowings/:id` | Publik | Mengambil peminjaman berdasarkan ID. |
| `POST` | `/borrowings` | Publik | Mengajukan peminjaman. |
| `PUT` | `/admin/borrowings/:id` | Admin | Memperbarui data peminjaman. |
| `DELETE` | `/admin/borrowings/:id` | Admin | Menghapus data peminjaman. |
| `PATCH` | `/admin/borrowings/:id/approve` | Admin | Memproses persetujuan peminjaman. |
| `PATCH` | `/admin/borrowings/:id/return` | Admin | Memproses pengembalian barang. |

### Audit kondisi barang

| Method | Path | Akses | Kegunaan |
| --- | --- | --- | --- |
| `GET` | `/condition-audits` | Publik | Mengambil audit; dapat difilter dengan `item_id` dan `user_id`. |
| `GET` | `/condition-audits/:id` | Publik | Mengambil audit berdasarkan ID. |
| `GET` | `/condition-audits/item/:item_id` | Publik | Mengambil riwayat audit sebuah item. |
| `POST` | `/admin/condition-audits` | Admin | Mencatat perubahan kondisi item. |
| `DELETE` | `/admin/condition-audits/:id` | Admin | Menghapus log audit. |

## Contoh Penggunaan

Gunakan `Content-Type: application/json` untuk endpoint yang memiliki request body. Contoh berikut menggunakan Bash, Git Bash, atau terminal yang mendukung `curl`.

### Register

`role_id` mengacu pada ID role di database. Seeder membuat role `admin`, `guru`, dan `siswa`; gunakan ID aktual dari database.

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Budi Santoso",
    "user_name": "budi",
    "password": "rahasia123",
    "role_id": 2
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "user_name": "budi",
    "password": "rahasia123"
  }'
```

Simpan nilai `token` dari response login dan kirimkan pada endpoint yang membutuhkan autentikasi.

### Membuat kategori sebagai admin

```bash
curl -X POST http://localhost:8080/api/admin/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -d '{
    "name": "Elektronik",
    "code": "ELK"
  }'
```

### Membuat lokasi sebagai admin

```bash
curl -X POST http://localhost:8080/api/admin/locations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -d '{
    "name": "Laboratorium Komputer",
    "code": "LAB-KOM-01",
    "building": "Gedung A",
    "floor": 2,
    "pic_user_id": 1
  }'
```

### Membuat item sebagai admin

```bash
curl -X POST http://localhost:8080/api/admin/items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -d '{
    "item_code": "INV-001",
    "category_id": 1,
    "location_id": 1,
    "name": "Laptop Lenovo",
    "source_of_funds": "BOS",
    "purchase_date": "2026-01-15",
    "purchase_price": 8500000,
    "condition": "baik",
    "status": "tersedia"
  }'
```

### Mengajukan peminjaman

```bash
curl -X POST http://localhost:8080/api/borrowings \
  -H "Content-Type: application/json" \
  -d '{
    "item_id": 1,
    "user_id": 2,
    "borrow_date": "2026-08-21 08:00:00",
    "due_date": "2026-08-23 16:00:00",
    "condition_before": "baik",
    "notes": "Digunakan untuk kegiatan sekolah"
  }'
```

### Menyetujui atau menolak peminjaman

Nilai `status` yang diterima endpoint approval adalah `disetujui`, `ditolak`, atau `dipinjam`.

```bash
curl -X PATCH http://localhost:8080/api/admin/borrowings/1/approve \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -d '{
    "status": "disetujui",
    "notes": "Peminjaman disetujui"
  }'
```

### Memproses pengembalian

```bash
curl -X PATCH http://localhost:8080/api/admin/borrowings/1/return \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -d '{
    "condition_after": "baik",
    "notes": "Barang kembali lengkap"
  }'
```

### Mencatat audit kondisi

```bash
curl -X POST http://localhost:8080/api/admin/condition-audits \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -d '{
    "item_id": 1,
    "user_id": 1,
    "condition_after": "rusak_ringan",
    "reason": "laporan_kerusakan",
    "notes": "Terdapat goresan pada casing"
  }'
```

### Filter audit

```bash
curl "http://localhost:8080/api/condition-audits?item_id=1&user_id=1"
```

## Struktur Project

```text
.
├── config/              # Koneksi database dan konfigurasi aplikasi
├── controllers/         # Handler HTTP dan aturan bisnis endpoint
├── database/seeders/    # Seeder role awal
├── middlewares/         # Middleware JWT dan validasi role admin
├── models/              # Model GORM dan struktur input JSON
├── routes/              # Registrasi route API
├── .env                 # Konfigurasi lokal, tidak di-commit
├── go.mod               # Module dan dependency Go
├── go.sum               # Checksum dependency
├── main.go              # Entry point aplikasi
└── README.md            # Dokumentasi project
```

## Catatan Implementasi

- Tidak ada endpoint health check khusus; keberhasilan startup dapat dicek dari log aplikasi dan request ke endpoint publik seperti `GET /api/items`.
- Database otomatis dimigrasikan saat aplikasi dijalankan. Pastikan user PostgreSQL memiliki izin membuat dan mengubah tabel.
- `item_code`, `category.code`, dan `location.code` bersifat unik.
- Saat pengembalian dengan kondisi `rusak_berat`, status item otomatis menjadi `perbaikan`. Kondisi lain mengembalikan status item menjadi `tersedia`.
- Saat peminjaman disetujui atau diubah menjadi `dipinjam`, status item menjadi `dipinjam`.
- Response error menggunakan object JSON dengan format umum `{ "error": "..." }`.
- Password disimpan menggunakan bcrypt dan tidak dikembalikan pada response API.

## Pengembangan

Format kode dan cek seluruh package dapat dijalankan dengan:

```bash
gofmt -w .
go test ./...
```

Saat ini repository belum memiliki file test khusus. Perintah `go test ./...` tetap berguna untuk memastikan seluruh package dapat dikompilasi.
