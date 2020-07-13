package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var (
	GlobalConfig *Config
	v            *viper.Viper
	BaseDir      string
)

type Config struct {
	ApiServerPort     int
	HeartBeatInterval int64
	RegisterInterval  int64
	ApiServerUrl      string
	Token             string
	ApiServerHost     string
	LogDirName        string
	LogFileName       string
	UuidFileName      string
	PidFileName       string
	PidFilePath       string
	RunEnv            string
	UUID              string
	PID               int
}

func NewConfig() *Config {
	// return a new config with default value above
	return &Config{
		ApiServerHost:     v.GetString("apiserver.host"),
		ApiServerPort:     v.GetInt("apiserver.port"),
		ApiServerUrl:      v.GetString("url"),
		Token:             v.GetString("token"),
		LogDirName:        v.GetString("log.log_dir_name"),
		LogFileName:       v.GetString("log.log_file_name"),
		UuidFileName:      v.GetString("uuid_file"),
		PidFileName:       v.GetString("pid_file"),
		RunEnv:            v.GetString("env"),
		HeartBeatInterval: v.GetInt64("heartbeat.interval"),
		RegisterInterval:  v.GetInt64("register.interval"),
	}
}

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

	v = viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(filepath.Join(BaseDir, "conf"))

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("Agent config file not found")
		} else {
			// Config file was found but another error was produced
			fmt.Println("other errors", err)
		}
		fmt.Println("Load Agent Config Failed:", err)
		os.Exit(1)
	}
	// set default config
	v.SetDefault("log.log_dir_name", "logs")
	v.SetDefault("log.log_file_name", "agent.log")
	v.SetDefault("uuid_file", "agent.uuid")
	v.SetDefault("pid_file", "agent.pid")
	v.SetDefault("env", "dev")
	v.SetDefault("heartbeat.interval", 10)
	v.SetDefault("register.interval", 10)

	GlobalConfig = NewConfig()
	GlobalConfig.UUID = GetUuid()
	PID := os.Getpid()
	GlobalConfig.PID = PID
	RecordPid(PID)
}
