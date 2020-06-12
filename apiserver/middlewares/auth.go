package middlewares

import (
	"github.com/gin-gonic/gin"
	"gogo-cmdb/apiserver/utils"
)


func GetSuperClaims(c *gin.Context) (*utils.AuthCustomClaims ,bool) {
	claims, ok := c.Get("claims")
	if !ok {
		utils.Logger.Info("上级认证信息缺失")
		utils.BadResponse(c, "上级认证信息缺失")
		return nil, false
	}
	cla, ok := claims.(*utils.AuthCustomClaims)
	if !ok {
		utils.BadResponse(c, "上级认证信息类型不符")
		return nil, false
	}
	return cla, true
}


func BaseAuth(c *gin.Context) {
	tokeStr, err := c.Cookie("authToken")
	if err != nil {
		utils.Logger.Error(err)
		utils.BadResponse(c, "认证信息缺失")
		c.Abort()
		return
	}

	claims := utils.ParseJwtAuthToken(tokeStr)

	if claims == nil {
		utils.Logger.Error(err)
		utils.BadResponse(c, "认证失败")
		c.Abort()
		return
	}

	c.Set("claims", claims)
	c.Next()
}



func AuthAdmin(c *gin.Context) {
	claims, ok := GetSuperClaims(c)
	if !ok {
		c.Abort()
		return
	}

	if claims.UserType != utils.AdminUser {
		utils.BadResponse(c, "只有管理员有此权限")
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

	userId := utils.GetUserIdFromUrl(c)
	if userId == utils.InvalidUid {
		c.Abort()
		return
	}

	if userId != currentUserId {
		utils.BadResponse(c, "只用当前用户有此权限")
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

	userId := utils.GetUserIdFromUrl(c)
	if userId == utils.InvalidUid {
		c.Abort()
		return
	}

	if userId != currentUserId && currentUseType != utils.AdminUser {
		utils.BadResponse(c, "只用当前用户和管理员有此权限")
		c.Abort()
		return
	}

	c.Next()

}
