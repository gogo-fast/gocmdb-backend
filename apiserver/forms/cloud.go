package forms

type DataDisk struct {
	Size     string
	Category string
}

type SystemDisk struct {
	Size     int
	Category string
}

type Tag struct {
	Key   string
	Value string
}

type InstanceTemplate struct {
	ClientToken     string // 随机生成
	RegionId        string //     动态查询
	ImageId         string // 手动查询
	InstanceType    string // 手动查询
	SecurityGroupId string // 手动查询

	VSwitchId               string //     动态查询
	Description             string // 可配置
	InternetMaxBandwidthIn  int    // 可配置
	InternetMaxBandwidthOut int    // 可配置

	HostName           string // 可配置
	Password           string // 可配置
	ZoneId             string // 可配置
	InternetChargeType string // PayByBandwidth：按固定带宽计费。	PayByTraffic（默认）

	SystemDisk   SystemDisk  // 可配置
	DataDisk     []*DataDisk // 可配置
	IoOptimized  string      // 可选 none / optimized
	SpotStrategy string      // 可选 NoSpot（默认）SpotWithPriceLimit：设置上限价格的抢占式实例。SpotAsPriceGo：系统自动出价，跟随当前市场实际价格。

	Tag                []*Tag // 可配置
	ResourceGroupId    string // 手动查询配置
	InstanceChargeType string // PrePaid：包年包月。PostPaid（默认）：按量付费。
	PrivateIpAddress   string // 实例私网IP地址。该IP地址必须为交换机（VSwitchId）网段的空闲地址。

	CreditSpecification string // Standard：标准模式 Unlimited：无性能约束模式
	Tenancy             string // default：在非专有宿主机上创建实例。 host：在专有宿主机上创建实例。若您不指定DedicatedHostId，则由阿里云自动选择专有宿主机部署实例 默认值：default。
}
