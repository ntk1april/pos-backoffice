package handler

import (
	"net/http"
	"strconv"

	"pos-backoffice/internal/models"
	"pos-backoffice/internal/repository"
	"pos-backoffice/pkg/response"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionRepo *repository.TransactionRepository
	productRepo     *repository.ProductRepository
}

func NewTransactionHandler(transactionRepo *repository.TransactionRepository, productRepo *repository.ProductRepository) *TransactionHandler {
	return &TransactionHandler{
		transactionRepo: transactionRepo,
		productRepo:     productRepo,
	}
}

// CreateTransaction creates a new stock transaction
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req models.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	// Validate: DECREASE requires store_id
	if req.TransactionType == "DECREASE" && req.StoreID == nil {
		response.Error(c, http.StatusBadRequest, "Store ID is required for DECREASE transactions", nil)
		return
	}

	// Validate: INCREASE should not have store_id
	if req.TransactionType == "INCREASE" && req.StoreID != nil {
		response.Error(c, http.StatusBadRequest, "Store ID should not be provided for INCREASE transactions", nil)
		return
	}

	// Check if product exists and has enough stock for DECREASE
	if req.TransactionType == "DECREASE" {
		product, err := h.productRepo.FindByID(req.ProductID)
		if err != nil {
			response.Error(c, http.StatusNotFound, "Product not found", err)
			return
		}

		if product.Stock < req.Quantity {
			response.Error(c, http.StatusBadRequest, "Insufficient stock", nil)
			return
		}
	}

	userID := c.GetInt64("user_id")

	transaction := &models.Transaction{
		TransactionType: req.TransactionType,
		ProductID:       req.ProductID,
		StoreID:         req.StoreID,
		Quantity:        req.Quantity,
		UnitPrice:       req.UnitPrice,
		TotalAmount:     req.UnitPrice * float64(req.Quantity),
		Notes:           req.Notes,
		CreatedBy:       userID,
	}

	err := h.transactionRepo.Create(transaction)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create transaction", err)
		return
	}

	response.Success(c, "Transaction created successfully", transaction)
}

// GetTransactions returns all transactions with pagination
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	transactions, total, err := h.transactionRepo.GetAll(page, limit)
	if err != nil {
		// Log the actual error for debugging
		println("Error fetching transactions:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to fetch transactions", err)
		return
	}

	result := map[string]interface{}{
		"transactions": transactions,
		"total":        total,
		"page":         page,
		"limit":        limit,
	}

	response.Success(c, "Transactions retrieved successfully", result)
}

// GetTransactionsByProduct returns transactions for a specific product
func (h *TransactionHandler) GetTransactionsByProduct(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Param("product_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid product ID", err)
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	transactions, err := h.transactionRepo.GetByProductID(productID, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch transactions", err)
		return
	}

	response.Success(c, "Transactions retrieved successfully", transactions)
}

// GetTransactionsByStore returns transactions for a specific store
func (h *TransactionHandler) GetTransactionsByStore(c *gin.Context) {
	storeID, err := strconv.ParseInt(c.Param("store_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid store ID", err)
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	transactions, err := h.transactionRepo.GetByStoreID(storeID, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch transactions", err)
		return
	}

	response.Success(c, "Transactions retrieved successfully", transactions)
}
