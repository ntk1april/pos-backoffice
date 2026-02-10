package service

import (
	"fmt"
	"math"

	"pos-backoffice/internal/models"
	"pos-backoffice/internal/repository"
)

type ProductService struct {
	productRepo *repository.ProductRepository
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

// GetProducts retrieves products with pagination and search
func (s *ProductService) GetProducts(page, pageSize int, search, status string) (*models.ProductListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	products, total, err := s.productRepo.FindAll(page, pageSize, search, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &models.ProductListResponse{
		Products:   products,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetProductByID retrieves a product by ID
func (s *ProductService) GetProductByID(id int64) (*models.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return product, nil
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(req *models.CreateProductRequest, userID int64) (*models.Product, error) {
	// Check if SKU already exists
	existing, err := s.productRepo.FindBySKU(req.SKU)
	if err != nil {
		return nil, fmt.Errorf("failed to check SKU: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("SKU already exists")
	}

	// Validate price > cost
	if req.Price < req.Cost {
		return nil, fmt.Errorf("price must be greater than or equal to cost")
	}

	product := &models.Product{
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Cost:        req.Cost,
		Stock:       req.Stock,
		Status:      "ACTIVE",
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}

	err = s.productRepo.Create(product)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(id int64, req *models.UpdateProductRequest, userID int64) (*models.Product, error) {
	// Check if product exists
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}

	// Validate price > cost
	if req.Price < req.Cost {
		return nil, fmt.Errorf("price must be greater than or equal to cost")
	}

	// Update fields
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Cost = req.Cost
	product.UpdatedBy = userID

	err = s.productRepo.Update(product)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

// DeleteProduct soft deletes a product
func (s *ProductService) DeleteProduct(id int64, userID int64) error {
	err := s.productRepo.Delete(id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}
