package middlewares

import (
	"github.com/gin-gonic/gin"
)

func TokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.Request.Header["Authorization"]

		if len(token) == 0 {
			c.AbortWithStatusJSON(401, gin.H{"message": "Authorization header missing"})
		}
		c.Next()
	}
}
