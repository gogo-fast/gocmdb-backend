package utils

type HostStatus int8

const (
	HOST_RUNNING HostStatus = iota
	HOST_STARTING
	HOST_STOPPING
	HOST_STOPPED
	HOST_MAINTAINING
	HOST_DELETED
	HOST_UNKNOWN
)

func HostStatusToInt(hostType HostStatus) int {
	switch hostType {
	case HOST_RUNNING:
		return 0
	case HOST_STARTING:
		return 1
	case HOST_STOPPING:
		return 2
	case HOST_STOPPED:
		return 3
	case HOST_MAINTAINING:
		return 4
	case HOST_DELETED:
		return 5
	case HOST_UNKNOWN:
		return 6
	default:
		return 6
	}
}

func IntToHostStatus(hostType int) HostStatus {
	switch hostType {
	case 0:
		return HOST_RUNNING
	case 1:
		return HOST_STARTING
	case 2:
		return HOST_STOPPING
	case 3:
		return HOST_STOPPED
	case 4:
		return HOST_MAINTAINING
	case 5:
		return HOST_DELETED
	case 6:
		return HOST_UNKNOWN
	default:
		return HOST_UNKNOWN
	}
}

func StrToHostStatus(hostType string) HostStatus {
	switch hostType {
	case "running":
		return HOST_RUNNING
	case "starting":
		return HOST_STARTING
	case "stopping":
		return HOST_STOPPING
	case "stopped":
		return HOST_STOPPED
	case "maintaining":
		return HOST_MAINTAINING
	case "deleted":
		return HOST_DELETED
	case "unknown":
		return HOST_UNKNOWN
	default:
		return HOST_UNKNOWN
	}
}

func HostStatusToStr(hostType HostStatus) string {
	switch hostType {
	case HOST_RUNNING:
		return "running"
	case HOST_STARTING:
		return "starting"
	case HOST_STOPPING:
		return "stopping"
	case HOST_STOPPED:
		return "stopped"
	case HOST_MAINTAINING:
		return "maintaining"
	case HOST_DELETED:
		return "deleted"
	case HOST_UNKNOWN:
		return "unknown"
	default:
		return "unknown"
	}
}
