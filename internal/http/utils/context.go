package utils

import "github.com/gin-gonic/gin"

func GetAccountIDFromContext(c *gin.Context) (uint, bool) {
	v, ok := c.Get("account_id")
	if !ok {
		return 0, false
	}
	id, ok := v.(uint)
	return id, ok
}
