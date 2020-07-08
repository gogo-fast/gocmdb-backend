package main

import (
	"fmt"
	_ "gogo-cmdb/apiserver/cloud/plugins"
	"gogo-cmdb/apiserver/handlers"
	_ "gogo-cmdb/apiserver/hosts"
	"gogo-cmdb/apiserver/utils"
	"os"
)

func main() {

	utils.InitLogger()

	port := utils.GlobalConfig.GetString("server.port")

	addr := fmt.Sprintf(":%s", port)

	utils.Logger.Info(fmt.Sprintf("Server Listening on [%s]", addr))
	err := handlers.Route.Run(addr)
	if err != nil {
		fmt.Println("start api server failed:", err)
		os.Exit(-1)
	}
}
