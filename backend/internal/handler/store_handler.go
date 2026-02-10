package handler

import (
	"net/http"
	"strconv"

	"pos-backoffice/internal/models"
	"pos-backoffice/internal/repository"
	"pos-backoffice/pkg/response"

	"github.com/gin-gonic/gin"
)

type StoreHandler struct {
	storeRepo *repository.StoreRepository
}

func NewStoreHandler(storeRepo *repository.StoreRepository) *StoreHandler {
	return &StoreHandler{storeRepo: storeRepo}
}

// GetStores returns all stores
func (h *StoreHandler) GetStores(c *gin.Context) {
	stores, err := h.storeRepo.GetAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch stores", err)
		return
	}

	response.Success(c, "Stores retrieved successfully", stores)
}

// GetStore returns a single store
func (h *StoreHandler) GetStore(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid store ID", err)
		return
	}

	store, err := h.storeRepo.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Store not found", err)
		return
	}

	response.Success(c, "Store retrieved successfully", store)
}

// CreateStore creates a new store
func (h *StoreHandler) CreateStore(c *gin.Context) {
	var req models.StoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	userID := c.GetInt64("user_id")

	store := &models.Store{
		Code:      req.Code,
		Name:      req.Name,
		Address:   req.Address,
		Phone:     req.Phone,
		Status:    req.Status,
		CreatedBy: userID,
		UpdatedBy: userID,
	}

	if store.Status == "" {
		store.Status = "ACTIVE"
	}

	err := h.storeRepo.Create(store)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create store", err)
		return
	}

	response.Success(c, "Store created successfully", store)
}

// UpdateStore updates an existing store
func (h *StoreHandler) UpdateStore(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid store ID", err)
		return
	}

	var req models.StoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	userID := c.GetInt64("user_id")

	store := &models.Store{
		ID:        id,
		Code:      req.Code,
		Name:      req.Name,
		Address:   req.Address,
		Phone:     req.Phone,
		Status:    req.Status,
		UpdatedBy: userID,
	}

	err = h.storeRepo.Update(store)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update store", err)
		return
	}

	response.Success(c, "Store updated successfully", store)
}

// DeleteStore soft deletes a store
func (h *StoreHandler) DeleteStore(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid store ID", err)
		return
	}

	err = h.storeRepo.Delete(id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete store", err)
		return
	}

	response.Success(c, "Store deleted successfully", nil)
}
