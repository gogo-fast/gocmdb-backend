package main

import (
	"fmt"
	"gogo-cmdb/apiserver/handlers"
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
		fmt.Println("start server failed:", err)
		os.Exit(-1)
	}
}
