package contextutils

import "github.com/gin-gonic/gin"

func GetUserID(c *gin.Context) (int64, bool) {
	val, exists := c.Get("userID") // Matches the c.Set("userID", ...) key in your middleware
	if !exists {
		return 0, false
	}
	id, ok := val.(int64)
	return id, ok
}