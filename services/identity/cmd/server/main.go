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

	"github.com/gondolia/gondolia/services/identity/internal/auth"
	"github.com/gondolia/gondolia/services/identity/internal/config"
	"github.com/gondolia/gondolia/services/identity/internal/domain"
	"github.com/gondolia/gondolia/services/identity/internal/handler"
	"github.com/gondolia/gondolia/services/identity/internal/middleware"
	"github.com/gondolia/gondolia/services/identity/internal/repository/postgres"
	"github.com/gondolia/gondolia/services/identity/internal/service"
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
	userRepo := postgres.NewUserRepository(db)
	companyRepo := postgres.NewCompanyRepository(db)
	roleRepo := postgres.NewRoleRepository(db)
	userCompanyRepo := postgres.NewUserCompanyRepository(db)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(db)
	passwordResetRepo := postgres.NewPasswordResetRepository(db)
	authLogRepo := postgres.NewAuthLogRepository(db)

	// Initialize JWT manager
	tokenConfig := auth.TokenConfig{
		AccessSecret:       cfg.JWTAccessSecret,
		RefreshSecret:      cfg.JWTRefreshSecret,
		AccessTokenExpiry:  cfg.JWTAccessTokenExpiry,
		RefreshTokenExpiry: cfg.JWTRefreshTokenExpiry,
		Issuer:             cfg.ServiceName,
	}
	jwtManager := auth.NewJWTManager(tokenConfig)

	// Initialize services
	authService := service.NewAuthService(
		userRepo,
		companyRepo,
		roleRepo,
		userCompanyRepo,
		refreshTokenRepo,
		passwordResetRepo,
		authLogRepo,
		jwtManager,
	)
	userService := service.NewUserService(userRepo, userCompanyRepo, roleRepo, authLogRepo)
	companyService := service.NewCompanyService(companyRepo, userCompanyRepo, roleRepo, userRepo)
	roleService := service.NewRoleService(roleRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, userService, cfg)
	userHandler := handler.NewUserHandler(userService)
	companyHandler := handler.NewCompanyHandler(companyService)
	roleHandler := handler.NewRoleHandler(roleService)

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

	// Public auth endpoints (no JWT required)
	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/refresh", authHandler.Refresh)
		authRoutes.POST("/forgot-password", authHandler.ForgotPassword)
		authRoutes.POST("/reset-password", authHandler.ResetPassword)
		authRoutes.GET("/invitations/:token", authHandler.ValidateInvitation)
		authRoutes.POST("/invitations/:token/accept", authHandler.AcceptInvitation)
	}

	// Protected routes (JWT required)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))

	// Protected auth endpoints
	protectedAuth := protected.Group("/auth")
	{
		protectedAuth.POST("/logout", authHandler.Logout)
		protectedAuth.GET("/me", authHandler.Me)
		protectedAuth.POST("/switch-company", authHandler.SwitchCompany)
	}

	// User endpoints
	users := protected.Group("/users")
	{
		users.GET("", userHandler.List)
		users.POST("", middleware.RequirePermission(domain.PermManageUsersAndRoles), userHandler.Create)
		users.GET("/:id", userHandler.Get)
		users.PUT("/:id", middleware.RequirePermission(domain.PermManageUsersAndRoles), userHandler.Update)
		users.DELETE("/:id", middleware.RequirePermission(domain.PermManageUsersAndRoles), userHandler.Delete)
		users.POST("/:id/activate", middleware.RequirePermission(domain.PermManageUsersAndRoles), userHandler.Activate)
		users.POST("/:id/deactivate", middleware.RequirePermission(domain.PermManageUsersAndRoles), userHandler.Deactivate)
		users.POST("/invite", middleware.RequirePermission(domain.PermManageUsersAndRoles), userHandler.Invite)
	}

	// Company endpoints
	companies := protected.Group("/companies")
	{
		companies.GET("", companyHandler.List)
		companies.POST("", middleware.RequirePermission(domain.PermManageCompany), companyHandler.Create)
		companies.GET("/:id", companyHandler.Get)
		companies.PUT("/:id", middleware.RequirePermission(domain.PermManageCompany), companyHandler.Update)
		companies.DELETE("/:id", middleware.RequirePermission(domain.PermManageCompany), companyHandler.Delete)
		companies.GET("/:id/users", companyHandler.ListUsers)
		companies.POST("/:id/users", middleware.RequirePermission(domain.PermManageUsersAndRoles), companyHandler.AddUser)
		companies.PUT("/:id/users/:userId", middleware.RequirePermission(domain.PermManageUsersAndRoles), companyHandler.UpdateUserRole)
		companies.DELETE("/:id/users/:userId", middleware.RequirePermission(domain.PermManageUsersAndRoles), companyHandler.RemoveUser)
	}

	// Role endpoints
	roles := protected.Group("/roles")
	{
		roles.GET("", roleHandler.List)
		roles.POST("", middleware.RequirePermission(domain.PermManageUsersAndRoles), roleHandler.Create)
		roles.GET("/:id", roleHandler.Get)
		roles.PUT("/:id", middleware.RequirePermission(domain.PermManageUsersAndRoles), roleHandler.Update)
		roles.DELETE("/:id", middleware.RequirePermission(domain.PermManageUsersAndRoles), roleHandler.Delete)
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
