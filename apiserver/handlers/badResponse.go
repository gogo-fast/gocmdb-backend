package handlers

import (
	"github.com/gin-gonic/gin"
	"gogo-cmdb/apiserver/cloud"
	"gogo-cmdb/apiserver/models"
)

func BadResponse(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"status": "error",
		"msg":    msg,
		"data":   gin.H{},
	})
}

func EmptyUsersResponse(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"status": "error",
		"msg":    msg,
		"data": gin.H{
			"total":          0,
			"users":          []*models.User{},
			"currentPageNum": -1,
		},
	})
}

func EmptyUserResponse(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"status": "error",
		"msg":    msg,
		"data":   []*models.User{},
	})
}

func EmptyUserDetailsResponse(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"status": "error",
		"msg":    msg,
		"data":   &models.UserFullInfo{},
	})
}

func EmptyHostResponse(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"status": "error",
		"msg":    msg,
		"data":   []*models.Host{},
	})
}

func EmptyHostsResponse(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"status": "error",
		"msg":    msg,
		"data": gin.H{
			"total":          0,
			"hosts":          []*models.Host{},
			"currentPageNum": -1,
		},
	})
}

func EmptyRegionsResponse(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"status": "error",
		"msg":    msg,
		"data":   []*cloud.Region{},
	})
}

func EmptyInstancesStatusResponse(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"status": "error",
		"msg":    msg,
		"data": gin.H{
			"total":           0,
			"instancesStatus": []*cloud.InstanceStaus{},
			"currentPageNum":  -1,
		},
	})
}

func EmptyInstancesResponse(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"status": "error",
		"msg":    msg,
		"data": gin.H{
			"total":          0,
			"instances":      []*cloud.Instance{},
			"currentPageNum": -1,
		},
	})
}

func EmptyInstanceResponse(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"status": "error",
		"msg":    msg,
		"data":   []*cloud.Instance{},
	})
}

func EmptySecurityGroupsResponse(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"status": "error",
		"msg":    msg,
		"data":   []*cloud.SecurityGroup{},
	})
}

func EmptySecurityZonesResponse(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"status": "error",
		"msg":    msg,
		"data":   []*cloud.Zone{},
	})
}
