package repository
import (
	"errors"
	"inventory-management-system/internal/model" // Sesuaikan dengan module path di go.mod kamu
	"time"
)

// ProductRepository mendefinisikan interface/kontrak fungsi apa saja yang harus dimiliki oleh repository.
// Ini memudahkan proses Testing (Mocking) di industri.
type ProductRepository interface {
	GetAll() ([]model.Product, error)
	GetByID(id int) (model.Product, error)
	Create(product *model.Product) error
}

// memoryRepository adalah implementasi konkrit dari ProductRepository yang menyimpan data di slice/memory.
type memoryRepository struct {
	products []model.Product
}

// NewMemoryRepository adalah Constructor Function untuk membuat instance baru dari repository.
func NewMemoryRepository() ProductRepository {
	return &memoryRepository{
		// Kita isi dummy data awal untuk testing
		products: []model.Product{
			{
				ID:        1,
				SKU:       "LAP-001",
				Name:      "Laptop Gaming",
				Stock:     10,
				Price:     15000000,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}
}

// GetAll mengambil seluruh daftar produk dari memory
func (r *memoryRepository) GetAll() ([]model.Product, error) {
	return r.products, nil
}

// GetByID mencari produk berdasarkan ID uniknya
func (r *memoryRepository) GetByID(id int) (model.Product, error) {
	for _, p := range r.products {
		if p.ID == id {
			return p, nil
		}
	}
	// Mengembalikan error standar jika barang tidak ditemukan
	return model.Product{}, errors.New("product not found")
}

// Create menambahkan produk baru ke dalam list memory
func (r *memoryRepository) Create(product *model.Product) error {
	// Generasi ID sederhana (Auto increment manual)
	product.ID = len(r.products) + 1
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	r.products = append(r.products, *product)
	return nil
}