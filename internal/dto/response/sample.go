// File: internal/dto/response/sample.go
// Clean version - NO duplicate CategoryResponse

package response

import "time"

type SampleResponse struct {
	ID                uint      `json:"id"`
	SKU               string    `json:"sku"`
	ProductNameID     uint      `json:"product_name_id"` // Keep for reference
	CategoryID        uint      `json:"category_id"`     // Keep for reference
	Description       string    `json:"description"`
	SampleType        string    `json:"sample_type"`
	Weight            float64   `json:"weight"`
	Width             float64   `json:"width"`
	Color             string    `json:"color"`
	ColorCode         string    `json:"color_code"`
	Quality           string    `json:"quality"`
	RemainingQuantity int       `json:"remaining_quantity"`
	FiberContent      string    `json:"fiber_content"`
	Source            string    `json:"source"`
	SampleLocation    string    `json:"sample_location"`
	Barcode           string    `json:"barcode"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	// Use existing structs - CategoryResponse is already defined in category.go
	ProductName *ProductNameResponse `json:"product_name,omitempty"`
	Category    *CategoryResponse    `json:"category,omitempty"` // This uses the one from category.go
}

type ProductNameResponse struct {
	ID            uint   `json:"id"`
	ProductNameVI string `json:"product_name_vi"`
	ProductNameEN string `json:"product_name_en"`
	SKUParent     string `json:"sku_parent"`
}

// Enhanced response for listings (optional)
type SampleListResponse struct {
	ID                uint      `json:"id"`
	SKU               string    `json:"sku"`
	ProductNameVI     string    `json:"product_name_vi"` // Direct field for easy display
	ProductNameEN     string    `json:"product_name_en"` // Direct field for easy display
	CategoryName      string    `json:"category_name"`   // Direct field for easy display
	SampleType        string    `json:"sample_type"`
	Weight            float64   `json:"weight"`
	Width             float64   `json:"width"`
	Color             string    `json:"color"`
	RemainingQuantity int       `json:"remaining_quantity"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type PaginatedResponse struct {
	Items      []interface{} `json:"items"`
	TotalItems int64         `json:"total_items"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int           `json:"total_pages"`
}
