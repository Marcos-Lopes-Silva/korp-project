package models

import (
	"stock/internal/apperrors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string         `json:"name" binding:"required"`
	SKU       string         `json:"sku" binding:"required"`
	Quantity  int            `json:"quantity" binding:"gte=0"`
	Price     float64        `json:"price" binding:"gt=0"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index" swaggerignore:"true"`
}

func (p *Product) ReduceStock(qty int) error {
	if qty <= 0 {
		return apperrors.ErrInvalidInput
	}
	if p.Quantity < qty {
		return apperrors.ErrInsufficientStock
	}
	p.Quantity -= qty
	return nil
}

func (p *Product) RestoreStock(qty int) error {
	if qty <= 0 {
		return apperrors.ErrInvalidInput
	}
	p.Quantity += qty
	return nil
}
