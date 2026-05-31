package netfinder

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"testing"
)

func TestFileName(t *testing.T) {
	fileS := File{
		Ip:       getLocalIp(),
		Port:     fmt.Sprintf("%v", getLocalTcpPort()),
		Id:       Id(),
		FileName: "E:\\GMStudioProject\\GMSender\about\\全国计算机技术与软件专业技术资格（水平）考试.pdf",
	}
	fmt.Println(publicSelfFileBytes(fileS))
}
func TestMask(t *testing.T) {
	fmt.Printf("%08b\n", afterMask)
	fmt.Printf("%08b\n", netOrderMask)
}
func TestNet(t *testing.T) {
	castAddr, _ := net.ResolveUDPAddr("", "255.255.255.255:1999") // 广播公用地址
	listenAddr, _ := net.ResolveUDPAddr("", ":1999")

	// 发送端
	sendConn, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		panic(err)
	}
	// 接收句柄
	cieveConn, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		panic(err)
	}

	// 接收
	go func() {
		buf := make([]byte, 128)
		n, err := cieveConn.Read(buf)
		fmt.Println("接受：", buf[:n], err)
	}()

	// 发送
	go func() {
		n, err := sendConn.WriteToUDP([]byte("hello"), castAddr)
		fmt.Println("发送：", n, err)
	}()

	<-make(chan bool)

}

func TestMulNet(t *testing.T) {
	castAddr, _ := net.ResolveUDPAddr("", "224.0.0.251:1999") // 广播公用地址
	// listenAddr, _ := net.ResolveUDPAddr("", ":1999")

	// 发送端
	sendConn, err := net.ListenMulticastUDP("udp", nil, castAddr)
	if err != nil {
		panic(err)
	}
	// // 接收句柄
	// cieveConn, err := net.ListenUDP("udp", listenAddr)
	// if err != nil {
	// 	panic(err)
	// }

	// 接收
	go func() {
		buf := make([]byte, 128)
		n, err := sendConn.Read(buf)
		fmt.Println("接受：", buf[:n], err)
	}()

	// 发送
	go func() {
		// 需要多个"连接"时，创建新的 UDPConn
		sendConn, _ := net.DialUDP("udp", nil, castAddr)
		n, err := sendConn.Write([]byte("hello")) // 专门用于组播发送
		fmt.Println("发送：", n, err)
	}()

	<-make(chan bool)

}

func TestCode(t *testing.T) {
	bs := askMasterBytes()
	fmt.Println("ask 长度:", len(bs), "B")
	for _, b := range bs {
		fmt.Printf("%08b ", b)
	}
	bss := masterAnswerBytes()
	fmt.Println("\nanswerk 长度:", len(bss), "B")
	for _, b := range bss {
		fmt.Printf("%08b ", b)
	}
	fmt.Println("")

	decode(bs)

	decode(bss)
}

func TestChan(t *testing.T) {
	m := make(chan bool, 1)
	// m = nil
	select {
	case <-m:
		fmt.Println("截断")
		return
	default:
	}
	// m <- false

	{
	loop:
		for i := 0; i < 2; i++ {
			fmt.Println("第", i, "次循环")
			select {
			case b, ok := <-m:
				if ok {
					fmt.Println("收到消息：", b)
				} else {
					fmt.Println("通道关闭")
					break loop
				}
			default:
			}
			if i == 1 {
				fmt.Println("成为master")
			}
		}
	}
	fmt.Println("done")

}

func TestUDP(t *testing.T) {
	lister, err := net.Listen("tcp", ":0")
	// 端口写 0，系统自动分配
	// addr, err := net.ResolveTCPAddr("", ":0")
	if err != nil {
		panic(err)
	}
	defer lister.Close()

	// conn, err := net.ListenTCP("", addr)
	// if err != nil {
	// 	panic(err)
	// }
	// defer conn.Close()

	// 获取实际分配到的端口
	// realAddr := conn.LocalAddr().(*net.UDPAddr)
	// realAddr.IP.
	fmt.Printf("绑定到的端口: %v\n", lister.Addr())

}

func TestBool(t *testing.T) {
	a := atomic.Bool{}

	fmt.Println(a.CompareAndSwap(false, true), a.Load()) //true,true
	fmt.Println(a.CompareAndSwap(false, true), a.Load()) //false,true
}

func TestBuffer(t *testing.T) {

	b := bytes.NewBuffer(make([]byte, 2))
	o := bytes.NewReader([]byte{1, 2, 3})
	b.Reset()
	io.Copy(b, o)
	fmt.Println(b.Bytes())
	oo := bytes.NewReader([]byte{4, 5})
	b.Reset()
	io.Copy(b, oo)
	fmt.Println(b.Bytes())
	fmt.Println(b.Len(), b.Cap())

}

func TestDownLoad(t *testing.T) {
	b := downLoadFileBytes("https://chat.deepseek.com/a/chat/s/6fe57427-bd57-493c-ab34-c5aaeca4b877")
	buf := bytes.NewBuffer(make([]byte, 10))
	fmt.Println(len(b), b)
	fmt.Println(decodeDownloadFileInfo(b))
	buf.Reset()
	io.Copy(buf, bytes.NewReader(b))

	fmt.Println(buf.Len(), buf)
}
