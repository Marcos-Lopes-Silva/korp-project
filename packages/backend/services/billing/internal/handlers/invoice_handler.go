package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"billing/internal/apperrors"
	"billing/internal/services"
)

type InvoiceController struct {
	service *services.InvoiceService
}

func NewInvoiceController(service *services.InvoiceService) *InvoiceController {
	return &InvoiceController{service: service}
}

func (h *InvoiceController) CreateInvoice(c *gin.Context) {
	invoice, err := h.service.CreateInvoice()
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, invoice)
}

func (h *InvoiceController) AddItem(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.HandleError(c, apperrors.ErrInvalidInput)
		return
	}

	var body struct {
		ProductID uuid.UUID `json:"product_id" binding:"required"`
		Quantity  int       `json:"quantity" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		apperrors.HandleError(c, apperrors.ErrInvalidInput)
		return
	}

	invoice, err := h.service.AddItem(c.Request.Context(), id, body.ProductID, int64(body.Quantity))
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, invoice)
}

func (h *InvoiceController) GetInvoice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.HandleError(c, apperrors.ErrInvalidInput)
		return
	}

	invoice, err := h.service.GetInvoice(id)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, invoice)
}

func (h *InvoiceController) PrintInvoice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.HandleError(c, apperrors.ErrInvalidInput)
		return
	}

	invoice, err := h.service.PrintInvoice(c.Request.Context(), id)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, invoice)
}
