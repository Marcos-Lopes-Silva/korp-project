package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Invoice struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name       string         `json:"name"`
	Address    string         `json:"address"`
	TotalPrice float64        `json:"totalPrice"`
	SeqCode    int            `gorm:"default:nextval('seq_code')" json:"seqCode"`
	Status     string         `json:"status"`
	Items      []*InvoiceItem `gorm:"foreignKey:InvoiceID" json:"items"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-" swaggerignore:"true"`
}

func CreateInvoice(name string, address string, invoiceItems []*InvoiceItem) *Invoice {
	return &Invoice{
		Name:    name,
		Address: address,
		Status:  "Aberta",
		Items:   invoiceItems,
	}
}

func (i *Invoice) CalculateInvoiceTotalAmount() float64 {
	var invoiceTotalAmount int64 = 0
	for _, data := range i.Items {
		amount := data.CalculateTotalAmount()
		invoiceTotalAmount += amount
	}

	totalAmount := float64(invoiceTotalAmount) / 100

	return totalAmount
}
