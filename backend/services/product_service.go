package services

import (
    "errors"
    "fmt"
    "strings"

    "go-ecommerce/database"
    "go-ecommerce/models"
)

type ProductService struct{}
func NewProductService() *ProductService {
	return &ProductService{}
}

type ProductFilter struct {
 Search     string  // search by name
    CategoryID uint    // filter by category
    MinPrice   float64 // minimum price
    MaxPrice   float64 // maximum price
    InStock    bool    // only show in-stock products
    Page       int     // page number (for pagination)
    Limit      int     // items per page
    SortBy     string  // field to sort by (price, name, created_at)
    SortOrder  string
	
}

type PaginatedResult struct {
    Data       interface{} `json:"data"`
    Total      int64       `json:"total"`      // total matching records
    Page       int         `json:"page"`
    Limit      int         `json:"limit"`
    TotalPages int         `json:"total_pages"`
}

func (s *ProductService) ListProducts (filter ProductFilter) (*PaginatedResult, error) {
	if filter.Page < 1 {
		filter.Page = 1 
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 20
	}
 query := database.DB.Model(&models.Product{}).Where("is_active = ?", true)

    // Apply filters — each is conditional
    if filter.Search != "" {
        // LIKE with % wildcards: finds products whose name contains the search term
        // The ? prevents SQL injection by parameterizing the value
        searchTerm := "%" + strings.ToLower(filter.Search) + "%"
        query = query.Where("LOWER(name) LIKE ?", searchTerm)
    }

}