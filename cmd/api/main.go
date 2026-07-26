package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"inventory-management-system/internal/handler"
	"inventory-management-system/internal/middleware"
	"inventory-management-system/internal/repository"
	"inventory-management-system/internal/service"

	_ "github.com/glebarez/go-sqlite"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load .env file (jika tidak ada, log warning tapi jangan panic agar tetap bisa jalan di environment server)
	if err := godotenv.Load(); err != nil {
		log.Println("[INFO] File .env tidak ditemukan, menggunakan environment variable sistem")
	}

	// 2. Baca Konfigurasi dari Environment Variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default fallback
	}

	dbDriver := os.Getenv("DB_DRIVER")
	if dbDriver == "" {
		dbDriver = "sqlite"
	}

	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		dbSource = "inventory.db"
	}

	// 3. Inisialisasi Database
	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("Gagal membuka database (%s): %v", dbDriver, err)
	}
	defer db.Close()

	// 4. Inisialisasi Repository, Service, & Handler
	repo, err := repository.NewSQLiteRepository(db)
	if err != nil {
		log.Fatalf("Gagal inisialisasi repository: %v", err)
	}

	svc := service.NewProductService(repo)
	productHandler := handler.NewProductHandler(svc)

	mux := http.NewServeMux()

	// 5. Serve Web Dashboard UI (Root /)
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		indexPath := filepath.Join(".", "web", "index.html")
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			http.Error(w, "File Web Dashboard (web/index.html) tidak ditemukan.", http.StatusNotFound)
			return
		}
		http.ServeFile(w, r, indexPath)
	})

	// 6. API Endpoints
	mux.HandleFunc("GET /products", productHandler.GetAllProducts)
	mux.HandleFunc("POST /products", productHandler.CreateProduct)
	mux.HandleFunc("PUT /products/{sku}", productHandler.UpdateProduct)
	mux.HandleFunc("DELETE /products/{sku}", productHandler.DeleteProduct)

	mux.HandleFunc("POST /products/{sku}/stock", productHandler.AdjustStock)
	mux.HandleFunc("GET /products/{sku}/logs", productHandler.GetStockLogs)
	mux.HandleFunc("GET /products/low-stock", productHandler.GetLowStockProducts)

	// 7. Middleware Chain
	handlerWithMiddleware := middleware.Recoverer(middleware.Logger(mux))

	// 8. Start Server
	serverAddr := fmt.Sprintf(":%s", port)
	log.Printf("Server berjalan di http://localhost%s (DB Driver: %s)", serverAddr, dbDriver)
	if err := http.ListenAndServe(serverAddr, handlerWithMiddleware); err != nil {
		log.Fatal(err)
	}
}