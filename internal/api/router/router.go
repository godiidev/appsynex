// File: internal/api/router/router.go
// Cập nhật tại: internal/api/router/router.go
// Mục đích: Fixed router với enhanced permission system (no route conflicts)

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

	// Global middlewares
	r.Use(middleware.CORS())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"version": "1.0.0",
		})
	})

	// Initialize repositories
	userRepo := mysql.NewUserRepository(db)
	roleRepo := mysql.NewRoleRepository(db)
	permissionRepo := mysql.NewPermissionRepository(db)
	productNameRepo := mysql.NewProductNameRepository(db)
	productCategoryRepo := mysql.NewProductCategoryRepository(db)
	sampleRepo := mysql.NewSampleRepository(db)

	// Initialize services
	jwtService := auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpiresIn)
	authService := services.NewAuthService(userRepo, roleRepo, permissionRepo, jwtService)
	userService := services.NewUserService(userRepo, roleRepo)
	permissionService := services.NewPermissionService(permissionRepo, roleRepo, userRepo)
	categoryService := services.NewCategoryService(productCategoryRepo)
	sampleService := services.NewSampleService(sampleRepo, productNameRepo, productCategoryRepo)

	// Initialize handlers
	authHandler := v1.NewAuthHandler(authService)
	userHandler := v1.NewUserHandler(userService)
	permissionHandler := v1.NewPermissionHandler(permissionService)
	categoryHandler := v1.NewCategoryHandler(categoryService)
	sampleHandler := v1.NewSampleHandler(sampleService)

	// Initialize permission middleware
	permMiddleware := middleware.NewPermissionMiddleware(permissionRepo)

	// API v1 routes
	api := r.Group("/api/v1")
	{
		// Public routes (no authentication required)
		public := api.Group("")
		{
			public.POST("/auth/login", authHandler.Login)
			// Add other public endpoints here (e.g., password reset request)
		}

		// Protected routes (authentication required)
		protected := api.Group("")
		protected.Use(middleware.Auth(jwtService))
		protected.Use(middleware.InjectPermissionMiddleware(permissionRepo))
		{
			// Permission Management Routes
			permissions := protected.Group("/permissions")
			{
				permissions.GET("", permMiddleware.RequirePermission("SYSTEM", "VIEW"), permissionHandler.GetAllPermissions)
				permissions.GET("/groups", permMiddleware.RequirePermission("SYSTEM", "VIEW"), permissionHandler.GetPermissionGroups)
				permissions.GET("/module/:module", permMiddleware.RequirePermission("SYSTEM", "VIEW"), permissionHandler.GetPermissionsByModule)
				permissions.POST("", permMiddleware.RequirePermission("SYSTEM", "CREATE"), permissionHandler.CreatePermission)
				permissions.PUT("/:id", permMiddleware.RequirePermission("SYSTEM", "UPDATE"), permissionHandler.UpdatePermission)
				permissions.DELETE("/:id", permMiddleware.RequirePermission("SYSTEM", "DELETE"), permissionHandler.DeletePermission)
				permissions.POST("/check", permMiddleware.RequirePermission("SYSTEM", "VIEW"), permissionHandler.CheckUserPermission)
				permissions.POST("/bulk-assign", permMiddleware.RequirePermission("ROLE", "ASSIGN_PERMISSIONS"), permissionHandler.BulkAssignPermissions)
				permissions.POST("/clone-role", permMiddleware.RequirePermission("ROLE", "ASSIGN_PERMISSIONS"), permissionHandler.CloneRolePermissions)
			}

			// User Management Routes
			users := protected.Group("/users")
			{
				users.GET("", permMiddleware.RequirePermission("USER", "VIEW"), userHandler.GetAll)
				users.POST("", permMiddleware.RequirePermission("USER", "CREATE"), userHandler.Create)
				users.GET("/:id", permMiddleware.RequireAnyPermission(
					middleware.PermissionCheck{Module: "USER", Action: "VIEW"},
					middleware.PermissionCheck{Module: "USER", Action: "VIEW_OWN"},
				), userHandler.GetByID)
				users.PUT("/:id", permMiddleware.RequirePermission("USER", "UPDATE"), userHandler.Update)
				users.DELETE("/:id", permMiddleware.RequirePermission("USER", "DELETE"), userHandler.Delete)

				// User role assignment
				users.POST("/:id/roles", permMiddleware.RequirePermission("USER", "ASSIGN_ROLES"), userHandler.AssignRoles)

				// User permission management - FIXED: use same :id parameter
				users.GET("/:id/permissions", permMiddleware.RequirePermission("USER", "VIEW"), permissionHandler.GetUserPermissions)
				users.GET("/:id/effective-permissions", permMiddleware.RequirePermission("USER", "VIEW"), permissionHandler.GetUserEffectivePermissions)
				users.POST("/:id/permissions", permMiddleware.RequirePermission("USER", "ASSIGN_PERMISSIONS"), permissionHandler.GrantUserPermission)
				users.DELETE("/:id/permissions/:permissionId", permMiddleware.RequirePermission("USER", "ASSIGN_PERMISSIONS"), permissionHandler.RevokeUserPermission)
			}

			// Role Management Routes - FIXED: Sử dụng consistent parameter name
			roles := protected.Group("/roles")
			{
				roles.GET("", permMiddleware.RequirePermission("ROLE", "VIEW"), func(c *gin.Context) {
					// TODO: Implement role handler
					c.JSON(200, gin.H{"message": "Role list endpoint"})
				})
				roles.POST("", permMiddleware.RequirePermission("ROLE", "CREATE"), func(c *gin.Context) {
					// TODO: Implement role creation
					c.JSON(200, gin.H{"message": "Role creation endpoint"})
				})
				roles.GET("/:id", permMiddleware.RequirePermission("ROLE", "VIEW"), func(c *gin.Context) {
					// TODO: Implement role details
					c.JSON(200, gin.H{"message": "Role details endpoint"})
				})
				roles.PUT("/:id", permMiddleware.RequirePermission("ROLE", "UPDATE"), func(c *gin.Context) {
					// TODO: Implement role update
					c.JSON(200, gin.H{"message": "Role update endpoint"})
				})
				roles.DELETE("/:id", permMiddleware.RequirePermission("ROLE", "DELETE"), func(c *gin.Context) {
					// TODO: Implement role deletion
					c.JSON(200, gin.H{"message": "Role deletion endpoint"})
				})

				// FIXED: Role permission management - sử dụng cùng parameter name ":id"
				roles.GET("/:id/permissions", permMiddleware.RequirePermission("ROLE", "VIEW"), permissionHandler.GetRolePermissions)
				roles.POST("/:id/permissions", permMiddleware.RequirePermission("ROLE", "ASSIGN_PERMISSIONS"), permissionHandler.AssignPermissionsToRole)
				roles.DELETE("/:id/permissions", permMiddleware.RequirePermission("ROLE", "ASSIGN_PERMISSIONS"), permissionHandler.RemovePermissionsFromRole)
			}

			// Product Category Management Routes
			categories := protected.Group("/categories")
			{
				categories.GET("", permMiddleware.RequirePermission("PRODUCT_CATEGORY", "VIEW"), categoryHandler.GetAll)
				categories.POST("", permMiddleware.RequirePermission("PRODUCT_CATEGORY", "CREATE"), categoryHandler.Create)
				categories.GET("/:id", permMiddleware.RequirePermission("PRODUCT_CATEGORY", "VIEW"), categoryHandler.GetByID)
				categories.PUT("/:id", permMiddleware.RequirePermission("PRODUCT_CATEGORY", "UPDATE"), categoryHandler.Update)
				categories.DELETE("/:id", permMiddleware.RequirePermission("PRODUCT_CATEGORY", "DELETE"), categoryHandler.Delete)
			}

			// Product Management Routes
			products := protected.Group("/products")
			{
				products.GET("", permMiddleware.RequirePermission("PRODUCT", "VIEW"), func(c *gin.Context) {
					// TODO: Implement product handler
					c.JSON(200, gin.H{"message": "Product list endpoint"})
				})
				products.POST("", permMiddleware.RequirePermission("PRODUCT", "CREATE"), func(c *gin.Context) {
					// TODO: Implement product creation
					c.JSON(200, gin.H{"message": "Product creation endpoint"})
				})
				products.GET("/:id", permMiddleware.RequirePermission("PRODUCT", "VIEW"), func(c *gin.Context) {
					// TODO: Implement product details
					c.JSON(200, gin.H{"message": "Product details endpoint"})
				})
				products.PUT("/:id", permMiddleware.RequirePermission("PRODUCT", "UPDATE"), func(c *gin.Context) {
					// TODO: Implement product update
					c.JSON(200, gin.H{"message": "Product update endpoint"})
				})
				products.DELETE("/:id", permMiddleware.RequirePermission("PRODUCT", "DELETE"), func(c *gin.Context) {
					// TODO: Implement product deletion
					c.JSON(200, gin.H{"message": "Product deletion endpoint"})
				})
			}

			// Sample Product Management Routes
			samples := protected.Group("/samples")
			{
				samples.GET("", permMiddleware.RequirePermission("SAMPLE", "VIEW"), sampleHandler.GetAll)
				samples.POST("", permMiddleware.RequirePermission("SAMPLE", "CREATE"), sampleHandler.Create)
				samples.GET("/:id", permMiddleware.RequirePermission("SAMPLE", "VIEW"), sampleHandler.GetByID)
				samples.PUT("/:id", permMiddleware.RequirePermission("SAMPLE", "UPDATE"), sampleHandler.Update)
				samples.DELETE("/:id", permMiddleware.RequirePermission("SAMPLE", "DELETE"), sampleHandler.Delete)

				// Additional sample operations
				samples.POST("/:id/dispatch", permMiddleware.RequirePermission("SAMPLE", "DISPATCH"), func(c *gin.Context) {
					// TODO: Implement sample dispatch
					c.JSON(200, gin.H{"message": "Sample dispatch endpoint"})
				})
				samples.GET("/:id/tracking", permMiddleware.RequirePermission("SAMPLE", "TRACK"), func(c *gin.Context) {
					// TODO: Implement sample tracking
					c.JSON(200, gin.H{"message": "Sample tracking endpoint"})
				})
			}

			// Customer Management Routes
			customers := protected.Group("/customers")
			{
				customers.GET("", permMiddleware.RequirePermission("CUSTOMER", "VIEW"), func(c *gin.Context) {
					// TODO: Implement customer handler
					c.JSON(200, gin.H{"message": "Customer list endpoint"})
				})
				customers.POST("", permMiddleware.RequirePermission("CUSTOMER", "CREATE"), func(c *gin.Context) {
					// TODO: Implement customer creation
					c.JSON(200, gin.H{"message": "Customer creation endpoint"})
				})
				customers.GET("/:id", permMiddleware.RequirePermission("CUSTOMER", "VIEW"), func(c *gin.Context) {
					// TODO: Implement customer details
					c.JSON(200, gin.H{"message": "Customer details endpoint"})
				})
				customers.PUT("/:id", permMiddleware.RequirePermission("CUSTOMER", "UPDATE"), func(c *gin.Context) {
					// TODO: Implement customer update
					c.JSON(200, gin.H{"message": "Customer update endpoint"})
				})
				customers.DELETE("/:id", permMiddleware.RequirePermission("CUSTOMER", "DELETE"), func(c *gin.Context) {
					// TODO: Implement customer deletion
					c.JSON(200, gin.H{"message": "Customer deletion endpoint"})
				})
				customers.GET("/:id/activity", permMiddleware.RequirePermission("CUSTOMER", "VIEW_ACTIVITY"), func(c *gin.Context) {
					// TODO: Implement customer activity logs
					c.JSON(200, gin.H{"message": "Customer activity endpoint"})
				})
			}

			// Order Management Routes
			orders := protected.Group("/orders")
			{
				orders.GET("", permMiddleware.RequirePermission("ORDER", "VIEW"), func(c *gin.Context) {
					// TODO: Implement order handler
					c.JSON(200, gin.H{"message": "Order list endpoint"})
				})
				orders.POST("", permMiddleware.RequirePermission("ORDER", "CREATE"), func(c *gin.Context) {
					// TODO: Implement order creation
					c.JSON(200, gin.H{"message": "Order creation endpoint"})
				})
				orders.GET("/:id", permMiddleware.RequirePermission("ORDER", "VIEW"), func(c *gin.Context) {
					// TODO: Implement order details
					c.JSON(200, gin.H{"message": "Order details endpoint"})
				})
				orders.PUT("/:id", permMiddleware.RequirePermission("ORDER", "UPDATE"), func(c *gin.Context) {
					// TODO: Implement order update
					c.JSON(200, gin.H{"message": "Order update endpoint"})
				})
				orders.DELETE("/:id", permMiddleware.RequirePermission("ORDER", "DELETE"), func(c *gin.Context) {
					// TODO: Implement order deletion
					c.JSON(200, gin.H{"message": "Order deletion endpoint"})
				})
				orders.POST("/:id/approve", permMiddleware.RequirePermission("ORDER", "APPROVE"), func(c *gin.Context) {
					// TODO: Implement order approval
					c.JSON(200, gin.H{"message": "Order approval endpoint"})
				})
				orders.POST("/:id/cancel", permMiddleware.RequirePermission("ORDER", "CANCEL"), func(c *gin.Context) {
					// TODO: Implement order cancellation
					c.JSON(200, gin.H{"message": "Order cancellation endpoint"})
				})
				orders.POST("/:id/ship", permMiddleware.RequirePermission("ORDER", "SHIP"), func(c *gin.Context) {
					// TODO: Implement order shipping
					c.JSON(200, gin.H{"message": "Order shipping endpoint"})
				})
			}

			// Warehouse Management Routes
			warehouse := protected.Group("/warehouse")
			{
				warehouse.GET("", permMiddleware.RequirePermission("WAREHOUSE", "VIEW"), func(c *gin.Context) {
					// TODO: Implement warehouse handler
					c.JSON(200, gin.H{"message": "Warehouse data endpoint"})
				})
				warehouse.POST("/transfer", permMiddleware.RequirePermission("WAREHOUSE", "TRANSFER"), func(c *gin.Context) {
					// TODO: Implement inventory transfer
					c.JSON(200, gin.H{"message": "Inventory transfer endpoint"})
				})
			}

			// Financial Management Routes
			finance := protected.Group("/finance")
			{
				finance.GET("", permMiddleware.RequirePermission("FINANCE", "VIEW"), func(c *gin.Context) {
					// TODO: Implement finance handler
					c.JSON(200, gin.H{"message": "Financial data endpoint"})
				})
				finance.POST("/approve", permMiddleware.RequirePermission("FINANCE", "APPROVE"), func(c *gin.Context) {
					// TODO: Implement financial approval
					c.JSON(200, gin.H{"message": "Financial approval endpoint"})
				})
			}

			// Reporting Routes
			reports := protected.Group("/reports")
			{
				reports.GET("", permMiddleware.RequirePermission("REPORT", "VIEW"), func(c *gin.Context) {
					// TODO: Implement reports handler
					c.JSON(200, gin.H{"message": "Reports list endpoint"})
				})
				reports.POST("", permMiddleware.RequirePermission("REPORT", "CREATE"), func(c *gin.Context) {
					// TODO: Implement report creation
					c.JSON(200, gin.H{"message": "Report creation endpoint"})
				})
				reports.GET("/:id/export", permMiddleware.RequirePermission("REPORT", "EXPORT"), func(c *gin.Context) {
					// TODO: Implement report export
					c.JSON(200, gin.H{"message": "Report export endpoint"})
				})
			}

			// System Administration Routes (Super Admin only)
			system := protected.Group("/system")
			system.Use(middleware.SuperAdminOnly())
			{
				system.GET("/logs", permMiddleware.RequirePermission("SYSTEM", "VIEW_LOGS"), func(c *gin.Context) {
					// TODO: Implement system logs
					c.JSON(200, gin.H{"message": "System logs endpoint"})
				})
				system.GET("/settings", permMiddleware.RequirePermission("SYSTEM", "MANAGE_SETTINGS"), func(c *gin.Context) {
					// TODO: Implement system settings
					c.JSON(200, gin.H{"message": "System settings endpoint"})
				})
				system.POST("/backup", permMiddleware.RequirePermission("SYSTEM", "BACKUP"), func(c *gin.Context) {
					// TODO: Implement system backup
					c.JSON(200, gin.H{"message": "System backup endpoint"})
				})
				system.POST("/restore", permMiddleware.RequirePermission("SYSTEM", "RESTORE"), func(c *gin.Context) {
					// TODO: Implement system restore
					c.JSON(200, gin.H{"message": "System restore endpoint"})
				})
			}
		}
	}

	return r
}
