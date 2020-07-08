package utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

var InvalidUid = -1


func GetUserIdFromUrl(c *gin.Context) int {
	userIdStr := c.Query("userId")
	if userIdStr == "" {
		Logger.Info("lose userId in url")
		return InvalidUid
	}
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		Logger.Error(err)
		return InvalidUid
	}
	return userId
}