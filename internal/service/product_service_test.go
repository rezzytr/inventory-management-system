package service_test

import (
	"testing"

	"inventory-management-system/internal/model"
	"inventory-management-system/internal/repository"
	"inventory-management-system/internal/service"
)

func TestCreateProduct_Success(t *testing.T) {
	mockRepo := repository.NewMockProductRepository()
	svc := service.NewProductService(mockRepo)

	newProduct := model.Product{
		SKU:   "TEST-001",
		Name:  "Kopi Testing",
		Stock: 10,
		Price: 15000,
	}

	err := svc.CreateProduct(newProduct)
	if err != nil {
		t.Fatalf("diharapkan sukses, tetapi dapat error: %v", err)
	}

	products, _ := svc.GetAllProducts()
	if len(products) != 1 {
		t.Errorf("diharapkan jumlah produk 1, dapat: %d", len(products))
	}
}

func TestCreateProduct_InvalidValidation(t *testing.T) {
	mockRepo := repository.NewMockProductRepository()
	svc := service.NewProductService(mockRepo)

	// Kasus 1: SKU Kosong
	invalidProduct := model.Product{
		SKU:   "",
		Name:  "Kopi Tanpa SKU",
		Stock: 10,
		Price: 15000,
	}

	err := svc.CreateProduct(invalidProduct)
	if err == nil {
		t.Errorf("diharapkan error karena SKU kosong, tetapi malah sukses")
	}

	// Kasus 2: Harga Negatif
	negativePriceProduct := model.Product{
		SKU:   "TEST-NEG",
		Name:  "Kopi Diskon Gila",
		Stock: 10,
		Price: -5000,
	}

	err = svc.CreateProduct(negativePriceProduct)
	if err == nil {
		t.Errorf("diharapkan error karena harga negatif, tetapi malah sukses")
	}
}

func TestAdjustStock_NegativeStockProtection(t *testing.T) {
	mockRepo := repository.NewMockProductRepository()
	svc := service.NewProductService(mockRepo)

	// Persiapkan produk dengan stok awal 10
	_ = svc.CreateProduct(model.Product{
		SKU:   "KOP-001",
		Name:  "Kopi Robusta",
		Stock: 10,
		Price: 20000,
	})

	// Kurangi stok melebihi sisa stok (-15) -> Harus Gagal!
	err := svc.AdjustStock("KOP-001", -15, "Penjualan Besar")
	if err == nil {
		t.Errorf("diharapkan error stok tidak mencukupi, tetapi transaksi berhasil")
	}

	// Kurangi stok yang valid (-5) -> Harus Sukses!
	err = svc.AdjustStock("KOP-001", -5, "Penjualan Valid")
	if err != nil {
		t.Fatalf("diharapkan sukses, tetapi dapat error: %v", err)
	}

	// Cek sisa stok (harus 5)
	products, _ := svc.GetAllProducts()
	if products[0].Stock != 5 {
		t.Errorf("diharapkan sisa stok 5, tetapi dapat: %d", products[0].Stock)
	}
}

func TestGetLowStockProducts(t *testing.T) {
	mockRepo := repository.NewMockProductRepository()
	svc := service.NewProductService(mockRepo)

	_ = svc.CreateProduct(model.Product{SKU: "A1", Name: "Barang Stok Tipis", Stock: 3, Price: 1000})
	_ = svc.CreateProduct(model.Product{SKU: "A2", Name: "Barang Stok Aman", Stock: 20, Price: 1000})

	lowStockList, err := svc.GetLowStockProducts(5)
	if err != nil {
		t.Fatalf("gagal mengambil data low stock: %v", err)
	}

	if len(lowStockList) != 1 {
		t.Errorf("diharapkan 1 produk stok tipis, dapat: %d", len(lowStockList))
	}
}