package service

import (
	"fmt"
	"inventory-management-system/internal/model"
	"inventory-management-system/internal/repository"
)

type ProductService interface {
	GetAllProducts() ([]model.Product, error)
	CreateProduct(product model.Product) error
	UpdateProduct(sku string, product model.Product) error
	DeleteProduct(sku string) error
	AdjustStock(sku string, amount int, reason string) error
	GetStockLogs(sku string) ([]model.StockTransaction, error)
	GetLowStockProducts(threshold int) ([]model.Product, error)
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetAllProducts() ([]model.Product, error) {
	return s.repo.GetAll()
}

func (s *productService) CreateProduct(product model.Product) error {
	if product.SKU == "" || product.Name == "" {
		return fmt.Errorf("SKU dan Nama Produk tidak boleh kosong")
	}
	if product.Price < 0 || product.Stock < 0 {
		return fmt.Errorf("Harga dan Stok tidak boleh bernilai negatif")
	}
	return s.repo.Create(&product)
}

func (s *productService) UpdateProduct(sku string, product model.Product) error {
	if product.Price < 0 || product.Stock < 0 {
		return fmt.Errorf("Harga dan Stok tidak boleh bernilai negatif")
	}
	return s.repo.Update(sku, product)
}

func (s *productService) DeleteProduct(sku string) error {
	return s.repo.Delete(sku)
}

func (s *productService) AdjustStock(sku string, amount int, reason string) error {
	if amount == 0 {
		return fmt.Errorf("jumlah perubahan stok tidak boleh nol (0)")
	}
	if reason == "" {
		reason = "Penyesuaian Stok Umum"
	}
	return s.repo.UpdateStock(sku, amount, reason)
}

func (s *productService) GetStockLogs(sku string) ([]model.StockTransaction, error) {
	if sku == "" {
		return nil, fmt.Errorf("SKU wajib diisi")
	}
	return s.repo.GetStockTransactions(sku)
}

func (s *productService) GetLowStockProducts(threshold int) ([]model.Product, error) {
	if threshold < 0 {
		return nil, fmt.Errorf("threshold tidak boleh negatif")
	}
	return s.repo.GetLowStockProducts(threshold)
}