package utils

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
)

var (
	GlobalConfig *ini.File
	BaseDir      string
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
	confPath := filepath.Join(BaseDir, "conf", "config.ini")
	GlobalConfig, err = ini.Load(confPath)

	if err != nil {
		fmt.Println("Load Config Failed:", err)
		os.Exit(1)
	}
	fmt.Println("Load Config Success")
}
