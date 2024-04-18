package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Produk struct {
	ID       int
	Nama     string
	Kategori string
	Harga    int
}

func main() {
	// 1 - menghubungkan golang dengan postgres
	connURI := "postgresql://postgres:postgres@localhost:5432?sslmode=disable"
	db, err := sql.Open("pgx", connURI)
	if err != nil {
		fmt.Printf("Gagal menghubungkan ke database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Printf("Terjadi kesalahan: %v\n", err)
		os.Exit(1)
	}

	// 2 - pengaturan lanjutan
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(15 * time.Minute)
	db.SetConnMaxLifetime(1 * time.Hour)

	fmt.Println("Berhasil terhubung ke database")

	// 3 - membuat tabel
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS produk (
	    id SERIAL PRIMARY KEY,
	    nama VARCHAR(255),
	    kategori VARCHAR(50),
	    harga INT
	);
	`)
	if err != nil {
		fmt.Printf("Gagal membuat tabel: %v\n", err)
		os.Exit(1)
	}

	// 4 - menambah data
	_, err = db.Exec(`
	INSERT INTO produk (nama, kategori, harga)
	VALUES ($1, $2, $3)
	`, "Kertas A4", "Kertas", 35000)
	if err != nil {
		fmt.Printf("Gagal mengisi tabel: %v\n", err)
		os.Exit(1)
	}

	// 5 - mengambil 1 data
	row := db.QueryRow(`SELECT id, nama, kategori, harga FROM produk WHERE id = $1`, 1)
	if row == nil {
		fmt.Println("Gagal mengambil data dari tabel")
	}

	var produk Produk
	err = row.Scan(&produk.ID, &produk.Nama, &produk.Kategori, &produk.Harga)
	if err != nil {
		fmt.Printf("Gagal membaca data baris: %v\n", err)
	}

	fmt.Println(produk)

	// 6 - mengambil banyak data
	rows, err := db.Query(`SELECT id, nama, kategori, harga FROM produk`)
	if err != nil || rows == nil {
		fmt.Printf("Gagal mengambil data: %v\n", err)
	}

	var produkSlice []Produk
	for rows.Next() {
		var produk Produk
		err = rows.Scan(&produk.ID, &produk.Nama, &produk.Kategori, &produk.Harga)
		if err != nil {
			fmt.Printf("Gagal membaca data baris: %v\n", err)
		}

		produkSlice = append(produkSlice, produk)
	}

	fmt.Println(produkSlice)

	// 7 - memperbarui data
	_, err = db.Exec(`
	UPDATE produk
	SET nama = $1, kategori = $2, harga = $3
	WHERE id = $4
	`, "New Kertas A5", "Kertas", 30000, 1)
	if err != nil {
		fmt.Printf("Gagal mengupdate data: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Data berhasil diupdate")

	// 8 - menghapus data
	_, err = db.Exec(`
	DELETE FROM produk
	WHERE id = $1
	`, 1)
	if err != nil {
		fmt.Printf("Gagal menghapus data: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Data berhasil dihapus")

	// 9 - transaction
	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("Gagal memulai transaction: %v\n", err)
		os.Exit(1)
	}

	_, err = tx.Exec(`DELETE FROM produk WHERE id = $1`, "a")
	if err != nil {
		fmt.Printf("Gagal menghapus data: %v\n", err)
		tx.Rollback()
		os.Exit(1)
	}

	tx.Commit()
	fmt.Println("Transaction selesai")
}
