package service
import (
	"errors"
	"inventory-management-system/internal/model"
	"inventory-management-system/internal/repository"
)

// ProductService adalah interface yang mendefinisikan business logic terkait produk
type ProductService interface {
	GetAllProducts() ([]model.Product, error)
	CreateProduct(product model.Product) error
}

type productService struct {
	// Service membutuhkan Repository untuk mengakses/menyimpan data
	repo repository.ProductRepository
}

// NewProductService mendefinisikan Constructor Service (Dependency Injection)
func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

// GetAllProducts mengambil semua data produk lewat repository
func (s *productService) GetAllProducts() ([]model.Product, error) {
	return s.repo.GetAll()
}

// CreateProduct memvalidasi input sebelum diteruskan ke repository
func (s *productService) CreateProduct(product model.Product) error {
	// --- BUSINESS LOGIC / VALIDASI ---
	
	// Validasi 1: Nama produk tidak boleh kosong
	if product.Name == "" {
		return errors.New("nama produk tidak boleh kosong")
	}

	// Validasi 2: Harga tidak boleh nol atau negatif
	if product.Price <= 0 {
		return errors.New("harga produk harus lebih besar dari 0")
	}

	// Validasi 3: Stok tidak boleh kurang dari 0
	if product.Stock < 0 {
		return errors.New("stok produk tidak boleh negatif")
	}

	// Jika semua validasi lolos, simpan ke repository
	return s.repo.Create(&product)
}