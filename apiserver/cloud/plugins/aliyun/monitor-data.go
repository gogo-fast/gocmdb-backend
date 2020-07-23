package aliyun

import (
	"apiserver/cloud"
	"apiserver/utils"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"strconv"
	"strings"
)

type DataPoint struct {
	Maximum    float64 `json:"Maximum"`
	Minimum    float64 `json:"Minimum"`
	Average    float64 `json:"Average"`
	Timestamp  float64 `json:"timestamp"`
	UserId     string  `json:"userId"`
	InstanceId string  `json:"instanceId"`
}

// GetMonitorDataOfInstances
// get monitor data per instance, different from tenCent cLoud.
func (m *AliMgr) GetMonitorDataOfInstances(regionId, metricName, startTime, endTime string, period int, instanceIDs []string) ([]*cloud.DataPoint, error) {
	client, err := cms.NewClientWithAccessKey(regionId, m.AccessKeyId, m.AccessKeySecret)
	if err != nil {
		utils.Logger.Error("init aliyun client failed")
		return nil, err
	}
	// reformat time format important
	// time format should be 2000-01-02 15:04:45+08:00
	// time format 2000-01-02T15:04:45+08:00 is not correct for aliyun.
	startTime = strings.Replace(startTime,"T", " ", 1)
	endTime = strings.Replace(endTime,"T", " ", 1)
	dataPointList := make([]*cloud.DataPoint, 0, len(instanceIDs))
	//fmt.Println(startTime, "  |   ", endTime)
	//fmt.Println(strings.Repeat("~", 30))
	for _, instanceId := range instanceIDs {
		request := cms.CreateDescribeMetricDataRequest()
		request.Scheme = "https"

		request.Namespace = "acs_ecs_dashboard"
		request.MetricName = metricName
		request.StartTime = startTime
		request.EndTime = endTime
		request.Dimensions = fmt.Sprintf(`{"instanceId": "%s"}`, instanceId)
		request.Period = strconv.Itoa(period)
		response, err := client.DescribeMetricData(request)

		if err != nil {
			utils.Logger.Error(fmt.Sprintf("An API error has returned on [%s]: %s", m.CloudType, err))
			return nil, err
		}

		var ds []*DataPoint  // you should use pointer of struct, or all the data are the last one.
		err = json.Unmarshal([]byte(response.Datapoints), &ds)
		if err != nil {
			utils.Logger.Error("parse response.Datapoints failed")
			return nil, err
		}

		dp := &cloud.DataPoint{
			MetricName: metricName,
			InstanceId: instanceId,
		}
		Timestamps := make([]*float64, 0, 10)
		Values := make([]*float64, 0, 10)
		for _, d := range ds {
			Timestamps = append(Timestamps, &d.Timestamp)
			Values = append(Values, &d.Maximum)
		}
		dp.Timestamps = Timestamps
		dp.Values = Values
		dataPointList = append(dataPointList, dp)

	}

	return dataPointList, nil
}
