package tencent

import (
	"apiserver/cloud"
	"apiserver/utils"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func (m *TenCentMgr) GetZones(regionId string) ([]*cloud.Zone, error) {
	client, err := cvm.NewClient(m.Credential, regionId, m.CvmClientProfile)
	if err != nil {
		utils.Logger.Error("init tencent client failed")
		return nil, err
	}

	request := cvm.NewDescribeZonesRequest()
	params := "{}"
	err = request.FromJsonString(params)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("params of get regions of [%s] invalid, err: %s", m.CloudType, err))
		return nil, err
	}

	response, err := client.DescribeZones(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		utils.Logger.Error(fmt.Sprintf("An API error has returned on [%s]: %s", m.CloudType, err))
		return nil, err
	}
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("An API error has returned on [%s]: %s", m.CloudType, err))
		return nil, err
	}
	zoneList := make([]*cloud.Zone, 0, 5)

	for _, z := range response.Response.ZoneSet {
		zone := &cloud.Zone{}
		zone.ZoneId = *z.ZoneId
		zone.ZoneName = *z.ZoneName
		zoneList = append(zoneList, zone)
	}
	return zoneList, nil
}
