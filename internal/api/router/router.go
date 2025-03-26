package router

import (
	"github.com/gin-gonic/gin"
	"github.com/godiidev/appsynex/config"
	v1 "github.com/godiidev/appsynex/internal/api/handlers/v1"
	"github.com/godiidev/appsynex/internal/api/middleware"
	"github.com/godiidev/appsynex/internal/domain/services"
	"github.com/godiidev/appsynex/internal/repository/mysql"
	"github.com/godiidev/appsynex/pkg/auth"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	// Set Gin mode
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Middlewares
	r.Use(middleware.CORS())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Repositories
	userRepo := mysql.NewUserRepository(db)
	roleRepo := mysql.NewRoleRepository(db)
	productNameRepo := mysql.NewProductNameRepository(db)
	productCategoryRepo := mysql.NewProductCategoryRepository(db)
	sampleRepo := mysql.NewSampleRepository(db)

	// Services
	jwtService := auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpiresIn)
	authService := services.NewAuthService(userRepo, roleRepo, jwtService)
	sampleService := services.NewSampleService(sampleRepo, productNameRepo, productCategoryRepo)

	// Handlers
	authHandler := v1.NewAuthHandler(authService)
	sampleHandler := v1.NewSampleHandler(sampleService)

	// API v1 routes
	api := r.Group("/api/v1")
	{
		// Auth routes (no auth required)
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.Auth(jwtService))
		{
			// User routes
			// Product routes
			// Sample routes
			samples := protected.Group("/samples")
			{
				samples.GET("", middleware.HasPermission("SAMPLE_VIEW"), sampleHandler.GetAll)
				samples.POST("", middleware.HasPermission("SAMPLE_CREATE"), sampleHandler.Create)
				samples.GET("/:id", middleware.HasPermission("SAMPLE_VIEW"), sampleHandler.GetByID)
				samples.PUT("/:id", middleware.HasPermission("SAMPLE_UPDATE"), sampleHandler.Update)
				samples.DELETE("/:id", middleware.HasPermission("SAMPLE_DELETE"), sampleHandler.Delete)
			}
			// Customer routes
			// Order routes
		}
	}

	return r
}
