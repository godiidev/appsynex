package request

type SampleFilterRequest struct {
	Page       int     `form:"page" json:"page"`
	Limit      int     `form:"limit" json:"limit"`
	Search     string  `form:"search" json:"search"`
	Category   string  `form:"category" json:"category"`
	SampleType string  `form:"sample_type" json:"sample_type"`
	WeightMin  float64 `form:"weight_min" json:"weight_min"`
	WeightMax  float64 `form:"weight_max" json:"weight_max"`
	WidthMin   float64 `form:"width_min" json:"width_min"`
	WidthMax   float64 `form:"width_max" json:"width_max"`
	Color      string  `form:"color" json:"color"`
}

type CreateSampleRequest struct {
	SKU               string  `json:"sku" binding:"required"`
	ProductNameID     uint    `json:"product_name_id" binding:"required"`
	CategoryID        uint    `json:"category_id" binding:"required"`
	Description       string  `json:"description"`
	SampleType        string  `json:"sample_type"`
	Weight            float64 `json:"weight"`
	Width             float64 `json:"width"`
	Color             string  `json:"color"`
	ColorCode         string  `json:"color_code"`
	Quality           string  `json:"quality"`
	RemainingQuantity int     `json:"remaining_quantity"`
	FiberContent      string  `json:"fiber_content"`
	Source            string  `json:"source"`
	SampleLocation    string  `json:"sample_location"`
	Barcode           string  `json:"barcode"`
}

type UpdateSampleRequest struct {
	SKU               string   `json:"sku"`
	ProductNameID     uint     `json:"product_name_id"`
	CategoryID        uint     `json:"category_id"`
	Description       *string  `json:"description"` // Pointer để phân biệt null với empty string
	SampleType        *string  `json:"sample_type"`
	Weight            *float64 `json:"weight"`
	Width             *float64 `json:"width"`
	Color             *string  `json:"color"`
	ColorCode         *string  `json:"color_code"`
	Quality           *string  `json:"quality"`
	RemainingQuantity *int     `json:"remaining_quantity"`
	FiberContent      *string  `json:"fiber_content"`
	Source            *string  `json:"source"`
	SampleLocation    *string  `json:"sample_location"`
	Barcode           *string  `json:"barcode"`
}
