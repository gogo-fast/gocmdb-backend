package middlewares

import (
	"github.com/gin-gonic/gin"
	"gogo-cmdb/apiserver/badresponse"
	"gogo-cmdb/apiserver/utils"
)

func AuthToken(c *gin.Context)  {
	token := c.PostForm("token")
	if token != utils.GlobalConfig.GetString("token") {
		badresponse.BadResponse(c, "Invalid token")
		c.Abort()
		return
	}
	c.Next()
}

