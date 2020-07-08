package utils

import (
	"fmt"
	"github.com/google/uuid"
	"gogo-cmdb/apiserver/utils"
	"io/ioutil"
	"os"
)

func GetUuid() (uuidStr string) {
	uuidFile := utils.GlobalConfig.GetString("uuid_file")
	uuidFilePath := fmt.Sprintf("%s/%s", utils.BaseDir, uuidFile)
	fUuid, err := os.Open(uuidFilePath)
	if err != nil {
		if notExist := os.IsNotExist(err); notExist {
			fUuid, err = os.Create(uuidFilePath)
			if err != nil {
				utils.Logger.Error("create uuid file failed")
				os.Exit(1)
			}
			defer fUuid.Close()
			uuidStr = uuid.New().String()
			_, err = fUuid.WriteString(uuidStr)
			if err != nil {
				utils.Logger.Error("write to uuid file failed")
				os.Exit(1)
			}
		}
	} else {
		defer fUuid.Close()
		b, err := ioutil.ReadAll(fUuid)
		if err != nil {
			utils.Logger.Error("read from uuid file failed")
			os.Exit(1)
		}
		uuidStr = string(b)
	}
	return
}
