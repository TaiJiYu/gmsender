package netfinder

import (
	"net"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

// 基础常量

var (
	multicastAddr, _ = net.ResolveUDPAddr("", "239.198.239.137:1999") // 组播公用地址
	id, _            = gonanoid.New(8)                                // 本机id
	ptopTcpAddr, _   = net.ResolveTCPAddr("", ":0")                   // 点对点通信的tcp地址
)
