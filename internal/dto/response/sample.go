package response

import "time"

type SampleResponse struct {
	ID                uint                 `json:"id"`
	SKU               string               `json:"sku"`
	ProductNameID     uint                 `json:"product_name_id"`
	CategoryID        uint                 `json:"category_id"`
	Description       string               `json:"description"`
	SampleType        string               `json:"sample_type"`
	Weight            float64              `json:"weight"`
	Width             float64              `json:"width"`
	Color             string               `json:"color"`
	ColorCode         string               `json:"color_code"`
	Quality           string               `json:"quality"`
	RemainingQuantity int                  `json:"remaining_quantity"`
	FiberContent      string               `json:"fiber_content"`
	Source            string               `json:"source"`
	SampleLocation    string               `json:"sample_location"`
	Barcode           string               `json:"barcode"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
	ProductName       *ProductNameResponse `json:"product_name,omitempty"`
	Category          *CategoryResponse    `json:"category,omitempty"`
}

type ProductNameResponse struct {
	ID            uint   `json:"id"`
	ProductNameVI string `json:"product_name_vi"`
	ProductNameEN string `json:"product_name_en"`
	SKUParent     string `json:"sku_parent"`
}

type CategoryResponse struct {
	ID           uint   `json:"id"`
	CategoryName string `json:"category_name"`
}

type PaginatedResponse struct {
	Items      []interface{} `json:"items"`
	TotalItems int64         `json:"total_items"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int           `json:"total_pages"`
}
