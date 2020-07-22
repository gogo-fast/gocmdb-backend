package tencent

import (
	"apiserver/cloud"
	"apiserver/utils"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"sync"
)

type TenCentMgr struct {
	mutex                sync.RWMutex
	CloudName            string
	CloudType            string
	DefaultRegionId      string
	AccessKeyId          string
	AccessKeySecret      string
	Page                 int
	Size                 int
	Credential           *common.Credential
	CvmClientProfile     *profile.ClientProfile
	MonitorClientProfile *profile.ClientProfile
}

func NewTenCentMgr() *TenCentMgr {
	var TenCentMgr = &TenCentMgr{}
	TenCentMgr.CloudName = TenCentMgr.Name()
	TenCentMgr.CloudType = TenCentMgr.PlatType()
	return TenCentMgr
}

func (m *TenCentMgr) Name() string {
	return "腾讯云"
}

func (m *TenCentMgr) PlatType() string {
	return "tencent"
}

func (m *TenCentMgr) Init(DefaultRegionId, AccessKeyId, AccessKeySecret string) {
	m.DefaultRegionId = DefaultRegionId
	m.AccessKeyId = AccessKeyId
	m.AccessKeySecret = AccessKeySecret
	credential := common.NewCredential(m.AccessKeyId, m.AccessKeySecret)
	m.Credential = credential
	cvmCpf := profile.NewClientProfile()
	cvmCpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"
	m.CvmClientProfile = cvmCpf
	monitorCpf := profile.NewClientProfile()
	monitorCpf.HttpProfile.Endpoint = "monitor.tencentcloudapi.com"
	m.MonitorClientProfile = monitorCpf
}

func (m *TenCentMgr) TestConn() error {
	_, err := m.GetRegions()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	defaultRegionId := utils.GlobalConfig.GetString("tencent.default_region_id")
	accessKeyId := utils.GlobalConfig.GetString("tencent.access_key_id")
	accessKeySecret := utils.GlobalConfig.GetString("tencent.access_key_secret")

	TenCentMgr := NewTenCentMgr()
	TenCentMgr.Init(defaultRegionId, accessKeyId, accessKeySecret)
	err := TenCentMgr.TestConn()
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("connect to [%s] failed, err: %s", TenCentMgr.CloudType, err))
	} else {
		utils.Logger.Info(fmt.Sprintf("connect to [%s] success", TenCentMgr.CloudType))
	}

	cloud.DefaultCloudMgr.Register(TenCentMgr)
}
