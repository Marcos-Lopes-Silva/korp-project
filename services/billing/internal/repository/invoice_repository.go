package repository

import (
	"billing/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

func (r *InvoiceRepository) CreateInvoice(invoice *models.Invoice) error {
	return r.db.Create(invoice).Error
}

func (r *InvoiceRepository) GetInvoice(id uuid.UUID) (*models.Invoice, error) {
	var invoice models.Invoice
	if err := r.db.Preload("Items").First(&invoice, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *InvoiceRepository) UpdateInvoice(invoice *models.Invoice) error {
	return r.db.Save(invoice).Error
}
