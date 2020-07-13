package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var (
	GlobalConfig *viper.Viper
	BaseDir string
)

func init() {
	absDir, err := filepath.Abs(os.Args[0])
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	// you should use the following base dir while binary was build into src/bin/
	BaseDir = filepath.Dir(filepath.Dir(absDir))
	// you should use the following base dir while binary was build into src/
	//BaseDir = filepath.Dir(absDir)

	GlobalConfig = viper.New()

	GlobalConfig.SetConfigName("config")
	GlobalConfig.SetConfigType("yaml")
	GlobalConfig.AddConfigPath(filepath.Join(BaseDir, "conf"))

	if err := GlobalConfig.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("api server config file not found")
		} else {
			// Config file was found but another error was produced
			fmt.Println("other errors", err)
		}
		fmt.Println("Load Api Server Config Failed:", err)
		os.Exit(1)
	}

	GlobalConfig.SetDefault("cors.allow_credentials", false)
	GlobalConfig.SetDefault("cors.max_age", 24)
	GlobalConfig.SetDefault("mysql.max_conn", 10)
	GlobalConfig.SetDefault("mysql.max_idle", 5)
	GlobalConfig.SetDefault("filesystem.max_multipart_memory", 20)
	GlobalConfig.SetDefault("server.default_page_num", 1)
	GlobalConfig.SetDefault("server.default_page_size", 5)
	GlobalConfig.SetDefault("heartbeat.interval", 10)
	GlobalConfig.SetDefault("server.websocket_update_interval", 5)

	fmt.Println("### Load Api Server Config Success")
}
