package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"CLOB/internal/http/handlers"
	"CLOB/internal/http/middleware"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	accountHandler := handlers.NewAccountHandler(db)
	orderHandler := handlers.NewOrderHandler(db)
	balanceHandler := handlers.NewBalanceHandler(db)
	tradeHandler := handlers.NewTradeHandler(db)

	// Auth
	r.POST("/account/register", accountHandler.CreateAccount)
	r.POST("/account/login", accountHandler.Login)

	authorized := r.Group("/")
	authorized.Use(middleware.AuthRequired())

	// Order
	authorized.POST("/orders", orderHandler.CreateOrder)
	authorized.DELETE("/orders/:id", orderHandler.CancelOrder)
	authorized.GET("/orders", orderHandler.ListOrders)
	authorized.GET("/orders/:id", orderHandler.GetOrder)

	//Balace
	authorized.POST("/balance/credit", balanceHandler.Credit)
	authorized.POST("/balance/debit", balanceHandler.Debit)
	authorized.GET("/balance", balanceHandler.List)
	authorized.GET("/balance/:id", balanceHandler.GetOne)

	//Trade
	authorized.GET("/trades", tradeHandler.ListTrades)
	authorized.GET("/trades/:id", tradeHandler.GetTrade)
	return r
}
