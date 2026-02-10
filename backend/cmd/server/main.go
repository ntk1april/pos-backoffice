package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"pos-backoffice/internal/config"
	"pos-backoffice/internal/database"
	"pos-backoffice/internal/handler"
	"pos-backoffice/internal/middleware"
	"pos-backoffice/internal/repository"
	"pos-backoffice/internal/service"
)

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Initialize repositories
	db := database.GetDB()
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	stockRepo := repository.NewStockRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo)
	productService := service.NewProductService(productRepo)
	stockService := service.NewStockService(db, productRepo, stockRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	productHandler := handler.NewProductHandler(productService)
	stockHandler := handler.NewStockHandler(stockService)

	// Setup Gin router
	router := setupRouter(authHandler, productHandler, stockHandler)

	// Create HTTP server
	srv := &http.Server{
		Addr:           ":" + config.AppConfig.ServerPort,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🚀 Server starting on port %s", config.AppConfig.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func setupRouter(authHandler *handler.AuthHandler, productHandler *handler.ProductHandler, stockHandler *handler.StockHandler) *gin.Engine {
	// Set Gin mode
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Apply CORS middleware
	router.Use(middleware.CORSMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Public routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Product routes
			products := protected.Group("/products")
			{
				products.GET("", productHandler.GetProducts)
				products.GET("/:id", productHandler.GetProduct)

				// Admin only routes
				adminProducts := products.Group("")
				adminProducts.Use(middleware.RequireRole("ADMIN"))
				{
					adminProducts.POST("", productHandler.CreateProduct)
					adminProducts.PUT("/:id", productHandler.UpdateProduct)
					adminProducts.DELETE("/:id", productHandler.DeleteProduct)
				}
			}

			// Stock routes
			stock := protected.Group("/stock")
			{
				stock.POST("/increase", stockHandler.IncreaseStock)
				stock.POST("/decrease", stockHandler.DecreaseStock)
				stock.GET("/logs/:product_id", stockHandler.GetStockLogs)
			}
		}
	}

	return router
}
