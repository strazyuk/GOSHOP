package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;not null" json:"name"`
	Slug        string `gorm:"uniqueIndex;not null" json:"slug"` // URL-friendly name, e.g. "mens-shoes"
	Description string `                            json:"description,omitempty"`
	ImageURL    string `                            json:"image_url,omitempty"`

	// has many: one category can have many products
	Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}
