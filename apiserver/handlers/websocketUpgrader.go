package handlers

import (
	"github.com/gorilla/websocket"
	"gogo-cmdb/apiserver/utils"
	"net/http"
)

var (
	upGrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	wsUpdateInterval = utils.GlobalConfig.GetInt64("server.websocket_update_interval")
)
