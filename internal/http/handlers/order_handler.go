package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"CLOB/internal/http/utils"
	"CLOB/internal/services"
	"CLOB/models"
)

type OrderHandler struct {
	db           *gorm.DB
	orderService *services.OrderService
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
	engine := services.NewMatchingEngine(db)
	orderService := services.NewOrderService(db, engine)

	return &OrderHandler{
		db:           db,
		orderService: orderService,
	}
}

type CreateOrderRequest struct {
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

	accountID, ok := utils.GetAccountIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ctx := c.Request.Context()

	order, match, err := h.orderService.CreateOrder(
		ctx,
		accountID,
		req.Instrument,
		req.Side,
		req.Price,
		req.Quantity,
	)
	if err != nil {
		msg := err.Error()

		if strings.Contains(msg, "side must be BUY or SELL") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "side must be BUY or SELL"})
			return
		}
		if strings.Contains(msg, "insufficient balance") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient balance"})
			return
		}
		if strings.Contains(msg, "invalid instrument") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid instrument"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"order":         order,
		"matched_order": match,
	})
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

	ctx := c.Request.Context()

	if err := h.orderService.CancelOrder(ctx, accountID, orderID); err != nil {
		msg := err.Error()

		switch {
		case strings.Contains(msg, "order not found"):
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		case strings.Contains(msg, "order already has executed trades"):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "order already has executed trades and cannot be canceled",
			})
			return
		case strings.Contains(msg, "locked balance too low for cancel"):
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "inconsistent locked balance for this order",
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cancel order"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "order canceled",
		"id":      orderID,
	})
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	_, ok := utils.GetAccountIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	side := c.Query("side") // "BUY" ou "SELL"

	query := h.db.Model(&models.Order{})

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
