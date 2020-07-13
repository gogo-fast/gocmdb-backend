package utils

import (
	"fmt"
	"os"
	"strconv"
)

func RecordPid(pid int) {
	pidFilePath := fmt.Sprintf("%s/%s", BaseDir, GlobalConfig.PidFileName)
	fPid, err := os.OpenFile(pidFilePath, os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		Logger.Error("create pid file failed")
		return
	}
	_, err = fPid.WriteString(strconv.Itoa(pid))
	if err != nil {
		Logger.Error("write to pid file failed")
		return
	}
	defer fPid.Close()
	GlobalConfig.PidFilePath = pidFilePath
}
