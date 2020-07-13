package hosts

import (
	"fmt"
	"gogo-cmdb/apiserver/models"
	"gogo-cmdb/apiserver/utils"
	"sync"
	"time"
)

type HostMgr struct {
	mutex              sync.RWMutex
	Hosts              map[string]*models.Host
	HostsChan          chan *models.Host
	OffsetLineHostChan chan *models.Host
}

func NewHostMgr() *HostMgr {
	return &HostMgr{
		mutex:              sync.RWMutex{},
		Hosts:              make(map[string]*models.Host, 1000),
		HostsChan:          make(chan *models.Host, 1000),
		OffsetLineHostChan: make(chan *models.Host, 1000),
	}
}

func (h *HostMgr) Register(rgMsg *models.RegisterMsg) {
	h.mutex.Lock()
	if _, ok := h.Hosts[rgMsg.UUID]; !ok {
		h.Hosts[rgMsg.UUID] = &models.Host{}
	}
	hostPtr := h.Hosts[rgMsg.UUID]
	h.mutex.Unlock()

	hostPtr.UUID = rgMsg.UUID
	hostPtr.Hostname = rgMsg.Hostname

	if rgMsg.OutBoundIP.String == "" {
		hostPtr.OutBoundIP.Valid = false
	} else {
		hostPtr.OutBoundIP.Valid = true
	}
	hostPtr.OutBoundIP.String = rgMsg.OutBoundIP.String

	if rgMsg.ClusterIP.String == "" {
		hostPtr.ClusterIP.Valid = false
	} else {
		hostPtr.ClusterIP.Valid = true
	}
	hostPtr.ClusterIP.String = rgMsg.ClusterIP.String

	hostPtr.IPs = rgMsg.IPs
	hostPtr.OS = rgMsg.OS
	hostPtr.Arch = rgMsg.Arch
	hostPtr.CpuCount = rgMsg.CpuCount
	hostPtr.CpuUsePercent = rgMsg.CpuUsePercent
	hostPtr.RamTotal = rgMsg.RamTotal
	hostPtr.RamUsed = rgMsg.RamUsed
	hostPtr.RamUsePercent = rgMsg.RamUsePercent
	hostPtr.Disks = rgMsg.Disks
	hostPtr.AvgLoad = rgMsg.AvgLoad
	hostPtr.BootTime = rgMsg.BootTime
}

var hbInterval = utils.GlobalConfig.GetInt64("heartbeat.interval")

func (h *HostMgr) HearBeat(hbMsg *models.HeartBeatMsg) {
	h.mutex.RLock()
	if host, ok := h.Hosts[hbMsg.UUID]; ok {
		host.HeartbeatTime = hbMsg.Timestamp
		if time.Now().Unix() > hbMsg.Timestamp.Int64+2*hbInterval {
			host.IsOnline = false
		} else {
			host.IsOnline = true
		}
	}
	h.mutex.RUnlock()
}

func (h *HostMgr) LoadHostsFromDB() {
	hosts, err := models.DefaultHostManager.GetAllHosts()
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("get host list failed, %s", err))
		return
	}
	for _, host := range hosts {
		h.HostsChan <- host
		h.mutex.Lock()
		h.Hosts[host.UUID] = host
		h.mutex.Unlock()
	}
}

func (h *HostMgr) DeleteHost(uuid string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	delete(h.Hosts, uuid)
}

func (h *HostMgr) Run() {

	interval := utils.GlobalConfig.GetInt64("update_host_record.interval")
	h.LoadHostsFromDB()

	go func() {
		for {
			h.mutex.RLock()
			for _, host := range h.Hosts {
				if time.Now().Unix() > host.HeartbeatTime.Int64+2*interval || time.Now().Unix() < host.HeartbeatTime.Int64-2*interval {
					host.Status = utils.HOST_UNKNOWN
					host.IsOnline = false
					h.OffsetLineHostChan <- host
				} else {
					host.Status = utils.HOST_RUNNING
					host.IsOnline = true
				}
			}
			h.mutex.RUnlock()
			time.Sleep(time.Second * time.Duration(interval))
		}
	}()

	go func() {
		for {
			host := <-h.OffsetLineHostChan
			err := models.DefaultHostManager.UpdateHostStatusByUUID(host.UUID, utils.HOST_UNKNOWN)
			if err != nil {
				utils.Logger.Error(err)
			}
		}
	}()

	go func(hc <-chan *models.Host) {
		for v := range hc {
			_, err := models.DefaultHostManager.GetHostRecordByUUID(v.UUID)
			if err != nil {
				_, err = models.DefaultHostManager.AddHostRecord(v)
				if err != nil {
					utils.Logger.Error(fmt.Sprintf("add host record failed"))
				}
			} else {
				err := models.DefaultHostManager.UpdateHostRecordByUUID(v)
				if err != nil {
					utils.Logger.Error(fmt.Sprintf("update host record failed"))
				}
			}
		}
	}(h.HostsChan)

	go func(hc chan<- *models.Host) {
		for {
			h.mutex.RLock()
			for _, v := range h.Hosts {
				hc <- v
			}
			h.mutex.RUnlock()
			time.Sleep(time.Second * time.Duration(interval))
		}
	}(h.HostsChan)
}

var DefaultHostManager = NewHostMgr()

func init() {
	DefaultHostManager.Run()
}
