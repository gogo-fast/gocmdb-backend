package cloud

import (
	"sync"
)

type CCloud interface {
	Name() string
	PlatType() string
	Init(DefaultRegionId, AccessKeyId, AccessKeySecret string)
	TestConn() error
	GetRegions() ([]*Region, error)
	GetZones(regionId string) ([]*Zone, error)
	GetSecurityGroups(regionId string) ([]*SecurityGroup, error)
	GetInstancesStatus(regionId string, instanceIds []string) ([]*InstanceStaus, error)
	GetAllInstancesStatus(regionId string) ([]*InstanceStaus, error)
	GetInstance(regionId, instanceId string) (*Instance, error)
	GetInstanceListPerPage(regionId string, page, size int) ([]*Instance, int, error)
	GetAllInstance(regionId string) ([]*Instance, error)
	StartInstance(regionId, instanceId string) error
	StopInstance(regionId, instanceId string) error
	RebootInstance(regionId, instanceId string) error
	DeleteInstance(regionId, instanceId string) error
	InstanceStatusTransform(string) string
	GetMonitorDataOfInstances(regionId, metricName, startTime, endTime string, period int, instanceIDs []string) ([]*DataPoint, error)
}

type CCloudMgr struct {
	rwLock sync.RWMutex
	Clouds map[string]CCloud
}

func (cm *CCloudMgr) Register(c CCloud) {
	cm.rwLock.Lock()
	defer cm.rwLock.Unlock()
	cm.Clouds[c.PlatType()] = c
}

func (cm *CCloudMgr) GetCloud(platType string) (CCloud, bool) {
	cm.rwLock.RLock()
	defer cm.rwLock.RUnlock()
	cloud, ok := cm.Clouds[platType]
	return cloud, ok
}

func NewCloudMgr() *CCloudMgr {
	return &CCloudMgr{
		Clouds: make(map[string]CCloud),
	}
}

var DefaultCloudMgr = NewCloudMgr()
