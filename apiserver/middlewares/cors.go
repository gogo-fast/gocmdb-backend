package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gogo-cmdb/apiserver/utils"
	"runtime"
	"strings"
	"time"
)

func AllowCors() gin.HandlerFunc {

	var sySep string
	switch runtime.GOOS {
	case "linux":
		sySep = "\n"
	case "windows":
		sySep = "\r\n"
	default:
		sySep = "\n"
	}


	allowOrigins := strings.Split(strings.TrimSpace(utils.GlobalConfig.Section("cors").Key("allow_origins").String()), sySep)
	allowMethods := strings.Split(strings.TrimSpace(utils.GlobalConfig.Section("cors").Key("allow_methods").String()), sySep)
	allowHeaders := strings.Split(strings.TrimSpace(utils.GlobalConfig.Section("cors").Key("allow_headers").String()), sySep)
	exposeHeaders := strings.Split(strings.TrimSpace(utils.GlobalConfig.Section("cors").Key("allow_origins").String()), sySep)
	allowCredentials := utils.GlobalConfig.Section("cors").Key("allow_credentials").MustBool(false)
	maxAge :=utils.GlobalConfig.Section("cors").Key("max_age").MustInt(24)

	//fmt.Println(allowOrigins,allowMethods,allowHeaders,exposeHeaders,allowCredentials,maxAge)

	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     allowMethods,
		AllowHeaders:     allowHeaders,
		ExposeHeaders:    exposeHeaders,
		AllowCredentials: allowCredentials,
		MaxAge: time.Duration(maxAge) * time.Hour,
	})
}
