// File: scripts/seed_enhanced_permissions.go
// Tạo tại: scripts/seed_enhanced_permissions.go
// Mục đích: Seed enhanced permission system data

package main

import (
	"log"
	"time"

	"github.com/godiidev/appsynex/config"
	"github.com/godiidev/appsynex/internal/domain/models"
	"github.com/godiidev/appsynex/internal/repository/mysql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := mysql.NewDBConnection(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Seed enhanced permission system
	if err := seedEnhancedPermissions(db); err != nil {
		log.Fatalf("Failed to seed enhanced permissions: %v", err)
	}

	// Seed roles with enhanced permissions
	if err := seedRolesWithPermissions(db); err != nil {
		log.Fatalf("Failed to seed roles with permissions: %v", err)
	}

	// Seed admin user with enhanced system
	if err := seedEnhancedAdmin(db); err != nil {
		log.Fatalf("Failed to seed enhanced admin user: %v", err)
	}

	// Seed test users with different roles
	if err := seedTestUsers(db); err != nil {
		log.Fatalf("Failed to seed test users: %v", err)
	}

	log.Println("Enhanced permission system seeding completed successfully")
}

func seedEnhancedPermissions(db *gorm.DB) error {
	log.Println("Seeding enhanced permissions...")

	// Check if permissions already exist
	var count int64
	if err := db.Model(&models.Permission{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Println("Permissions already exist, skipping...")
		return nil
	}

	// Insert all predefined permissions
	for _, permission := range models.PreDefinedPermissions {
		var existingPerm models.Permission
		err := db.Where("permission_name = ?", permission.PermissionName).First(&existingPerm).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&permission).Error; err != nil {
					return err
				}
				log.Printf("Created permission: %s", permission.PermissionName)
			} else {
				return err
			}
		} else {
			log.Printf("Permission already exists: %s", permission.PermissionName)
		}
	}

	// Insert permission groups
	for _, group := range models.PreDefinedPermissionGroups {
		var existingGroup models.PermissionGroup
		err := db.Where("group_name = ?", group.GroupName).First(&existingGroup).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&group).Error; err != nil {
					return err
				}
				log.Printf("Created permission group: %s", group.GroupName)
			} else {
				return err
			}
		} else {
			log.Printf("Permission group already exists: %s", group.GroupName)
		}
	}

	return nil
}

func seedRolesWithPermissions(db *gorm.DB) error {
	log.Println("Seeding roles with enhanced permissions...")

	// Define role permissions mapping
	rolePermissions := map[string][]string{
		"ADMIN": {
			// Full access to everything
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

	// Get admin user for granted_by field
	var adminUser models.User
	if err := db.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		return err
	}

	for roleName, permissionNames := range rolePermissions {
		// Get role
		var role models.Role
		if err := db.Where("role_name = ?", roleName).First(&role).Error; err != nil {
			log.Printf("Role %s not found, skipping...", roleName)
			continue
		}

		// Clear existing permissions
		if err := db.Where("role_id = ?", role.ID).Delete(&models.RolePermission{}).Error; err != nil {
			return err
		}

		// Assign new permissions
		for _, permName := range permissionNames {
			var permission models.Permission
			if err := db.Where("permission_name = ?", permName).First(&permission).Error; err != nil {
				log.Printf("Permission %s not found, skipping...", permName)
				continue
			}

			rolePermission := models.RolePermission{
				RoleID:       role.ID,
				PermissionID: permission.ID,
				GrantedBy:    adminUser.ID,
				GrantedAt:    time.Now(),
				IsActive:     true,
			}

			if err := db.Create(&rolePermission).Error; err != nil {
				return err
			}
		}

		log.Printf("Assigned %d permissions to role %s", len(permissionNames), roleName)
	}

	return nil
}

func seedEnhancedAdmin(db *gorm.DB) error {
	log.Println("Seeding enhanced admin user...")

	// Check if admin user exists
	var adminUser models.User
	err := db.Where("username = ?", "admin").First(&adminUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create admin user
			passwordHash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
			if err != nil {
				return err
			}

			adminUser = models.User{
				Username:      "admin",
				PasswordHash:  string(passwordHash),
				Email:         "admin@appsynex.vn",
				AccountStatus: "active",
			}

			if err := db.Create(&adminUser).Error; err != nil {
				return err
			}

			log.Println("Created admin user")
		} else {
			return err
		}
	} else {
		log.Println("Admin user already exists")
	}

	// Assign admin role
	var adminRole models.Role
	if err := db.Where("role_name = ?", "ADMIN").First(&adminRole).Error; err != nil {
		return err
	}

	// Check if role assignment exists
	var userRole models.UserRole
	err = db.Where("user_id = ? AND role_id = ?", adminUser.ID, adminRole.ID).First(&userRole).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			userRole = models.UserRole{
				UserID:    adminUser.ID,
				RoleID:    adminRole.ID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if err := db.Create(&userRole).Error; err != nil {
				return err
			}

			log.Println("Assigned admin role to admin user")
		} else {
			return err
		}
	} else {
		log.Println("Admin role already assigned to admin user")
	}

	return nil
}

func seedTestUsers(db *gorm.DB) error {
	log.Println("Seeding test users...")

	testUsers := []struct {
		Username string
		Email    string
		Password string
		RoleName string
	}{
		{"manager", "manager@appsynex.vn", "manager123", "MANAGER"},
		{"staff1", "staff1@appsynex.vn", "staff123", "STAFF"},
		{"staff2", "staff2@appsynex.vn", "staff123", "STAFF"},
	}

	for _, userData := range testUsers {
		// Check if user exists
		var user models.User
		err := db.Where("username = ?", userData.Username).First(&user).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// Create user
				passwordHash, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
				if err != nil {
					return err
				}

				user = models.User{
					Username:      userData.Username,
					PasswordHash:  string(passwordHash),
					Email:         userData.Email,
					AccountStatus: "active",
				}

				if err := db.Create(&user).Error; err != nil {
					return err
				}

				log.Printf("Created user: %s", userData.Username)
			} else {
				return err
			}
		} else {
			log.Printf("User already exists: %s", userData.Username)
			continue
		}

		// Assign role
		var role models.Role
		if err := db.Where("role_name = ?", userData.RoleName).First(&role).Error; err != nil {
			log.Printf("Role %s not found for user %s", userData.RoleName, userData.Username)
			continue
		}

		userRole := models.UserRole{
			UserID:    user.ID,
			RoleID:    role.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := db.Create(&userRole).Error; err != nil {
			return err
		}

		log.Printf("Assigned role %s to user %s", userData.RoleName, userData.Username)
	}

	return nil
}