package utils

import (
	"fmt"
	"gogo-cmdb/apiserver/utils"
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
	serverHost :=  utils.GlobalConfig.GetString("server.host")
	serverPort :=  utils.GlobalConfig.GetString("server.port")
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%s", serverHost, serverPort))
	if err != nil {
		Logger.Error(err)
		return "", nil
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx], nil
}
