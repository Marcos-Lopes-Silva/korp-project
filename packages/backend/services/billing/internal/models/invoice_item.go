package models

import (
	"errors"

	"github.com/google/uuid"
)

type InvoiceItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	InvoiceID uuid.UUID `gorm:"type:uuid" json:"invoiceId"`
	ProductID uuid.UUID `json:"productId"`
	Quantity  int64     `json:"quantity"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
}

func (i *InvoiceItem) CalculateTotalAmount() int64 {
	totalAmount := i.Quantity * i.Price
	return totalAmount
}

func (i *InvoiceItem) ReturnItemTotalAmount() float64 {
	totalAmount := i.CalculateTotalAmount()
	converted := float64(totalAmount) / 100
	return converted
}

func (i *InvoiceItem) ReturnItemPrice() float64 {
	converted := float64(i.Price) / 100
	return converted
}

func NewInvoiceData(title string, qty int64, price interface{}) (*InvoiceItem, error) {
	var convertedPrice int64

	switch priceValue := price.(type) {
	case int64:
		convertedPrice = priceValue * 100
	case int:
		convertedPrice = int64(priceValue * 100)
	case float32:
		convertedPrice = int64(priceValue * 100)
	case float64:
		convertedPrice = int64(priceValue * 100)
	default:
		return nil, errors.New("type not permitted")
	}

	return &InvoiceItem{
		Name:     title,
		Quantity: qty,
		Price:    convertedPrice,
	}, nil
}
