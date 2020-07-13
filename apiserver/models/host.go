package models

import (
	"apiserver/utils"
	"time"
)

type Host struct {
	Id            int              `db:"id" json:"id,string"`
	IsOnline      bool             `db:"is_online" json:"isOnline"`
	Status        utils.HostStatus `db:"host_status" json:"hostStatus"`
	UUID          string           `db:"uuid" json:"uuid"`
	Hostname      string           `db:"hostname" json:"hostname"`
	OutBoundIP    NullString       `db:"out_bound_ip" json:"outBoundIp"`
	ClusterIP     NullString       `db:"cluster_ip" json:"clusterIp"`
	IPs           string           `db:"ips" json:"ips"`
	OS            string           `db:"os" json:"os"`
	Arch          string           `db:"arch" json:"arch"`
	CpuCount      int32            `db:"cpu_count" json:"cpuCount"`
	CpuUsePercent float64          `db:"cpu_usage_percent" json:"cpuUsePercent"`
	RamTotal      int64            `db:"ram_total" json:"ramTotal"` // MB
	RamUsed       int64            `db:"ram_usage" json:"ramUsed"`  // MB
	RamUsePercent float64          `db:"ram_usage_percent" json:"ramPercent"`
	Disks         string           `db:"disks" json:"disks"`
	AvgLoad       string           `db:"avg_load" json:"avgLoad"` // do not use "load" in db, use avg_load instead.
	BootTime      NullInt64        `db:"boot_time" json:"bootTime"`
	CreateTime    NullInt64        `db:"create_time" json:"createTime"`
	UpdateTime    NullInt64        `db:"update_time" json:"updateTime"`
	DeleteTime    NullInt64        `db:"delete_time" json:"deleteTime"`
	HeartbeatTime NullInt64        `db:"heartbeat_time" json:"-"`
}

type HostManager struct{}

func NewHostManager() *HostManager {
	return &HostManager{}
}

var DefaultHostManager *HostManager

func init() {
	DefaultHostManager = NewHostManager()
}

func (h *HostManager) GetAllHosts() ([]*Host, error) {
	sqlHosts := `select id, is_online, host_status, uuid, hostname, out_bound_ip, cluster_ip, ips, os, arch, cpu_count, 
			cpu_usage_percent, ram_total, ram_usage, ram_usage_percent, disks,
			avg_load, boot_time, heartbeat_time, create_time, update_time from hosts where host_status != ?;`
	var hostList []*Host

	err := db.Select(&hostList, sqlHosts, utils.HOST_DELETED)
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	return hostList, nil
}

func (h *HostManager) GetHostRecordList(page, size int) (int, []*Host, *utils.Pagination, error) {
	sqlCount := `select count(uuid) from hosts where host_status != ?;`
	sqlHosts := `select id, is_online, host_status, uuid, hostname, out_bound_ip, cluster_ip, ips, os, arch, cpu_count, 
			cpu_usage_percent, ram_total, ram_usage, ram_usage_percent, disks,
			avg_load, boot_time, heartbeat_time, create_time, update_time from hosts where host_status != ? order by id desc limit ?,?;`
	var total int
	var hostList []*Host

	err := db.Get(&total, sqlCount, utils.HOST_DELETED)
	if err != nil {
		utils.Logger.Error(err)
		return 0, nil, nil, err
	}
	pagination, err := utils.NewPagination(total, page, size, 5)
	if err != nil {
		utils.Logger.Error(err)
		return 0, nil, nil, err
	}
	err = db.Select(&hostList, sqlHosts, utils.HOST_DELETED, pagination.CurrentPage.Offset, pagination.CurrentPage.Limit)
	if err != nil {
		utils.Logger.Error(err)
		return 0, nil, nil, err
	}
	return total, hostList, pagination, nil
}

func (h *HostManager) GetHostRecordByUUID(uuid string) (*Host, error) {
	sql := `select id, is_online, host_status, uuid, hostname, out_bound_ip, cluster_ip, ips, os, arch, cpu_count, 
			cpu_usage_percent, ram_total, ram_usage, ram_usage_percent, disks,
			avg_load, boot_time, heartbeat_time, create_time, update_time from hosts where uuid = ?`
	var host Host
	err := db.Get(&host, sql, uuid)
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	return &host, nil
}

func (h *HostManager) AddHostRecord(host *Host) (int64, error) {
	sql := `insert into hosts(is_online, host_status, uuid, hostname, out_bound_ip, cluster_ip, ips, os, arch, cpu_count, 
						cpu_usage_percent, ram_total, ram_usage, ram_usage_percent, disks,
						avg_load, boot_time, heartbeat_time, create_time) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`
	result, err := db.Exec(sql, host.IsOnline, utils.HOST_RUNNING, host.UUID, host.Hostname,
		host.OutBoundIP.String, host.ClusterIP.String, host.IPs, host.OS, host.Arch, host.CpuCount,
		host.CpuUsePercent, host.RamTotal, host.RamUsed, host.RamUsePercent, host.Disks,
		host.AvgLoad, host.BootTime.Int64, host.HeartbeatTime.Int64, time.Now().Unix())
	if err != nil {
		utils.Logger.Error(err)
		return -1, err
	}
	return result.LastInsertId()
}

func (h *HostManager) UpdateHostRecordByUUID(host *Host) error {
	sql := `update hosts set is_online=?, host_status=?, uuid=?, hostname=?, out_bound_ip=?, cluster_ip=?, ips=?, os=?, arch=?, cpu_count=?, 
						cpu_usage_percent=?, ram_total=?, ram_usage=?, ram_usage_percent=?, disks=?,
						avg_load=?, boot_time=?, heartbeat_time=?, update_time=? where uuid=?`
	_, err := db.Exec(sql, host.IsOnline, host.Status, host.UUID, host.Hostname,
		host.OutBoundIP.String, host.ClusterIP.String, host.IPs, host.OS, host.Arch, host.CpuCount,
		host.CpuUsePercent, host.RamTotal, host.RamUsed, host.RamUsePercent, host.Disks, host.AvgLoad,
		host.BootTime.Int64, host.HeartbeatTime.Int64, time.Now().Unix(), host.UUID)
	if err != nil {
		utils.Logger.Error(err)
		return err
	}
	return nil
}

func (h *HostManager) UpdateHostStatusByUUID(uuid string, status utils.HostStatus) error {
	sql := `update hosts set host_status=?, update_time=? where uuid=?`
	_, err := db.Exec(sql, status, time.Now().Unix(), uuid)
	if err != nil {
		utils.Logger.Error(err)
		return err
	}
	return nil
}

func (h *HostManager) DeleteHostRecordByUUID(uuid string) error {
	sql := `update hosts set host_status=?, delete_time=? where uuid=?`
	_, err := db.Exec(sql, utils.HOST_DELETED, time.Now().Unix(), uuid)
	if err != nil {
		utils.Logger.Error(err)
		return err
	}
	return nil
}
