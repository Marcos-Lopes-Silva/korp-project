package handlers

import (
	"errors"
	"net/http"
	"stock/internal/apperrors"
	"stock/internal/models"
	"stock/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductController struct {
	service *services.ProductService
}

type MessageResponse struct {
	Message string `json:"message"`
}

func NewProductController(service *services.ProductService) *ProductController {
	return &ProductController{service: service}
}

// GetAllProducts godoc
// @Summary Get all products
// @Description Get a list of all products
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {array} models.Product
// @Failure 500 {object} apperrors.ErrorResponse
// @Router /products [get]
func (pc *ProductController) GetAllProducts(c *gin.Context) {
	products, err := pc.service.GetProducts()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar os produtos",
		})
		return
	}

	c.JSON(http.StatusOK, products)
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the provided details
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product to create"
// @Success 200 {object} models.Product
// @Failure 400 {object} apperrors.ErrorResponse
// @Failure 409 {object} apperrors.ErrorResponse
// @Failure 500 {object} apperrors.ErrorResponse
// @Router /products [post]
func (pc *ProductController) CreateProduct(c *gin.Context) {
	var newProduct models.Product

	if err := c.ShouldBindJSON(&newProduct); err != nil {
		apperrors.HandleError(c, apperrors.ErrInvalidInput)
		return
	}

	product, err := pc.service.CreateProduct(&newProduct)

	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetProductByID godoc
// @Summary Get a product by ID
// @Description Get a product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product
// @Failure 400 {object} apperrors.ErrorResponse
// @Failure 404 {object} apperrors.ErrorResponse
// @Router /products/{id} [get]
func (pc *ProductController) GetProductByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		apperrors.HandleError(c, apperrors.ErrInvalidInput)
		return
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		apperrors.HandleError(c, apperrors.ErrInvalidInput)
		return
	}

	product, err := pc.service.GetProductByID(uuid)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct godoc
// @Summary Update a product by ID
// @Description Update a product by its ID with the provided details
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body models.Product true "Product to update"
// @Success 200 {object} models.Product
// @Failure 400 {object} apperrors.ErrorResponse
// @Failure 404 {object} apperrors.ErrorResponse
// @Failure 409 {object} apperrors.ErrorResponse
// @Router /products/{id} [put]
func (pc *ProductController) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		apperrors.HandleError(c, apperrors.ErrInvalidInput)
		return
	}

	var updatedProduct models.Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		apperrors.HandleError(c, apperrors.ErrInvalidInput)
		return
	}

	idUUID, err := uuid.Parse(id)

	if err != nil {
		apperrors.HandleError(c, apperrors.ErrInvalidInput)
		return
	}

	product, err := pc.service.UpdateProduct(idUUID, &updatedProduct)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			apperrors.HandleError(c, apperrors.ErrNotFound)
			return
		}
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary Delete a product by ID
// @Description Delete a product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 204 {object} MessageResponse
// @Failure 400 {object} apperrors.ErrorResponse
// @Failure 404 {object} apperrors.ErrorResponse
// @Router /products/{id} [delete]
func (pc *ProductController) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		apperrors.HandleError(c, apperrors.ErrInvalidInput)
		return
	}

	idUUID, err := uuid.Parse(id)

	if err != nil {
		apperrors.HandleError(c, apperrors.ErrInvalidInput)
		return
	}

	err = pc.service.DeleteProduct(idUUID)

	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Produto deletado com sucesso"})
}

// ReduceStock godoc
// @Summary Reduce stock of a product by ID
// @Description Reduce stock of a product by its ID with the provided quantity
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param quantity query int true "Quantity to reduce"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} apperrors.ErrorResponse
// @Failure 404 {object} apperrors.ErrorResponse
// @Failure 409 {object} apperrors.ErrorResponse
// @Router /products/{id}/reduce-stock [post]
func (pc *ProductController) ReduceStock(c *gin.Context) {
	id := c.Param("id")

	productID, err := uuid.Parse(id)
	if err != nil {
		apperrors.HandleError(c, apperrors.ErrInvalidInput)
		return
	}

	var body struct {
		Quantity int64 `json:"quantity" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		apperrors.HandleError(c, apperrors.ErrInvalidInput)
		return
	}

	if err := pc.service.ReduceStock(productID, body.Quantity); err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Estoque reduzido com sucesso"})
}

// RestoreStock godoc
// @Summary Restore stock of a product by ID
// @Description Restore stock of a product by its ID with the provided quantity
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param quantity query int true "Quantity to restore"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} apperrors.ErrorResponse
// @Failure 404 {object} apperrors.ErrorResponse
// @Router /products/{id}/restore-stock [post]
func (pc *ProductController) RestoreStock(c *gin.Context) {
	id := c.Param("id")

	productID, err := uuid.Parse(id)
	if err != nil {
		apperrors.HandleError(c, apperrors.ErrInvalidInput)
		return
	}

	var body struct {
		Quantity int64 `json:"quantity" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		apperrors.HandleError(c, apperrors.ErrInvalidInput)
		return
	}

	if err := pc.service.RestoreStock(productID, body.Quantity); err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Estoque restaurado com sucesso"})
}

func (pc *ProductController) VerifyAvailability(c *gin.Context) {
	id := c.Param("id")
	quantityStr := c.Query("quantity")

	if id == "" {
		apperrors.HandleError(c, apperrors.ErrInvalidInput)
		return
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		apperrors.HandleError(c, apperrors.ErrInvalidInput)
		return
	}

	quantity, err := strconv.ParseInt(quantityStr, 10, 10)
	if err != nil || quantity <= 0 {
		apperrors.HandleError(c, apperrors.ErrInvalidInput)
		return
	}

	isAvailable, err := pc.service.VerifyAvailability(uuid, quantity)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"available": isAvailable})
}
