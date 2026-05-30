package netfinder

import (
	"net"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

// 基础常量
const (
	broadcastPort = "1999" // 广播端口
)

var (
	broadcastListenAddr, _ = net.ResolveUDPAddr("", ":"+broadcastPort)                // 本地监听
	broadcastAddr, _       = net.ResolveUDPAddr("", "255.255.255.255:"+broadcastPort) // 广播公用地址
	id, _                  = gonanoid.New(8)                                          // 本机id
	ptopTcpAddr, _         = net.ResolveTCPAddr("", ":0")                             // 点对点通信的tcp地址
)
