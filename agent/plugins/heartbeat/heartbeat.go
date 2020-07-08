package heartbeat

import (
	"fmt"
	"github.com/imroc/req"
	"gogo-cmdb/agent/utils"
	"gogo-cmdb/commons"
	"time"
)

func Run() {
	url := fmt.Sprintf("%s/%s", utils.GlobalConfig.GetString("url"), "heartbeat")
	tokenStr := utils.GlobalConfig.GetString("token")
	interval := utils.GlobalConfig.GetInt64("heartbeat.interval")

	heartbeat := commons.HeartBeatMsg{}
	for {
		t := time.Now().Unix()
		params := req.Param{"token": tokenStr}
		heartbeat.UUID = utils.GetUuid()
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
