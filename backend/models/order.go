package models

import (
    "time"

    "gorm.io/gorm"
)

// OrderStatus tracks where an order is in the fulfillment process
type OrderStatus string

const (
    StatusPending    OrderStatus = "pending"
    StatusConfirmed  OrderStatus = "confirmed"
    StatusShipping   OrderStatus = "shipping"
    StatusDelivered  OrderStatus = "delivered"
    StatusCancelled  OrderStatus = "cancelled"
)

type Order struct {
    gorm.Model
    UserID          uint        `gorm:"not null"              json:"user_id"`
    User            User        `gorm:"foreignKey:UserID"     json:"user,omitempty"`
    Status          OrderStatus `gorm:"default:pending"       json:"status"`
    TotalAmount     float64     `gorm:"not null"              json:"total_amount"`

    // Shipping address — stored as flat fields for simplicity
    // (in a real app you might use a separate Address model)
    ShippingName    string `json:"shipping_name"`
    ShippingAddress string `json:"shipping_address"`
    ShippingCity    string `json:"shipping_city"`
    ShippingZip     string `json:"shipping_zip"`
    ShippingCountry string `json:"shipping_country"`

    // Notes from the customer
    Notes string `json:"notes,omitempty"`

    // When the order was confirmed/shipped/delivered
    ConfirmedAt *time.Time `json:"confirmed_at,omitempty"`
    ShippedAt   *time.Time `json:"shipped_at,omitempty"`
    DeliveredAt *time.Time `json:"delivered_at,omitempty"`

    // Line items — the actual products ordered and at what price
    Items []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

// OrderItem is a snapshot of a product at the time of purchase.
// We store price separately (not just a foreign key to Product)
// because prices can change later — we want the historical price.
type OrderItem struct {
    gorm.Model
    OrderID   uint    `gorm:"not null"               json:"order_id"`
    ProductID uint    `gorm:"not null"               json:"product_id"`
    Product   Product `gorm:"foreignKey:ProductID"   json:"product,omitempty"`

    // Snapshot values at time of purchase
    ProductName string  `gorm:"not null"              json:"product_name"`
    UnitPrice   float64 `gorm:"not null"              json:"unit_price"`
    Quantity    int     `gorm:"not null"              json:"quantity"`
    Subtotal    float64 `gorm:"not null"              json:"subtotal"`
}
