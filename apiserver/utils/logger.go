package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

// Create a new instance of the logger. You can have any number of instances.
var Logger = logrus.New()

func InitLogger() {
	Logger = logrus.New()
	// The API for setting attributes is a little different than the package level
	// exported logger. See Godoc.

	env := GlobalConfig.GetString("env.env")
	logDirName := GlobalConfig.GetString("log.log_dir_name")
	logFileName := GlobalConfig.GetString("log.log_file_name")

	Logger.Info(fmt.Sprintf("Env: [%s]", env))
	switch env {
	case "dev":
		Logger.Out = os.Stdout
		Logger.SetLevel(logrus.InfoLevel)
		Logger.SetReportCaller(true)
		//Logger.Formatter = &logrus.JSONFormatter{}
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = Logger.Out
	case "product":
		logDir := filepath.Join(BaseDir, logDirName)
		_, err := os.Stat(logDir)
		if err != nil {
			exist := os.IsExist(err)
			if !exist {
				err = os.MkdirAll(logDir, 0755)
				if err != nil {
					fmt.Println("create log dir failed")
					os.Exit(1)
				}
			}
		}

		logPath := filepath.Join(logDir, logFileName)

		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err == nil {
			Logger.Out = file
			Logger.SetLevel(logrus.InfoLevel)
			//Logger.Formatter = &logrus.JSONFormatter{}
			gin.SetMode(gin.ReleaseMode)
			gin.DefaultWriter = Logger.Out
		} else {
			Logger.Info("Failed to log to file, using default stderr")
		}
	default:
		Logger.Out = os.Stdout
		Logger.SetLevel(logrus.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = Logger.Out
	}

}
