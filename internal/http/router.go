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

	r.POST("/account/register", accountHandler.CreateAccount)
	r.POST("/account/login", accountHandler.Login)

	authorized := r.Group("/")
	authorized.Use(middleware.AuthRequired())

	authorized.POST("/orders", orderHandler.CreateOrder)
	authorized.POST("/orders/:id/cancel", orderHandler.CancelOrder)

	return r
}
