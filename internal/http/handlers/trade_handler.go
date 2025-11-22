package handlers

import (
	"CLOB/internal/http/utils"
	"CLOB/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TradeHandler struct {
	db *gorm.DB
}

func NewTradeHandler(db *gorm.DB) *TradeHandler {
	return &TradeHandler{db: db}
}

func (h *TradeHandler) ListTrades(c *gin.Context) {
	accountID, ok := utils.GetAccountIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var trades []models.Trade

	err := h.db.
		Joins("JOIN orders ob ON ob.id = trades.buy_order_id").
		Joins("JOIN orders os ON os.id = trades.sell_order_id").
		Where("ob.account_id = ? OR os.account_id = ?", accountID, accountID).
		Order("trades.created_at DESC").
		Find(&trades).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list trades"})
		return
	}

	c.JSON(http.StatusOK, trades)
}

func (h *TradeHandler) GetTrade(c *gin.Context) {
	accountID, ok := utils.GetAccountIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")

	var trade models.Trade

	err := h.db.
		Joins("JOIN orders ob ON ob.id = trades.buy_order_id").
		Joins("JOIN orders os ON os.id = trades.sell_order_id").
		Where("trades.id = ?", id).
		Where("ob.account_id = ? OR os.account_id = ?", accountID, accountID).
		First(&trade).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "trade not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch trade"})
		return
	}

	c.JSON(http.StatusOK, trade)
}
