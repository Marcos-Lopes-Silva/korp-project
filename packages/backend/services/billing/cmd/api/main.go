package main

import (
	"os"

	"billing/internal/config"
	"billing/internal/config/clients"
	"billing/internal/handlers"
	"billing/internal/models"
	"billing/internal/repository"
	"billing/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.Invoice{}, &models.InvoiceItem{})

	stockClient := clients.NewStockClient(os.Getenv("STOCK_SERVICE_URL"))

	repo := repository.NewInvoiceRepository(config.DB)
	service := services.NewInvoiceService(repo, stockClient)
	handler := handlers.NewInvoiceController(service)

	router := gin.Default()

	router.GET("/invoices", handler.GetInvoices)
	router.POST("/invoices", handler.CreateInvoice)
	router.GET("/invoices/:id", handler.GetInvoice)
	router.POST("/invoices/:id/items", handler.AddItem)
	router.POST("/invoices/:id/print", handler.PrintInvoice)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "UP",
		})
	})
	router.Run(":8082")
}
