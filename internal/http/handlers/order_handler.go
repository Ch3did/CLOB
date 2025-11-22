package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"CLOB/internal/http/utils"
	"CLOB/models"
)

type OrderHandler struct {
	db *gorm.DB
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{db: db}
}

type CreateOrderRequest struct {
	Instrument string      `json:"instrument" binding:"required"`
	Side       models.Side `json:"side" binding:"required"`
	Price      float64     `json:"price" binding:"required"`
	Quantity   float64     `json:"quantity" binding:"required"`
}

type CancelOrderRequest struct {
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accountID, exists := c.Get("account_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if req.Side != models.SideBuy && req.Side != models.SideSell {
		c.JSON(http.StatusBadRequest, gin.H{"error": "side must be BUY or SELL"})
		return
	}

	order := models.Order{
		AccountID:  accountID.(uint),
		Instrument: req.Instrument,
		Side:       req.Side,
		Price:      req.Price,
		Quantity:   req.Quantity,
	}

	if err := h.db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	// TODO:
	// aqui depois entra a lógica do matching engine

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	accountID, ok := utils.GetAccountIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idParam := c.Param("id")
	idUint64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	orderID := uint(idUint64)

	var order models.Order
	if err := h.db.
		Where("id = ? AND account_id = ?", orderID, accountID).
		First(&order).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch order"})
		return
	}

	var tradeCount int64
	if err := h.db.
		Model(&models.Trade{}).
		Where("buy_order_id = ? OR sell_order_id = ?", order.ID, order.ID).
		Count(&tradeCount).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check trades for order"})
		return
	}

	if tradeCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "order already has executed trades and cannot be canceled",
		})
		return
	}

	if err := h.db.Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cancel order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "order canceled",
		"id":      orderID,
	})
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	accountID, ok := utils.GetAccountIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	side := c.Query("side") // "BUY" ou "SELL"

	query := h.db.Model(&models.Order{}).
		Where("account_id = ?", accountID)

	if side != "" {
		query = query.Where("side = ?", side)
	}

	var orders []models.Order
	if err := query.Order("created_at DESC").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	accountID, ok := utils.GetAccountIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idParam := c.Param("id")
	idUint64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}
	orderID := uint(idUint64)

	var order models.Order
	if err := h.db.
		Where("id = ? AND account_id = ?", orderID, accountID).
		First(&order).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch order"})
		return
	}

	c.JSON(http.StatusOK, order)
}
