package main

import (
	"log"

	"github.com/godiidev/appsynex/config"
	"github.com/godiidev/appsynex/internal/domain/models"
	"github.com/godiidev/appsynex/internal/repository/mysql"
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

	// Seed product categories
	if err := seedProductCategories(db); err != nil {
		log.Fatalf("Failed to seed product categories: %v", err)
	}

	// Seed product names
	if err := seedProductNames(db); err != nil {
		log.Fatalf("Failed to seed product names: %v", err)
	}

	// Seed sample products
	if err := seedSampleProducts(db); err != nil {
		log.Fatalf("Failed to seed sample products: %v", err)
	}

	log.Println("Sample data seeding completed successfully!")
}

func seedProductCategories(db *gorm.DB) error {
	categories := []models.ProductCategory{
		{
			CategoryName: "Vải Thun Cotton",
			Description:  "Các loại vải cotton và vải cotton pha",
		},
		{
			CategoryName: "Vải Thun Polyester",
			Description:  "Các loại vải polyester",
		},
		{
			CategoryName: "Vải Khaki & Kaki",
			Description:  "Các loại vải khaki dùng cho quần áo công sở và bình thường",
		},
		{
			CategoryName: "Vải thun pha",
			Description:  "Các loại vải pha cotton, pha polyester, pha spandex",
		},
		{
			CategoryName: "Vải in",
			Description:  "Các loại vải in ấn",
		},
	}

	for _, category := range categories {
		var existingCategory models.ProductCategory
		err := db.Where("category_name = ?", category.CategoryName).First(&existingCategory).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&category).Error; err != nil {
					return err
				}
				log.Printf("Created category: %s", category.CategoryName)
			} else {
				return err
			}
		} else {
			log.Printf("Category already exists: %s", category.CategoryName)
		}
	}

	return nil
}

func seedProductNames(db *gorm.DB) error {
	productNames := []models.ProductName{
		{
			ProductNameVI: "Vải Cotton 2 chiều",
			ProductNameEN: "Cotton Single Jersey",
			SKUParent:     "SY1015",
		},
		{
			ProductNameVI: "Vải bo rib Cotton 2 chiều ",
			ProductNameEN: "Cotton Rib 1x1",
			SKUParent:     "SY1021",
		},
		{
			ProductNameVI: "Vải cá sấu CVC 2 chiều",
			ProductNameEN: "CVC Lacoste",
			SKUParent:     "SY1041",
		},
		{
			ProductNameVI: "Cotton Khaki thun 2/1",
			ProductNameEN: "Cotton Khaki 2/1 Span",
			SKUParent:     "SY1362",
		},
		{
			ProductNameVI: "Cotton Canvas không giãn",
			ProductNameEN: "Cotton Canvas",
			SKUParent:     "SY1359",
		},
		{
			ProductNameVI: "Poly Pique (mè kim)",
			ProductNameEN: "Vải mè kim",
			SKUParent:     "SY1266",
		},
		{
			ProductNameVI: "Vải Pique Cotton back",
			ProductNameEN: "Cotton Back Pique",
			SKUParent:     "SY1290",
		},
		{
			ProductNameVI: "Vải cá sấu Tici 2 chiều ",
			ProductNameEN: "TC Pique (cá mập)",
			SKUParent:     "SY1056",
		},
	}

	for _, productName := range productNames {
		var existingProductName models.ProductName
		err := db.Where("sku_parent = ?", productName.SKUParent).First(&existingProductName).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&productName).Error; err != nil {
					return err
				}
				log.Printf("Created product name: %s", productName.ProductNameEN)
			} else {
				return err
			}
		} else {
			log.Printf("Product name already exists: %s", productName.ProductNameEN)
		}
	}

	return nil
}

func seedSampleProducts(db *gorm.DB) error {
	// Get category IDs
	var cottonCategory models.ProductCategory
	var polyCategory models.ProductCategory
	var khakiCategory models.ProductCategory
	var mixCategory models.ProductCategory
	var printedCategory models.ProductCategory

	if err := db.Where("category_name = ?", "Vải Thun Cotton").First(&cottonCategory).Error; err != nil {
		return err
	}
	if err := db.Where("category_name = ?", "Vải Thun Polyester").First(&polyCategory).Error; err != nil {
		return err
	}
	if err := db.Where("category_name = ?", "Vải Khaki & Kaki").First(&khakiCategory).Error; err != nil {
		return err
	}
	if err := db.Where("category_name = ?", "Vải thun pha").First(&mixCategory).Error; err != nil {
		return err
	}
	if err := db.Where("category_name = ?", "Vải In").First(&printedCategory).Error; err != nil {
		return err
	}

	// Get product name IDs
	var cottonName models.ProductName
	var cottonPolyName models.ProductName
	var khakiName models.ProductName
	var cvcName models.ProductName
	var polyName models.ProductName

	if err := db.Where("sku_parent = ?", "SY1015").First(&cottonName).Error; err != nil {
		return err
	}
	if err := db.Where("sku_parent = ?", "SY1056").First(&cottonPolyName).Error; err != nil {
		return err
	}
	if err := db.Where("sku_parent = ?", "SY1290").First(&cvcName).Error; err != nil {
		return err
	}
	if err := db.Where("sku_parent = ?", "SY1266").First(&polyName).Error; err != nil {
		return err
	}
	if err := db.Where("sku_parent = ?", "SY1021").First(&cottonName).Error; err != nil {
		return err
	}
	if err := db.Where("sku_parent = ?", "SY1362").First(&khakiName).Error; err != nil {
		return err
	}

	// Create sample products
	sampleProducts := []models.SampleProduct{
		{
			SKU:               "SY1015205185-WHT",
			ProductNameID:     cottonName.ID,
			CategoryID:        cottonCategory.ID,
			Description:       "Mẫu vải cotton trắng 100%",
			SampleType:        "Vải mét",
			Weight:            205.0,
			Width:             185.0,
			Color:             "Trắng",
			ColorCode:         "250304A",
			Quality:           "Hàng bền màu 4",
			RemainingQuantity: 10,
			FiberContent:      "100% Cotton",
			Source:            "CMP20W mua việt thắng, dệt anh soạn",
			SampleLocation:    "Kệ A-12",
			Barcode:           "1234567890",
		},
		{
			SKU:               "SY1015220180-BLK",
			ProductNameID:     cottonName.ID,
			CategoryID:        cottonCategory.ID,
			Description:       "Mẫu vải cotton đen 100%",
			SampleType:        "Vải cây",
			Weight:            180.0,
			Width:             220.0,
			Color:             "Đen",
			ColorCode:         "250307A",
			Quality:           "Hàng bền màu 4",
			RemainingQuantity: 8,
			FiberContent:      "100% Cotton",
			Source:            "Sợi CM30 mua việt thắng, dệt anh soạn",
			SampleLocation:    "Kệ A-12",
			Barcode:           "1234567891",
		},
		{
			SKU:               "SY1056245195-3108A",
			ProductNameID:     cottonPolyName.ID,
			CategoryID:        mixCategory.ID,
			Description:       "Mẫu vải cotton pha polyester pique",
			SampleType:        "Bền màu phối",
			Weight:            160.0,
			Width:             150.0,
			Color:             "Xanh đen",
			ColorCode:         "3108A",
			Quality:           "Hàng cvc chất lượng tốt phối vải trắng",
			RemainingQuantity: 12,
			FiberContent:      "66% Cotton, 34% Polyester",
			Source:            "CM36 + Poly 150 dệt bên anh thuấn",
			SampleLocation:    "Kệ B-5",
			Barcode:           "1234567892",
		},
		{
			SKU:               "SY1021260160-250307A",
			ProductNameID:     cottonName.ID,
			CategoryID:        cottonCategory.ID,
			Description:       "Mẫu vải rib cotton span trắng",
			SampleType:        "Vải cây",
			Weight:            260.0,
			Width:             160.0,
			Color:             "Trắng kem",
			ColorCode:         "250307A",
			Quality:           "Cao cấp",
			RemainingQuantity: 5,
			FiberContent:      "95% Cotton, 5% Spandex",
			Source:            "CM30 + SP40 dệt bên anh sử",
			SampleLocation:    "Kệ D-1",
			Barcode:           "1234567893",
		},
		{
			SKU:               "SY1362270150-NV04",
			ProductNameID:     khakiName.ID,
			CategoryID:        khakiCategory.ID,
			Description:       "Mẫu vải kaki màu Navy",
			SampleType:        "Vải mét",
			Weight:            240.0,
			Width:             150.0,
			Color:             "Navy",
			ColorCode:         "NV04",
			Quality:           "Hàng bên màu cơ bản, không phối",
			RemainingQuantity: 15,
			FiberContent:      "100% Cotton",
			Source:            "Việt Nam",
			SampleLocation:    "Kệ C-8",
			Barcode:           "1234567894",
		},
		{
			SKU:               "SY1266135160-4610TCX",
			ProductNameID:     polyName.ID,
			CategoryID:        polyCategory.ID,
			Description:       "Mẫu vải mè poly giá rẻ",
			SampleType:        "Vải mét",
			Weight:            135.0,
			Width:             160.0,
			Color:             "Trắng",
			ColorCode:         "4610TCX",
			Quality:           "Hàng bền màu 3",
			RemainingQuantity: 7,
			FiberContent:      "100% Polyester",
			Source:            "Hàng mè 75/72 dệt bên nam thành",
			SampleLocation:    "Kệ E-3",
			Barcode:           "1234567895",
		},
		{
			SKU:               "SY1266200160-4610TCX",
			ProductNameID:     polyName.ID,
			CategoryID:        polyCategory.ID,
			Description:       "Mẫu vải poly sọc ngang",
			SampleType:        "Vải mét",
			Weight:            200.0,
			Width:             160.0,
			Color:             "Trắng",
			ColorCode:         "4610TCX",
			Quality:           "Hàng bền màu 3",
			RemainingQuantity: 7,
			FiberContent:      "100% Polyester",
			Source:            "Sợi Poly 150D dệt máy 2 giường kim",
			SampleLocation:    "Kệ E-3",
			Barcode:           "1234567895",
		},
	}

	for _, sample := range sampleProducts {
		var existingSample models.SampleProduct
		err := db.Where("sku = ?", sample.SKU).First(&existingSample).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&sample).Error; err != nil {
					return err
				}
				log.Printf("Created sample product: %s", sample.SKU)
			} else {
				return err
			}
		} else {
			log.Printf("Sample product already exists: %s", sample.SKU)
		}
	}

	return nil
}
