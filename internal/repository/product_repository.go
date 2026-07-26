package repository

import "inventory-management-system/internal/model"

type ProductRepository interface {
	GetAll() ([]model.Product, error)
	Create(product *model.Product) error
	Update(sku string, product model.Product) error
	Delete(sku string) error
	UpdateStock(sku string, amount int, reason string) error
	GetStockTransactions(sku string) ([]model.StockTransaction, error)
	GetLowStockProducts(threshold int) ([]model.Product, error)
}