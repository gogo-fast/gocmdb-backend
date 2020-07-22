package tencent

import (
	"apiserver/cloud"
	"apiserver/utils"
	"encoding/json"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
)

type Params struct {
	Namespace  string
	MetricName string
	Period     uint64
	StartTime  string
	EndTime    string
	Instances  []Instance
}

type Instance struct {
	Dimensions []Dimension
}

type Dimension struct {
	Name  string
	Value string
}

func (m *TenCentMgr) GetMonitorDataOfInstances(regionId, metricName, startTime, endTime string, period int, instanceIDs []string) ([]*cloud.DataPoint, error) {
	client, err := monitor.NewClient(m.Credential, regionId, m.MonitorClientProfile)
	if err != nil {
		utils.Logger.Error("init tencent client failed")
		return nil, err
	}

	request := monitor.NewGetMonitorDataRequest()
	p := Params{
		Namespace:  "QCE/CVM",
		MetricName: metricName,
		Period:     uint64(period),
		StartTime:  startTime,
		EndTime:    endTime,
	}
	for _, id := range instanceIDs {
		dm := Dimension{
			Name:  "InstanceId",
			Value: id,
		}
		p.Instances = append(p.Instances, Instance{Dimensions: []Dimension{dm}})
	}
	bytes, err := json.Marshal(p)
	err = request.FromJsonString(fmt.Sprint(string(bytes)))
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("An API error has returned on [%s]: %s", m.CloudType, err))
		return nil, err
	}
	response, err := client.GetMonitorData(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		utils.Logger.Error(fmt.Sprintf("An API error has returned on [%s]: %s", m.CloudType, err))
		return nil, err
	}
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("An API error has returned on [%s]: %s", m.CloudType, err))
		return nil, err
	}
	dataPointList := make([]*cloud.DataPoint, 0, len(instanceIDs))
	for _, v := range response.Response.DataPoints {
		dp := &cloud.DataPoint{}
		dp.Timestamps = v.Timestamps
		dp.Values = v.Values
		dp.MetricName = metricName
		dp.InstanceId = *v.Dimensions[0].Value
		dataPointList = append(dataPointList, dp)
	}

	return dataPointList, nil
}
