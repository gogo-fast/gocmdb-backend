package handlers

import (
	"apiserver/utils"
	"github.com/gorilla/websocket"
	"net/http"
)

var (
	upGrader websocket.Upgrader

	wsUpdateInterval int64
)

func InitWs() {
	upGrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	wsUpdateInterval = utils.GlobalConfig.GetInt64("server.websocket_update_interval")
}
