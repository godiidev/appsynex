package models

import (
	"time"

	"gorm.io/gorm"
)

type SampleProduct struct {
	ID                uint            `gorm:"primaryKey" json:"id"`
	SKU               string          `gorm:"size:100;uniqueIndex" json:"sku"`
	ProductNameID     uint            `json:"product_name_id"`
	CategoryID        uint            `json:"category_id"`
	Description       string          `gorm:"type:text" json:"description"`
	SampleType        string          `gorm:"size:255" json:"sample_type"`
	Weight            float64         `json:"weight"`
	Width             float64         `json:"width"`
	Color             string          `gorm:"size:255" json:"color"`
	ColorCode         string          `gorm:"size:50" json:"color_code"`
	Quality           string          `gorm:"type:text" json:"quality"`
	RemainingQuantity int             `json:"remaining_quantity"`
	FiberContent      string          `gorm:"size:255" json:"fiber_content"`
	Source            string          `gorm:"size:255" json:"source"`
	SampleLocation    string          `gorm:"size:255" json:"sample_location"`
	Barcode           string          `gorm:"size:255" json:"barcode"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	DeletedAt         gorm.DeletedAt  `gorm:"index" json:"-"`
	ProductName       ProductName     `gorm:"foreignKey:ProductNameID" json:"product_name,omitempty"`
	Category          ProductCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}
