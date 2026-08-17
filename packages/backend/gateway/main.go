package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

func newProxy(targetURL string) *httputil.ReverseProxy {
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatalf("URL inválida: %v", err)
	}
	return httputil.NewSingleHostReverseProxy(target)
}

func main() {
	stockURL := os.Getenv("STOCK_SERVICE_URL")
	billingURL := os.Getenv("BILLING_SERVICE_URL")

	stockProxy := newProxy(stockURL)
	billingProxy := newProxy(billingURL)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	r.Any("/products/*proxyPath", func(c *gin.Context) {
		stockProxy.ServeHTTP(c.Writer, c.Request)
	})

	r.Any("/invoices/*proxyPath", func(c *gin.Context) {
		billingProxy.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.Run(":8000")
}
