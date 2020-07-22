package aliyun

import (
	"apiserver/cloud"
	"fmt"
)

func (m *AliMgr) GetMonitorDataOfInstances(regionId, metricName, startTime, endTime string, period int, instanceIDs []string) ([]*cloud.DataPoint, error) {
	return nil, fmt.Errorf("not implement yet")
}
