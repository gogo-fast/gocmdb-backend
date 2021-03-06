package register

import (
	"agent/models"
	"agent/utils"
	"encoding/json"
	"fmt"
	"github.com/imroc/req"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"strings"
	"time"
)

func NewRegister(uuid string) (*models.RegisterMsg, error) {

	register := &models.RegisterMsg{}
	var _cpuCount int32 = 0
	_ips := make([]string, 0, 10)
	_disks := map[string]*models.Disk{}
	_cpuInfo, err := cpu.Info()
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	_cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	_interfaceStats, err := net.Interfaces()
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	_hostInfo, err := host.Info()
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	_memory, err := mem.VirtualMemory()
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	_partitionStats, err := disk.Partitions(true)
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	// not implemented yet on windows
	_loadAvgStat, err := load.Avg()
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	_loads, err := json.Marshal(_loadAvgStat)
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	for _, c := range _cpuInfo {
		_cpuCount += c.Cores
	}
	for _, InterfaceStat := range _interfaceStats {
		for _, addr := range InterfaceStat.Addrs {
			if strings.Index(addr.Addr, ":") >= 0 {
				continue
			}
			if strings.Index(addr.Addr, "127.") == 0 {
				continue
			}
			ip := strings.Split(addr.Addr, "/")[0]
			_ips = append(_ips, ip)
		}
	}
	ips, err := json.Marshal(_ips)
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	for _, pts := range _partitionStats {
		// windows user pts.Device, linux use pts.Mountpoint
		usageStats, err := disk.Usage(pts.Mountpoint)
		if err != nil {
			utils.Logger.Error(err)
			return nil, err
		}
		_disk := models.Disk{
			Total:      int64(usageStats.Total / 1024 / 1024 / 1024), // GB
			Used:       int64(usageStats.Used / 1024 / 1024 / 1024),  // GB
			UsePercent: usageStats.UsedPercent,
		}
		_disks[usageStats.Path] = &_disk
	}
	disks, err := json.Marshal(_disks)
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}

	register.UUID = uuid
	register.Hostname = _hostInfo.Hostname
	register.IPs = string(ips)
	register.OS = _hostInfo.OS
	register.CpuCount = _cpuCount
	register.CpuUsePercent = _cpuPercent[0]
	register.BootTime.Int64 = int64(_hostInfo.BootTime)
	register.BootTime.Valid = true
	register.Arch = _hostInfo.KernelArch
	register.RamTotal = int64(_memory.Total / 1024 / 1024) // MB
	register.RamUsed = int64(_memory.Used / 1024 / 1024)   // MB
	register.RamUsePercent = _memory.UsedPercent
	register.Disks = string(disks)
	register.AvgLoad = string(_loads)
	return register, nil
}

func Run() {
	url := fmt.Sprintf("%s/%s", utils.GlobalConfig.ApiServerUrl, "register")
	tokenStr := utils.GlobalConfig.Token
	uuid := utils.GlobalConfig.UUID
	interval := utils.GlobalConfig.RegisterInterval

	outBoundIp := utils.GetOutBoundIp()
	clusterIp := utils.GetAgentIp()

	for {
		register, err:= NewRegister(uuid)
		if err != nil {
			utils.Logger.Error(fmt.Sprintf("get register error, %s",err))
			time.Sleep(time.Second * time.Duration(interval))
			continue
		}
		register.OutBoundIP.String = outBoundIp
		if outBoundIp == "" {
			register.OutBoundIP.Valid = false
		} else {
			register.OutBoundIP.Valid = true
		}
		register.ClusterIP.String = clusterIp
		if clusterIp == "" {
			register.ClusterIP.Valid = false
		} else {
			register.ClusterIP.Valid = true
		}

		params := req.Param{"token": tokenStr}
		resp, err := req.Post(url, params, req.BodyJSON(&register))
		if err != nil {
			utils.Logger.Error(fmt.Sprintf("register agent failed, %s", err))
		} else {
			result := map[string]interface{}{}
			resp.ToJSON(&result)
			utils.Logger.Info(result["msg"])
		}
		time.Sleep(time.Second * time.Duration(interval))
	}

}
