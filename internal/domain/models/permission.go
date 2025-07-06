// File: internal/domain/models/permission.go
// Tạo tại: internal/domain/models/permission.go
// Mục đích: Enhanced permission system với module và action granular control

package models

import (
	"time"

	"gorm.io/gorm"
)

// Permission defines granular permissions
type Permission struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Module         string         `gorm:"size:100;not null" json:"module"`         // USER, PRODUCT, SAMPLE, ORDER, CUSTOMER, etc.
	Action         string         `gorm:"size:50;not null" json:"action"`          // VIEW, CREATE, UPDATE, DELETE, EXPORT, IMPORT
	Resource       string         `gorm:"size:100" json:"resource"`                // Specific resource if needed
	PermissionName string         `gorm:"size:200;not null" json:"permission_name"` // Computed: MODULE_ACTION or MODULE_ACTION_RESOURCE
	Description    string         `gorm:"type:text" json:"description"`
	IsActive       bool           `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Roles          []Role         `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
}

// PermissionGroup defines logical grouping of permissions
type PermissionGroup struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	GroupName   string         `gorm:"size:100;not null" json:"group_name"`   // USER_MANAGEMENT, PRODUCT_MANAGEMENT, etc.
	DisplayName string         `gorm:"size:200;not null" json:"display_name"` // "User Management", "Product Management"
	Description string         `gorm:"type:text" json:"description"`
	Module      string         `gorm:"size:100;not null" json:"module"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// RolePermission junction table with additional metadata
type RolePermission struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RoleID       uint      `gorm:"not null" json:"role_id"`
	PermissionID uint      `gorm:"not null" json:"permission_id"`
	GrantedBy    uint      `json:"granted_by"`    // User who granted this permission
	GrantedAt    time.Time `json:"granted_at"`
	ExpiresAt    *time.Time `json:"expires_at"`   // Optional expiration
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
	Role       Role       `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Permission Permission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
	GrantedByUser User    `gorm:"foreignKey:GrantedBy" json:"granted_by_user,omitempty"`
}

// UserPermission for direct user permissions (overrides role permissions)
type UserPermission struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	UserID       uint       `gorm:"not null" json:"user_id"`
	PermissionID uint       `gorm:"not null" json:"permission_id"`
	GrantType    string     `gorm:"size:20;not null" json:"grant_type"` // GRANT, DENY
	GrantedBy    uint       `json:"granted_by"`
	GrantedAt    time.Time  `json:"granted_at"`
	ExpiresAt    *time.Time `json:"expires_at"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	Reason       string     `gorm:"type:text" json:"reason"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	
	User          User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Permission    Permission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
	GrantedByUser User       `gorm:"foreignKey:GrantedBy" json:"granted_by_user,omitempty"`
}

// PreDefinedPermissions - system permissions
var PreDefinedPermissions = []Permission{
	// User Management
	{Module: "USER", Action: "VIEW", PermissionName: "USER_VIEW", Description: "View users"},
	{Module: "USER", Action: "CREATE", PermissionName: "USER_CREATE", Description: "Create new users"},
	{Module: "USER", Action: "UPDATE", PermissionName: "USER_UPDATE", Description: "Update user information"},
	{Module: "USER", Action: "DELETE", PermissionName: "USER_DELETE", Description: "Delete users"},
	{Module: "USER", Action: "ASSIGN_ROLES", PermissionName: "USER_ASSIGN_ROLES", Description: "Assign roles to users"},
	{Module: "USER", Action: "RESET_PASSWORD", PermissionName: "USER_RESET_PASSWORD", Description: "Reset user passwords"},
	
	// Role Management
	{Module: "ROLE", Action: "VIEW", PermissionName: "ROLE_VIEW", Description: "View roles"},
	{Module: "ROLE", Action: "CREATE", PermissionName: "ROLE_CREATE", Description: "Create new roles"},
	{Module: "ROLE", Action: "UPDATE", PermissionName: "ROLE_UPDATE", Description: "Update role information"},
	{Module: "ROLE", Action: "DELETE", PermissionName: "ROLE_DELETE", Description: "Delete roles"},
	{Module: "ROLE", Action: "ASSIGN_PERMISSIONS", PermissionName: "ROLE_ASSIGN_PERMISSIONS", Description: "Assign permissions to roles"},
	
	// Product Management
	{Module: "PRODUCT", Action: "VIEW", PermissionName: "PRODUCT_VIEW", Description: "View products"},
	{Module: "PRODUCT", Action: "CREATE", PermissionName: "PRODUCT_CREATE", Description: "Create new products"},
	{Module: "PRODUCT", Action: "UPDATE", PermissionName: "PRODUCT_UPDATE", Description: "Update product information"},
	{Module: "PRODUCT", Action: "DELETE", PermissionName: "PRODUCT_DELETE", Description: "Delete products"},
	{Module: "PRODUCT", Action: "EXPORT", PermissionName: "PRODUCT_EXPORT", Description: "Export product data"},
	{Module: "PRODUCT", Action: "IMPORT", PermissionName: "PRODUCT_IMPORT", Description: "Import product data"},
	
	// Product Category Management
	{Module: "PRODUCT_CATEGORY", Action: "VIEW", PermissionName: "PRODUCT_CATEGORY_VIEW", Description: "View product categories"},
	{Module: "PRODUCT_CATEGORY", Action: "CREATE", PermissionName: "PRODUCT_CATEGORY_CREATE", Description: "Create product categories"},
	{Module: "PRODUCT_CATEGORY", Action: "UPDATE", PermissionName: "PRODUCT_CATEGORY_UPDATE", Description: "Update product categories"},
	{Module: "PRODUCT_CATEGORY", Action: "DELETE", PermissionName: "PRODUCT_CATEGORY_DELETE", Description: "Delete product categories"},
	
	// Sample Management
	{Module: "SAMPLE", Action: "VIEW", PermissionName: "SAMPLE_VIEW", Description: "View samples"},
	{Module: "SAMPLE", Action: "CREATE", PermissionName: "SAMPLE_CREATE", Description: "Create new samples"},
	{Module: "SAMPLE", Action: "UPDATE", PermissionName: "SAMPLE_UPDATE", Description: "Update sample information"},
	{Module: "SAMPLE", Action: "DELETE", PermissionName: "SAMPLE_DELETE", Description: "Delete samples"},
	{Module: "SAMPLE", Action: "DISPATCH", PermissionName: "SAMPLE_DISPATCH", Description: "Dispatch samples to customers"},
	{Module: "SAMPLE", Action: "TRACK", PermissionName: "SAMPLE_TRACK", Description: "Track sample status"},
	
	// Customer Management
	{Module: "CUSTOMER", Action: "VIEW", PermissionName: "CUSTOMER_VIEW", Description: "View customers"},
	{Module: "CUSTOMER", Action: "CREATE", PermissionName: "CUSTOMER_CREATE", Description: "Create new customers"},
	{Module: "CUSTOMER", Action: "UPDATE", PermissionName: "CUSTOMER_UPDATE", Description: "Update customer information"},
	{Module: "CUSTOMER", Action: "DELETE", PermissionName: "CUSTOMER_DELETE", Description: "Delete customers"},
	{Module: "CUSTOMER", Action: "VIEW_ACTIVITY", PermissionName: "CUSTOMER_VIEW_ACTIVITY", Description: "View customer activity logs"},
	
	// Order Management
	{Module: "ORDER", Action: "VIEW", PermissionName: "ORDER_VIEW", Description: "View orders"},
	{Module: "ORDER", Action: "CREATE", PermissionName: "ORDER_CREATE", Description: "Create new orders"},
	{Module: "ORDER", Action: "UPDATE", PermissionName: "ORDER_UPDATE", Description: "Update order information"},
	{Module: "ORDER", Action: "DELETE", PermissionName: "ORDER_DELETE", Description: "Delete orders"},
	{Module: "ORDER", Action: "APPROVE", PermissionName: "ORDER_APPROVE", Description: "Approve orders"},
	{Module: "ORDER", Action: "CANCEL", PermissionName: "ORDER_CANCEL", Description: "Cancel orders"},
	{Module: "ORDER", Action: "SHIP", PermissionName: "ORDER_SHIP", Description: "Ship orders"},
	
	// Warehouse Management
	{Module: "WAREHOUSE", Action: "VIEW", PermissionName: "WAREHOUSE_VIEW", Description: "View warehouse data"},
	{Module: "WAREHOUSE", Action: "CREATE", PermissionName: "WAREHOUSE_CREATE", Description: "Create warehouse entries"},
	{Module: "WAREHOUSE", Action: "UPDATE", PermissionName: "WAREHOUSE_UPDATE", Description: "Update warehouse data"},
	{Module: "WAREHOUSE", Action: "DELETE", PermissionName: "WAREHOUSE_DELETE", Description: "Delete warehouse entries"},
	{Module: "WAREHOUSE", Action: "TRANSFER", PermissionName: "WAREHOUSE_TRANSFER", Description: "Transfer inventory"},
	
	// Financial Management
	{Module: "FINANCE", Action: "VIEW", PermissionName: "FINANCE_VIEW", Description: "View financial data"},
	{Module: "FINANCE", Action: "CREATE", PermissionName: "FINANCE_CREATE", Description: "Create financial records"},
	{Module: "FINANCE", Action: "UPDATE", PermissionName: "FINANCE_UPDATE", Description: "Update financial data"},
	{Module: "FINANCE", Action: "DELETE", PermissionName: "FINANCE_DELETE", Description: "Delete financial records"},
	{Module: "FINANCE", Action: "APPROVE", PermissionName: "FINANCE_APPROVE", Description: "Approve financial transactions"},
	
	// Reporting
	{Module: "REPORT", Action: "VIEW", PermissionName: "REPORT_VIEW", Description: "View reports"},
	{Module: "REPORT", Action: "CREATE", PermissionName: "REPORT_CREATE", Description: "Create custom reports"},
	{Module: "REPORT", Action: "EXPORT", PermissionName: "REPORT_EXPORT", Description: "Export reports"},
	
	// System Administration
	{Module: "SYSTEM", Action: "VIEW_LOGS", PermissionName: "SYSTEM_VIEW_LOGS", Description: "View system logs"},
	{Module: "SYSTEM", Action: "MANAGE_SETTINGS", PermissionName: "SYSTEM_MANAGE_SETTINGS", Description: "Manage system settings"},
	{Module: "SYSTEM", Action: "BACKUP", PermissionName: "SYSTEM_BACKUP", Description: "Perform system backup"},
	{Module: "SYSTEM", Action: "RESTORE", PermissionName: "SYSTEM_RESTORE", Description: "Restore system from backup"},
}

var PreDefinedPermissionGroups = []PermissionGroup{
	{GroupName: "USER_MANAGEMENT", DisplayName: "User Management", Module: "USER", SortOrder: 1},
	{GroupName: "ROLE_MANAGEMENT", DisplayName: "Role & Permission Management", Module: "ROLE", SortOrder: 2},
	{GroupName: "PRODUCT_MANAGEMENT", DisplayName: "Product Management", Module: "PRODUCT", SortOrder: 3},
	{GroupName: "CATEGORY_MANAGEMENT", DisplayName: "Category Management", Module: "PRODUCT_CATEGORY", SortOrder: 4},
	{GroupName: "SAMPLE_MANAGEMENT", DisplayName: "Sample Management", Module: "SAMPLE", SortOrder: 5},
	{GroupName: "CUSTOMER_MANAGEMENT", DisplayName: "Customer Management", Module: "CUSTOMER", SortOrder: 6},
	{GroupName: "ORDER_MANAGEMENT", DisplayName: "Order Management", Module: "ORDER", SortOrder: 7},
	{GroupName: "WAREHOUSE_MANAGEMENT", DisplayName: "Warehouse Management", Module: "WAREHOUSE", SortOrder: 8},
	{GroupName: "FINANCIAL_MANAGEMENT", DisplayName: "Financial Management", Module: "FINANCE", SortOrder: 9},
	{GroupName: "REPORTING", DisplayName: "Reports & Analytics", Module: "REPORT", SortOrder: 10},
	{GroupName: "SYSTEM_ADMINISTRATION", DisplayName: "System Administration", Module: "SYSTEM", SortOrder: 11},
}