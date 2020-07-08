package handlers

import (
	"github.com/gin-gonic/gin"
)

func StopHost(c *gin.Context) {
	//win.ExitWindowsEx(win.EWX_SHUTDOWN, 0)
	c.JSON(200, gin.H{})
}
