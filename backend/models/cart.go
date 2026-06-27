package models

import "gorm.io/gorm"

// Cart belongs to a User (one-to-one)
type Cart struct {
    gorm.Model
    UserID uint       `gorm:"uniqueIndex;not null" json:"user_id"`  // unique: one cart per user
    Items  []CartItem `gorm:"foreignKey:CartID"    json:"items,omitempty"`
}

// CartItem is one product inside a cart
type CartItem struct {
    gorm.Model
    CartID    uint    `gorm:"not null"              json:"cart_id"`
    ProductID uint    `gorm:"not null"              json:"product_id"`
    Product   Product `gorm:"foreignKey:ProductID"  json:"product,omitempty"`
    Quantity  int     `gorm:"not null;default:1"    json:"quantity"`
}

// CartResponse is what we send to the client — includes computed totals
type CartResponse struct {
    ID         uint               `json:"id"`
    UserID     uint               `json:"user_id"`
    Items      []CartItemResponse `json:"items"`
    TotalItems int                `json:"total_items"`   // sum of all quantities
    TotalPrice float64            `json:"total_price"`   // sum of price * quantity
}

type CartItemResponse struct {
    ID        uint    `json:"id"`
    ProductID uint    `json:"product_id"`
    Name      string  `json:"name"`
    Price     float64 `json:"price"`
    ImageURL  string  `json:"image_url"`
    Quantity  int     `json:"quantity"`
    Subtotal  float64 `json:"subtotal"`  // price * quantity
}
