package main

import (
	"agent/handlers"
	"agent/plugins/heartbeat"
	"agent/plugins/register"
	"agent/utils"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	utils.InitLogger()

	utils.Logger.WithFields(
		logrus.Fields{
			"PID":  utils.GlobalConfig.PID,
			"UUID": utils.GlobalConfig.UUID,
		}).Info("Agent Started")
	go heartbeat.Run()
	go register.Run()

	quitChan := make(chan os.Signal, 1)
	go func(<-chan os.Signal) {
		signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
		<-quitChan
		utils.Logger.WithFields(
			logrus.Fields{
				"PID":  utils.GlobalConfig.PID,
				"UUID": utils.GlobalConfig.UUID,
			}).Info("Agent Stopped")

		os.Remove(utils.GlobalConfig.PidFilePath)
		os.Exit(1)
	}(quitChan)

	utils.Logger.Info(fmt.Sprintf("agent http server running at [%d]", utils.GlobalConfig.AgentHttpServerPort))
	err := handlers.Route.Run(fmt.Sprintf(":%d", utils.GlobalConfig.AgentHttpServerPort))
	if err != nil {
		fmt.Println("start agent httpServer failed:", err)
		os.Exit(-1)
	}
}
