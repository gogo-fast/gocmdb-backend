package utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

var NoneUserId = ""
var UnknownUserName = ""
var InvalidUid = -1

func GetUserIdFromUrl(c *gin.Context) int {
	userIdStr := c.DefaultQuery("userId", NoneUserId)
	if userIdStr == NoneUserId {
		Logger.Info(" userId needed within url")
		BadResponse(c, "用户ID不能为空")
		return InvalidUid
	}
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		Logger.Error(err)
		BadResponse(c, "非法用户ID")
		return InvalidUid
	}
	return userId
}
