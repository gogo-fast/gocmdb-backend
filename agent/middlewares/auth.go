package middlewares

import (
	"github.com/gin-gonic/gin"
	"gogo-cmdb/agent/utils"
)

func AuthToken(c *gin.Context)  {
	token := c.PostForm("token")
	if token != utils.GlobalConfig.Token {
		BadResponse(c, "Invalid token")
		c.Abort()
		return
	}
	c.Next()
}

