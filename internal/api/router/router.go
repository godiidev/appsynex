// File: internal/api/router/router.go  
// Cập nhật tại: internal/api/router/router.go
// Mục đích: Thêm routes cho User và Category APIs

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
	userService := services.NewUserService(userRepo, roleRepo)
	categoryService := services.NewCategoryService(productCategoryRepo)
	sampleService := services.NewSampleService(sampleRepo, productNameRepo, productCategoryRepo)

	// Handlers
	authHandler := v1.NewAuthHandler(authService)
	userHandler := v1.NewUserHandler(userService)
	categoryHandler := v1.NewCategoryHandler(categoryService)
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
			users := protected.Group("/users")
			{
				users.GET("", middleware.HasPermission("USER_VIEW"), userHandler.GetAll)
				users.POST("", middleware.HasPermission("USER_CREATE"), userHandler.Create)
				users.GET("/:id", middleware.HasPermission("USER_VIEW"), userHandler.GetByID)
				users.PUT("/:id", middleware.HasPermission("USER_UPDATE"), userHandler.Update)
				users.DELETE("/:id", middleware.HasPermission("USER_DELETE"), userHandler.Delete)
				users.POST("/:id/roles", middleware.HasPermission("USER_UPDATE"), userHandler.AssignRoles)
			}

			// Category routes
			categories := protected.Group("/categories")
			{
				categories.GET("", middleware.HasPermission("PRODUCT_VIEW"), categoryHandler.GetAll)
				categories.POST("", middleware.HasPermission("PRODUCT_CREATE"), categoryHandler.Create)
				categories.GET("/:id", middleware.HasPermission("PRODUCT_VIEW"), categoryHandler.GetByID)
				categories.PUT("/:id", middleware.HasPermission("PRODUCT_UPDATE"), categoryHandler.Update)
				categories.DELETE("/:id", middleware.HasPermission("PRODUCT_DELETE"), categoryHandler.Delete)
			}

			// Sample routes
			samples := protected.Group("/samples")
			{
				samples.GET("", middleware.HasPermission("SAMPLE_VIEW"), sampleHandler.GetAll)
				samples.POST("", middleware.HasPermission("SAMPLE_CREATE"), sampleHandler.Create)
				samples.GET("/:id", middleware.HasPermission("SAMPLE_VIEW"), sampleHandler.GetByID)
				samples.PUT("/:id", middleware.HasPermission("SAMPLE_UPDATE"), sampleHandler.Update)
				samples.DELETE("/:id", middleware.HasPermission("SAMPLE_DELETE"), sampleHandler.Delete)
			}

			// Future routes: Orders, Customers, Warehouses, etc.
		}
	}

	return r
}