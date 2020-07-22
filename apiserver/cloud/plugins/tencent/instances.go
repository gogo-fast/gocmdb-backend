package tencent

import (
	"apiserver/cloud"
	"apiserver/utils"
	"encoding/json"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func (m *TenCentMgr) GetInstancesStatus(regionId string, instanceIds []string) ([]*cloud.InstanceStaus, error) {
	var (
		total              int64 = 100
		offset             int64 = 0
		limit              int64 = 100 //  max page size of DescribeInstancesStatus of tencent is 100.
		InstanceStatusList       = make([]*cloud.InstanceStaus, 0, len(instanceIds))
	)

	client, err := cvm.NewClient(m.Credential, regionId, m.CvmClientProfile)
	if err != nil {
		utils.Logger.Error("init tencent client failed")
		return nil, err
	}

	request := cvm.NewDescribeInstancesStatusRequest()
	ids := make([]*string, len(instanceIds))

	for i, v := range instanceIds {
		ids[i] = &v
	}
	request.InstanceIds = ids

	for offset <= total {
		request.Offset = &offset
		request.Limit = &limit
		response, err := client.DescribeInstancesStatus(request)
		if err != nil {
			utils.Logger.Error(fmt.Sprintf("An API error has returned on [%s]: %s", m.CloudType, err))
			return InstanceStatusList, err
		}
		total = *response.Response.TotalCount
		offset += limit
		for _, v := range response.Response.InstanceStatusSet {
			s := &cloud.InstanceStaus{
				InstanceId:    *v.InstanceId,
				InstanceState: m.InstanceStatusTransform(*v.InstanceState),
			}
			InstanceStatusList = append(InstanceStatusList, s)
		}

	}
	return InstanceStatusList, nil
}

func (m *TenCentMgr) GetAllInstancesStatus(regionId string) ([]*cloud.InstanceStaus, error) {
	var (
		total              int64 = 100
		offset             int64 = 0
		limit              int64 = 100 //  max page size of DescribeInstancesStatus of tencent is 100.
		InstanceStatusList       = make([]*cloud.InstanceStaus, 0, 100)
	)

	client, err := cvm.NewClient(m.Credential, regionId, m.CvmClientProfile)
	if err != nil {
		utils.Logger.Error("init tencent client failed")
		return nil, err
	}

	request := cvm.NewDescribeInstancesStatusRequest()

	for offset <= total {
		request.Offset = &offset
		request.Limit = &limit
		response, err := client.DescribeInstancesStatus(request)
		if err != nil {
			utils.Logger.Error(fmt.Sprintf("An API error has returned on [%s]: %s", m.CloudType, err))
			return InstanceStatusList, err
		}
		total = *response.Response.TotalCount
		offset += limit
		for _, v := range response.Response.InstanceStatusSet {
			s := &cloud.InstanceStaus{
				InstanceId:    *v.InstanceId,
				InstanceState: m.InstanceStatusTransform(*v.InstanceState),
			}
			InstanceStatusList = append(InstanceStatusList, s)
		}

	}
	return InstanceStatusList, nil
}

func (m *TenCentMgr) GetInstance(regionId, instanceId string) (*cloud.Instance, error) {
	client, err := cvm.NewClient(m.Credential, regionId, m.CvmClientProfile)
	if err != nil {
		utils.Logger.Error("init tencent client failed")
		return nil, err
	}

	request := cvm.NewDescribeInstancesRequest()
	request.InstanceIds = []*string{&instanceId}
	response, err := client.DescribeInstances(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		utils.Logger.Error(fmt.Sprintf("An API error has returned on [%s]: %s", m.CloudType, err))
		return nil, err
	}
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("An API error has returned on [%s]: %s", m.CloudType, err))
		return nil, err
	}
	instanceSet := response.Response.InstanceSet
	if len(instanceSet) == 0 {
		utils.Logger.Error(fmt.Sprintf("No instance returned on [%s]: %s", m.CloudType, err))
		// you should return a error by hand, or err is nil, nil err with nil cloud.Instance is not a correct return.
		return nil, fmt.Errorf("instance with id [%s] not exist", instanceId)
	}
	v := instanceSet[0]
	instance := &cloud.Instance{}

	instance.InstanceId = *v.InstanceId
	instance.Uuid = *v.Uuid
	instance.HostName = ""
	instance.RegionId = regionId
	instance.ZoneId = *v.Placement.Zone
	instance.Status = m.InstanceStatusTransform(*v.InstanceState)
	instance.OSName = *v.OsName
	instance.Cpu = int(*v.CPU)
	instance.Memory = int(*v.Memory) * 1024
	instance.InstanceType = *v.InstanceType
	instance.CreatedTime = *v.CreatedTime
	instance.Description = ""
	instance.InternetChargeType = *v.InternetAccessible.InternetChargeType
	instance.VpcId = *v.VirtualPrivateCloud.VpcId
	privateIps, _ := json.Marshal(v.PrivateIpAddresses)
	instance.PrivateIpAddress = string(privateIps)

	publicIps, _ := json.Marshal(v.PublicIpAddresses)
	instance.PublicIpAddress = string(publicIps)

	instance.InternetMaxBandwidthOut = int(*v.InternetAccessible.InternetMaxBandwidthOut)
	instance.InternetMaxBandwidthIn = instance.InternetMaxBandwidthOut

	return instance, nil
}

func (m *TenCentMgr) GetInstanceListPerPage(regionId string, page, size int) ([]*cloud.Instance, int, error) {

	client, err := cvm.NewClient(m.Credential, regionId, m.CvmClientProfile)
	if err != nil {
		utils.Logger.Error("init tencent client failed")
		return nil, 0, err
	}

	request := cvm.NewDescribeInstancesRequest()

	p := int64(page)
	s := int64(size)
	offset := (p - 1) * s
	limit := s
	request.Offset = &offset
	request.Limit = &limit

	response, err := client.DescribeInstances(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		utils.Logger.Error(fmt.Sprintf("An API error has returned on [%s]: %s", m.CloudType, err))
		return nil, 0, err
	}
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("An API error has returned on [%s]: %s", m.CloudType, err))
		return nil, 0, err
	}

	instances := response.Response.InstanceSet
	instanceList := make([]*cloud.Instance, 0, len(instances))
	for _, v := range instances {
		instance := &cloud.Instance{}

		instance.InstanceId = *v.InstanceId
		instance.Uuid = *v.Uuid
		instance.HostName = ""
		instance.RegionId = regionId
		instance.ZoneId = *v.Placement.Zone
		instance.Status = m.InstanceStatusTransform(*v.InstanceState)
		instance.OSName = *v.OsName
		instance.Cpu = int(*v.CPU)
		instance.Memory = int(*v.Memory) * 1024
		instance.InstanceType = *v.InstanceType
		instance.CreatedTime = *v.CreatedTime
		instance.Description = ""
		instance.InternetChargeType = *v.InternetAccessible.InternetChargeType
		instance.VpcId = *v.VirtualPrivateCloud.VpcId
		privateIps, _ := json.Marshal(v.PrivateIpAddresses)
		instance.PrivateIpAddress = string(privateIps)

		publicIps, _ := json.Marshal(v.PublicIpAddresses)
		instance.PublicIpAddress = string(publicIps)

		instance.InternetMaxBandwidthOut = int(*v.InternetAccessible.InternetMaxBandwidthOut)
		instance.InternetMaxBandwidthIn = instance.InternetMaxBandwidthOut

		instanceList = append(instanceList, instance)
	}
	return instanceList, int(*response.Response.TotalCount), nil
}

func (m *TenCentMgr) GetAllInstance(regionId string) ([]*cloud.Instance, error) {
	var (
		err          error
		total        int
		p            = 1
		s            = 100 //  max page size of DescribeInstances of tencent is 100.
		instanceList []*cloud.Instance
		instances    []*cloud.Instance
	)

	instanceList = make([]*cloud.Instance, 0, 100)
	for {
		instances, total, err = m.GetInstanceListPerPage(regionId, p, s)
		instanceList = append(instanceList, instances...)
		if p*s >= total {
			break
		}
		p += 1
	}
	return instanceList, err
}

func (m *TenCentMgr) StartInstance(regionId, instanceId string) error {
	client, err := cvm.NewClient(m.Credential, regionId, m.CvmClientProfile)
	if err != nil {
		utils.Logger.Error("init tencent client failed")
		return err
	}

	request := cvm.NewStartInstancesRequest()

	request.InstanceIds = []*string{&instanceId}

	_, err = client.StartInstances(request)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("start instance [%s] failed", instanceId))
		return err
	}
	return nil
}

func (m *TenCentMgr) StopInstance(regionId, instanceId string) error {
	client, err := cvm.NewClient(m.Credential, regionId, m.CvmClientProfile)
	if err != nil {
		utils.Logger.Error("init tencent client failed")
		return err
	}

	request := cvm.NewStopInstancesRequest()

	request.InstanceIds = []*string{&instanceId}

	_, err = client.StopInstances(request)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("stop instance [%s] failed", instanceId))
		return err
	}
	return nil
}

func (m *TenCentMgr) RebootInstance(regionId, instanceId string) error {
	client, err := cvm.NewClient(m.Credential, regionId, m.CvmClientProfile)
	if err != nil {
		utils.Logger.Error("init tencent client failed")
		return err
	}

	request := cvm.NewRebootInstancesRequest()

	request.InstanceIds = []*string{&instanceId}

	_, err = client.RebootInstances(request)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("reboot instance [%s] failed", instanceId))
		return err
	}
	return nil
}

func (m *TenCentMgr) DeleteInstance(regionId, instanceId string) error {
	client, err := cvm.NewClient(m.Credential, regionId, m.CvmClientProfile)
	if err != nil {
		utils.Logger.Error("init tencent client failed")
		return err
	}

	request := cvm.NewTerminateInstancesRequest()

	request.InstanceIds = []*string{&instanceId}

	_, err = client.TerminateInstances(request)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("delete instance [%s] failed", instanceId))
		return err
	}
	return nil
}

func (m *TenCentMgr) InstanceStatusTransform(status string) string {
	InstanceStatusMap := map[string]string{
		"PENDING":       cloud.StatusPending,
		"CREATE_FAILED": cloud.StatusLaunchFailed,
		"RUNNING":       cloud.StatusRunning,
		"STOPPED":       cloud.StatusStopped,
		"STARTING":      cloud.StatusStarting,
		"STOPPING":      cloud.StatusStopping,
		"REBOOTING":     cloud.StatusRebooting,
		"SHUTDOWN":      cloud.StatusShutdown,
		"DELETING":      cloud.StatusDeleting,
	}
	if s, ok := InstanceStatusMap[status]; ok {
		return s
	}
	return cloud.StatusUnknown
}
