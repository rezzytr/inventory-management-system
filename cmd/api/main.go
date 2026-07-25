package main

import (
	"fmt"
	"log"

	// Import package repository dan model sesuai path module di go.mod
	"inventory-management-system/internal/model"
	"inventory-management-system/internal/repository"
)

func main() {
	fmt.Println("=== Inventory Management System API ===")

	// 1. Inisialisasi layer repository
	repo := repository.NewMemoryRepository()

	// 2. Tambah produk baru
	newProduct := model.Product{
		SKU:   "MOU-002",
		Name:  "Wireless Mouse",
		Stock: 25,
		Price: 250000,
	}

	err := repo.Create(&newProduct)
	if err != nil {
		log.Fatalf("Gagal menambahkan produk: %v", err)
	}
	fmt.Println("✅ Berhasil menambahkan produk baru!")

	// 3. Ambil dan tampilkan semua data produk
	products, err := repo.GetAll()
	if err != nil {
		log.Fatalf("Gagal mengambil data produk: %v", err)
	}

	fmt.Println("\n--- Daftar Inventaris ---")
	for _, p := range products {
		fmt.Printf("[%d] SKU: %s | Name: %s | Stock: %d | Price: Rp%.0f\n",
			p.ID, p.SKU, p.Name, p.Stock, p.Price)
	}
}