package models

import (
	"time"

	"gorm.io/gorm"
)

// Role is a custom string type for user roles
// Using a custom type prevents typos — you can't accidentally write "Amdin"
type Role string

const (
	RoleCustomer Role = "customer"
	RoleAdmin    Role = "admin"
)

type User struct {
	// gorm.Model gives us:
	//   ID         uint (primary key, auto-increment)
	//   CreatedAt  time.Time
	//   UpdatedAt  time.Time
	//   DeletedAt  gorm.DeletedAt (for soft deletes — records aren't truly deleted)
	gorm.Model

	// uniqueIndex ensures no two users have the same email
	Email string `gorm:"uniqueIndex;not null" json:"email"`

	// Password is stored as a bcrypt hash, never plain text
	// json:"-" means this field is NEVER included in JSON output
	Password string `gorm:"not null"             json:"-"`

	FirstName string `gorm:"not null"             json:"first_name"`
	LastName  string `gorm:"not null"             json:"last_name"`
	Role      Role   `gorm:"default:customer"     json:"role"`
	Phone     string `                            json:"phone,omitempty"`

	// has many: one user can have many orders
	Orders []Order `gorm:"foreignKey:UserID"    json:"orders,omitempty"`

	// has one: one user has one cart
	Cart Cart `gorm:"foreignKey:UserID"    json:"cart,omitempty"`
}

// UserResponse is what we send back in API responses
// We define this separately so we never accidentally expose the password
type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      Role      `json:"role"`
	Phone     string    `json:"phone,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts a User to a UserResponse (strips the password)
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
		Phone:     u.Phone,
		CreatedAt: u.CreatedAt,
	}
}
