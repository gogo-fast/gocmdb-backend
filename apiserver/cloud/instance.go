package cloud

const (
	StatusPending = "PENDING"
	StatusLaunchFailed = "CREATE_FAILED"
	StatusRunning = "RUNNING"
	StatusStopped = "STOPPED"
	StatusStarting = "STARTING"
	StatusStopping = "STOPPING"
	StatusRebooting = "REBOOTING"
	StatusShutdown = "SHUTDOWN"
	StatusDeleting = "DELETING"
	StatusUnknown = "UNKNOWN"
)


type Region struct {
	RegionId   string
	RegionName string
}

type Zone struct {
	ZoneId   string
	ZoneName string
}

type DataDisk struct {
	Size     string
	Category string
}

type Tag struct {
	Key   string
	Value string
}

type SystemDisk struct {
	Size     int
	Category string
}

type SecurityGroup struct {
	SecurityGroupId   string
	SecurityGroupName string
	VpcId             string
}

type Instance struct {
	InstanceId              string
	Uuid                    string
	HostName                string
	RegionId                string
	ZoneId                  string
	Status                  string
	OSName                  string
	Cpu                     int // vCpu count
	Memory                  int
	InstanceType            string
	CreatedTime             string
	Description             string
	InternetChargeType      string
	VpcId                   string // aliyun in VpcAttributes
	PrivateIpAddress        string // aliyun in VpcAttributes
	PublicIpAddress         string
	InternetMaxBandwidthIn  int
	InternetMaxBandwidthOut int
}

type InstanceType struct {
	InstanceType        string
	Memory              int
	GPUAmount           int
	Cores               int
	CpuCoreCount        int
	MemorySize          int
	InstanceBandwidthRx int
	InstanceBandwidthTx int
	InstanceFamilyLevel string
	InstanceTypeFamily  string
}

type InstanceStaus struct {
	InstanceId    string
	InstanceState string
}


