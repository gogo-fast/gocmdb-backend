package handlers

import (
	"github.com/gin-gonic/gin"
	"gogo-cmdb/apiserver/utils"
)

func GetRegions(c *gin.Context) {



}

func GetZones(c *gin.Context) {
	regionId := c.DefaultQuery("regionId", "")
	if regionId == "" {
		utils.BadResponse(c, "Need regionId")
		return
	}

}



