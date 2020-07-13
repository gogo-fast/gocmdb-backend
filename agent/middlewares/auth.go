package middlewares

import (
	"agent/utils"
	"github.com/gin-gonic/gin"
)

func AuthToken(c *gin.Context) {
	token := c.PostForm("token")
	if token != utils.GlobalConfig.Token {
		BadResponse(c, "Invalid token")
		c.Abort()
		return
	}
	c.Next()
}
