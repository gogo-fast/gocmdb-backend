package utils

import (
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
)

func GetUuid() (uuidStr string) {
	uuidFilePath := fmt.Sprintf("%s/%s", BaseDir, GlobalConfig.UuidFileName)
	fUuid, err := os.Open(uuidFilePath)
	if err != nil {
		if notExist := os.IsNotExist(err); notExist {
			fUuid, err = os.Create(uuidFilePath)
			if err != nil {
				Logger.Error("create uuid file failed")
				os.Exit(1)
			}
			defer fUuid.Close()
			uuidStr = uuid.New().String()
			_, err = fUuid.WriteString(uuidStr)
			if err != nil {
				Logger.Error("write to uuid file failed")
				os.Exit(1)
			}
		}
	} else {
		defer fUuid.Close()
		b, err := ioutil.ReadAll(fUuid)
		if err != nil {
			Logger.Error("read from uuid file failed")
			os.Exit(1)
		}
		uuidStr = string(b)
	}
	return
}
