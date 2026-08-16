package apperrors

import "errors"

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	ErrNotFound           = errors.New("recurso não encontrado")
	ErrInvalidInput       = errors.New("dados inválidos")
	ErrServiceUnavailable = errors.New("serviço indisponível")
	ErrDuplicated         = errors.New("recurso já existe")
	ErrInvoiceNotOpened   = errors.New("fatura não aberta")
)
