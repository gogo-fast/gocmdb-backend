package handlers

import (
	"agent/middlewares"
	"github.com/gin-gonic/gin"
)

var Route *gin.Engine

func init() {
	Route = gin.Default()

	agent := Route.Group("/host", middlewares.AuthToken)
	{
		agent.POST("/stop", StopHost)
	}
}
