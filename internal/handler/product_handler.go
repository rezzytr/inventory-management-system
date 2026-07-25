package handler

import (
	"encoding/json"
	"inventory-management-system/internal/model"
	"inventory-management-system/internal/service"
	"net/http"
)

// ProductHandler menangani request HTTP dan mengirimkan response HTTP
type ProductHandler struct {
	service service.ProductService
}

// NewProductHandler mendefinisikan Constructor untuk Handler
func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

// GetProducts mendefinisikan handler untuk GET /products
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	// 1. Set Response Header menjadi JSON
	w.Header().Set("Content-Type", "application/json")

	// 2. Ambil data dari Service Layer
	products, err := h.service.GetAllProducts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// 3. Kirim response OK (200) beserta data produk dalam format JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

// CreateProduct mendefinisikan handler untuk POST /products
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Decode body JSON dari client ke struct Product
	var req model.Product
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Format JSON tidak valid"})
		return
	}

	// 2. Panggil Service Layer
	err = h.service.CreateProduct(req)
	if err != nil {
		// Jika ada error validasi dari service, kirim status Bad Request (400)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// 3. Response Created (201)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Produk berhasil dibuat"})
}