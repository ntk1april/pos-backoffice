package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"pos-backoffice/internal/middleware"
	"pos-backoffice/internal/models"
	"pos-backoffice/internal/service"
	"pos-backoffice/pkg/response"
)

type StockHandler struct {
	stockService *service.StockService
}

func NewStockHandler(stockService *service.StockService) *StockHandler {
	return &StockHandler{
		stockService: stockService,
	}
}

// IncreaseStock increases product stock
// @Summary Increase stock
// @Description Increase product stock quantity
// @Tags stock
// @Accept json
// @Produce json
// @Param request body models.StockAdjustmentRequest true "Stock adjustment data"
// @Success 200 {object} response.Response
// @Router /api/stock/increase [post]
func (h *StockHandler) IncreaseStock(c *gin.Context) {
	var req models.StockAdjustmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	userID := middleware.GetUserID(c)
	err := h.stockService.IncreaseStock(&req, userID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Stock increased successfully", nil)
}

// DecreaseStock decreases product stock
// @Summary Decrease stock
// @Description Decrease product stock quantity
// @Tags stock
// @Accept json
// @Produce json
// @Param request body models.StockAdjustmentRequest true "Stock adjustment data"
// @Success 200 {object} response.Response
// @Router /api/stock/decrease [post]
func (h *StockHandler) DecreaseStock(c *gin.Context) {
	var req models.StockAdjustmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	userID := middleware.GetUserID(c)
	err := h.stockService.DecreaseStock(&req, userID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Stock decreased successfully", nil)
}

// GetStockLogs retrieves stock logs for a product
// @Summary Get stock logs
// @Description Get stock transaction history for a product
// @Tags stock
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} response.Response{data=models.StockLogResponse}
// @Router /api/stock/logs/{product_id} [get]
func (h *StockHandler) GetStockLogs(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Param("product_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid product ID", err)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.stockService.GetStockLogs(productID, page, pageSize)
	if err != nil {
		response.InternalServerError(c, "Failed to get stock logs", err)
		return
	}

	response.Success(c, "Stock logs retrieved successfully", result)
}
