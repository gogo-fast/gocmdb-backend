package aliyun

import (
	"apiserver/cloud"
	"apiserver/utils"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

func (m *AliMgr) GetZones(regionId string) ([]*cloud.Zone, error) {
	client, err := ecs.NewClientWithAccessKey(regionId, m.AccessKeyId, m.AccessKeySecret)
	if err != nil {
		utils.Logger.Error("init aliyun client failed")
		return nil, err
	}
	request := ecs.CreateDescribeZonesRequest()
	request.Scheme = "https"

	response, err := client.DescribeZones(request)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("An API error has returned on [%s]: %s", m.CloudType, err))
		return nil, err
	}
	zoneList := make([]*cloud.Zone, 0, 5)

	for _, z := range response.Zones.Zone {
		zone := &cloud.Zone{}
		zone.ZoneId = z.ZoneId
		zone.ZoneName = z.LocalName
		zoneList = append(zoneList, zone)
	}
	return zoneList, nil
}