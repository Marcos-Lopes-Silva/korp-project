package main

import (
	"stock/internal/config"
	_ "stock/internal/docs"
	"stock/internal/handlers"
	"stock/internal/models"
	"stock/internal/repository"
	"stock/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.Product{})
	repo := repository.NewProductRepository(config.DB)
	service := services.NewProductService(repo)
	handler := handlers.NewProductController(service)

	router := gin.Default()

	router.GET("/products", handler.GetAllProducts)
	router.POST("/products", handler.CreateProduct)
	router.GET("/products/:id", handler.GetProductByID)
	router.PUT("/products/:id", handler.UpdateProduct)
	router.DELETE("/products/:id", handler.DeleteProduct)
	router.POST("/products/:id/reduce-stock", handler.ReduceStock)
	router.POST("/products/:id/restore-stock", handler.RestoreStock)
	router.GET("/products/:id/availability", handler.VerifyAvailability)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "UP",
		})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":8080")
}
