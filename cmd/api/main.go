package main

import (
	"fmt"
	"log"

	"inventory-management-system/internal/model"
	"inventory-management-system/internal/repository"
	"inventory-management-system/internal/service"
)

func main() {
	fmt.Println("=== Inventory Management System API ===")

	// 1. Init Repository
	repo := repository.NewMemoryRepository()

	// 2. Init Service (memasukkan Repository ke dalam Service)
	productService := service.NewProductService(repo)

	// 3. Skenario 1: Coba tambah produk yang VALID
	validProduct := model.Product{
		SKU:   "MOU-002",
		Name:  "Wireless Mouse",
		Stock: 15,
		Price: 150000,
	}

	err := productService.CreateProduct(validProduct)
	if err != nil {
		log.Printf("❌ Gagal membuat produk valid: %v\n", err)
	} else {
		fmt.Println("✅ Produk berhasil ditambahkan!")
	}

	// 4. Skenario 2: Coba tambah produk INVALID (Harga 0) -> Harusnya Kena Validasi!
	invalidProduct := model.Product{
		SKU:   "KEY-001",
		Name:  "Mechanical Keyboard",
		Stock: 5,
		Price: 0, // Invalid!
	}

	err = productService.CreateProduct(invalidProduct)
	if err != nil {
		fmt.Printf("⚠️ Validasi Bekerja! Error: %v\n", err)
	}

	// 5. Tampilkan semua produk yang ada saat ini
	products, _ := productService.GetAllProducts()
	fmt.Println("\n--- Daftar Produk Saat Ini ---")
	for _, p := range products {
		fmt.Printf("[%d] SKU: %s | %s | Stok: %d | Rp%.0f\n", p.ID, p.SKU, p.Name, p.Stock, p.Price)
	}
}