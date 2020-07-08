package handlers

import (
	"github.com/gin-gonic/gin"
	"gogo-cmdb/agent/middlewares"
)

var Route *gin.Engine

func init() {
	Route = gin.Default()

	agent := Route.Group("/host", middlewares.AuthToken)
	{
		agent.POST("/stop", StopHost)
	}
}
