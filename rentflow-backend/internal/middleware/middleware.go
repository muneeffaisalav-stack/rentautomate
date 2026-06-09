package middleware

import "github.com/gin-gonic/gin"

// Auth is a placeholder for an authentication middleware
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement authentication logic
		c.Next()
	}
}
