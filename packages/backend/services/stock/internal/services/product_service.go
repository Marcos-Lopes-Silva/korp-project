package services

import (
	"stock/internal/apperrors"
	"stock/internal/models"
	"stock/internal/repository"

	"github.com/google/uuid"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetProducts() ([]models.Product, error) {
	return s.repo.GetAllProducts()
}

func (s *ProductService) CreateProduct(product *models.Product) (*models.Product, error) {
	if product.Name == "" || product.SKU == "" || product.Quantity < 0 || product.Price <= 0 {
		return nil, apperrors.ErrInvalidInput
	}

	existing, _ := s.repo.GetProductBySKU(product.SKU)

	if existing != nil {
		return nil, apperrors.ErrDuplicated
	}

	product, err := s.repo.CreateProduct(product)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) GetProductByID(id uuid.UUID) (*models.Product, error) {
	product, err := s.repo.GetProductByID(id)

	if err != nil {
		return nil, apperrors.ErrNotFound
	}

	return product, nil
}

func (s *ProductService) UpdateProduct(id uuid.UUID, product *models.Product) (*models.Product, error) {
	if product.Name == "" || product.SKU == "" || product.Quantity < 0 || product.Price <= 0 {
		return nil, apperrors.ErrInvalidInput
	}

	existing, _ := s.repo.GetProductBySKU(product.SKU)

	if existing != nil && existing.ID != id {
		return nil, apperrors.ErrDuplicated
	}

	updatedProduct, err := s.repo.UpdateProduct(id, product)

	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (s *ProductService) DeleteProduct(id uuid.UUID) error {
	err := s.repo.DeleteProduct(id)

	if err != nil {
		return apperrors.ErrNotFound
	}

	return nil
}

func (s *ProductService) ReduceStock(id uuid.UUID, quantity int64) error {
	product, err := s.repo.GetProductByID(id)
	if err != nil {
		return apperrors.ErrNotFound
	}

	if err := product.ReduceStock(quantity); err != nil {
		return err
	}

	_, err = s.repo.UpdateProduct(id, product)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProductService) RestoreStock(id uuid.UUID, quantity int64) error {
	product, err := s.repo.GetProductByID(id)
	if err != nil {
		return apperrors.ErrNotFound
	}

	if err := product.RestoreStock(quantity); err != nil {
		return err
	}

	_, err = s.repo.UpdateProduct(id, product)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProductService) VerifyAvailability(uuid uuid.UUID, quantity int64) (bool, error) {
	product, err := s.repo.GetProductByID(uuid)
	if err != nil {
		return false, apperrors.ErrNotFound
	}

	return product.Quantity >= quantity, nil
}
