package aliyun

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"apiserver/cloud"
	"apiserver/utils"
	"sync"
)

type AliMgr struct {
	mutex           sync.RWMutex
	CloudName       string
	CloudType       string
	DefaultRegionId string
	AccessKeyId     string
	AccessKeySecret string
	Page            int
	Size            int
}

func NewAliMgr() *AliMgr {
	var aliMgr = &AliMgr{}
	aliMgr.CloudName = aliMgr.Name()
	aliMgr.CloudType = aliMgr.PlatType()
	return aliMgr
}

func (m *AliMgr) Name() string {
	return "阿里云"
}

func (m *AliMgr) PlatType() string {
	return "aliyun"
}

func (m *AliMgr) Init(DefaultRegionId, AccessKeyId, AccessKeySecret string) {
	m.DefaultRegionId = DefaultRegionId
	m.AccessKeyId = AccessKeyId
	m.AccessKeySecret = AccessKeySecret
}

func (m *AliMgr) TestConn() error {
	_, err := m.GetRegions()
	if err != nil {
		return err
	}
	return nil
}

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

func (m *AliMgr) GetInstancesStatus(regionId string, instanceIds []string) ([]*cloud.InstanceStaus, error) {
	var (
		total              int = 50
		p                  int = 1
		s                  int = 50 //  max page size of DescribeInstanceStatus of aliyun is 50.
		InstanceStatusList     = make([]*cloud.InstanceStaus, 0, len(instanceIds))
	)

	client, err := ecs.NewClientWithAccessKey(regionId, m.AccessKeyId, m.AccessKeySecret)
	if err != nil {
		utils.Logger.Error("init aliyun client failed")
		return nil, err
	}

	request := ecs.CreateDescribeInstanceStatusRequest()
	request.Scheme = "https"
	request.InstanceId = &instanceIds
	for p*s <= total {
		request.PageSize = requests.NewInteger(s)
		request.PageNumber = requests.NewInteger(p)
		response, err := client.DescribeInstanceStatus(request)
		if err != nil {
			utils.Logger.Error(fmt.Sprintf("An API error has returned on [%s]: %s", m.CloudType, err))
			return InstanceStatusList, err
		}
		total = response.TotalCount
		p = response.PageNumber
		s = response.PageSize
		p ++
		for _, v := range response.InstanceStatuses.InstanceStatus {
			s := &cloud.InstanceStaus{
				InstanceId:    v.InstanceId,
				InstanceState: m.InstanceStatusTransform(v.Status),
			}
			InstanceStatusList = append(InstanceStatusList, s)
		}
	}
	return InstanceStatusList, nil
}

func (m *AliMgr) GetAllInstancesStatus(regionId string) ([]*cloud.InstanceStaus, error) {
	var (
		total              = 50
		p                  = 1
		s                  = 50 //  max page size of DescribeInstanceStatus of aliyun is 50.
		InstanceStatusList = make([]*cloud.InstanceStaus, 0, 50)
	)

	client, err := ecs.NewClientWithAccessKey(regionId, m.AccessKeyId, m.AccessKeySecret)
	if err != nil {
		utils.Logger.Error("init aliyun client failed")
		return nil, err
	}

	request := ecs.CreateDescribeInstanceStatusRequest()
	request.Scheme = "https"

	for p*s <= total {
		request.PageNumber = requests.NewInteger(p)
		request.PageSize = requests.NewInteger(s)
		response, err := client.DescribeInstanceStatus(request)
		if err != nil {
			utils.Logger.Error(fmt.Sprintf("An API error has returned on [%s]: %s", m.CloudType, err))
			return InstanceStatusList, err
		}
		total = response.TotalCount
		p = response.PageNumber
		s = response.PageSize
		p ++
		for _, v := range response.InstanceStatuses.InstanceStatus {
			s := &cloud.InstanceStaus{
				InstanceId:    v.InstanceId,
				InstanceState: m.InstanceStatusTransform(v.Status),
			}
			InstanceStatusList = append(InstanceStatusList, s)
		}
	}
	fmt.Println("InstanceStatusList:", InstanceStatusList)
	return InstanceStatusList, nil
}

func (m *AliMgr) GetInstance(regionId, instanceId string) (*cloud.Instance, error) {
	client, err := ecs.NewClientWithAccessKey(regionId, m.AccessKeyId, m.AccessKeySecret)
	if err != nil {
		utils.Logger.Error("init aliyun client failed")
		return nil, err
	}
	request := ecs.CreateDescribeInstanceAttributeRequest()
	request.Scheme = "https"
	request.InstanceId = instanceId
	response, err := client.DescribeInstanceAttribute(request)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("An API error has returned on [%s]: %s", m.CloudType, err))
		return nil, err
	}

	instance := &cloud.Instance{}

	instance.InstanceId = response.InstanceId
	instance.Uuid = response.SerialNumber
	instance.HostName = response.HostName
	instance.RegionId = response.RegionId
	instance.ZoneId = response.ZoneId
	instance.Status = m.InstanceStatusTransform(response.Status)
	instance.OSName = "Api UNSupport"
	instance.Cpu = response.Cpu
	instance.Memory = response.Memory
	instance.InstanceType = response.InstanceType
	instance.CreatedTime = response.CreationTime
	instance.Description = response.Description
	instance.InternetChargeType = response.InternetChargeType
	instance.VpcId = response.VpcAttributes.VpcId
	privateIps, _ := json.Marshal(response.VpcAttributes.PrivateIpAddress.IpAddress)
	instance.PrivateIpAddress = string(privateIps)

	publicIpList := response.PublicIpAddress.IpAddress
	allPublicIpList := make([]string, 0, len(publicIpList)+1)
	if response.EipAddress.IpAddress != "" {
		allPublicIpList = append(allPublicIpList, response.EipAddress.IpAddress)
	}
	allPublicIpList = append(allPublicIpList, publicIpList...)
	publicIps, _ := json.Marshal(allPublicIpList)
	instance.PublicIpAddress = string(publicIps)

	instance.InternetMaxBandwidthIn = response.InternetMaxBandwidthIn
	instance.InternetMaxBandwidthOut = response.InternetMaxBandwidthOut

	return instance, nil
}

func (m *AliMgr) GetInstanceListPerPage(regionId string, page, size int) ([]*cloud.Instance, int, error) {
	client, err := ecs.NewClientWithAccessKey(regionId, m.AccessKeyId, m.AccessKeySecret)
	if err != nil {
		utils.Logger.Error("init aliyun client failed")
		return nil, 0, err
	}
	request := ecs.CreateDescribeInstancesRequest()
	request.Scheme = "https"

	request.PageNumber = requests.NewInteger(page)
	request.PageSize = requests.NewInteger(size)
	//fmt.Println("request.PageNumber:", request.PageNumber, "request.PageSize:" , request.PageSize)
	response, err := client.DescribeInstances(request)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("An API error has returned on [%s]: %s", m.CloudType, err))
		return nil, 0, err
	}

	instances := response.Instances.Instance
	instanceList := make([]*cloud.Instance, 0, len(instances))
	for _, v := range instances {
		instance := &cloud.Instance{}

		instance.InstanceId = v.InstanceId
		instance.Uuid = v.SerialNumber
		instance.HostName = v.HostName
		instance.RegionId = v.RegionId
		instance.ZoneId = v.ZoneId
		instance.Status = m.InstanceStatusTransform(v.Status)
		instance.OSName = v.OSName
		instance.Cpu = v.Cpu
		instance.Memory = v.Memory
		instance.InstanceType = v.InstanceType
		instance.CreatedTime = v.CreationTime
		instance.Description = v.Description
		instance.InternetChargeType = v.InternetChargeType
		instance.VpcId = v.VpcAttributes.VpcId
		privateIps, _ := json.Marshal(v.VpcAttributes.PrivateIpAddress.IpAddress)
		instance.PrivateIpAddress = string(privateIps)

		publicIpList := v.PublicIpAddress.IpAddress
		allPublicIpList := make([]string, 0, len(publicIpList)+1)
		if v.EipAddress.IpAddress != "" {
			allPublicIpList = append(allPublicIpList, v.EipAddress.IpAddress)
		}
		allPublicIpList = append(allPublicIpList, publicIpList...)
		publicIps, _ := json.Marshal(allPublicIpList)
		instance.PublicIpAddress = string(publicIps)

		instance.InternetMaxBandwidthIn = v.InternetMaxBandwidthIn
		instance.InternetMaxBandwidthOut = v.InternetMaxBandwidthOut

		instanceList = append(instanceList, instance)
	}
	return instanceList, response.TotalCount, nil
}

func (m *AliMgr) GetAllInstance(regionId string) ([]*cloud.Instance, error) {
	var (
		err          error
		total        int
		p            = 1
		s            = 100 // max page size of DescribeInstances of aliyun is 100.
		instanceList []*cloud.Instance
		instances    []*cloud.Instance
	)

	instanceList = make([]*cloud.Instance, 0, 100)
	for {
		instances, total, err = m.GetInstanceListPerPage(regionId, p, s)
		if err != nil {
			break
		}
		instanceList = append(instanceList, instances...)
		if p*s >= total {
			break
		}
		p += 1
	}
	return instanceList, err
}

func (m *AliMgr) StartInstance(regionId, instanceId string) error {
	client, err := ecs.NewClientWithAccessKey(regionId, m.AccessKeyId, m.AccessKeySecret)
	if err != nil {
		utils.Logger.Error("init aliyun client failed")
		return err
	}
	request := ecs.CreateStartInstanceRequest()
	request.Scheme = "https"

	request.InstanceId = instanceId

	_, err = client.StartInstance(request)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("An API error has returned on [%s]: %s", m.CloudType, err))
		return err
	}
	return nil
}

func (m *AliMgr) StopInstance(regionId, instanceId string) error {
	client, err := ecs.NewClientWithAccessKey(regionId, m.AccessKeyId, m.AccessKeySecret)
	if err != nil {
		utils.Logger.Error("init aliyun client failed")
		return err
	}
	request := ecs.CreateStopInstanceRequest()
	request.Scheme = "https"

	request.InstanceId = instanceId

	_, err = client.StopInstance(request)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("An API error has returned on [%s]: %s", m.CloudType, err))
		return err
	}
	return nil
}

func (m *AliMgr) RebootInstance(regionId, instanceId string) error {
	client, err := ecs.NewClientWithAccessKey(regionId, m.AccessKeyId, m.AccessKeySecret)
	if err != nil {
		utils.Logger.Error("init aliyun client failed")
		return err
	}
	request := ecs.CreateRebootInstanceRequest()
	request.Scheme = "https"

	request.InstanceId = instanceId

	_, err = client.RebootInstance(request)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("An API error has returned on [%s]: %s", m.CloudType, err))
		return err
	}
	return nil
}

func (m *AliMgr) DeleteInstance(regionId, instanceId string) error {
	client, err := ecs.NewClientWithAccessKey(regionId, m.AccessKeyId, m.AccessKeySecret)
	if err != nil {
		utils.Logger.Error("init aliyun client failed")
		return err
	}
	request := ecs.CreateDeleteInstanceRequest()
	request.Scheme = "https"

	request.InstanceId = instanceId

	_, err = client.DeleteInstance(request)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("An API error has returned on [%s]: %s", m.CloudType, err))
		return err
	}
	return nil
}

func (m *AliMgr) InstanceStatusTransform(status string) string {
	InstanceStatusMap := map[string]string{
		"Pending":  cloud.StatusPending,
		"Running":  cloud.StatusRunning,
		"Stopped":  cloud.StatusStopped,
		"Starting": cloud.StatusStarting,
		"Stopping": cloud.StatusStopping,
	}
	if s, ok := InstanceStatusMap[status]; ok {
		return s
	}
	return cloud.StatusUnknown
}

func init() {
	defaultRegionId := utils.GlobalConfig.GetString("aliyun.default_region_id")
	accessKeyId := utils.GlobalConfig.GetString("aliyun.access_key_id")
	accessKeySecret := utils.GlobalConfig.GetString("aliyun.access_key_secret")
	aliMgr := NewAliMgr()
	aliMgr.Init(defaultRegionId, accessKeyId, accessKeySecret)
	err := aliMgr.TestConn()
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("connect to [%s] failed, err: %s", aliMgr.CloudType, err))
	} else {
		utils.Logger.Info(fmt.Sprintf("connect to [%s] success", aliMgr.CloudType))
	}

	cloud.DefaultCloudMgr.Register(aliMgr)
}
