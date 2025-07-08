// File: scripts/setup.go
// Tạo tại: scripts/setup.go
// Mục đích: Complete setup với enhanced permission system

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/godiidev/appsynex/config"
	"github.com/godiidev/appsynex/internal/domain/models"
	"github.com/godiidev/appsynex/internal/repository/mysql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	log.Println("🚀 Setting up AppSynex with enhanced permission system...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := mysql.NewDBConnection(&cfg.Database)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	log.Println("✅ Connected to database")

	// Auto migrate all models
	if err := autoMigrate(db); err != nil {
		log.Fatalf("❌ Auto migration failed: %v", err)
	}

	// Seed enhanced permissions
	if err := seedEnhancedPermissions(db); err != nil {
		log.Fatalf("❌ Failed to seed enhanced permissions: %v", err)
	}

	// Seed roles with permissions
	if err := seedRolesWithPermissions(db); err != nil {
		log.Fatalf("❌ Failed to seed roles with permissions: %v", err)
	}

	// Seed admin user
	if err := seedAdminUser(db); err != nil {
		log.Fatalf("❌ Failed to seed admin user: %v", err)
	}

	// Seed sample data
	if err := seedSampleData(db); err != nil {
		log.Fatalf("❌ Failed to seed sample data: %v", err)
	}

	log.Println("")
	log.Println("🎉 Setup completed successfully!")
	log.Println("📝 Default admin credentials:")
	log.Println("   Username: admin")
	log.Println("   Password: admin123")
	log.Println("")
	log.Println("🔗 Start server with: go run cmd/api/main.go")
	log.Println("🌐 API will be available at: http://localhost:8080")
}

func autoMigrate(db *gorm.DB) error {
	log.Println("📋 Running auto migrations...")

	// Define all models to migrate in correct order
	models := []interface{}{
		&models.Role{},
		&models.User{},
		&models.UserRole{},
		&models.Permission{},
		&models.PermissionGroup{},
		&models.RolePermission{},
		&models.UserPermission{},
		&models.ProductCategory{},
		&models.ProductName{},
		&models.Product{},
		&models.SampleProduct{},
	}

	// Run auto migration
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %w", model, err)
		}
		log.Printf("✅ Migrated: %T", model)
	}

	return nil
}

func seedEnhancedPermissions(db *gorm.DB) error {
	log.Println("🔑 Seeding enhanced permissions...")

	// Check if permissions already exist
	var count int64
	if err := db.Model(&models.Permission{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Println("📝 Permissions already exist, skipping...")
		return nil
	}

	// Insert all predefined permissions
	for _, permission := range models.PreDefinedPermissions {
		if err := db.Create(&permission).Error; err != nil {
			return fmt.Errorf("failed to create permission %s: %w", permission.PermissionName, err)
		}
	}
	log.Printf("✅ Created %d permissions", len(models.PreDefinedPermissions))

	// Insert permission groups
	for _, group := range models.PreDefinedPermissionGroups {
		if err := db.Create(&group).Error; err != nil {
			return fmt.Errorf("failed to create permission group %s: %w", group.GroupName, err)
		}
	}
	log.Printf("✅ Created %d permission groups", len(models.PreDefinedPermissionGroups))

	return nil
}

func seedRolesWithPermissions(db *gorm.DB) error {
	log.Println("👥 Seeding roles with permissions...")

	// Create default roles if they don't exist
	defaultRoles := []models.Role{
		{RoleName: "SUPER_ADMIN", Description: "Super Administrator with full system access"},
		{RoleName: "ADMIN", Description: "Administrator with full system access"},
		{RoleName: "MANAGER", Description: "Manager with limited administrative access"},
		{RoleName: "STAFF", Description: "Staff with basic access"},
	}

	for _, role := range defaultRoles {
		var existingRole models.Role
		err := db.Where("role_name = ?", role.RoleName).First(&existingRole).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&role).Error; err != nil {
					return fmt.Errorf("failed to create role %s: %w", role.RoleName, err)
				}
				log.Printf("✅ Created role: %s", role.RoleName)
			} else {
				return err
			}
		} else {
			log.Printf("📝 Role already exists: %s", role.RoleName)
		}
	}

	// Define role permissions mapping
	rolePermissions := map[string][]string{
		"SUPER_ADMIN": {
			// All permissions
			"USER_VIEW", "USER_CREATE", "USER_UPDATE", "USER_DELETE", "USER_ASSIGN_ROLES", "USER_RESET_PASSWORD", "USER_ASSIGN_PERMISSIONS",
			"ROLE_VIEW", "ROLE_CREATE", "ROLE_UPDATE", "ROLE_DELETE", "ROLE_ASSIGN_PERMISSIONS",
			"PRODUCT_VIEW", "PRODUCT_CREATE", "PRODUCT_UPDATE", "PRODUCT_DELETE", "PRODUCT_EXPORT", "PRODUCT_IMPORT",
			"PRODUCT_CATEGORY_VIEW", "PRODUCT_CATEGORY_CREATE", "PRODUCT_CATEGORY_UPDATE", "PRODUCT_CATEGORY_DELETE",
			"SAMPLE_VIEW", "SAMPLE_CREATE", "SAMPLE_UPDATE", "SAMPLE_DELETE", "SAMPLE_DISPATCH", "SAMPLE_TRACK",
			"CUSTOMER_VIEW", "CUSTOMER_CREATE", "CUSTOMER_UPDATE", "CUSTOMER_DELETE", "CUSTOMER_VIEW_ACTIVITY",
			"ORDER_VIEW", "ORDER_CREATE", "ORDER_UPDATE", "ORDER_DELETE", "ORDER_APPROVE", "ORDER_CANCEL", "ORDER_SHIP",
			"WAREHOUSE_VIEW", "WAREHOUSE_CREATE", "WAREHOUSE_UPDATE", "WAREHOUSE_DELETE", "WAREHOUSE_TRANSFER",
			"FINANCE_VIEW", "FINANCE_CREATE", "FINANCE_UPDATE", "FINANCE_DELETE", "FINANCE_APPROVE",
			"REPORT_VIEW", "REPORT_CREATE", "REPORT_EXPORT",
			"SYSTEM_VIEW", "SYSTEM_VIEW_LOGS", "SYSTEM_MANAGE_SETTINGS", "SYSTEM_BACKUP", "SYSTEM_RESTORE",
		},
		"ADMIN": {
			// Admin level access (same as SUPER_ADMIN except some system functions)
			"USER_VIEW", "USER_CREATE", "USER_UPDATE", "USER_DELETE", "USER_ASSIGN_ROLES", "USER_RESET_PASSWORD",
			"ROLE_VIEW", "ROLE_CREATE", "ROLE_UPDATE", "ROLE_DELETE", "ROLE_ASSIGN_PERMISSIONS",
			"PRODUCT_VIEW", "PRODUCT_CREATE", "PRODUCT_UPDATE", "PRODUCT_DELETE", "PRODUCT_EXPORT", "PRODUCT_IMPORT",
			"PRODUCT_CATEGORY_VIEW", "PRODUCT_CATEGORY_CREATE", "PRODUCT_CATEGORY_UPDATE", "PRODUCT_CATEGORY_DELETE",
			"SAMPLE_VIEW", "SAMPLE_CREATE", "SAMPLE_UPDATE", "SAMPLE_DELETE", "SAMPLE_DISPATCH", "SAMPLE_TRACK",
			"CUSTOMER_VIEW", "CUSTOMER_CREATE", "CUSTOMER_UPDATE", "CUSTOMER_DELETE", "CUSTOMER_VIEW_ACTIVITY",
			"ORDER_VIEW", "ORDER_CREATE", "ORDER_UPDATE", "ORDER_DELETE", "ORDER_APPROVE", "ORDER_CANCEL", "ORDER_SHIP",
			"WAREHOUSE_VIEW", "WAREHOUSE_CREATE", "WAREHOUSE_UPDATE", "WAREHOUSE_DELETE", "WAREHOUSE_TRANSFER",
			"FINANCE_VIEW", "FINANCE_CREATE", "FINANCE_UPDATE", "FINANCE_DELETE", "FINANCE_APPROVE",
			"REPORT_VIEW", "REPORT_CREATE", "REPORT_EXPORT",
			"SYSTEM_VIEW",
		},
		"MANAGER": {
			// Manager level access
			"USER_VIEW", "USER_CREATE", "USER_UPDATE", "USER_ASSIGN_ROLES",
			"PRODUCT_VIEW", "PRODUCT_CREATE", "PRODUCT_UPDATE", "PRODUCT_EXPORT",
			"PRODUCT_CATEGORY_VIEW", "PRODUCT_CATEGORY_CREATE", "PRODUCT_CATEGORY_UPDATE",
			"SAMPLE_VIEW", "SAMPLE_CREATE", "SAMPLE_UPDATE", "SAMPLE_DISPATCH", "SAMPLE_TRACK",
			"CUSTOMER_VIEW", "CUSTOMER_CREATE", "CUSTOMER_UPDATE", "CUSTOMER_VIEW_ACTIVITY",
			"ORDER_VIEW", "ORDER_CREATE", "ORDER_UPDATE", "ORDER_APPROVE", "ORDER_CANCEL", "ORDER_SHIP",
			"WAREHOUSE_VIEW", "WAREHOUSE_CREATE", "WAREHOUSE_UPDATE", "WAREHOUSE_TRANSFER",
			"FINANCE_VIEW", "FINANCE_CREATE", "FINANCE_UPDATE",
			"REPORT_VIEW", "REPORT_CREATE", "REPORT_EXPORT",
		},
		"STAFF": {
			// Basic staff access
			"USER_VIEW_OWN",
			"PRODUCT_VIEW",
			"PRODUCT_CATEGORY_VIEW",
			"SAMPLE_VIEW", "SAMPLE_CREATE", "SAMPLE_UPDATE", "SAMPLE_TRACK",
			"CUSTOMER_VIEW", "CUSTOMER_CREATE", "CUSTOMER_UPDATE",
			"ORDER_VIEW", "ORDER_CREATE", "ORDER_UPDATE",
			"WAREHOUSE_VIEW", "WAREHOUSE_CREATE", "WAREHOUSE_UPDATE",
			"REPORT_VIEW",
		},
	}

	// Get system user ID for granted_by field
	systemUserID := uint(1) // Will be created admin user

	// Assign permissions to roles
	for roleName, permissionNames := range rolePermissions {
		var role models.Role
		if err := db.Where("role_name = ?", roleName).First(&role).Error; err != nil {
			log.Printf("⚠️  Role %s not found, skipping...", roleName)
			continue
		}

		// Clear existing permissions for this role
		if err := db.Where("role_id = ?", role.ID).Delete(&models.RolePermission{}).Error; err != nil {
			return fmt.Errorf("failed to clear existing permissions for role %s: %w", roleName, err)
		}

		// Assign new permissions
		assignedCount := 0
		for _, permName := range permissionNames {
			var permission models.Permission
			if err := db.Where("permission_name = ?", permName).First(&permission).Error; err != nil {
				log.Printf("⚠️  Permission %s not found, skipping...", permName)
				continue
			}

			rolePermission := models.RolePermission{
				RoleID:       role.ID,
				PermissionID: permission.ID,
				GrantedBy:    systemUserID,
				GrantedAt:    time.Now(),
				IsActive:     true,
			}

			if err := db.Create(&rolePermission).Error; err != nil {
				log.Printf("⚠️  Failed to assign permission %s to role %s: %v", permName, roleName, err)
				continue
			}
			assignedCount++
		}

		log.Printf("✅ Assigned %d permissions to role %s", assignedCount, roleName)
	}

	return nil
}

func seedAdminUser(db *gorm.DB) error {
	log.Println("👤 Seeding admin user...")

	// Check if admin user exists
	var adminUser models.User
	err := db.Where("username = ?", "admin").First(&adminUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create admin user
			passwordHash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
			if err != nil {
				return fmt.Errorf("failed to hash password: %w", err)
			}

			adminUser = models.User{
				Username:      "admin",
				PasswordHash:  string(passwordHash),
				Email:         "admin@appsynex.vn",
				AccountStatus: "active",
			}

			if err := db.Create(&adminUser).Error; err != nil {
				return fmt.Errorf("failed to create admin user: %w", err)
			}

			log.Println("✅ Created admin user")
		} else {
			return fmt.Errorf("error checking admin user: %w", err)
		}
	} else {
		log.Println("📝 Admin user already exists")
	}

	// Assign SUPER_ADMIN role
	var superAdminRole models.Role
	if err := db.Where("role_name = ?", "SUPER_ADMIN").First(&superAdminRole).Error; err != nil {
		// Fallback to ADMIN role
		if err := db.Where("role_name = ?", "ADMIN").First(&superAdminRole).Error; err != nil {
			return fmt.Errorf("no admin role found: %w", err)
		}
	}

	// Check if role assignment exists
	var userRole models.UserRole
	err = db.Where("user_id = ? AND role_id = ?", adminUser.ID, superAdminRole.ID).First(&userRole).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			userRole = models.UserRole{
				UserID:    adminUser.ID,
				RoleID:    superAdminRole.ID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if err := db.Create(&userRole).Error; err != nil {
				return fmt.Errorf("failed to assign role to admin user: %w", err)
			}

			log.Printf("✅ Assigned %s role to admin user", superAdminRole.RoleName)
		} else {
			return fmt.Errorf("error checking role assignment: %w", err)
		}
	} else {
		log.Printf("📝 %s role already assigned to admin user", superAdminRole.RoleName)
	}

	return nil
}

func seedSampleData(db *gorm.DB) error {
	log.Println("📦 Seeding sample data...")

	// Seed product categories
	categories := []models.ProductCategory{
		{CategoryName: "Vải Thun Cotton", Description: "Các loại vải cotton và vải cotton pha"},
		{CategoryName: "Vải Thun Polyester", Description: "Các loại vải polyester"},
		{CategoryName: "Vải Khaki & Kaki", Description: "Các loại vải khaki dùng cho quần áo công sở và bình thường"},
		{CategoryName: "Vải thun pha", Description: "Các loại vải pha cotton, pha polyester, pha spandex"},
		{CategoryName: "Vải in", Description: "Các loại vải in ấn"},
	}

	for _, category := range categories {
		var existingCategory models.ProductCategory
		err := db.Where("category_name = ?", category.CategoryName).First(&existingCategory).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&category).Error; err != nil {
					log.Printf("⚠️  Failed to create category %s: %v", category.CategoryName, err)
				} else {
					log.Printf("✅ Created category: %s", category.CategoryName)
				}
			}
		} else {
			log.Printf("📝 Category already exists: %s", category.CategoryName)
		}
	}

	// Seed product names
	productNames := []models.ProductName{
		{ProductNameVI: "Vải Cotton 2 chiều", ProductNameEN: "Cotton Single Jersey", SKUParent: "SY1015"},
		{ProductNameVI: "Vải bo rib Cotton 2 chiều", ProductNameEN: "Cotton Rib 1x1", SKUParent: "SY1021"},
		{ProductNameVI: "Vải cá sấu CVC 2 chiều", ProductNameEN: "CVC Lacoste", SKUParent: "SY1041"},
		{ProductNameVI: "Cotton Khaki thun 2/1", ProductNameEN: "Cotton Khaki 2/1 Span", SKUParent: "SY1362"},
		{ProductNameVI: "Cotton Canvas không giãn", ProductNameEN: "Cotton Canvas", SKUParent: "SY1359"},
	}

	for _, productName := range productNames {
		var existingProductName models.ProductName
		err := db.Where("sku_parent = ?", productName.SKUParent).First(&existingProductName).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&productName).Error; err != nil {
					log.Printf("⚠️  Failed to create product name %s: %v", productName.ProductNameEN, err)
				} else {
					log.Printf("✅ Created product name: %s", productName.ProductNameEN)
				}
			}
		} else {
			log.Printf("📝 Product name already exists: %s", productName.ProductNameEN)
		}
	}

	// Create a few sample products
	var cottonCategory models.ProductCategory
	var cottonSingleJersey models.ProductName

	if err := db.Where("category_name = ?", "Vải Thun Cotton").First(&cottonCategory).Error; err != nil {
		log.Printf("⚠️  Could not find cotton category: %v", err)
		return nil // Not critical, skip sample products
	}

	if err := db.Where("sku_parent = ?", "SY1015").First(&cottonSingleJersey).Error; err != nil {
		log.Printf("⚠️  Could not find cotton single jersey: %v", err)
		return nil // Not critical, skip sample products
	}

	sampleProducts := []models.SampleProduct{
		{
			SKU:               "SY1015205185-WHT",
			ProductNameID:     cottonSingleJersey.ID,
			CategoryID:        cottonCategory.ID,
			Description:       "Mẫu vải cotton trắng 100%",
			SampleType:        "Vải mét",
			Weight:            205.0,
			Width:             185.0,
			Color:             "Trắng",
			ColorCode:         "250304A",
			Quality:           "Hàng bền màu 4",
			RemainingQuantity: 10,
			FiberContent:      "100% Cotton",
			Source:            "CMP20W mua việt thắng, dệt anh soạn",
			SampleLocation:    "Kệ A-12",
			Barcode:           "1234567890",
		},
		{
			SKU:               "SY1015220180-BLK",
			ProductNameID:     cottonSingleJersey.ID,
			CategoryID:        cottonCategory.ID,
			Description:       "Mẫu vải cotton đen 100%",
			SampleType:        "Vải cây",
			Weight:            180.0,
			Width:             220.0,
			Color:             "Đen",
			ColorCode:         "250307A",
			Quality:           "Hàng bền màu 4",
			RemainingQuantity: 8,
			FiberContent:      "100% Cotton",
			Source:            "Sợi CM30 mua việt thắng, dệt anh soạn",
			SampleLocation:    "Kệ A-12",
			Barcode:           "1234567891",
		},
	}

	for _, sample := range sampleProducts {
		var existingSample models.SampleProduct
		err := db.Where("sku = ?", sample.SKU).First(&existingSample).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&sample).Error; err != nil {
					log.Printf("⚠️  Failed to create sample product %s: %v", sample.SKU, err)
				} else {
					log.Printf("✅ Created sample product: %s", sample.SKU)
				}
			}
		} else {
			log.Printf("📝 Sample product already exists: %s", sample.SKU)
		}
	}

	return nil
}
