package hosts

import (
	"gogo-cmdb/commons"
)

var HeartbeatObj *HeartBeat

type HeartBeat struct {
	HeartBeatChan chan *commons.HeartBeatMsg
}

func NewHeartBeat() *HeartBeat {
	return &HeartBeat{
		HeartBeatChan: make(chan *commons.HeartBeatMsg, 1000),
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
