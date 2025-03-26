package main

import (
	"fmt"
	"log"
	"strings"
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

	// Run seeds
	if err := seedPermissions(db); err != nil {
		log.Fatalf("Failed to seed permissions: %v", err)
	}

	if err := seedAdmin(db); err != nil {
		log.Fatalf("Failed to seed admin user: %v", err)
	}

	log.Println("Seed completed successfully")
}

func seedPermissions(db *gorm.DB) error {
	// Define all permissions
	permissions := []string{
		// User permissions
		"USER_VIEW", "USER_CREATE", "USER_UPDATE", "USER_DELETE",
		// Role permissions
		"ROLE_VIEW", "ROLE_CREATE", "ROLE_UPDATE", "ROLE_DELETE",
		// Product permissions
		"PRODUCT_VIEW", "PRODUCT_CREATE", "PRODUCT_UPDATE", "PRODUCT_DELETE",
		// Sample permissions
		"SAMPLE_VIEW", "SAMPLE_CREATE", "SAMPLE_UPDATE", "SAMPLE_DELETE",
		// Customer permissions
		"CUSTOMER_VIEW", "CUSTOMER_CREATE", "CUSTOMER_UPDATE", "CUSTOMER_DELETE",
		// Order permissions
		"ORDER_VIEW", "ORDER_CREATE", "ORDER_UPDATE", "ORDER_DELETE",
	}

	// Insert permissions
	for _, permName := range permissions {
		var perm models.Permission
		err := db.Where("permission_name = ?", permName).First(&perm).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// Create permission if it doesn't exist
				description := fmt.Sprintf("Permission to %s", strings.ToLower(strings.Replace(permName, "_", " ", -1)))
				perm = models.Permission{
					PermissionName: permName,
					Description:    description,
				}
				if err := db.Create(&perm).Error; err != nil {
					return err
				}
				log.Printf("Created permission: %s", permName)
			} else {
				return err
			}
		} else {
			log.Printf("Permission already exists: %s", permName)
		}
	}

	// Get admin role
	var adminRole models.Role
	if err := db.Where("role_name = ?", "ADMIN").First(&adminRole).Error; err != nil {
		return err
	}

	// Get all permissions
	var allPermissions []models.Permission
	if err := db.Find(&allPermissions).Error; err != nil {
		return err
	}

	// Assign all permissions to admin role
	for _, perm := range allPermissions {
		var rolePermission models.RolePermission
		err := db.Where("role_id = ? AND permission_id = ?", adminRole.ID, perm.ID).First(&rolePermission).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// Get module name from permission name (e.g., USER_VIEW -> USER)
				module := strings.Split(perm.PermissionName, "_")[0]

				// Create role permission
				rolePermission = models.RolePermission{
					RoleID:       adminRole.ID,
					PermissionID: perm.ID,
					Module:       module,
				}
				if err := db.Create(&rolePermission).Error; err != nil {
					return err
				}
				log.Printf("Assigned permission %s to admin role", perm.PermissionName)
			} else {
				return err
			}
		}
	}

	return nil
}

func seedAdmin(db *gorm.DB) error {
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
				LastLogin:     nil,
			}

			if err := db.Create(&adminUser).Error; err != nil {
				return err
			}

			log.Println("Created admin user")

			// Assign admin role to admin user
			var adminRole models.Role
			if err := db.Where("role_name = ?", "ADMIN").First(&adminRole).Error; err != nil {
				return err
			}

			userRole := models.UserRole{
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
		log.Println("Admin user already exists")
	}

	return nil
}
