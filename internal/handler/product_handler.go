package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"inventory-management-system/internal/model"
	"inventory-management-system/internal/service"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(s service.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	products, err := h.service.GetAllProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cegah return 'null' saat database kosong, ubah jadi array kosong '[]'
	if products == nil {
		products = []model.Product{}
	}

	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Format JSON tidak valid", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateProduct(product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sku := r.PathValue("sku")

	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Format JSON tidak valid", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateProduct(sku, product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Produk berhasil diperbarui"})
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sku := r.PathValue("sku")

	if err := h.service.DeleteProduct(sku); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Produk berhasil dihapus"})
}

func (h *ProductHandler) AdjustStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sku := r.PathValue("sku")

	var req model.StockAdjustment
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Format JSON tidak valid", http.StatusBadRequest)
		return
	}

	if err := h.service.AdjustStock(sku, req.Amount, req.Reason); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Stok berhasil diperbarui & dicatat ke log"})
}

func (h *ProductHandler) GetStockLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sku := r.PathValue("sku")

	logs, err := h.service.GetStockLogs(sku)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if logs == nil {
		logs = []model.StockTransaction{}
	}

	json.NewEncoder(w).Encode(logs)
}

func (h *ProductHandler) GetLowStockProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	threshold := 5
	thresholdStr := r.URL.Query().Get("threshold")
	if thresholdStr != "" {
		parsed, err := strconv.Atoi(thresholdStr)
		if err != nil || parsed < 0 {
			http.Error(w, "Query parameter threshold harus berupa angka positif", http.StatusBadRequest)
			return
		}
		threshold = parsed
	}

	products, err := h.service.GetLowStockProducts(threshold)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if products == nil {
		products = []model.Product{}
	}

	json.NewEncoder(w).Encode(products)
}