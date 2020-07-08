package middlewares

import "github.com/gin-gonic/gin"

func BadResponse(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"status": "error",
		"msg":    msg,
		"data":   gin.H{},
	})
}
