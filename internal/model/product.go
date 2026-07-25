package model

import "time"

// Product mendefinisikan struktur data barang dalam sistem inventaris.
// Di dunia kerja, struct ini merepresentasikan tabel di database.
type Product struct {
	// ID unik untuk setiap produk (Primary Key)
	ID int `json:"id"`

	// SKU (Stock Keeping Unit) adalah kode unik barang di gudang (misal: "PROD-001")
	SKU string `json:"sku"`

	// Nama barang
	Name string `json:"name"`

	// Jumlah stok yang tersedia di gudang
	Stock int `json:"stock"`

	// Harga barang per unit
	Price float64 `json:"price"`

	// Waktu kapan data produk ini dibuat dan terakhir diubah
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}