package models

import (
	"errors"

	"github.com/google/uuid"
)

type InvoiceItem struct {
	ProductID uuid.UUID
	Name      string
	Quantity  int64
	Price     int64
	SKU       string
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
