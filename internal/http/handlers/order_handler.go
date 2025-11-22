package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"CLOB/models"
)

type OrderHandler struct {
	db *gorm.DB
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{db: db}
}

type CreateOrderRequest struct {
	AccountID  uint        `json:"account_id" binding:"required"`
	Instrument string      `json:"instrument" binding:"required"`
	Side       models.Side `json:"side" binding:"required"`
	Price      float64     `json:"price" binding:"required"`
	Quantity   float64     `json:"quantity" binding:"required"`
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Side != models.SideBuy && req.Side != models.SideSell {
		c.JSON(http.StatusBadRequest, gin.H{"error": "side must be BUY or SELL"})
		return
	}

	order := models.Order{
		AccountID:  req.AccountID,
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
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var order models.Order
	if err := h.db.First(&order, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch order"})
		return
	}

	if err := h.db.Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cancel order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order canceled"})
}
