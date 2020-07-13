package hosts

import (
	"gogo-cmdb/apiserver/models"
)

var HeartbeatObj *HeartBeat

type HeartBeat struct {
	HeartBeatChan chan *models.HeartBeatMsg
}

func NewHeartBeat() *HeartBeat {
	return &HeartBeat{
		HeartBeatChan: make(chan *models.HeartBeatMsg, 1000),
	}
}

func (h *HeartBeat) Run() {
	go func() {
		for {
			hbMsg := <-h.HeartBeatChan
			DefaultHostManager.HearBeat(hbMsg)
		}
	}()

}

func init() {
	HeartbeatObj = NewHeartBeat()
	HeartbeatObj.Run()
}
