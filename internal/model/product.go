package model

import "time"

type Product struct {
	ID        int64     `json:"id"`
	SKU       string    `json:"sku"`
	Name      string    `json:"name"`
	Stock     int       `json:"stock"`
	Price     int64     `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StockAdjustment struct {
	Amount int    `json:"amount"` // Nilai positif (+10) atau negatif (-5)
	Reason string `json:"reason"` // Contoh: "Restock Suplier", "Penjualan Toko", "Barang Rusak"
}

type StockTransaction struct {
	ID        int64     `json:"id"`
	SKU       string    `json:"sku"`
	Type      string    `json:"type"`   // "IN" atau "OUT"
	Amount    int       `json:"amount"` // Selalu bernilai positif untuk jumlah transaksi
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}