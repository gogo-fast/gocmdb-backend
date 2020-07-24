package handlers

import (
	"apiserver/hosts"
	"apiserver/models"
	"apiserver/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/imroc/req"
	"strconv"
	"time"
)

func Heartbeat(c *gin.Context) {
	hbMsg := models.HeartBeatMsg{}
	err := c.ShouldBindJSON(&hbMsg)
	if err != nil {
		utils.Logger.Error("invalid heartbeat msg")
		return
	}
	//fmt.Println(hbMsg)
	if hbMsg.UUID == "" || hbMsg.Timestamp.Int64 == 0 {
		BadResponse(c, "invalid heartbeat msg")
		return
	}
	hosts.HeartbeatObj.HeartBeatChan <- &hbMsg
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "heartbeat success",
		"data":   gin.H{},
	})
}

func Register(c *gin.Context) {
	rgMsg := models.RegisterMsg{}
	err := c.ShouldBindJSON(&rgMsg)
	if err != nil {
		fmt.Println("err:", err)
		utils.Logger.Error("invalid register msg")
		return
	}
	//fmt.Println(hbMsg)
	if rgMsg.UUID == "" {
		BadResponse(c, "invalid register msg")
		return
	}
	hosts.RegisterObj.RegisterChan <- &rgMsg
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "register success",
		"data":   gin.H{},
	})
}

func GetHost(c *gin.Context) {

	uuid := c.DefaultQuery("uuid", "")
	if uuid == "" {
		utils.Logger.Info("uuid must be set")
		EmptyHostResponse(c, "主机查询条件不能为空，需要制定uuid")
		return
	}

	host, err := models.DefaultHostManager.GetHostRecordByUUID(uuid)
	if err != nil {
		utils.Logger.Error(err)
		EmptyHostResponse(c, "所查主机不存在")
		return
	}
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "获取主机成功",
		"data":   []*models.Host{host},
	})

}

func GetHostList(c *gin.Context) {

	var (
		err  error
		p, s int
	)

	defaultPageNum := utils.GlobalConfig.GetInt("server.default_page_num")
	defaultPageSize := utils.GlobalConfig.GetInt("server.default_page_size")
	page := c.Query("page")
	size := c.Query("size")
	p, err = strconv.Atoi(page)
	if err != nil {
		p = defaultPageNum
	}
	s, err = strconv.Atoi(size)
	if err != nil {
		s = defaultPageSize
	}

	total, hostList, pagination, err := models.DefaultHostManager.GetHostRecordList(p, s)
	if err != nil {
		utils.Logger.Error(err)
		EmptyHostsResponse(c, "获取主机列表失败")
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "获取主机列表成功",
		"data": gin.H{
			"total":          total,
			"hosts":          hostList,
			"currentPageNum": pagination.CurrentPageNum,
		},
	})

}

func WsGetHost(c *gin.Context) {
	var (
		wsConn *websocket.Conn
		err    error
		data   []byte
	)
	uuid := c.Query("uuid")
	if uuid == "" {
		utils.Logger.Warning("uuid must be set")
		EmptyHostResponse(c, "主机查询条件不能为空，需要制定uuid")
		return
	}

	// get host once before establish websocket connection,
	// if error happened, do not establish websocket connection.
	_, err = models.DefaultHostManager.GetHostRecordByUUID(uuid)
	if err != nil {
		utils.Logger.Error(err)
		EmptyHostResponse(c, "所查主机不存在")
		return
	}

	// Upgrade: websocket
	wsConn, err = upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.Logger.Error(err)
		EmptyHostResponse(c, fmt.Sprintf("获取主机[%s]的websocket连接建立失败", uuid))
		return
	}
	conn := utils.InitConnection(wsConn)
	utils.Logger.Info(fmt.Sprintf("getting host [%s] websocket success", uuid))

	go func(conn *utils.Connection, uuid string) {
		for {
			host, err := models.DefaultHostManager.GetHostRecordByUUID(uuid)
			if err != nil {
				utils.Logger.Error(err)
				//utils.BadResponse(c, "所查主机不存在")
				data, err = json.Marshal(gin.H{
					"status": "error",
					"msg":    "获取主机失败",
					"data":   []*models.Host{},
				})
			} else {
				data, err = json.Marshal(gin.H{
					"status": "ok",
					"msg":    "获取主机成功",
					"data":   []*models.Host{host},
				})
			}
			err = conn.WriteMessage(data)
			if err != nil {
				conn.Close()
				return
			}
			time.Sleep(time.Second * time.Duration(wsUpdateInterval))
		}
	}(conn, uuid)

}

func WsGetHostList(c *gin.Context) {
	var (
		wsConn *websocket.Conn
		err    error
		data   []byte
		p, s   int
	)

	// Upgrade: websocket
	wsConn, err = upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.Logger.Error(err)
		return
	}
	conn := utils.InitConnection(wsConn)

	defaultPageNum := utils.GlobalConfig.GetInt("server.default_page_num")
	defaultPageSize := utils.GlobalConfig.GetInt("server.default_page_size")
	page := c.DefaultQuery("page", "none")
	size := c.DefaultQuery("size", "none")
	p, err = strconv.Atoi(page)
	if err != nil {
		p = defaultPageNum
	}
	s, err = strconv.Atoi(size)
	if err != nil {
		s = defaultPageSize
	}

	go func(conn *utils.Connection, p, s int) {
		for {
			total, hostList, pagination, err := models.DefaultHostManager.GetHostRecordList(p, s)
			if err != nil {
				utils.Logger.Error(err)
				data, err = json.Marshal(gin.H{
					"status": "error",
					"msg":    "获取主机列表失败",
					"data": gin.H{
						"total":          0,
						"hosts":          []*models.Host{},
						"currentPageNum": -1,
					},
				})
			} else {
				data, err = json.Marshal(gin.H{
					"status": "ok",
					"msg":    "获取主机列表成功",
					"data": gin.H{
						"total":          total,
						"hosts":          hostList,
						"currentPageNum": pagination.CurrentPageNum,
					},
				})
			}
			err = conn.WriteMessage(data)
			if err != nil {
				conn.Close()
				return
			}
			time.Sleep(time.Second * time.Duration(wsUpdateInterval))
		}
	}(conn, p, s)

}

func DeleteHost(c *gin.Context) {
	uuid := c.Query("uuid")
	userId := c.Query("userId")
	host, err := models.DefaultHostManager.GetHostRecordByUUID(uuid)
	if err != nil {
		BadResponse(c, "host not exist")
		return
	}
	if host.IsOnline {
		utils.Logger.Error("stop host before delete")
		BadResponse(c, "stop host before delete")
		return
	}
	err = models.DefaultHostManager.DeleteHostRecordByUUID(uuid)
	if err != nil {
		utils.Logger.Error("delete host failed")
		BadResponse(c, "delete host failed")
		return
	}
	hosts.DefaultHostManager.DeleteHost(uuid)
	utils.Logger.Info(fmt.Sprintf("delete host %s, userId %s", uuid, userId))
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "delete host success",
		"data":   gin.H{},
	})
}

func StopHost(c *gin.Context) {
	uuid := c.Query("uuid")
	clusterIp := c.Query("clusterIp")
	tokenStr := utils.GlobalConfig.GetString("agent.token")
	agentPort := utils.GlobalConfig.GetInt("agent.port")
	if clusterIp != "" {
		params := req.Param{
			"token": tokenStr,
		}
		_, err := req.Post(fmt.Sprintf("http://%s:%d/host/stop", clusterIp, agentPort), params)
		if err != nil {
			utils.Logger.Error("connect to agent failed")
			c.JSON(200, gin.H{
				"status": "error",
				"msg":    fmt.Sprintf("connect to agent %s failed", uuid),
				"data":   gin.H{},
			})
		} else {
			utils.Logger.Info(fmt.Sprintf("stop host %s", uuid))
			c.JSON(200, gin.H{
				"status": "ok",
				"msg":    fmt.Sprintf("stopping host %s", uuid),
				"data":   gin.H{},
			})
		}
	} else {
		c.JSON(200, gin.H{
			"status": "error",
			"msg":    "invalid params",
			"data":   gin.H{},
		})
	}
}
