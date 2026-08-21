package services

import (
	"billing/internal/apperrors"
	"billing/internal/config/clients"
	"billing/internal/models"
	"billing/internal/repository"
	"context"

	"github.com/google/uuid"
)

type InvoiceService struct {
	repo        *repository.InvoiceRepository
	stockClient *clients.StockClient
	pdfService  *PDFService
}

func NewInvoiceService(repo *repository.InvoiceRepository, stockClient *clients.StockClient) *InvoiceService {
	return &InvoiceService{repo: repo, stockClient: stockClient}
}

func (s *InvoiceService) CreateInvoice(name string, address string) (*models.Invoice, error) {
	if name == "" || address == "" {
		return nil, apperrors.ErrInvalidInput
	}

	invoice := models.CreateInvoice(name, address, []*models.InvoiceItem{})

	if err := s.repo.CreateInvoice(invoice); err != nil {
		return nil, apperrors.ErrServiceUnavailable
	}

	return invoice, nil
}

func (s *InvoiceService) GetInvoices() ([]models.Invoice, error) {
	return s.repo.GetAllInvoices()
}

func (s *InvoiceService) AddItem(context context.Context, invoiceID, productID uuid.UUID, quantity int64) (*models.Invoice, error) {
	invoice, err := s.repo.GetInvoice(invoiceID)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}

	product, err := s.stockClient.GetProduct(context, productID)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}

	err = s.stockClient.ReduceStock(context, productID, quantity)

	if err != nil {
		return nil, apperrors.ErrNotFound
	}

	item := &models.InvoiceItem{
		ProductID: productID,
		Quantity:  quantity,
		Name:      product.Name,
		Price:     int64(product.Price * 100),
	}
	invoice.Items = append(invoice.Items, item)

	if err := s.repo.UpdateInvoice(invoice); err != nil {
		return nil, apperrors.ErrServiceUnavailable
	}

	product.Quantity -= quantity
	if err := s.stockClient.RestoreStock(context, productID, quantity); err != nil {
		return nil, apperrors.ErrServiceUnavailable
	}

	return invoice, nil
}

func (s *InvoiceService) GetInvoice(id uuid.UUID) (*models.Invoice, error) {
	invoice, err := s.repo.GetInvoice(id)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	return invoice, nil
}

func (s *InvoiceService) PrintInvoice(ctx context.Context, id uuid.UUID) ([]byte, error) {
	invoice, err := s.repo.GetInvoice(id)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}

	// if invoice.Status != "Aberta" {
	// 	return nil, apperrors.ErrInvoiceNotOpened
	// }

	// for _, item := range invoice.Items {
	// 	if err := s.stockClient.ReduceStock(ctx, item.ProductID, item.Quantity); err != nil {
	// 		return nil, err
	// 	}
	// }

	invoice.Status = "Fechada"
	invoice.TotalPrice = invoice.CalculateInvoiceTotalAmount()

	if err := s.repo.UpdateInvoice(invoice); err != nil {
		return nil, apperrors.ErrServiceUnavailable
	}

	pdfBytes, err := s.pdfService.GenerateInvoicePDF(invoice)
	if err != nil {
		return nil, apperrors.ErrServiceUnavailable
	}

	return pdfBytes, nil
}
