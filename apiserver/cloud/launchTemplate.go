package cloud

import "sync"

type LaunchTemplateSet struct {
	mutex        sync.RWMutex
	TemplateSets []*LaunchTemplate
}

type LaunchTemplate struct {
	VersionNumber               int
	LaunchTemplateId            string
	LaunchTemplateName          string
	UpdatedTime                 string
	CreatedTime                 string
	ImageId                     string
	SecurityGroupId             string
	Tags                        string
	DataDisks                   string
	EnhancedService string
	InternetChargeType          string
	InstanceType                string
	InstanceChargeType          string
	HostName                    string
	InternetMaxBandwidthOut     int
	SystemDisk                  string
	InstanceName                string
	Description                 string
}
