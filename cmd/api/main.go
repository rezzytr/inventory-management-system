package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"inventory-management-system/internal/handler"
	"inventory-management-system/internal/middleware"
	"inventory-management-system/internal/repository"
	"inventory-management-system/internal/service"

	_ "github.com/glebarez/go-sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "inventory.db")
	if err != nil {
		log.Fatalf("Gagal membuka database: %v", err)
	}
	defer db.Close()

	repo, err := repository.NewSQLiteRepository(db)
	if err != nil {
		log.Fatalf("Gagal inisialisasi repository: %v", err)
	}

	svc := service.NewProductService(repo)
	productHandler := handler.NewProductHandler(svc)

	mux := http.NewServeMux()

	// INV-12: Serve Web Dashboard UI (Root /)
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		// Cari lokasi file web/index.html relatif dari Current Working Directory
		indexPath := filepath.Join(".", "web", "index.html")

		// Cek apakah file benar-benar ada di disk
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			dir, _ := os.Getwd()
			log.Printf("[ERROR] File index.html tidak ditemukan di path: %s (Working Dir: %s)", indexPath, dir)
			http.Error(w, "File Web Dashboard (web/index.html) tidak ditemukan. Pastikan folder 'web' berada di root project.", http.StatusNotFound)
			return
		}

		http.ServeFile(w, r, indexPath)
	})

	// API Endpoints
	mux.HandleFunc("GET /products", productHandler.GetAllProducts)
	mux.HandleFunc("POST /products", productHandler.CreateProduct)
	mux.HandleFunc("PUT /products/{sku}", productHandler.UpdateProduct)
	mux.HandleFunc("DELETE /products/{sku}", productHandler.DeleteProduct)

	// INV-07, INV-08, & INV-09 Endpoints
	mux.HandleFunc("POST /products/{sku}/stock", productHandler.AdjustStock)
	mux.HandleFunc("GET /products/{sku}/logs", productHandler.GetStockLogs)
	mux.HandleFunc("GET /products/low-stock", productHandler.GetLowStockProducts)

	// INV-10 Middleware Chain
	handlerWithMiddleware := middleware.Recoverer(middleware.Logger(mux))

	log.Println("Server berjalan di http://localhost:8080")
	if err := http.ListenAndServe(":8080", handlerWithMiddleware); err != nil {
		log.Fatal(err)
	}
}