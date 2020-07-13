package middlewares

import (
	"apiserver/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func AllowCors() gin.HandlerFunc {

	allowOrigins := utils.GlobalConfig.GetStringSlice("cors.allow_origins")
	allowMethods := utils.GlobalConfig.GetStringSlice("cors.allow_methods")
	allowHeaders := utils.GlobalConfig.GetStringSlice("cors.allow_headers")
	exposeHeaders := utils.GlobalConfig.GetStringSlice("cors.expose_headers")
	allowCredentials := utils.GlobalConfig.GetBool("cors.allow_credentials")
	maxAge := utils.GlobalConfig.GetInt("cors.max_age")

	//fmt.Println(allowOrigins,allowMethods,allowHeaders,exposeHeaders,allowCredentials,maxAge)

	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     allowMethods,
		AllowHeaders:     allowHeaders,
		ExposeHeaders:    exposeHeaders,
		AllowCredentials: allowCredentials,
		MaxAge:           time.Duration(maxAge) * time.Hour,
	})
}
