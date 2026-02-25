package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/gondolia/gondolia/provider/pim"
	"github.com/gondolia/gondolia/provider/search"
	"github.com/gondolia/gondolia/services/catalog/internal/config"
	"github.com/gondolia/gondolia/services/catalog/internal/handler"
	"github.com/gondolia/gondolia/services/catalog/internal/middleware"
	"github.com/gondolia/gondolia/services/catalog/internal/repository/postgres"
	"github.com/gondolia/gondolia/services/catalog/internal/service"
)

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize database
	ctx := context.Background()
	db, err := postgres.NewDB(ctx, cfg.DatabaseURL())
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Initialize repositories
	tenantRepo := postgres.NewTenantRepository(db)
	productRepo := postgres.NewProductRepository(db)
	categoryRepo := postgres.NewCategoryRepository(db)
	priceRepo := postgres.NewPriceRepository(db)
	attrTransRepo := postgres.NewAttributeTranslationRepository(db)

	// Initialize providers (mock for now)
	var pimProvider pim.PIMProvider
	var searchProvider search.SearchProvider

	// TODO: Initialize actual providers based on config
	// For now, these would be nil and services would handle gracefully
	_ = pimProvider
	_ = searchProvider

	// Initialize parametric repositories
	parametricPricingRepo := postgres.NewParametricPricingRepository(db)
	axisOptionRepo := postgres.NewAxisOptionRepository(db)
	skuMappingRepo := postgres.NewSKUMappingRepository(db)

	// Initialize bundle repository
	bundleRepo := postgres.NewBundleRepository(db)

	// Initialize services
	productService := service.NewProductService(productRepo, priceRepo, attrTransRepo)
	variantService := service.NewVariantService(productRepo, priceRepo)
	categoryService := service.NewCategoryService(categoryRepo, productRepo)
	priceService := service.NewPriceService(priceRepo, productRepo)
	attrTransService := service.NewAttributeTranslationService(attrTransRepo)
	parametricService := service.NewParametricService(productRepo, parametricPricingRepo, axisOptionRepo, skuMappingRepo)
	bundleService := service.NewBundleService(bundleRepo, productRepo, priceRepo, parametricService)
	
	var syncService *service.SyncService
	var searchService *service.SearchService
	
	if pimProvider != nil && searchProvider != nil {
		syncService = service.NewSyncService(productRepo, categoryRepo, pimProvider, searchProvider)
		searchService = service.NewSearchService(searchProvider)
	}

	// Initialize handlers
	productHandler := handler.NewProductHandler(productService)
	variantHandler := handler.NewVariantHandler(variantService)
	variantHandler.SetParametricService(parametricService)
	categoryHandler := handler.NewCategoryHandler(categoryService, productService)
	priceHandler := handler.NewPriceHandler(priceService)
	parametricHandler := handler.NewParametricHandler(parametricService)
	bundleHandler := handler.NewBundleHandler(bundleService)
	attrTransHandler := handler.NewAttributeTranslationHandler(attrTransService)
	
	var searchHandler *handler.SearchHandler
	if searchService != nil && syncService != nil {
		searchHandler = handler.NewSearchHandler(searchService, syncService)
	}

	// Initialize HTTP server (REST API)
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// CORS middleware - must be before all routes
	router.Use(middleware.CORSMiddleware(middleware.CORSConfig{
		AllowedOrigins: cfg.AllowedOrigins,
	}))

	// Health endpoints
	router.GET("/health/live", handler.LivenessHandler)
	router.GET("/health/ready", handler.ReadinessHandler)
	router.GET("/metrics", handler.MetricsHandler)

	// API routes
	api := router.Group("/api/v1")

	// Apply tenant middleware to all API routes
	api.Use(middleware.TenantMiddleware(tenantRepo))

	// Product endpoints
	products := api.Group("/products")
	{
		products.GET("", productHandler.List)
		products.POST("", productHandler.Create)
		products.GET("/:id", variantHandler.GetProductWithVariants) // Enhanced to include variants
		products.PUT("/:id", productHandler.Update)
		products.DELETE("/:id", productHandler.Delete)
		products.PATCH("/:id/status", productHandler.UpdateStatus) // PIM: Status management

		// Variant endpoints
		products.GET("/:id/variants", variantHandler.ListVariants)
		products.POST("/:id/variants", variantHandler.CreateVariant)
		products.GET("/:id/variants/select", variantHandler.SelectVariant)
		products.GET("/:id/variants/available", variantHandler.GetAvailableAxisValues)

		// Price endpoints for product
		products.GET("/:id/prices", priceHandler.ListByProduct)
		products.POST("/:id/prices", priceHandler.Create)

		// Attribute endpoints (PIM)
		products.POST("/:id/attributes", productHandler.AddAttribute)
		products.PUT("/:id/attributes/:key", productHandler.UpdateAttribute)
		products.DELETE("/:id/attributes/:key", productHandler.DeleteAttribute)

		// Parametric endpoints
		products.POST("/:id/calculate-price", parametricHandler.CalculatePrice)

		// Bundle endpoints
		products.GET("/:id/bundle-components", bundleHandler.GetComponents)
		products.PUT("/:id/bundle-components", bundleHandler.SetComponents)
		products.POST("/:id/bundle-components", bundleHandler.AddComponent)         // PIM: Add component
		products.PUT("/:id/bundle-components/:compId", bundleHandler.UpdateComponent) // PIM: Update component
		products.DELETE("/:id/bundle-components/:compId", bundleHandler.DeleteComponent) // PIM: Delete component
	}

	// Bundle endpoints (storefront)
	bundles := api.Group("/bundles")
	{
		bundles.POST("/:id/calculate-price", bundleHandler.CalculatePrice)
	}

	// Category endpoints
	categories := api.Group("/categories")
	{
		categories.GET("", categoryHandler.GetTree) // Returns tree by default
		categories.GET("/list", categoryHandler.List) // Paginated list
		categories.POST("", categoryHandler.Create)
		categories.GET("/:id", categoryHandler.Get)
		categories.PUT("/:id", categoryHandler.Update)
		categories.DELETE("/:id", categoryHandler.Delete)
		categories.PATCH("/:id/sort", categoryHandler.UpdateSortOrder) // PIM: Sort order management
		categories.GET("/:id/products", categoryHandler.GetProducts) // Products by category
		categories.POST("/:id/products", categoryHandler.AddProduct)  // PIM: Assign product to category
		categories.DELETE("/:id/products/:productId", categoryHandler.RemoveProduct) // PIM: Remove product from category
	}

	// Price endpoints
	prices := api.Group("/prices")
	{
		prices.PUT("/:id", priceHandler.Update)
		prices.DELETE("/:id", priceHandler.Delete)
	}

	// Attribute translation endpoints
	attrTrans := api.Group("/attribute-translations")
	{
		attrTrans.GET("", attrTransHandler.List)
		attrTrans.GET("/by-locale/:locale", attrTransHandler.GetByLocale)
		attrTrans.POST("", attrTransHandler.Create)
		attrTrans.PUT("/:id", attrTransHandler.Update)
		attrTrans.DELETE("/:id", attrTransHandler.Delete)
	}

	// Search endpoints (if available)
	if searchHandler != nil {
		api.GET("/search", searchHandler.Search)
		api.POST("/sync/pim", searchHandler.SyncPIM)
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.HTTPPort),
		Handler: router,
	}

	// Initialize gRPC server
	grpcServer := grpc.NewServer()
	// Register gRPC services here

	// Start servers
	go func() {
		logger.Info("Starting HTTP server", zap.String("port", cfg.HTTPPort))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("HTTP server error", zap.Error(err))
		}
	}()

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
		if err != nil {
			logger.Fatal("Failed to listen for gRPC", zap.Error(err))
		}
		logger.Info("Starting gRPC server", zap.String("port", cfg.GRPCPort))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("gRPC server error", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down servers...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	grpcServer.GracefulStop()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server shutdown error", zap.Error(err))
	}

	logger.Info("Servers stopped")
}
