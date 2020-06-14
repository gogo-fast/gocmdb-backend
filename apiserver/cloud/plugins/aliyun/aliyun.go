package aliyun

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"gogo-cmdb/apiserver/utils"
)

type Instance struct {
	RequestId               string
	InstanceId              string
	HostName                string
	RegionId                string
	ZoneId                  string
	Status                  string
	Cpu                     int
	Memory                  int
	InstanceType            string
	CreationTime            string
	Description             string
	InternetChargeType      string
	VpcAttributes           string
	StoppedMode             string
	SerialNumber            string
	VlanId                  string
	EipAddress              string
	InnerIpAddress          string
	ImageId                 string
	PublicIpAddress         string
	InstanceNetworkType     string
	InternetMaxBandwidthIn  string
	InternetMaxBandwidthOut string
}

type AliMgr struct {
	PlatName string
	Regions  []Region
}

type Region struct {
	RegionId       string
	Status         string
	LocalName      string
	RegionEndpoint string
	Zones          []*Zone
}

type Zone struct {
	ZoneId    string
	LocalName string
}

func (m *AliMgr) Name() string {
	return "阿里云"
}
func (m *AliMgr) PlatType() string {
	return "aliyun"
}
func (m *AliMgr) TestConn() error {
	deafult_region_id := utils.GlobalConfig.Section("aliyun").Key("default_region_id").String()
	client, err := ecs.NewClientWithAccessKey(deafult_region_id, "LTAI4G5eagkrwTs1Rcd8pZHz", "8F1k3W2CSkW09DfvKHeS9QmhigMtrF")
	if err != nil {
		fmt.Println("init client failed")
		return err
	}

	request := ecs.CreateDescribeRegionsRequest()
	request.Scheme = "https"

	_, err = client.DescribeRegions(request)
	if err != nil {
		fmt.Print(err.Error())
		return err
	}
	return nil
}
func (m *AliMgr) Init(AccessKeyId, AccessKeySecret string) {

}
func (m *AliMgr) CreateInstance(regionId, uuid string) (string, error) {

}
func (m *AliMgr) RunInstance(regionId, uuid string) (string, error) {

}
func (m *AliMgr) GetInstanceAttribute(regionId, instanceId string) (*Instance, error) {

}
func (m *AliMgr) GetInstanceList(regionId string) ([]*Instance, error) {

}

func (i *Instance) StartInstance(instanceId string) error {

}
func (i *Instance) StopInstance(instanceId string) error {

}
func (i *Instance) RebootInstance(instanceId string) error {

}
func (i *Instance) DeleteInstance(instanceId string) error {

}
