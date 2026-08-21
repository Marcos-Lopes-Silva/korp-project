package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func newProxy(targetURL string) *httputil.ReverseProxy {
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatalf("URL inválida: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = &http.Transport{
		ResponseHeaderTimeout: 5 * time.Second,
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("erro ao rotear pra %s: %v", targetURL, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"error":"serviço temporariamente indisponível"}`))
	}

	return proxy
}

func main() {
	stockURL := os.Getenv("STOCK_SERVICE_URL")
	billingURL := os.Getenv("BILLING_SERVICE_URL")

	stockProxy := newProxy(stockURL)
	billingProxy := newProxy(billingURL)

	r := gin.Default()
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false

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

	stockHealthURL, _ := url.Parse(os.Getenv("STOCK_SERVICE_URL"))
	billingHealthURL, _ := url.Parse(os.Getenv("BILLING_SERVICE_URL"))

	r.GET("/stock/health", func(c *gin.Context) {
		proxy := httputil.NewSingleHostReverseProxy(stockHealthURL)
		c.Request.URL.Path = "/health"
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/billing/health", func(c *gin.Context) {
		proxy := httputil.NewSingleHostReverseProxy(billingHealthURL)
		c.Request.URL.Path = "/health"
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	r.Any("/products", func(c *gin.Context) {
		stockProxy.ServeHTTP(c.Writer, c.Request)
	})
	r.Any("/products/*proxyPath", func(c *gin.Context) {
		stockProxy.ServeHTTP(c.Writer, c.Request)
	})

	r.Any("/invoices", func(c *gin.Context) {
		billingProxy.ServeHTTP(c.Writer, c.Request)
	})
	r.Any("/invoices/*proxyPath", func(c *gin.Context) {
		billingProxy.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.Run(":8080")
}
