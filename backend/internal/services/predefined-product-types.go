package services

import (
	"github.com/petrejonn/naytife/internal/api/models"
)

// PredefinedProductTypeService handles predefined product type templates
type PredefinedProductTypeService struct{}

// NewPredefinedProductTypeService creates a new service instance
func NewPredefinedProductTypeService() *PredefinedProductTypeService {
	return &PredefinedProductTypeService{}
}

// GetPredefinedProductTypes returns all available predefined product type templates
func (s *PredefinedProductTypeService) GetPredefinedProductTypes() []models.PredefinedProductType {
	return []models.PredefinedProductType{
		{
			ID:           "apparel",
			Title:        "Apparel (Clothing)",
			Description:  "Clothing items including shirts, pants, dresses, and other garments",
			SkuSubstring: "APP",
			Shippable:    true,
			Digital:      false,
			Category:     "Fashion",
			Icon:         "👕",
			Attributes: []models.PredefinedAttributeTemplate{
				// Product Attributes
				{Title: "Brand", DataType: "Text", Required: true, AppliesTo: "Product"},
				{Title: "Material", DataType: "Option", Required: true, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "Cotton"}, {Value: "Polyester"}, {Value: "Wool"}, {Value: "Silk"},
					{Value: "Linen"}, {Value: "Denim"}, {Value: "Leather"}, {Value: "Synthetic Blend"},
				}},
				{Title: "Care Instructions", DataType: "Option", Required: true, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "Machine Wash Cold"}, {Value: "Hand Wash Only"}, {Value: "Dry Clean Only"}, {Value: "Machine Wash Warm"},
				}},
				{Title: "Season", DataType: "Option", Required: false, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "Spring"}, {Value: "Summer"}, {Value: "Fall"}, {Value: "Winter"}, {Value: "All Season"},
				}},
				// Variant Attributes
				{Title: "Size", DataType: "Option", Required: true, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "XS"}, {Value: "S"}, {Value: "M"}, {Value: "L"}, {Value: "XL"}, {Value: "XXL"}, {Value: "3XL"},
				}},
				{Title: "Color", DataType: "Color", Required: true, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "#000000"}, {Value: "#FFFFFF"}, {Value: "#FF0000"}, {Value: "#0000FF"},
					{Value: "#008000"}, {Value: "#FFFF00"}, {Value: "#FFA500"}, {Value: "#800080"},
				}},
				{Title: "Fit", DataType: "Option", Required: false, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "Regular"}, {Value: "Slim"}, {Value: "Relaxed"}, {Value: "Oversized"},
				}},
			},
		},
		{
			ID:           "footwear",
			Title:        "Footwear",
			Description:  "Shoes, boots, sandals, and other footwear items",
			SkuSubstring: "FOOT",
			Shippable:    true,
			Digital:      false,
			Category:     "Fashion",
			Icon:         "👟",
			Attributes: []models.PredefinedAttributeTemplate{
				// Product Attributes
				{Title: "Brand", DataType: "Text", Required: true, AppliesTo: "Product"},
				{Title: "Material", DataType: "Option", Required: true, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "Leather"}, {Value: "Canvas"}, {Value: "Synthetic"}, {Value: "Rubber"},
					{Value: "Suede"}, {Value: "Mesh"}, {Value: "Patent Leather"},
				}},
				{Title: "Type", DataType: "Option", Required: true, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "Sneakers"}, {Value: "Boots"}, {Value: "Sandals"}, {Value: "Dress Shoes"},
					{Value: "Athletic"}, {Value: "Casual"}, {Value: "Formal"},
				}},
				{Title: "Heel Height", DataType: "Option", Required: false, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "Flat"}, {Value: "Low (1-2 inches)"}, {Value: "Medium (2-3 inches)"}, {Value: "High (3+ inches)"},
				}},
				// Variant Attributes
				{Title: "Size", DataType: "Option", Required: true, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "5"}, {Value: "5.5"}, {Value: "6"}, {Value: "6.5"}, {Value: "7"}, {Value: "7.5"},
					{Value: "8"}, {Value: "8.5"}, {Value: "9"}, {Value: "9.5"}, {Value: "10"}, {Value: "10.5"},
					{Value: "11"}, {Value: "11.5"}, {Value: "12"}, {Value: "13"}, {Value: "14"},
				}},
				{Title: "Color", DataType: "Color", Required: true, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "#000000"}, {Value: "#FFFFFF"}, {Value: "#8B4513"}, {Value: "#000080"},
					{Value: "#808080"}, {Value: "#FF0000"}, {Value: "#00FF00"},
				}},
				{Title: "Width", DataType: "Option", Required: false, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "Narrow"}, {Value: "Medium"}, {Value: "Wide"}, {Value: "Extra Wide"},
				}},
			},
		},
		{
			ID:           "books",
			Title:        "Books",
			Description:  "Physical and digital books, textbooks, novels, and other publications",
			SkuSubstring: "BOOK",
			Shippable:    true,
			Digital:      true,
			Category:     "Media",
			Icon:         "📚",
			Attributes: []models.PredefinedAttributeTemplate{
				// Product Attributes
				{Title: "Author", DataType: "Text", Required: true, AppliesTo: "Product"},
				{Title: "Publisher", DataType: "Text", Required: true, AppliesTo: "Product"},
				{Title: "ISBN", DataType: "Text", Required: false, AppliesTo: "Product"},
				{Title: "Publication Date", DataType: "Date", Required: false, AppliesTo: "Product"},
				{Title: "Genre", DataType: "Option", Required: true, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "Fiction"}, {Value: "Non-Fiction"}, {Value: "Mystery"}, {Value: "Romance"},
					{Value: "Science Fiction"}, {Value: "Biography"}, {Value: "History"}, {Value: "Self-Help"},
					{Value: "Children's"}, {Value: "Educational"}, {Value: "Reference"},
				}},
				{Title: "Language", DataType: "Option", Required: true, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "English"}, {Value: "Spanish"}, {Value: "French"}, {Value: "German"}, {Value: "Chinese"},
				}},
				{Title: "Pages", DataType: "Number", Required: false, AppliesTo: "Product"},
				// Variant Attributes
				{Title: "Format", DataType: "Option", Required: true, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "Hardcover"}, {Value: "Paperback"}, {Value: "eBook"}, {Value: "Audiobook"},
				}},
				{Title: "Condition", DataType: "Option", Required: false, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "New"}, {Value: "Like New"}, {Value: "Good"}, {Value: "Fair"},
				}},
			},
		},
		{
			ID:           "electronics",
			Title:        "Electronics",
			Description:  "Consumer electronics, gadgets, and electronic devices",
			SkuSubstring: "ELEC",
			Shippable:    true,
			Digital:      false,
			Category:     "Technology",
			Icon:         "📱",
			Attributes: []models.PredefinedAttributeTemplate{
				// Product Attributes
				{Title: "Brand", DataType: "Text", Required: true, AppliesTo: "Product"},
				{Title: "Model", DataType: "Text", Required: true, AppliesTo: "Product"},
				{Title: "Warranty Period", DataType: "Option", Required: true, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "1 Year"}, {Value: "2 Years"}, {Value: "3 Years"}, {Value: "Lifetime"},
				}},
				{Title: "Power Source", DataType: "Option", Required: false, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "Battery"}, {Value: "AC Adapter"}, {Value: "USB"}, {Value: "Solar"}, {Value: "Rechargeable"},
				}},
				{Title: "Connectivity", DataType: "Option", Required: false, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "WiFi"}, {Value: "Bluetooth"}, {Value: "USB"}, {Value: "Ethernet"}, {Value: "Wireless"},
				}},
				// Variant Attributes
				{Title: "Color", DataType: "Color", Required: true, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "#000000"}, {Value: "#FFFFFF"}, {Value: "#C0C0C0"}, {Value: "#FFD700"},
					{Value: "#0000FF"}, {Value: "#FF0000"}, {Value: "#008000"},
				}},
				{Title: "Storage Capacity", DataType: "Option", Required: false, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "32GB"}, {Value: "64GB"}, {Value: "128GB"}, {Value: "256GB"},
					{Value: "512GB"}, {Value: "1TB"}, {Value: "2TB"},
				}},
				{Title: "Memory", DataType: "Option", Required: false, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "4GB"}, {Value: "8GB"}, {Value: "16GB"}, {Value: "32GB"}, {Value: "64GB"},
				}},
			},
		},
		{
			ID:           "beauty",
			Title:        "Beauty & Personal Care",
			Description:  "Cosmetics, skincare, haircare, and personal hygiene products",
			SkuSubstring: "BEAU",
			Shippable:    true,
			Digital:      false,
			Category:     "Health & Beauty",
			Icon:         "💄",
			Attributes: []models.PredefinedAttributeTemplate{
				// Product Attributes
				{Title: "Brand", DataType: "Text", Required: true, AppliesTo: "Product"},
				{Title: "Category", DataType: "Option", Required: true, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "Skincare"}, {Value: "Makeup"}, {Value: "Haircare"}, {Value: "Fragrance"},
					{Value: "Personal Hygiene"}, {Value: "Tools & Accessories"},
				}},
				{Title: "Skin Type", DataType: "Option", Required: false, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "All Skin Types"}, {Value: "Dry"}, {Value: "Oily"}, {Value: "Combination"},
					{Value: "Sensitive"}, {Value: "Normal"},
				}},
				{Title: "Ingredients", DataType: "Text", Required: false, AppliesTo: "Product"},
				{Title: "Expiration Date", DataType: "Date", Required: false, AppliesTo: "Product"},
				// Variant Attributes
				{Title: "Shade/Color", DataType: "Color", Required: false, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "#FDB5A6"}, {Value: "#F4C2A1"}, {Value: "#E8B685"}, {Value: "#D4A574"},
					{Value: "#C09458"}, {Value: "#A67C4A"}, {Value: "#8B6332"},
				}},
				{Title: "Size/Volume", DataType: "Option", Required: true, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "Travel Size (15ml)"}, {Value: "Regular (30ml)"}, {Value: "Large (50ml)"},
					{Value: "Value Size (100ml)"}, {Value: "Jumbo (200ml)"},
				}},
				{Title: "Finish", DataType: "Option", Required: false, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "Matte"}, {Value: "Glossy"}, {Value: "Satin"}, {Value: "Shimmer"}, {Value: "Natural"},
				}},
			},
		},
		{
			ID:           "home_appliances",
			Title:        "Home Appliances",
			Description:  "Kitchen appliances, home electronics, and household equipment",
			SkuSubstring: "HOME",
			Shippable:    true,
			Digital:      false,
			Category:     "Home & Garden",
			Icon:         "🏠",
			Attributes: []models.PredefinedAttributeTemplate{
				// Product Attributes
				{Title: "Brand", DataType: "Text", Required: true, AppliesTo: "Product"},
				{Title: "Model Number", DataType: "Text", Required: true, AppliesTo: "Product"},
				{Title: "Energy Rating", DataType: "Option", Required: false, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "A+++"}, {Value: "A++"}, {Value: "A+"}, {Value: "A"}, {Value: "B"}, {Value: "C"},
				}},
				{Title: "Warranty", DataType: "Option", Required: true, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "1 Year"}, {Value: "2 Years"}, {Value: "3 Years"}, {Value: "5 Years"},
				}},
				{Title: "Power Consumption", DataType: "Number", Unit: stringPtr("Watts"), Required: false, AppliesTo: "Product"},
				// Variant Attributes
				{Title: "Color", DataType: "Color", Required: true, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "#FFFFFF"}, {Value: "#000000"}, {Value: "#C0C0C0"}, {Value: "#FF0000"}, {Value: "#0000FF"},
				}},
				{Title: "Capacity", DataType: "Option", Required: false, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "Small (1-3L)"}, {Value: "Medium (4-6L)"}, {Value: "Large (7-10L)"}, {Value: "Extra Large (10L+)"},
				}},
				{Title: "Installation Type", DataType: "Option", Required: false, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "Countertop"}, {Value: "Built-in"}, {Value: "Freestanding"}, {Value: "Wall Mount"},
				}},
			},
		},
		{
			ID:           "furniture",
			Title:        "Furniture",
			Description:  "Home and office furniture including chairs, tables, beds, and storage",
			SkuSubstring: "FURN",
			Shippable:    true,
			Digital:      false,
			Category:     "Home & Garden",
			Icon:         "🪑",
			Attributes: []models.PredefinedAttributeTemplate{
				// Product Attributes
				{Title: "Brand", DataType: "Text", Required: true, AppliesTo: "Product"},
				{Title: "Material", DataType: "Option", Required: true, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "Wood"}, {Value: "Metal"}, {Value: "Plastic"}, {Value: "Glass"},
					{Value: "Fabric"}, {Value: "Leather"}, {Value: "Composite"},
				}},
				{Title: "Style", DataType: "Option", Required: true, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "Modern"}, {Value: "Traditional"}, {Value: "Contemporary"}, {Value: "Rustic"},
					{Value: "Industrial"}, {Value: "Minimalist"}, {Value: "Vintage"},
				}},
				{Title: "Room", DataType: "Option", Required: true, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "Living Room"}, {Value: "Bedroom"}, {Value: "Dining Room"}, {Value: "Office"},
					{Value: "Kitchen"}, {Value: "Bathroom"}, {Value: "Outdoor"},
				}},
				{Title: "Assembly Required", DataType: "Option", Required: true, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "Yes"}, {Value: "No"}, {Value: "Minimal"},
				}},
				// Variant Attributes
				{Title: "Color", DataType: "Color", Required: true, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "#8B4513"}, {Value: "#000000"}, {Value: "#FFFFFF"}, {Value: "#808080"},
					{Value: "#D2B48C"}, {Value: "#A0522D"}, {Value: "#CD853F"},
				}},
				{Title: "Size", DataType: "Option", Required: true, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "Small"}, {Value: "Medium"}, {Value: "Large"}, {Value: "Extra Large"}, {Value: "Custom"},
				}},
				{Title: "Dimensions", DataType: "Text", Required: false, AppliesTo: "ProductVariation"},
			},
		},
		{
			ID:           "jewelry",
			Title:        "Jewelry",
			Description:  "Rings, necklaces, earrings, bracelets, and other jewelry items",
			SkuSubstring: "JEWL",
			Shippable:    true,
			Digital:      false,
			Category:     "Fashion",
			Icon:         "💍",
			Attributes: []models.PredefinedAttributeTemplate{
				// Product Attributes
				{Title: "Brand", DataType: "Text", Required: true, AppliesTo: "Product"},
				{Title: "Type", DataType: "Option", Required: true, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "Ring"}, {Value: "Necklace"}, {Value: "Earrings"}, {Value: "Bracelet"},
					{Value: "Watch"}, {Value: "Pendant"}, {Value: "Brooch"}, {Value: "Anklet"},
				}},
				{Title: "Metal Type", DataType: "Option", Required: true, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "Gold"}, {Value: "Silver"}, {Value: "Platinum"}, {Value: "Titanium"},
					{Value: "Stainless Steel"}, {Value: "Rose Gold"}, {Value: "White Gold"},
				}},
				{Title: "Gemstone", DataType: "Option", Required: false, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "Diamond"}, {Value: "Ruby"}, {Value: "Sapphire"}, {Value: "Emerald"},
					{Value: "Pearl"}, {Value: "Amethyst"}, {Value: "Topaz"}, {Value: "None"},
				}},
				{Title: "Occasion", DataType: "Option", Required: false, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "Wedding"}, {Value: "Engagement"}, {Value: "Anniversary"}, {Value: "Birthday"},
					{Value: "Everyday"}, {Value: "Special Occasion"},
				}},
				// Variant Attributes
				{Title: "Size", DataType: "Option", Required: true, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "5"}, {Value: "5.5"}, {Value: "6"}, {Value: "6.5"}, {Value: "7"}, {Value: "7.5"},
					{Value: "8"}, {Value: "8.5"}, {Value: "9"}, {Value: "9.5"}, {Value: "10"},
					{Value: "One Size"}, {Value: "Adjustable"},
				}},
				{Title: "Finish", DataType: "Option", Required: false, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "Polished"}, {Value: "Matte"}, {Value: "Brushed"}, {Value: "Textured"},
				}},
			},
		},
		{
			ID:           "watches",
			Title:        "Watches",
			Description:  "Wristwatches, smartwatches, and timepieces",
			SkuSubstring: "WTCH",
			Shippable:    true,
			Digital:      false,
			Category:     "Fashion",
			Icon:         "⌚",
			Attributes: []models.PredefinedAttributeTemplate{
				// Product Attributes
				{Title: "Brand", DataType: "Text", Required: true, AppliesTo: "Product"},
				{Title: "Model", DataType: "Text", Required: true, AppliesTo: "Product"},
				{Title: "Movement Type", DataType: "Option", Required: true, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "Quartz"}, {Value: "Automatic"}, {Value: "Manual"}, {Value: "Digital"}, {Value: "Smartwatch"},
				}},
				{Title: "Case Material", DataType: "Option", Required: true, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "Stainless Steel"}, {Value: "Aluminum"}, {Value: "Titanium"}, {Value: "Gold"},
					{Value: "Silver"}, {Value: "Plastic"}, {Value: "Ceramic"},
				}},
				{Title: "Water Resistance", DataType: "Option", Required: true, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "30m"}, {Value: "50m"}, {Value: "100m"}, {Value: "200m"}, {Value: "300m+"},
				}},
				{Title: "Gender", DataType: "Option", Required: true, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "Men's"}, {Value: "Women's"}, {Value: "Unisex"}, {Value: "Kids"},
				}},
				// Variant Attributes
				{Title: "Band Material", DataType: "Option", Required: true, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "Leather"}, {Value: "Metal"}, {Value: "Rubber"}, {Value: "Silicone"},
					{Value: "Fabric"}, {Value: "Ceramic"}, {Value: "Nylon"},
				}},
				{Title: "Band Color", DataType: "Color", Required: true, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "#000000"}, {Value: "#8B4513"}, {Value: "#C0C0C0"}, {Value: "#FFD700"},
					{Value: "#0000FF"}, {Value: "#FF0000"}, {Value: "#FFFFFF"},
				}},
				{Title: "Case Size", DataType: "Option", Required: true, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "38mm"}, {Value: "40mm"}, {Value: "42mm"}, {Value: "44mm"}, {Value: "46mm"},
				}},
			},
		},
		{
			ID:           "mobile_tablets",
			Title:        "Mobile Phones & Tablets",
			Description:  "Smartphones, tablets, and mobile accessories",
			SkuSubstring: "MOB",
			Shippable:    true,
			Digital:      false,
			Category:     "Technology",
			Icon:         "📱",
			Attributes: []models.PredefinedAttributeTemplate{
				// Product Attributes
				{Title: "Brand", DataType: "Text", Required: true, AppliesTo: "Product"},
				{Title: "Model", DataType: "Text", Required: true, AppliesTo: "Product"},
				{Title: "Operating System", DataType: "Option", Required: true, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "iOS"}, {Value: "Android"}, {Value: "Windows"}, {Value: "Other"},
				}},
				{Title: "Screen Size", DataType: "Option", Required: true, AppliesTo: "Product", Options: []models.PredefinedAttributeOption{
					{Value: "5.0-5.4 inches"}, {Value: "5.5-5.9 inches"}, {Value: "6.0-6.4 inches"},
					{Value: "6.5+ inches"}, {Value: "7-8 inches"}, {Value: "9-10 inches"}, {Value: "11+ inches"},
				}},
				{Title: "Camera", DataType: "Text", Required: false, AppliesTo: "Product"},
				{Title: "Battery Life", DataType: "Text", Required: false, AppliesTo: "Product"},
				{Title: "Processor", DataType: "Text", Required: false, AppliesTo: "Product"},
				// Variant Attributes
				{Title: "Storage", DataType: "Option", Required: true, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "32GB"}, {Value: "64GB"}, {Value: "128GB"}, {Value: "256GB"},
					{Value: "512GB"}, {Value: "1TB"},
				}},
				{Title: "Color", DataType: "Color", Required: true, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "#000000"}, {Value: "#FFFFFF"}, {Value: "#FF0000"}, {Value: "#0000FF"},
					{Value: "#008000"}, {Value: "#FFD700"}, {Value: "#FFC0CB"}, {Value: "#800080"},
				}},
				{Title: "RAM", DataType: "Option", Required: false, AppliesTo: "ProductVariation", Options: []models.PredefinedAttributeOption{
					{Value: "4GB"}, {Value: "6GB"}, {Value: "8GB"}, {Value: "12GB"}, {Value: "16GB"},
				}},
			},
		},
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
