# SQL dengan Golang
## Cara Menjalankan
1. Jalankan postgres dengan docker
```
docker run --name postgresql -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=database -p 5432:5432 -d postgres:16
```

2. Install module yang dibutuhkan
```
go mod tidy
```

3. Jalankan program
```
go run main.go
```

## Konten
Dalam file `main.go` terdapat 9 bagian contoh:
- Menghubungkan Golang dengan Postgres
- Pengaturan lanjutan
- Membuat tabel
- Menambah data
- Mengambil 1 data
- Mengambil banyak data
- Memperbarui data
- Menghapus data
- Transaction