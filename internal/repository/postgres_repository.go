package repository

import (
	"database/sql"
	"fmt"
	"time"

	"inventory-management-system/internal/model"

	_ "github.com/lib/pq"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) (ProductRepository, error) {
	repo := &postgresRepository{db: db}

	queryProducts := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		sku VARCHAR(100) UNIQUE NOT NULL,
		name VARCHAR(255) NOT NULL,
		stock INT NOT NULL DEFAULT 0,
		price BIGINT NOT NULL DEFAULT 0,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(queryProducts); err != nil {
		return nil, fmt.Errorf("gagal inisialisasi tabel products postgres: %w", err)
	}

	queryTransactions := `
	CREATE TABLE IF NOT EXISTS stock_transactions (
		id SERIAL PRIMARY KEY,
		sku VARCHAR(100) NOT NULL,
		type VARCHAR(10) NOT NULL,
		amount INT NOT NULL,
		reason TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (sku) REFERENCES products(sku)
	);`

	if _, err := db.Exec(queryTransactions); err != nil {
		return nil, fmt.Errorf("gagal inisialisasi tabel stock_transactions postgres: %w", err)
	}

	return repo, nil
}

func (r *postgresRepository) GetAll() ([]model.Product, error) {
	query := `SELECT id, sku, name, stock, price, created_at, updated_at FROM products`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data dari postgres: %w", err)
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Stock, &p.Price, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("gagal membaca baris postgres: %w", err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterasi baris produk postgres: %w", err)
	}

	return products, nil
}

func (r *postgresRepository) Create(product *model.Product) error {
	now := time.Now()
	query := `
		INSERT INTO products (sku, name, stock, price, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err := r.db.QueryRow(query, product.SKU, product.Name, product.Stock, product.Price, now, now).Scan(&product.ID)
	if err != nil {
		return fmt.Errorf("gagal menyimpan produk ke postgres: %w", err)
	}

	product.CreatedAt = now
	product.UpdatedAt = now

	return nil
}

func (r *postgresRepository) Update(sku string, product model.Product) error {
	now := time.Now()
	query := `UPDATE products SET name = $1, stock = $2, price = $3, updated_at = $4 WHERE sku = $5`
	result, err := r.db.Exec(query, product.Name, product.Stock, product.Price, now, sku)
	if err != nil {
		return fmt.Errorf("gagal mengupdate produk di postgres: %w", err)
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

func (r *postgresRepository) Delete(sku string) error {
	query := `DELETE FROM products WHERE sku = $1`
	result, err := r.db.Exec(query, sku)
	if err != nil {
		return fmt.Errorf("gagal menghapus produk di postgres: %w", err)
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

func (r *postgresRepository) UpdateStock(sku string, amount int, reason string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi database: %w", err)
	}
	defer tx.Rollback()

	updateQuery := `
		UPDATE products 
		SET stock = stock + $1, updated_at = $2 
		WHERE sku = $3 AND (stock + $4) >= 0
	`
	now := time.Now()
	result, err := tx.Exec(updateQuery, amount, now, sku, amount)
	if err != nil {
		return fmt.Errorf("gagal memperbarui stok di postgres: %w", err)
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
		VALUES ($1, $2, $3, $4, $5)
	`
	if _, err := tx.Exec(logQuery, sku, txType, absAmount, reason, now); err != nil {
		return fmt.Errorf("gagal mencatat riwayat transaksi stok di postgres: %w", err)
	}

	return tx.Commit()
}

func (r *postgresRepository) GetStockTransactions(sku string) ([]model.StockTransaction, error) {
	query := `
		SELECT id, sku, type, amount, COALESCE(reason, ''), created_at 
		FROM stock_transactions 
		WHERE sku = $1 
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, sku)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil log transaksi dari postgres: %w", err)
	}
	defer rows.Close()

	var transactions []model.StockTransaction
	for rows.Next() {
		var t model.StockTransaction
		if err := rows.Scan(&t.ID, &t.SKU, &t.Type, &t.Amount, &t.Reason, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("gagal membaca data log postgres: %w", err)
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterasi baris log postgres: %w", err)
	}

	return transactions, nil
}

func (r *postgresRepository) GetLowStockProducts(threshold int) ([]model.Product, error) {
	query := `
		SELECT id, sku, name, stock, price, created_at, updated_at 
		FROM products 
		WHERE stock <= $1 
		ORDER BY stock ASC
	`
	rows, err := r.db.Query(query, threshold)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data produk stok tipis postgres: %w", err)
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Stock, &p.Price, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("gagal membaca baris produk stok tipis postgres: %w", err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterasi baris produk stok tipis postgres: %w", err)
	}

	return products, nil
}