package models

import (
	"github.com/google/uuid"
)

type Invoice struct {
	ID         uuid.UUID
	Name       string
	Address    string
	TotalPrice float64
	SeqCode    int
	Status     string
	Items      []*InvoiceItem
}

func CreateInvoice(name string, address string, invoiceItems []*InvoiceItem) *Invoice {
	return &Invoice{
		Name:    name,
		Address: address,
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
