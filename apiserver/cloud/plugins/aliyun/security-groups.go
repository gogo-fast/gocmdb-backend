package aliyun

import (
	"apiserver/cloud"
	"apiserver/utils"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

func (m *AliMgr) GetSecurityGroups(regionId string) ([]*cloud.SecurityGroup, error) {
	var (
		limit       = 50 //  max page size of DescribeSecurityGroups of aliyun is 50.
		offset      = 1
		currentPage = 1
		total       = 2
	)
	client, err := ecs.NewClientWithAccessKey(regionId, m.AccessKeyId, m.AccessKeySecret)
	if err != nil {
		fmt.Println("init aliyun client failed")
		return nil, err
	}

	var sgs = make([]*cloud.SecurityGroup, 0, limit)
	for offset < total {
		request := ecs.CreateDescribeSecurityGroupsRequest()
		request.Scheme = "https"

		request.PageNumber = requests.NewInteger(currentPage)
		request.PageSize = requests.NewInteger(limit)
		response, err := client.DescribeSecurityGroups(request)
		if err != nil {
			utils.Logger.Error(fmt.Sprintf("An API error has returned on [%s]: %s", m.CloudType, err))
			return nil, err
		}

		total = response.TotalCount
		offset += limit
		for _, sg := range response.SecurityGroups.SecurityGroup {
			sgs = append(sgs, &cloud.SecurityGroup{
				SecurityGroupId:   sg.SecurityGroupId,
				SecurityGroupName: sg.SecurityGroupName,
				VpcId:             sg.VpcId,
			})
		}
	}
	return sgs, nil
}
