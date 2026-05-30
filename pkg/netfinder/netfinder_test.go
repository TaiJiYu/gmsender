package netfinder

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"testing"
)

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

func TestBuffer(t *testing.T) {

	b := bytes.NewBuffer(make([]byte, 10))

	b.Reset()
	fmt.Println(b.Len(), b.Cap())
	b.Write(make([]byte, 20))
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
