package services

import (
	"billing/internal/models"
	"bytes"
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

type PDFService struct{}

func NewPDFService() *PDFService {
	return &PDFService{}
}

func (s *PDFService) GenerateInvoicePDF(invoice *models.Invoice) ([]byte, error) {

	pdf := gofpdf.New("P", "mm", "A4", "")
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.AddPage()
	pdf.SetMargins(15, 15, 15)

	pdf.SetFont("Arial", "B", 18)
	pdf.SetTextColor(15, 23, 42)

	pdf.CellFormat(100, 10, "KORP Notas Fiscais", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(5, 150, 105)
	pdf.CellFormat(80, 10, fmt.Sprintf("NOTA #%d", invoice.SeqCode), "", 1, "R", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(100, 116, 139)
	pdf.CellFormat(100, 5, "CNPJ: 00.000.000/0001-00", "", 0, "L", false, 0, "")
	pdf.CellFormat(80, 5, fmt.Sprintf("Status: %s", invoice.Status), "", 1, "R", false, 0, "")

	pdf.Ln(8)
	pdf.SetDrawColor(226, 232, 240)
	pdf.Line(15, pdf.GetY(), 195, pdf.GetY())
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(15, 23, 42)
	pdf.Cell(180, 6, "DADOS DO CLIENTE")
	pdf.Ln(6)

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(51, 65, 85)
	pdf.Cell(180, 5, fmt.Sprintf("Cliente: %s", invoice.Name))
	pdf.Ln(5)
	pdf.Cell(180, 5, tr(fmt.Sprintf("Endereço: %s", invoice.Address)))
	pdf.Ln(10)

	pdf.SetFillColor(15, 23, 42)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 9)

	pdf.CellFormat(90, 8, "  PRODUTO", "1", 0, "L", true, 0, "")
	pdf.CellFormat(25, 8, "QTD", "1", 0, "C", true, 0, "")
	pdf.CellFormat(32, 8, tr("PREÇO UNIT."), "1", 0, "R", true, 0, "")
	pdf.CellFormat(33, 8, "TOTAL", "1", 1, "R", true, 0, "")

	pdf.SetTextColor(30, 41, 59)
	pdf.SetFont("Arial", "", 9)
	fill := false

	for _, item := range invoice.Items {
		if fill {
			pdf.SetFillColor(248, 250, 252)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		unitPrice := float64(item.Price) / 100.0
		totalItemPrice := unitPrice * float64(item.Quantity)

		pdf.CellFormat(90, 7, fmt.Sprintf("  %s", item.Name), "1", 0, "L", fill, 0, "")
		pdf.CellFormat(25, 7, fmt.Sprintf("%d", item.Quantity), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(32, 7, fmt.Sprintf("R$ %.2f", unitPrice), "1", 0, "R", fill, 0, "")
		pdf.CellFormat(33, 7, fmt.Sprintf("R$ %.2f", totalItemPrice), "1", 1, "R", fill, 0, "")

		fill = !fill
	}

	pdf.Ln(4)
	pdf.SetFont("Arial", "B", 11)
	totalAmount := float64(invoice.TotalPrice) / 100.0
	pdf.CellFormat(147, 8, "TOTAL DA NOTA:", "", 0, "R", false, 0, "")
	pdf.CellFormat(33, 8, fmt.Sprintf("R$ %.2f", totalAmount), "", 1, "R", false, 0, "")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
