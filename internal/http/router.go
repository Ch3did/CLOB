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

	authHandler := handlers.NewAuthHandler(db)
	accountHandler := handlers.NewAccountHandler(db)
	orderHandler := handlers.NewOrderHandler(db)

	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)

	authorized := r.Group("/")
	authorized.Use(middleware.AuthRequired())

	authorized.POST("/accounts", accountHandler.CreateAccount)
	authorized.GET("/accounts", accountHandler.ListAccounts)

	authorized.POST("/orders", orderHandler.CreateOrder)
	authorized.POST("/orders/:id/cancel", orderHandler.CancelOrder)

	return r
}
