package services

import (
	"backend/backend/models"
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
    if filter.CategoryID > 0 {
        query = query.Where("category_id = ?", filter.CategoryID)
    }

    if filter.MinPrice > 0 {
        query = query.Where("price >= ?", filter.MinPrice)
    }

    if filter.MaxPrice > 0 {
        query = query.Where("price <= ?", filter.MaxPrice)
    }

    if filter.InStock {
        query = query.Where("stock > 0")

}

var total int64
query.Count(&total)
allowedSortFields := map[string]bool{"price" : true, "name": true, "created_at",true}
sortField := "created_at"
if allowedSortFields[filter.SortBy] {
	sortField = filter.SortBy
}

sortOrder := "desc"
if filter.SortOrder == "asc"{
	sortOrder = "asc"
}
Query = Query.Order (fmt.Sprintf("%s %s", sortField, sortOrder))
offset := (filter.Page - 1) * filter.Limit
query = query.Offset(offset).Limit(filter.Limit)
var products []models.Product
if err := query.Preload("Category").Find(&products).Error; err != nil {
    return nil, err
}
totalPages := int((total + int64(filter.Limit) - 1) / int64(filter.Limit))
return &PaginatedResult{
Data:       products,
Total:      total,
Page:       filter.Page,
Limit:      filter.Limit,
TotalPages: totalPages
},nil
}
func (s *ProductService) GetProductBySlug(slug string) (*models.Product, error) {
    var product models.Product
    if err := database.DB.Preload("Category").Where("slug = ? AND is_active = ?", slug, true).First(&product).Error; err != nil {
        return nil, errors.New("product not found")
    }
    return &product, nil
}

type CreateProductInput struct {
    Name        string  `json:"name"        binding:"required"`
    Description string  `json:"description"`
    Price       float64 `json:"price"       binding:"required,gt=0"`  // gt=0 means greater than 0
    Stock       int     `json:"stock"       binding:"min=0"`
    CategoryID  uint    `json:"category_id" binding:"required"`
    ImageURL    string  `json:"image_url"`
}

func (s *ProductService) CreateProduct(input CreateProductInput) (*models.Product, error) {
    var category models.Category
    if err := database.DB.First(&category, input.CategoryID).Error; err != nil {
        return nil, errors.New("category not found")
    }

    product := models.Product{
        Name:        input.Name,
        Slug:        generateSlug(input.Name),
        Description: input.Description,
        Price:       input.Price,
        Stock:       input.Stock,
        CategoryID:  input.CategoryID,
        ImageURL:    input.ImageURL,
        IsActive:    true,
    }

    if err := database.DB.Create(&product).Error; err != nil {
        return nil, errors.New("failed to create product")
    }

    database.DB.Preload("Category").First(&product, product.ID)
    return &product, nil
}
func (s *ProductService) DeleteProduct(id uint) error {
    if err := database.DB.Delete(&models.Product{}, id).Error; err != nil {
        return errors.New("failed to delete product")
    }
    return nil
}
func generateSlug(name string) string {
    slug := strings.ToLower(name)
    slug = strings.ReplaceAll(slug, " ", "-")
    return slug
}