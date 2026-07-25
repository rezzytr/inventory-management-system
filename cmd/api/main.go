package main

import (
	"fmt"
	"log"
	"net/http"

	"inventory-management-system/internal/handler"
	"inventory-management-system/internal/repository"
	"inventory-management-system/internal/service"
)

func main() {
	// 1. Inisialisasi Layer (Repository -> Service -> Handler)
	repo := repository.NewMemoryRepository()
	productService := service.NewProductService(repo)
	productHandler := handler.NewProductHandler(productService)

	// 2. Routing HTTP (Penting: Fitur Mux baru di Go 1.22+)
	http.HandleFunc("GET /products", productHandler.GetProducts)
	http.HandleFunc("POST /products", productHandler.CreateProduct)

	// 3. Jalankan HTTP Server di port 8080
	port := ":8080"
	fmt.Printf("🚀 Server berjalan di http://localhost%s\n", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}