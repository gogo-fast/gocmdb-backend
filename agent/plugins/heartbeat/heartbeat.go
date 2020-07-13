package heartbeat

import (
	"agent/models"
	"agent/utils"
	"fmt"
	"github.com/imroc/req"

	"time"
)

func Run() {
	url := fmt.Sprintf("%s/%s", utils.GlobalConfig.ApiServerUrl, "heartbeat")
	tokenStr := utils.GlobalConfig.Token
	interval := utils.GlobalConfig.HeartBeatInterval

	heartbeat := models.HeartBeatMsg{}
	for {
		t := time.Now().Unix()
		params := req.Param{"token": tokenStr}
		heartbeat.UUID = utils.GlobalConfig.UUID
		heartbeat.Timestamp.Int64 = t
		heartbeat.Timestamp.Valid = true
		resp, err := req.Post(url, params, req.BodyJSON(&heartbeat))
		if err != nil {
			utils.Logger.Error("send heartbeat failed")
		} else {
			result := map[string]interface{}{}
			resp.ToJSON(&result)
			utils.Logger.Info(result["msg"])
		}
		time.Sleep(time.Second * time.Duration(interval))
	}
}
