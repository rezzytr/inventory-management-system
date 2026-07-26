package repository

import (
	"fmt"
	"inventory-management-system/internal/model"
)

type MockProductRepository struct {
	Products     map[string]model.Product
	Transactions map[string][]model.StockTransaction
}

func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{
		Products:     make(map[string]model.Product),
		Transactions: make(map[string][]model.StockTransaction),
	}
}

func (m *MockProductRepository) GetAll() ([]model.Product, error) {
	var list []model.Product
	for _, p := range m.Products {
		list = append(list, p)
	}
	return list, nil
}

func (m *MockProductRepository) Create(product *model.Product) error {
	if _, exists := m.Products[product.SKU]; exists {
		return fmt.Errorf("SKU sudah ada")
	}
	m.Products[product.SKU] = *product
	return nil
}

func (m *MockProductRepository) Update(sku string, product model.Product) error {
	if _, exists := m.Products[sku]; !exists {
		return fmt.Errorf("produk tidak ditemukan")
	}
	m.Products[sku] = product
	return nil
}

func (m *MockProductRepository) Delete(sku string) error {
	if _, exists := m.Products[sku]; !exists {
		return fmt.Errorf("produk tidak ditemukan")
	}
	delete(m.Products, sku)
	return nil
}

func (m *MockProductRepository) UpdateStock(sku string, amount int, reason string) error {
	p, exists := m.Products[sku]
	if !exists {
		return fmt.Errorf("produk tidak ditemukan")
	}

	newStock := p.Stock + amount
	if newStock < 0 {
		return fmt.Errorf("produk tidak ditemukan atau jumlah stok tidak mencukupi")
	}

	p.Stock = newStock
	m.Products[sku] = p

	txType := "IN"
	absAmount := amount
	if amount < 0 {
		txType = "OUT"
		absAmount = -amount
	}

	m.Transactions[sku] = append(m.Transactions[sku], model.StockTransaction{
		SKU:    sku,
		Type:   txType,
		Amount: absAmount,
		Reason: reason,
	})

	return nil
}

func (m *MockProductRepository) GetStockTransactions(sku string) ([]model.StockTransaction, error) {
	return m.Transactions[sku], nil
}

func (m *MockProductRepository) GetLowStockProducts(threshold int) ([]model.Product, error) {
	var list []model.Product
	for _, p := range m.Products {
		if p.Stock <= threshold {
			list = append(list, p)
		}
	}
	return list, nil
}