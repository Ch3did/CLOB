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

	// Auth
	r.POST("/account/register", accountHandler.CreateAccount)
	r.POST("/account/login", accountHandler.Login)

	authorized := r.Group("/")
	authorized.Use(middleware.AuthRequired())

	// Order
	authorized.POST("/orders", orderHandler.CreateOrder)
	authorized.DELETE("/orders/:id", orderHandler.CancelOrder)

	//Balace
	authorized.POST("/balances/credit", balanceHandler.Credit)
	authorized.POST("/balances/debit", balanceHandler.Debit)
	authorized.GET("/balances", balanceHandler.List)
	authorized.GET("/balances/:id", balanceHandler.GetOne)

	return r
}
