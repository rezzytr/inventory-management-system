package repository

import (
	"database/sql"
	"fmt"
	"time"

	"inventory-management-system/internal/model"

	_ "github.com/glebarez/go-sqlite"
)

type sqliteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) (ProductRepository, error) {
	repo := &sqliteRepository{db: db}

	queryProducts := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sku TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		stock INTEGER NOT NULL DEFAULT 0,
		price INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(queryProducts); err != nil {
		return nil, fmt.Errorf("gagal inisialisasi tabel products: %w", err)
	}

	queryTransactions := `
	CREATE TABLE IF NOT EXISTS stock_transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sku TEXT NOT NULL,
		type TEXT NOT NULL,
		amount INTEGER NOT NULL,
		reason TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (sku) REFERENCES products(sku)
	);`

	if _, err := db.Exec(queryTransactions); err != nil {
		return nil, fmt.Errorf("gagal inisialisasi tabel stock_transactions: %w", err)
	}

	return repo, nil
}

func (r *sqliteRepository) GetAll() ([]model.Product, error) {
	query := `SELECT id, sku, name, stock, price, created_at, updated_at FROM products`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data: %w", err)
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Stock, &p.Price, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("gagal membaca baris: %w", err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterasi baris produk: %w", err)
	}

	return products, nil
}

func (r *sqliteRepository) Create(product *model.Product) error {
	now := time.Now()
	query := `
		INSERT INTO products (sku, name, stock, price, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, product.SKU, product.Name, product.Stock, product.Price, now, now)
	if err != nil {
		return fmt.Errorf("gagal menyimpan produk: %w", err)
	}

	id, err := result.LastInsertId()
	if err == nil {
		product.ID = id
		product.CreatedAt = now
		product.UpdatedAt = now
	}

	return nil
}

func (r *sqliteRepository) Update(sku string, product model.Product) error {
	now := time.Now()
	query := `UPDATE products SET name = ?, stock = ?, price = ?, updated_at = ? WHERE sku = ?`
	result, err := r.db.Exec(query, product.Name, product.Stock, product.Price, now, sku)
	if err != nil {
		return fmt.Errorf("gagal mengupdate produk: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("produk dengan SKU %s tidak ditemukan", sku)
	}

	return nil
}

func (r *sqliteRepository) Delete(sku string) error {
	query := `DELETE FROM products WHERE sku = ?`
	result, err := r.db.Exec(query, sku)
	if err != nil {
		return fmt.Errorf("gagal menghapus produk: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("produk dengan SKU %s tidak ditemukan", sku)
	}

	return nil
}

func (r *sqliteRepository) UpdateStock(sku string, amount int, reason string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi database: %w", err)
	}
	defer tx.Rollback()

	updateQuery := `
		UPDATE products 
		SET stock = stock + ?, updated_at = ? 
		WHERE sku = ? AND (stock + ?) >= 0
	`
	now := time.Now()
	result, err := tx.Exec(updateQuery, amount, now, sku, amount)
	if err != nil {
		return fmt.Errorf("gagal memperbarui stok: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("produk tidak ditemukan atau jumlah stok tidak mencukupi")
	}

	txType := "IN"
	absAmount := amount
	if amount < 0 {
		txType = "OUT"
		absAmount = -amount
	}

	logQuery := `
		INSERT INTO stock_transactions (sku, type, amount, reason, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	if _, err := tx.Exec(logQuery, sku, txType, absAmount, reason, now); err != nil {
		return fmt.Errorf("gagal mencatat riwayat transaksi stok: %w", err)
	}

	return tx.Commit()
}

func (r *sqliteRepository) GetStockTransactions(sku string) ([]model.StockTransaction, error) {
	query := `
		SELECT id, sku, type, amount, COALESCE(reason, ''), created_at 
		FROM stock_transactions 
		WHERE sku = ? 
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, sku)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil log transaksi: %w", err)
	}
	defer rows.Close()

	var transactions []model.StockTransaction
	for rows.Next() {
		var t model.StockTransaction
		if err := rows.Scan(&t.ID, &t.SKU, &t.Type, &t.Amount, &t.Reason, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("gagal membaca data log: %w", err)
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterasi baris log: %w", err)
	}

	return transactions, nil
}

func (r *sqliteRepository) GetLowStockProducts(threshold int) ([]model.Product, error) {
	query := `
		SELECT id, sku, name, stock, price, created_at, updated_at 
		FROM products 
		WHERE stock <= ? 
		ORDER BY stock ASC
	`
	rows, err := r.db.Query(query, threshold)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data produk stok tipis: %w", err)
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Stock, &p.Price, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("gagal membaca baris produk stok tipis: %w", err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterasi baris produk stok tipis: %w", err)
	}

	return products, nil
}