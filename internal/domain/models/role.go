package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	RoleName    string         `gorm:"size:100;uniqueIndex" json:"role_name"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Users       []User         `gorm:"many2many:user_roles;" json:"users,omitempty"`
	Permissions []Permission   `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

type Permission struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	PermissionName string         `gorm:"size:100;uniqueIndex" json:"permission_name"`
	Description    string         `gorm:"type:text" json:"description"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Roles          []Role         `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
}

type RolePermission struct {
	RoleID       uint   `gorm:"primaryKey"`
	PermissionID uint   `gorm:"primaryKey"`
	Module       string `gorm:"size:100"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
