package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Username      string         `gorm:"size:100;uniqueIndex" json:"username"`
	PasswordHash  string         `gorm:"size:255" json:"-"`
	Email         string         `gorm:"size:255" json:"email"`
	Phone         *string        `gorm:"size:50" json:"phone"`
	LastLogin     *time.Time     `json:"last_login"`
	AccountStatus string         `gorm:"size:50;default:active" json:"account_status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Roles         []Role         `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

type UserRole struct {
	UserID    uint `gorm:"primaryKey"`
	RoleID    uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
