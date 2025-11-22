package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"CLOB/internal/http/utils"
	"CLOB/internal/services"
)

type BalanceHandler struct {
	svc *services.BalanceService
}

func NewBalanceHandler(db *gorm.DB) *BalanceHandler {
	return &BalanceHandler{
		svc: services.NewBalanceService(db),
	}
}

type balanceChangeRequest struct {
	Asset  string  `json:"asset" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}

func (h *BalanceHandler) Credit(c *gin.Context) {
	accountID, ok := utils.GetAccountIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req balanceChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.Credit(c.Request.Context(), accountID, req.Asset, req.Amount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *BalanceHandler) Debit(c *gin.Context) {
	accountID, ok := utils.GetAccountIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req balanceChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.Debit(c.Request.Context(), accountID, req.Asset, req.Amount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *BalanceHandler) List(c *gin.Context) {
	accountID, ok := utils.GetAccountIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	balances, err := h.svc.ListByAccount(c.Request.Context(), accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, balances)
}

func (h *BalanceHandler) GetOne(c *gin.Context) {
	accountID, ok := utils.GetAccountIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	asset := c.Param("asset")

	bal, err := h.svc.GetOne(c.Request.Context(), accountID, asset)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bal)
}
