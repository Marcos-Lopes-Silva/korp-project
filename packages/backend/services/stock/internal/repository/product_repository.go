package repository

import (
	"stock/internal/models"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) CreateProduct(newProduct *models.Product) (*models.Product, error) {

	if err := r.db.Create(newProduct).Error; err != nil {
		return nil, err
	}

	return newProduct, nil
}

func (r *ProductRepository) GetProductBySKU(sku string) (*models.Product, error) {
	var product models.Product
	if err := r.db.Where("sku = ?", sku).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetProductByID(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) UpdateProduct(id uuid.UUID, product *models.Product) (*models.Product, error) {
	if err := r.db.Model(&models.Product{}).Where("id = ?", id).Updates(product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductRepository) DeleteProduct(id uuid.UUID) error {
	if err := r.db.Delete(&models.Product{}, id).Error; err != nil {
		return err
	}
	return nil
}
