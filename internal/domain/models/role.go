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

// UserRole junction table - defined here to avoid duplication
type UserRole struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	RoleID    uint      `gorm:"primaryKey" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
