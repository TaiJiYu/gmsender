package netfinder

import (
	"fmt"
	"net"
)

var (
	selfIp      = ""
	selfTcpPort = 0 // 用于通信的端口
)

// 由系统分配
func getLocalTcpPort() int {
	return selfTcpPort
}

// 获取本机ip
func getLocalIp() string {
	if selfIp != "" {
		return selfIp
	}
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return backGetLoackIp()
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	if localAddr.IP.IsPrivate() {
		selfIp = localAddr.IP.String()
		return selfIp
	}
	return backGetLoackIp()
}

// 备用方案获取本机ip
func backGetLoackIp() string {
	// 备用方案
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		// 检查 IP 地址，过滤掉回环和链路本地地址
		if ipNet, ok := addr.(*net.IPNet); ok && ipNet.IP.IsPrivate() {
			selfIp = ipNet.IP.String()
			return selfIp
		}
	}
	panic(fmt.Errorf("没有可用ip"))
}
