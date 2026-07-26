package service

import (
	"errors"
	"inventory-management-system/internal/model"
	"inventory-management-system/internal/repository"
)

type ProductService interface {
	GetAllProducts() ([]model.Product, error)
	CreateProduct(product model.Product) error
	UpdateProduct(sku string, product model.Product) error // <-- Pastikan ini ada
	DeleteProduct(sku string) error                       // <-- Pastikan ini ada
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
	if product.Price <= 0 {
		return errors.New("harga produk harus lebih besar dari 0")
	}
	if product.Stock < 0 {
		return errors.New("stok produk tidak boleh negatif")
	}
	return s.repo.Create(product)
}

func (s *productService) UpdateProduct(sku string, product model.Product) error {
	if product.Price <= 0 {
		return errors.New("harga produk harus lebih besar dari 0")
	}
	if product.Stock < 0 {
		return errors.New("stok produk tidak boleh negatif")
	}
	return s.repo.Update(sku, product)
}

func (s *productService) DeleteProduct(sku string) error {
	return s.repo.Delete(sku)
}