package models

import "gorm.io/gorm"

type Product struct {
    gorm.Model

    Name        string  `gorm:"not null;index"   json:"name"`        // index speeds up search queries
    Slug        string  `gorm:"uniqueIndex"      json:"slug"`
    Description string  `                        json:"description,omitempty"`
    Price       float64 `gorm:"not null"         json:"price"`
    Stock       int     `gorm:"default:0"        json:"stock"`
    ImageURL    string  `                        json:"image_url,omitempty"`
    IsActive    bool    `gorm:"default:true"     json:"is_active"`

    // belongs to: each product belongs to one category
    // CategoryID is the foreign key column in the database
    CategoryID uint     `gorm:"not null"         json:"category_id"`
    Category   Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// ProductResponse is the API-facing shape of a product
type ProductResponse struct {
    ID          uint             `json:"id"`
    Name        string           `json:"name"`
    Slug        string           `json:"slug"`
    Description string           `json:"description"`
    Price       float64          `json:"price"`
    Stock       int              `json:"stock"`
    ImageURL    string           `json:"image_url"`
    IsActive    bool             `json:"is_active"`
    CategoryID  uint             `json:"category_id"`
    Category    *CategorySummary `json:"category,omitempty"`
}

// CategorySummary is a small view of a category embedded in product responses
type CategorySummary struct {
    ID   uint   `json:"id"`
    Name string `json:"name"`
    Slug string `json:"slug"`
}
