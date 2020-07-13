package main

import (
	_ "apiserver/cloud/plugins"
	"apiserver/handlers"
	_ "apiserver/hosts"
	"apiserver/utils"
	"fmt"
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
