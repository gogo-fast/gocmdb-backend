package cloud

type PlatformMgr interface {
	Name() string
	PlatType() string
	TestConn() error
	Init(AccessKeyId, AccessKeySecret string)
	CreateInstance(regionId, uuid string) (string, error)
	RunInstance(regionId, uuid string) (string, error)
	GetInstanceAttribute(regionId, instanceId string) (*Instance, error)
	GetInstanceList(regionId string) ([]*Instance, error)
}

type Instance interface {
	StartInstance(instanceId string) error
	StopInstance(instanceId string) error
	RebootInstance(instanceId string) error
	DeleteInstance(instanceId string) error
}
