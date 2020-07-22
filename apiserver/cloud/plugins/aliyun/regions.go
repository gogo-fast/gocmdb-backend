package aliyun

import (
	"apiserver/cloud"
	"apiserver/utils"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

func (m *AliMgr) GetRegions() ([]*cloud.Region, error) {
	client, err := ecs.NewClientWithAccessKey(m.DefaultRegionId, m.AccessKeyId, m.AccessKeySecret)
	if err != nil {
		utils.Logger.Error("init aliyun client failed")
		return nil, err
	}

	request := ecs.CreateDescribeRegionsRequest()
	request.Scheme = "https"

	response, err := client.DescribeRegions(request)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("An API error has returned on [%s]: %s", m.CloudType, err))
		return nil, err
	}

	var regionList = make([]*cloud.Region, 0, 30)
	for _, r := range response.Regions.Region {
		rg := &cloud.Region{}
		rg.RegionId = r.RegionId
		rg.RegionName = r.LocalName
		regionList = append(regionList, rg)
	}
	return regionList, nil
}
