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

	"github.com/gondolia/gondolia/provider"
	"github.com/gondolia/gondolia/provider/pim"
	"github.com/gondolia/gondolia/provider/search"
	_ "github.com/gondolia/gondolia/provider/search/meilisearch" // Register meilisearch provider
	_ "github.com/gondolia/gondolia/provider/search/noop"        // Register noop provider
	"github.com/gondolia/gondolia/services/catalog/internal/config"
	"github.com/gondolia/gondolia/services/catalog/internal/domain"
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

	// Initialize PIM provider
	var pimProvider pim.PIMProvider
	// TODO: Initialize PIM provider based on config
	_ = pimProvider

	// Initialize search provider
	searchProvider, err := initSearchProvider(ctx, cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize search provider", zap.Error(err))
	}
	if searchProvider != nil {
		logger.Info("Search provider initialized", zap.String("provider", cfg.SearchProvider))
	}

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
	} else if searchProvider != nil {
		// Create a minimal sync service for search indexing only
		syncService = service.NewSyncService(productRepo, categoryRepo, nil, searchProvider)
		searchService = service.NewSearchService(searchProvider)
	}

	// Bulk index all products on startup if search provider is available
	if searchProvider != nil && cfg.SearchProvider == "meilisearch" {
		go func() {
			logger.Info("Starting bulk product indexing...")
			if err := bulkIndexProducts(ctx, productRepo, searchProvider, tenantRepo, logger); err != nil {
				logger.Error("Bulk indexing failed", zap.Error(err))
			} else {
				logger.Info("Bulk product indexing completed")
			}
		}()
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

// initSearchProvider initializes the search provider based on configuration
func initSearchProvider(ctx context.Context, cfg *config.Config, logger *zap.Logger) (search.SearchProvider, error) {
	providerType := cfg.SearchProvider
	if providerType == "" || providerType == "mock" {
		providerType = "noop"
	}

	var providerConfig map[string]any

	if providerType == "meilisearch" {
		host := cfg.SearchURL
		if host == "" {
			host = os.Getenv("MEILI_HOST")
		}
		if host == "" {
			logger.Warn("Meilisearch host not configured, using noop provider")
			providerType = "noop"
		}

		apiKey := cfg.SearchAPIKey
		if apiKey == "" {
			apiKey = os.Getenv("MEILI_API_KEY")
		}
		if apiKey == "" {
			logger.Warn("Meilisearch API key not configured, using noop provider")
			providerType = "noop"
		}

		if providerType == "meilisearch" {
			providerConfig = map[string]any{
				"host":    host,
				"api_key": apiKey,
				"timeout": 30,
			}
		}
	}

	// Get provider factory from registry
	factory, err := provider.Get[search.SearchProvider]("search", providerType)
	if err != nil {
		return nil, fmt.Errorf("search provider '%s' not found: %w", providerType, err)
	}

	searchProv, err := factory(providerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create search provider: %w", err)
	}

	// Test connection
	if err := searchProv.Health(ctx); err != nil {
		logger.Warn("Search provider health check failed", zap.Error(err))
		// Don't fail startup, just log the warning
	}

	// Configure products index if meilisearch
	if providerType == "meilisearch" {
		if err := configureProductsIndex(ctx, searchProv, logger); err != nil {
			logger.Error("Failed to configure products index", zap.Error(err))
			// Don't fail startup
		}
	}

	return searchProv, nil
}

// configureProductsIndex configures the Meilisearch products index
func configureProductsIndex(ctx context.Context, provider search.SearchProvider, logger *zap.Logger) error {
	indexName := "products"

	// Try to create the index (will fail if it already exists, which is fine)
	err := provider.CreateIndex(ctx, indexName, "id")
	if err != nil {
		logger.Debug("Index creation skipped (may already exist)", zap.String("index", indexName))
	}

	// Configure index settings
	indexConfig := search.IndexConfig{
		SearchableAttributes: []string{
			"name.de",
			"name.en",
			"description.de",
			"description.en",
			"sku",
		},
		FilterableAttributes: []string{
			"tenant_id",
			"status",
			"product_type",
			"category_ids",
		},
		SortableAttributes: []string{
			"name.de",
			"sku",
			"created_at",
			"updated_at",
		},
		// TypoTolerance: use Meilisearch defaults (enabled by default)
	}

	if err := provider.ConfigureIndex(ctx, indexName, indexConfig); err != nil {
		return fmt.Errorf("failed to configure index: %w", err)
	}

	logger.Info("Products index configured successfully", zap.String("index", indexName))
	return nil
}

// bulkIndexProducts indexes all products into Meilisearch
func bulkIndexProducts(ctx context.Context, productRepo *postgres.ProductRepository, searchProv search.SearchProvider, tenantRepo *postgres.TenantRepository, logger *zap.Logger) error {
	// Get all tenants by trying known codes
	tenant, err := tenantRepo.GetByCode(ctx, "demo")
	if err != nil {
		return fmt.Errorf("failed to get demo tenant: %w", err)
	}

	filter := domain.ProductFilter{
		TenantID: tenant.ID,
		Limit:    10000,
	}
	products, _, err := productRepo.List(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to list products: %w", err)
	}

	if len(products) == 0 {
		logger.Info("No products to index")
		return nil
	}

	docs := make([]search.Document, 0, len(products))
	for _, p := range products {
		doc := search.Document{
			"id":           p.ID.String(),
			"tenant_id":    p.TenantID.String(),
			"sku":          p.SKU,
			"name":         p.Name,
			"description":  p.Description,
			"product_type": string(p.ProductType),
			"status":       string(p.Status),
			"category_ids": p.CategoryIDs,
			"created_at":   p.CreatedAt.Unix(),
			"updated_at":   p.UpdatedAt.Unix(),
		}
		docs = append(docs, doc)
	}

	if _, err := searchProv.IndexDocuments(ctx, "products", docs); err != nil {
		return fmt.Errorf("failed to index products: %w", err)
	}

	logger.Info("Bulk indexing complete", zap.Int("total", len(products)))
	return nil
}

