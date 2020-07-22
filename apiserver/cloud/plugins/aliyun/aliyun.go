package aliyun

import (
	"apiserver/cloud"
	"apiserver/utils"
	"fmt"
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
