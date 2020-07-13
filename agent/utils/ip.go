package utils

import (
	"fmt"
	"net"
	"strings"
)

func GetOutBoundIp() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		Logger.Error(err)
		return "", nil
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx], nil
}

func GetAgentIp() (string, error) {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", GlobalConfig.ApiServerHost, GlobalConfig.ApiServerPort))
	if err != nil {
		Logger.Error(err)
		return "", nil
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx], nil
}
