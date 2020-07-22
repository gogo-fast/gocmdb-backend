package handlers

import (
	"apiserver/cloud"
	"apiserver/forms"
	"apiserver/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"strconv"
	"strings"
	"time"
)

func GetRegions(c *gin.Context) {
	platType := c.Query("platType")
	//fmt.Printf("plat form type: %s\n", platType)
	if platType == "" {
		EmptyRegionsResponse(c, "need platType")
		return
	}
	cCloud, ok := cloud.DefaultCloudMgr.GetCloud(platType)
	if !ok {
		EmptyRegionsResponse(c, "get cloud platform failed")
		return
	}
	regions, err := cCloud.GetRegions()
	if err != nil {
		EmptyRegionsResponse(c, fmt.Sprintf("get regions of [%s] failed", cCloud.Name()))
		return
	}
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    fmt.Sprintf("connect to %s success", cCloud.Name()),
		"data":   regions,
	})
}

func GetZones(c *gin.Context) {
	platType := c.Query("platType")
	regionId := c.Query("regionId")
	if platType == "" {
		EmptyZonesResponse(c, "need platType")
		return
	}
	if regionId == "" {
		EmptyZonesResponse(c, "need regionId")
		return
	}

	cCloud, ok := cloud.DefaultCloudMgr.GetCloud(platType)
	if !ok {
		utils.Logger.Warning(fmt.Sprintf("get %s cloud manager failed", platType))
		EmptyZonesResponse(c, fmt.Sprintf("do not support cloud %s", platType))
		return
	}
	zones, err := cCloud.GetZones(regionId)
	if err != nil {
		EmptyZonesResponse(c, fmt.Sprintf("get zones from %s failed", cCloud.PlatType()))
		return
	}
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "get zones success",
		"data":   zones,
	})

}

func GetSecurityGroups(c *gin.Context) {
	platType := c.Query("platType")
	regionId := c.Query("regionId")
	if platType == "" {
		EmptySecurityGroupsResponse(c, "need platType")
		return
	}
	if regionId == "" {
		EmptySecurityGroupsResponse(c, "need regionId")
		return
	}

	cCloud, ok := cloud.DefaultCloudMgr.GetCloud(platType)
	if !ok {
		utils.Logger.Warning(fmt.Sprintf("get %s cloud manager failed", platType))
		EmptySecurityGroupsResponse(c, fmt.Sprintf("do not support cloud %s", platType))
		return
	}
	sgs, err := cCloud.GetSecurityGroups(regionId)
	if err != nil {
		EmptySecurityGroupsResponse(c, fmt.Sprintf("get SecurityGroups from %s failed", cCloud.PlatType()))
		return
	}
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "get SecurityGroups success",
		"data":   sgs,
	})

}

func LoadInstanceStatusList(c *gin.Context) {
	platType := c.Query("platType")
	regionId := c.Query("regionId")
	if platType == "" {
		EmptyInstancesStatusResponse(c, "need platType")
		return
	}
	if regionId == "" {
		EmptyInstancesStatusResponse(c, "need regionId")
		return
	}

	instanceIds := c.PostFormArray("instanceIds")
	if len(instanceIds) == 0 {
		EmptyInstancesStatusResponse(c, "need instanceIds in form-data")
		return
	}
	cCloud, ok := cloud.DefaultCloudMgr.GetCloud(platType)
	if !ok {
		utils.Logger.Warning(fmt.Sprintf("get %s cloud manager failed", platType))
		EmptyInstancesStatusResponse(c, fmt.Sprintf("do not support cloud %s", platType))
		return
	}
	statusList, err := cCloud.GetInstancesStatus(regionId, instanceIds)
	if err != nil {
		utils.Logger.Warning(fmt.Sprintf("get instances status on [%s] failed", platType))
		EmptyInstancesStatusResponse(c, fmt.Sprintf("get instances status on [%s] failed", platType))
		return
	}
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    fmt.Sprintf("get instances status set in region [%s] success", regionId),
		"data": gin.H{
			"total":           len(statusList),
			"instancesStatus": statusList,
			"currentPageNum":  -1,
		},},
	)
}

func GetAllInstanceStatusList(c *gin.Context) {
	platType := c.Query("platType")
	regionId := c.Query("regionId")
	if platType == "" {
		EmptyInstancesStatusResponse(c, "need platType")
		return
	}
	if regionId == "" {
		EmptyInstancesStatusResponse(c, "need regionId")
		return
	}

	cCloud, ok := cloud.DefaultCloudMgr.GetCloud(platType)
	if !ok {
		utils.Logger.Warning(fmt.Sprintf("get %s cloud manager failed", platType))
		EmptyInstancesStatusResponse(c, fmt.Sprintf("do not support cloud %s", platType))
		return
	}
	statusList, err := cCloud.GetAllInstancesStatus(regionId)
	if err != nil {
		utils.Logger.Warning(fmt.Sprintf("get all instances status on [%s] failed, err: %s", platType, err))
		EmptyInstancesStatusResponse(c, fmt.Sprintf("get all instances status on [%s] failed", platType))
		return
	}
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    fmt.Sprintf("get all instances status set in region [%s] success", regionId),
		"data": gin.H{
			"total":           len(statusList),
			"instancesStatus": statusList,
			"currentPageNum":  -1,
		},},
	)
}

func WsGetAllInstanceStatusList(c *gin.Context) {
	var (
		wsConn *websocket.Conn
		err    error
		data   []byte
	)

	platType := c.Query("platType")
	regionId := c.Query("regionId")
	if platType == "" {
		EmptyInstancesStatusResponse(c, "need platType")
		return
	}
	if regionId == "" {
		EmptyInstancesStatusResponse(c, "need regionId")
		return
	}

	cCloud, ok := cloud.DefaultCloudMgr.GetCloud(platType)
	if !ok {
		utils.Logger.Warning(fmt.Sprintf("get %s cloud manager failed", platType))
		EmptyInstancesStatusResponse(c, fmt.Sprintf("do not support cloud %s", platType))
		return
	}

	// get instance once before establish websocket connection,
	// if error happened, do not establish websocket connection.
	_, err = cCloud.GetAllInstancesStatus(regionId)
	if err != nil {
		utils.Logger.Error(err)
		return
	}

	// Upgrade: websocket
	wsConn, err = upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.Logger.Error(err)
		EmptyInstancesStatusResponse(c, fmt.Sprintf("获取[%s]实例状态列表的websocket连接建立失败", platType))
		return
	}
	conn := utils.InitConnection(wsConn)
	utils.Logger.Info(fmt.Sprintf("getting instances status websocket of [%s] success", platType))

	for {
		statusList, err := cCloud.GetAllInstancesStatus(regionId)
		if err != nil {
			utils.Logger.Error(err)
			data, err = json.Marshal(gin.H{
				"status": "error",
				"msg":    fmt.Sprintf("获取 [%s] [%s] 云主机状态列表失败", platType, regionId),
				"data": gin.H{
					"total":           0,
					"instancesStatus": []*cloud.InstanceStaus{},
					"currentPageNum":  -1,
				},
			})
		} else {
			data, err = json.Marshal(gin.H{
				"status": "ok",
				"msg":    fmt.Sprintf("获取 [%s] [%s] 云主机状态列表成功", platType, regionId),
				"data": gin.H{
					"total":           len(statusList),
					"instancesStatus": statusList,
					"currentPageNum":  -1,
				},
			})
		}
		//fmt.Println(instanceList)
		err = conn.WriteMessage(data)
		if err != nil {
			conn.Close()
			return
		}
		time.Sleep(time.Second * time.Duration(wsUpdateInterval))
	}
}

func GetInstanceList(c *gin.Context) {
	var (
		err  error
		p, s int
	)

	platType := c.Query("platType")
	regionId := c.Query("regionId")
	if platType == "" {
		EmptyInstancesResponse(c, "need platType")
		return
	}
	if regionId == "" {
		EmptyInstancesResponse(c, "need regionId")
		return
	}

	cCloud, ok := cloud.DefaultCloudMgr.GetCloud(platType)
	if !ok {
		utils.Logger.Warning(fmt.Sprintf("get %s cloud manager failed", platType))
		EmptyInstancesResponse(c, fmt.Sprintf("do not support cloud %s", platType))
		return
	}

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

	instanceList, total, err := cCloud.GetInstanceListPerPage(regionId, p, s)

	if err != nil {
		utils.Logger.Error(err)
		EmptyInstancesResponse(c, fmt.Sprintf("获取 [%s] 云主机列表失败", platType))
	} else {
		c.JSON(200, gin.H{
			"status": "ok",
			"msg":    "获取云主机列表成功",
			"data": gin.H{
				"total":          total,
				"instances":      instanceList,
				"currentPageNum": p,
			},
		})
	}
}

func WsGetInstanceList(c *gin.Context) {
	var (
		wsConn *websocket.Conn
		err    error
		data   []byte
		p, s   int
	)

	platType := c.Query("platType")
	regionId := c.Query("regionId")
	if platType == "" {
		EmptyInstancesResponse(c, "need platType")
		return
	}
	if regionId == "" {
		EmptyInstancesResponse(c, "need regionId")
		return
	}

	cCloud, ok := cloud.DefaultCloudMgr.GetCloud(platType)
	if !ok {
		utils.Logger.Warning(fmt.Sprintf("get %s cloud manager failed", platType))
		EmptyInstancesResponse(c, fmt.Sprintf("do not support cloud %s", platType))
		return
	}

	defaultPageNum := utils.GlobalConfig.GetInt("server.default_page_num")
	// aliyun max page size is 100
	defaultPageSize := 100
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

	// get instance once before establish websocket connection,
	// if error happened, do not establish websocket connection.
	_, _, err = cCloud.GetInstanceListPerPage(regionId, p, s)
	if err != nil {
		utils.Logger.Error(err)
		return
	}

	// Upgrade: websocket
	wsConn, err = upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.Logger.Error(err)
		EmptyInstancesResponse(c, fmt.Sprintf("获取[%s]实例列表的websocket连接建立失败", platType))
		return
	}
	conn := utils.InitConnection(wsConn)
	utils.Logger.Info(fmt.Sprintf("getting instances websocket of [%s] success", platType))

	for {
		instanceList, total, err := cCloud.GetInstanceListPerPage(regionId, p, s)
		if err != nil {
			utils.Logger.Error(err)
			data, err = json.Marshal(gin.H{
				"status": "error",
				"msg":    fmt.Sprintf("获取 [%s] 云主机列表失败", platType),
				"data": gin.H{
					"total":          0,
					"instances":      []*cloud.Instance{},
					"currentPageNum": -1,
				},
			})
		} else {
			data, err = json.Marshal(gin.H{
				"status": "ok",
				"msg":    fmt.Sprintf("获取 [%s] 云主机列表成功", platType),
				"data": gin.H{
					"total":          total,
					"instances":      instanceList,
					"currentPageNum": p,
				},
			})
		}
		//fmt.Println(instanceList)
		err = conn.WriteMessage(data)
		if err != nil {
			conn.Close()
			return
		}
		time.Sleep(time.Second * time.Duration(wsUpdateInterval))
	}
}

func GetAllInstance(c *gin.Context) {
	platType := c.Query("platType")
	regionId := c.Query("regionId")
	if platType == "" {
		EmptyInstancesResponse(c, "need platType")
		return
	}
	if regionId == "" {
		EmptyInstancesResponse(c, "need regionId")
		return
	}

	cCloud, ok := cloud.DefaultCloudMgr.GetCloud(platType)
	if !ok {
		utils.Logger.Warning(fmt.Sprintf("get %s cloud manager failed", platType))
		EmptyInstancesResponse(c, fmt.Sprintf("do not support cloud %s", platType))
		return
	}

	instances, err := cCloud.GetAllInstance(regionId)
	if err != nil {
		EmptyInstancesResponse(c, fmt.Sprintf("get all instances from %s failed", cCloud.PlatType()))
		return
	}
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "get all instances success",
		"data": gin.H{
			"total":          len(instances),
			"instances":      instances,
			"currentPageNum": -1,
		},
	})
}

func GetInstance(c *gin.Context) {
	platType := c.Query("platType")
	regionId := c.Query("regionId")
	instanceId := c.Query("instanceId")
	if platType == "" {
		EmptyInstanceResponse(c, "need platType")
		return
	}
	if regionId == "" {
		EmptyInstanceResponse(c, "need regionId")
		return
	}
	if instanceId == "" {
		EmptyInstanceResponse(c, "need instanceId")
		return
	}

	cCloud, ok := cloud.DefaultCloudMgr.GetCloud(platType)
	if !ok {
		utils.Logger.Warning(fmt.Sprintf("get %s cloud manager failed", platType))
		EmptyInstanceResponse(c, fmt.Sprintf("do not support cloud %s", platType))
		return
	}
	instance, err := cCloud.GetInstance(regionId, instanceId)
	if err != nil {
		EmptyInstanceResponse(c, fmt.Sprintf("get instance from %s failed", cCloud.PlatType()))
		return
	}
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "get instance success",
		"data":   []*cloud.Instance{instance},
	})
}

func WsGetInstance(c *gin.Context) {
	var (
		wsConn *websocket.Conn
		err    error
		data   []byte
	)

	platType := c.Query("platType")
	regionId := c.Query("regionId")
	instanceId := c.Query("instanceId")
	if platType == "" {
		EmptyInstanceResponse(c, "need platType")
		return
	}
	if regionId == "" {
		EmptyInstanceResponse(c, "need regionId")
		return
	}
	if instanceId == "" {
		EmptyInstanceResponse(c, "need instanceId")
		return
	}
	cCloud, ok := cloud.DefaultCloudMgr.GetCloud(platType)
	if !ok {
		utils.Logger.Warning(fmt.Sprintf("get %s cloud manager failed", platType))
		EmptyInstanceResponse(c, fmt.Sprintf("do not support cloud %s", platType))
		return
	}

	// get instance once before establish websocket connection,
	// if error happened, do not establish websocket connection.
	_, err = cCloud.GetInstance(regionId, instanceId)
	if err != nil {
		utils.Logger.Error(err)
		return
	}

	// Upgrade: websocket
	wsConn, err = upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.Logger.Error(err)
		EmptyInstanceResponse(c, fmt.Sprintf("获取[%s]实例的websocket连接建立失败", platType))
		return
	}
	conn := utils.InitConnection(wsConn)
	utils.Logger.Info(fmt.Sprintf("getting instance [%s] websocket of [%s] success", instanceId, platType))

	for {
		instance, err := cCloud.GetInstance(regionId, instanceId)
		if err != nil {
			utils.Logger.Error(err)
			data, err = json.Marshal(gin.H{
				"status": "error",
				"msg":    fmt.Sprintf("获取 [%s] 云主机失败", platType),
				"data":   []*cloud.Instance{},
			})
		} else {
			data, err = json.Marshal(gin.H{
				"status": "ok",
				"msg":    fmt.Sprintf("获取 [%s] 云主机成功", platType),
				"data":   []*cloud.Instance{instance},
			})
		}
		//fmt.Println(instanceList)
		err = conn.WriteMessage(data)
		if err != nil {
			conn.Close()
			return
		}
		time.Sleep(time.Second * time.Duration(wsUpdateInterval))
	}
}

func StopInstance(c *gin.Context) {
	platType := c.Query("platType")
	regionId := c.Query("regionId")
	instanceId := c.Query("instanceId")
	if platType == "" {
		BadResponse(c, "need platType")
		return
	}
	if regionId == "" {
		BadResponse(c, "need regionId")
		return
	}
	if instanceId == "" {
		BadResponse(c, "need instanceId")
		return
	}

	cCloud, ok := cloud.DefaultCloudMgr.GetCloud(platType)
	if !ok {
		utils.Logger.Warning(fmt.Sprintf("get %s cloud manager failed", platType))
		BadResponse(c, fmt.Sprintf("do not support cloud %s", platType))
		return
	}
	err := cCloud.StopInstance(regionId, instanceId)
	if err != nil {
		BadResponse(c, fmt.Sprintf("stop instance [%s] on [%s] failed", instanceId, cCloud.PlatType()))
		return
	}
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "stop instance success",
		"data":   gin.H{},
	})
}

func StartInstance(c *gin.Context) {
	platType := c.Query("platType")
	regionId := c.Query("regionId")
	instanceId := c.Query("instanceId")
	if platType == "" {
		BadResponse(c, "need platType")
		return
	}
	if regionId == "" {
		BadResponse(c, "need regionId")
		return
	}
	if instanceId == "" {
		BadResponse(c, "need instanceId")
		return
	}

	cCloud, ok := cloud.DefaultCloudMgr.GetCloud(platType)
	if !ok {
		utils.Logger.Warning(fmt.Sprintf("get %s cloud manager failed", platType))
		BadResponse(c, fmt.Sprintf("do not support cloud %s", platType))
		return
	}
	err := cCloud.StartInstance(regionId, instanceId)
	if err != nil {
		BadResponse(c, fmt.Sprintf("start instance [%s] on [%s] failed", instanceId, cCloud.PlatType()))
		return
	}
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "start instance success",
		"data":   gin.H{},
	})
}

func RebootInstance(c *gin.Context) {
	platType := c.Query("platType")
	regionId := c.Query("regionId")
	instanceId := c.Query("instanceId")
	if platType == "" {
		BadResponse(c, "need platType")
		return
	}
	if regionId == "" {
		BadResponse(c, "need regionId")
		return
	}
	if instanceId == "" {
		BadResponse(c, "need instanceId")
		return
	}

	cCloud, ok := cloud.DefaultCloudMgr.GetCloud(platType)
	if !ok {
		utils.Logger.Warning(fmt.Sprintf("get %s cloud manager failed", platType))
		BadResponse(c, fmt.Sprintf("do not support cloud %s", platType))
		return
	}
	err := cCloud.RebootInstance(regionId, instanceId)
	if err != nil {
		BadResponse(c, fmt.Sprintf("reboot instance [%s] on [%s] failed", instanceId, cCloud.PlatType()))
		return
	}
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "reboot instance success",
		"data":   gin.H{},
	})
}

func DeleteInstance(c *gin.Context) {
	platType := c.Query("platType")
	regionId := c.Query("regionId")
	instanceId := c.Query("instanceId")
	if platType == "" {
		BadResponse(c, "need platType")
		return
	}
	if regionId == "" {
		BadResponse(c, "need regionId")
		return
	}
	if instanceId == "" {
		BadResponse(c, "need instanceId")
		return
	}

	cCloud, ok := cloud.DefaultCloudMgr.GetCloud(platType)
	if !ok {
		utils.Logger.Warning(fmt.Sprintf("get %s cloud manager failed", platType))
		BadResponse(c, fmt.Sprintf("do not support cloud %s", platType))
		return
	}
	err := cCloud.DeleteInstance(regionId, instanceId)
	if err != nil {
		BadResponse(c, fmt.Sprintf("delete instance [%s] on [%s] failed", instanceId, cCloud.PlatType()))
		return
	}
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "delete instance success",
		"data":   gin.H{},
	})
}

func LoadInstanceMonitorData(c *gin.Context) {
	var f forms.InstanceMonitorPostForm
	err := c.ShouldBindJSON(&f)
	if err != nil {
		utils.Logger.Warning(fmt.Sprintf("invalid post data"))
		return
	}
	platType := f.PlatType
	regionId := f.RegionId
	metricName := f.MetricName
	startTime := f.StartTime
	endTime := f.EndTime
	period := f.Period
	instanceIds := f.InstanceIds

	cCloud, ok := cloud.DefaultCloudMgr.GetCloud(platType)
	if !ok {
		utils.Logger.Warning(fmt.Sprintf("get %s cloud manager failed", platType))
		EmptyDataPointsResponse(c, fmt.Sprintf("do not support cloud %s", platType))
		return
	}
	dataPoints, err := cCloud.GetMonitorDataOfInstances(regionId, metricName, startTime, endTime, period, instanceIds)
	if err != nil {
		EmptyDataPointsResponse(c, fmt.Sprintf("get [%s] from [%s] failed", metricName, cCloud.PlatType()))
		return
	}
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "get monitor data success",
		"data":   dataPoints,
	})

}

func SevMonitorTask(conn *utils.Connection) {

	go func(conn *utils.Connection) {
		for {
			bytes, err := conn.ReadMessage()

			if err != nil {
				utils.Logger.Error(fmt.Sprintf("read from websocket of [%s] failed", conn.WsConn.RemoteAddr().String()))
				conn.Close()
				return
			}

			var (
				data []byte
				//err  error
				ok bool
				//task *MonitorTask
			)

			var f forms.InstanceMonitorPostForm
			err = json.Unmarshal(bytes, &f)
			if err != nil {
				utils.Logger.Warningf("parse ws data failed, %s", err)
				return
			}

			platType := f.PlatType
			regionId := f.RegionId
			metricName := f.MetricName
			startTime := f.StartTime
			endTime := f.EndTime
			period := f.Period
			instanceIds := f.InstanceIds
			DurationType := f.DurationType

			ss := strings.Split(startTime, "+")
			es := strings.Split(endTime, "+")
			var (
				st, et time.Time
			)

			// calculate time delta for fixed.
			st, err = time.Parse("2006-01-02T15:04:05", ss[0])
			if err != nil {
				utils.Logger.Error(fmt.Sprintf("invalid time format [%s]", st))
				return
			}
			et, err = time.Parse("2006-01-02T15:04:05", es[0])
			if err != nil {
				utils.Logger.Error(fmt.Sprintf("invalid time format [%s]", et))
				return
			}
			timeDelta := et.Sub(st) // time.Duration

			cCloud, ok := cloud.DefaultCloudMgr.GetCloud(platType)
			if !ok {
				utils.Logger.Warningf("get %s cloud manager failed", platType)
				return
			}
			if DurationType == "fixed" {
				startTime = fmt.Sprintf("%s%s", time.Now().Add(-timeDelta).Format("2006-01-02T15:04:05"), "+08:00")
				endTime = fmt.Sprintf("%s%s", time.Now().Format("2006-01-02T15:04:05"), "+08:00")
			}

			go func() {
				dataPoints, err := cCloud.GetMonitorDataOfInstances(regionId, metricName, startTime, endTime, period, instanceIds)
				if err != nil {
					utils.Logger.Error(err)
					data, err = json.Marshal(gin.H{
						"status": "error",
						"msg":    fmt.Sprintf("获取 [%s] [%s] 数据失败", platType, metricName),
						"data":   []*cloud.DataPoint{},
					})
				} else {
					data, err = json.Marshal(gin.H{
						"status": "ok",
						"msg":    fmt.Sprintf("获取 [%s] [%s] 数据成功", platType, metricName),
						"data":   dataPoints,
					})
				}
				//fmt.Printf("%s data: %v\n", metricName, string(data))
				err = conn.WriteMessage(data)
				if err != nil {
					fmt.Println("写入失败")
					conn.Close()
					return
				}
			}()
		}
	}(conn)
}

func WsGetInstanceMonitorData(c *gin.Context) {
	var (
		wsConn *websocket.Conn
		err    error
	)

	// Upgrade: websocket
	wsConn, err = upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.Logger.Error(err)
		return
	}

	conn := utils.InitConnection(wsConn)
	utils.Logger.Info(fmt.Sprintf("getting new websocket [%s]", wsConn.RemoteAddr().String()))
	SevMonitorTask(conn)
}
