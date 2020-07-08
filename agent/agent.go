package main

import (
	"fmt"
	"gogo-cmdb/agent/handlers"
	"gogo-cmdb/agent/plugins/heartbeat"
	"gogo-cmdb/agent/plugins/register"
	"gogo-cmdb/agent/utils"
	"os"
)

func main() {
	utils.InitLogger()
	go heartbeat.Run()
	go register.Run()
	err := handlers.Route.Run(":8010")
	if err != nil {
		fmt.Println("start api server failed:", err)
		os.Exit(-1)
	}
}
