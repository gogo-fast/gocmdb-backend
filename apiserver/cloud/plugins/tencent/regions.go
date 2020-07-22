package tencent

import (
	"apiserver/cloud"
	"apiserver/utils"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func (m *TenCentMgr) GetRegions() ([]*cloud.Region, error) {
	client, err := cvm.NewClient(m.Credential, "", m.CvmClientProfile)
	if err != nil {
		utils.Logger.Error("init tencent client failed")
		return nil, err
	}

	request := cvm.NewDescribeRegionsRequest()
	params := "{}"
	err = request.FromJsonString(params)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("params of get regions of [%s] invalid, err: %s", m.CloudType, err))
		return nil, err
	}

	response, err := client.DescribeRegions(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		utils.Logger.Error(fmt.Sprintf("An API error has returned on [%s]: %s", m.CloudType, err))
		return nil, err
	}
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("An API error has returned on [%s]: %s", m.CloudType, err))
		return nil, err
	}
	var regionList = make([]*cloud.Region, 0, 30)
	for _, v := range response.Response.RegionSet {
		rg := &cloud.Region{}
		rg.RegionId = *v.Region
		rg.RegionName = *v.RegionName
		regionList = append(regionList, rg)
	}
	return regionList, nil
}
