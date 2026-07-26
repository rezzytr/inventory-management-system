package repository

import (
	"errors"
	"inventory-management-system/internal/model"
)

type ProductRepository interface {
	GetAll() ([]model.Product, error)
	GetBySKU(sku string) (model.Product, error)
	Create(product model.Product) error
	Update(sku string, product model.Product) error
	Delete(sku string) error
}

type memoryRepository struct {
	products map[string]model.Product
}

func NewMemoryRepository() ProductRepository {
	return &memoryRepository{
		products: make(map[string]model.Product),
	}
}

func (r *memoryRepository) GetAll() ([]model.Product, error) {
	var list []model.Product
	for _, p := range r.products {
		list = append(list, p)
	}
	return list, nil
}

func (r *memoryRepository) GetBySKU(sku string) (model.Product, error) {
	p, exists := r.products[sku]
	if !exists {
		return model.Product{}, errors.New("produk tidak ditemukan")
	}
	return p, nil
}

func (r *memoryRepository) Create(product model.Product) error {
	r.products[product.SKU] = product
	return nil
}

func (r *memoryRepository) Update(sku string, product model.Product) error {
	if _, exists := r.products[sku]; !exists {
		return errors.New("produk tidak ditemukan")
	}
	// Pertahankan SKU agar tidak berubah
	product.SKU = sku
	r.products[sku] = product
	return nil
}

func (r *memoryRepository) Delete(sku string) error {
	if _, exists := r.products[sku]; !exists {
		return errors.New("produk tidak ditemukan")
	}
	delete(r.products, sku)
	return nil
}