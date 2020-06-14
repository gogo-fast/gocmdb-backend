package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gogo-cmdb/apiserver/middlewares"
	"gogo-cmdb/apiserver/utils"
	"net/http"
)

var Route *gin.Engine

func init() {
	Route = gin.Default()

	mmy := utils.GlobalConfig.GetInt("filesystem.max_multipart_memory")
	Route.MaxMultipartMemory = 8 << mmy

	userStaticPath := utils.GlobalConfig.GetString(("filesystem.user_static_dir_name"))
	userStaticUrl := utils.GlobalConfig.GetString(("filesystem.user_static_url"))

	img := Route.Group("/img", middlewares.BaseAuth)
	{
		img.StaticFS(userStaticUrl, http.Dir(fmt.Sprintf("%s/upload/%s", utils.BaseDir, userStaticPath)))
	}

	host := utils.GlobalConfig.GetString("filesystem.host")
	port := utils.GlobalConfig.GetString("server.port")
	utils.Logger.Info(fmt.Sprintf("File System Served On: http://%s:%s/img%s", host, port, userStaticUrl))

	v1 := Route.Group("/v1", middlewares.AllowCors())
	{
		v1.POST("/login", Login)
		v1.OPTIONS("/login", Login)

		user := v1.Group("/user", middlewares.BaseAuth)
		{
			/*
				由于PUT方法的请求是非简单请求(HEAD,GET,POST之外的方法)，
				发送的时候会先有预检请求(OPTIONS方法)，所以这里要增加OPTIONS方法的路由。
			*/
			user.GET("/list", GetUserList)

			user.GET("", GetUser)
			user.POST("", middlewares.AuthAdmin, CreateUser)
			user.OPTIONS("", middlewares.AuthAdmin, CreateUser)

			user.GET("/detail", middlewares.AuthCurrentAndAdmin, GetUserDetailById)
			user.PUT("/detail", middlewares.AuthCurrentAndAdmin, UpdateDetailById)
			user.OPTIONS("/detail", middlewares.AuthCurrentAndAdmin, UpdateDetailById)

			user.PUT("/password", middlewares.AuthCurrentUser, UpdatePasswordById)
			user.OPTIONS("/password", middlewares.AuthCurrentUser, UpdatePasswordById)

			user.PUT("/status", middlewares.AuthAdmin, UpdateUserStatusById)
			user.OPTIONS("/status", middlewares.AuthAdmin, UpdateUserStatusById)

			user.PUT("/type", middlewares.AuthAdmin, UpdateUserTypeById)
			user.OPTIONS("/type", middlewares.AuthAdmin, UpdateUserTypeById)

			user.POST("/avatar", middlewares.AuthCurrentUser, UploadAvatar)
			user.OPTIONS("/avatar", middlewares.AuthCurrentUser, UploadAvatar)
		}

		cloud := v1.Group("/cloud", middlewares.BaseAuth)
		{
			cloud.GET("/regions", GetRegions)
			cloud.GET("/zones", GetZones)
		}
	}
}
