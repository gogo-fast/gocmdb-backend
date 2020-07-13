package middlewares

import (
	"github.com/gin-gonic/gin"
	"gogo-cmdb/apiserver/models"
	"gogo-cmdb/apiserver/utils"
)

func GetSuperClaims(c *gin.Context) (*models.AuthCustomClaims, bool) {
	claims, ok := c.Get("claims")
	if !ok {
		utils.Logger.Info("上级认证信息缺失")
		BadResponse(c, "上级认证信息缺失")
		return nil, false
	}
	cla, ok := claims.(*models.AuthCustomClaims)
	if !ok {
		BadResponse(c, "上级认证信息类型不符")
		return nil, false
	}
	return cla, true
}

func BaseAuth(c *gin.Context) {
	tokeStr, err := c.Cookie("authToken")
	if err != nil {
		utils.Logger.Error(err)
		BadResponse(c, "认证信息缺失,请重新登录")
		c.Abort()
		return
	}

	claims := models.ParseJwtAuthToken(tokeStr)

	if claims == nil {
		utils.Logger.Error(err)
		BadResponse(c, "认证失败")
		c.Abort()
		return
	}

	c.Set("claims", claims)
	c.Next()
}

func AuthAgent(c *gin.Context) {
	token := c.Query("token")
	if token != utils.GlobalConfig.GetString("agent.token") {
		BadResponse(c, "Invalid token")
		c.Abort()
		return
	}
	c.Next()
}

func AuthAdmin(c *gin.Context) {
	claims, ok := GetSuperClaims(c)
	if !ok {
		c.Abort()
		return
	}

	if claims.UserType != utils.AdminUser {
		BadResponse(c, "只有管理员有此权限")
		c.Abort()
		return
	}

	c.Next()
}

func AuthCurrentUser(c *gin.Context) {
	claims, ok := GetSuperClaims(c)
	if !ok {
		c.Abort()
		return
	}
	currentUserId := claims.UserId
	UserIdFromUrl := utils.GetUserIdFromUrl(c)

	if UserIdFromUrl == utils.InvalidUid {
		BadResponse(c, "用户id非法或缺失")
		c.Abort()
		return
	}
	if UserIdFromUrl != currentUserId {
		BadResponse(c, "只用当前用户有此权限")
		c.Abort()
		return
	}

	c.Next()

}

func AuthCurrentAndAdmin(c *gin.Context) {
	claims, ok := GetSuperClaims(c)
	if !ok {
		c.Abort()
		return
	}

	currentUserId := claims.UserId
	currentUseType := claims.UserType
	UserIdFromUrl := utils.GetUserIdFromUrl(c)

	if UserIdFromUrl == utils.InvalidUid {
		BadResponse(c, "当前用户id非法或缺失")
		c.Abort()
		return
	}

	if UserIdFromUrl != currentUserId && currentUseType != utils.AdminUser {
		BadResponse(c, "只用当前用户和管理员有此权限")
		c.Abort()
		return
	}

	c.Next()

}
