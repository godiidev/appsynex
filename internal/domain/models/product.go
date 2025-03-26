package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductName struct {
	ID            uint            `gorm:"primaryKey" json:"id"`
	ProductNameVI string          `gorm:"size:255" json:"product_name_vi"`
	ProductNameEN string          `gorm:"size:255" json:"product_name_en"`
	SKUParent     string          `gorm:"size:50" json:"sku_parent"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`
	Products      []Product       `json:"products,omitempty"`
	Samples       []SampleProduct `json:"samples,omitempty"`
}

type ProductCategory struct {
	ID               uint              `gorm:"primaryKey" json:"id"`
	CategoryName     string            `gorm:"size:255" json:"category_name"`
	ParentCategoryID *uint             `json:"parent_category_id"`
	Description      string            `gorm:"type:text" json:"description"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	DeletedAt        gorm.DeletedAt    `gorm:"index" json:"-"`
	Products         []Product         `json:"products,omitempty"`
	Samples          []SampleProduct   `json:"samples,omitempty"`
	ParentCategory   *ProductCategory  `gorm:"foreignKey:ParentCategoryID" json:"parent_category,omitempty"`
	ChildCategories  []ProductCategory `gorm:"foreignKey:ParentCategoryID" json:"child_categories,omitempty"`
}

type Product struct {
	ID             uint            `gorm:"primaryKey" json:"id"`
	ProductNameID  uint            `json:"product_name_id"`
	CategoryID     uint            `json:"category_id"`
	SKU            string          `gorm:"size:100;uniqueIndex" json:"sku"`
	SKUVariant     string          `gorm:"size:100" json:"sku_variant"`
	Description    string          `gorm:"type:text" json:"description"`
	FabricType     string          `gorm:"size:255" json:"fabric_type"`
	Weight         float64         `json:"weight"`
	Width          float64         `json:"width"`
	Color          string          `gorm:"size:150" json:"color"`
	Quality        string          `gorm:"size:50" json:"quality"`
	FiberContent   string          `gorm:"size:255" json:"fiber_content"`
	AdditionalInfo string          `gorm:"type:json" json:"additional_info"`
	Price          float64         `json:"price"`
	SalesPrice     float64         `json:"sales_price"`
	StockQuantity  float64         `json:"stock_quantity"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`
	ProductName    ProductName     `gorm:"foreignKey:ProductNameID" json:"product_name,omitempty"`
	Category       ProductCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}
