package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"pos-backoffice/internal/middleware"
	"pos-backoffice/internal/models"
	"pos-backoffice/internal/service"
	"pos-backoffice/pkg/response"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// GetProducts retrieves products with pagination and search
// @Summary List products
// @Description Get paginated list of products with optional search
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search term"
// @Param status query string false "Product status" Enums(ACTIVE, INACTIVE)
// @Success 200 {object} response.Response{data=models.ProductListResponse}
// @Router /api/products [get]
func (h *ProductHandler) GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")
	status := c.DefaultQuery("status", "ACTIVE")

	result, err := h.productService.GetProducts(page, pageSize, search, status)
	if err != nil {
		response.InternalServerError(c, "Failed to get products", err)
		return
	}

	response.Success(c, "Products retrieved successfully", result)
}

// GetProduct retrieves a product by ID
// @Summary Get product
// @Description Get product details by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.Response{data=models.Product}
// @Router /api/products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid product ID", err)
		return
	}

	product, err := h.productService.GetProductByID(id)
	if err != nil {
		response.NotFound(c, "Product not found")
		return
	}

	response.Success(c, "Product retrieved successfully", product)
}

// CreateProduct creates a new product
// @Summary Create product
// @Description Create a new product (ADMIN only)
// @Tags products
// @Accept json
// @Produce json
// @Param request body models.CreateProductRequest true "Product data"
// @Success 201 {object} response.Response{data=models.Product}
// @Router /api/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	userID := middleware.GetUserID(c)
	product, err := h.productService.CreateProduct(&req, userID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Product created successfully", product)
}

// UpdateProduct updates an existing product
// @Summary Update product
// @Description Update product details (ADMIN only)
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param request body models.UpdateProductRequest true "Product data"
// @Success 200 {object} response.Response{data=models.Product}
// @Router /api/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid product ID", err)
		return
	}

	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	userID := middleware.GetUserID(c)
	product, err := h.productService.UpdateProduct(id, &req, userID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Product updated successfully", product)
}

// DeleteProduct soft deletes a product
// @Summary Delete product
// @Description Soft delete a product (ADMIN only)
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.Response
// @Router /api/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid product ID", err)
		return
	}

	userID := middleware.GetUserID(c)
	err = h.productService.DeleteProduct(id, userID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Product deleted successfully", nil)
}
