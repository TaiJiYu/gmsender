package netfinder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// 服务发现
type finder struct {
	chanCloseOnce            sync.Once
	askMasterChan            chan time.Duration
	multicastListenerOutChan chan struct{}
	ptopTcpOutChan           chan struct{} // tcp链接管道
	isClose                  atomic.Bool

	lastWaitSec time.Duration // 剩余询问等待时间

	masterInfo baseInfo
	masterAddr *net.UDPAddr

	selfIsMaster atomic.Bool // 自己是否是master

	multicastCoon *net.UDPConn // 组播监听链接

	// nodes []baseInfo // 所有节点的信息，包含自身和master，所有节点与master同步
	files []file // 公开的文件列表
}

var (
	finderCli  *finder
	finderOnce sync.Once
)

func defaultFinder() *finder {
	finderOnce.Do(func() {
		finderCli = &finder{
			askMasterChan:            make(chan time.Duration, 1),
			multicastListenerOutChan: make(chan struct{}, 1),
			ptopTcpOutChan:           make(chan struct{}, 1),
			lastWaitSec:              3 * time.Second,
			files:                    make([]file, 0),
		}
		finderCli.pTopTcpListener()
		finderCli.multicastListener()
		finderCli.askMaster()
	})
	return finderCli
}

func (f *finder) closeNetFinder() {
	close(f.ptopTcpOutChan)
	if f.multicastCoon != nil {
		f.multicastCoon.Close()
	}
}

// 自身的点对点通信监听
func (f *finder) pTopTcpListener() {
	// 端口写 0，系统自动分配
	lister, err := net.Listen("tcp", fmt.Sprintf("%s:0", getLocalIp()))
	if err != nil {
		panic(err)
	}

	selfTcpPort = lister.Addr().(*net.TCPAddr).Port
	go func() {
		for {
			select {
			case <-f.ptopTcpOutChan:
				// 收到退出请求
				lister.Close()
				return
			default:
			}
			conn, err := lister.Accept()
			if err != nil {
				// 建联失败
				continue
			}
			go f.handlerDownLoad(conn)

		}
	}()
}

// 收到其他用户的下载请求
func (f *finder) handlerDownLoad(conn net.Conn) {
	// buf := make([]byte, 32*1024)
	buf := bytes.NewBuffer(make([]byte, 128))
	buf.Reset()

	if _, err := io.Copy(buf, conn); err != nil {
		conn.Close()
		return
	}

	// 收到了文件请求
	file := decodeDownloadFileInfo(buf.Bytes())

	fileS, err := os.Open(file)
	if err != nil {
		// 文件错误
		conn.Close()
		return
	}
	io.Copy(conn, fileS)
	conn.Close()
	fileS.Close()
}

const (
	retryTimesMax = 5 // 最多重试5次
)

// 下载文件
// dstFileName要下载的目标文件名称
// saveFileName 保存文件名
func (f *finder) downLoadFile(dstinfo baseInfo, dstFileName, saveFileName string, errChan chan error) {
	go func() {
		for i := 0; i < retryTimesMax; i++ {
			// 点对点链接
			conn, err := net.Dial("tcp", dstinfo.addr())
			if err != nil {
				// 链接出错，等待1-3秒后重试
				time.Sleep(time.Duration(rand.IntN(3)+1) * time.Second)
				continue
			}

			// 请求下载文件
			for j := 0; j < retryTimesMax; j++ {
				message := downLoadFileBytes(dstFileName)
				_, err = conn.Write(message)
				if err != nil {
					// 接受失败，等会重试
					time.Sleep(time.Duration(rand.IntN(3)+1) * time.Second)
					continue
				}

				// 接收文件
				file, err := os.Create(saveFileName)
				if err != nil {
					errChan <- err
					break
				}
				if _, err := io.Copy(file, conn); err != nil {
					file.Close()
					os.Remove(saveFileName)
				} else {
					file.Close()
				}
				break
			}

			conn.Close()
			break
		}
	}()
}

// 组播监听
func (f *finder) multicastListener() {
	// 绑定本地端口接收广播
	conn, err := net.ListenMulticastUDP("udp", nil, multicastAddr)
	if err != nil {
		panic(err)
	}
	f.multicastCoon = conn
	// 监听线程
	go func() {
		buf := make([]byte, 128)
		{
		loop:
			for {
				n, _, err := conn.ReadFromUDP(buf)
				if err != nil {
					continue
				}
				if decode(buf[:n]) {
					// 收到master回复并退出监听，自己不是master，但保留通信组播监听，用于接受之后的组播消息
					f.closeAsk()
					f.beNode()
					break
				}
				select {
				case <-f.multicastListenerOutChan:
					// 没收到master回复但是收到了退出信号，自己成为master
					f.closeAsk()
					break loop
				default:
				}
			} // loop结束
		}
	}()
}

// 询问master线程
func (f *finder) askMaster() {
	// 询问线程
	go func() {
		{
		loop:
			for i := 0; i < 3; i++ {
				f.multicastCoon.Write(askMasterBytes())
				time.Sleep(time.Second)
				f.lastWaitSec -= time.Second
				select {
				case waitTime, ok := <-f.askMasterChan:
					if ok {
						time.Sleep(waitTime)
					} else {
						// 通道已经关闭，被监听线程关闭
						break loop
					}
				default:
				}
				if i == 2 {
					// 最后一次循环，如果是最后一次循环，则自己成为master
					f.closeAsk()
					f.beMaster()
				}
			}
		} // loop标签结束
	}()
}

// 停止询问环节
func (f *finder) closeAsk() {
	f.chanCloseOnce.Do(func() {
		close(f.askMasterChan)
		close(f.multicastListenerOutChan)
		f.isClose.Store(true)
	})
}
func (f *finder) saveMasterInfo(info baseInfo) {
	f.masterInfo = info
	f.masterAddr, _ = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", info.Ip, info.Port)) // 组播公用地址
}

// 作为一个节点监听信息
func (f *finder) beNode() {
	f.selfIsMaster.Store(false)
	// 监听线程
	go func() {
		// 询问当前的公开文件,3次重试
		for i := 0; i < 3; i++ {
			if _, err := f.multicastCoon.Write(askFilesBytes()); err != nil {
				time.Sleep(time.Second)
				continue
			}
			break
		}

		buf := make([]byte, 32*1024) // 32kb缓存

		for {
			n, _, err := f.multicastCoon.ReadFromUDP(buf)
			if err != nil {
				continue
			}
			f.nodeDecode(buf[:n])
		}

	}()
}

// 成为master
func (f *finder) beMaster() {
	f.saveMasterInfo(readSelfBaseInfo())
	f.selfIsMaster.Store(true)

	// 成为master后需要继续监听组播消息，给其他人回复master信息
	// 监听线程
	go func() {
		buf := make([]byte, 128)
		for {
			n, _, err := f.multicastCoon.ReadFromUDP(buf)
			if err != nil {
				continue
			}
			f.masterDecode(buf[:n])
		}
	}()
}

// 内部接口
// 读取剩余等待秒数
func readWaitSec() int {
	return max(0, int(defaultFinder().lastWaitSec.Seconds()))
}

// 要求等待线程延长等待时间
func askWait(t time.Duration) {
	if !defaultFinder().isClose.Load() {
		defaultFinder().askMasterChan <- t
	}
}

// 保存master信息
func receivedMasterAnwer(info baseInfo) {
	defaultFinder().saveMasterInfo(info)
}

// 读取自身基础结构体信息
func readSelfBaseInfo() baseInfo {
	return baseInfo{
		Ip:   getLocalIp(),
		Port: fmt.Sprintf("%v", getLocalTcpPort()),
		Id:   id,
	}
}

// 读取自身基础结构体信息
func readSelfBaseInfoBytes() []byte {
	p, _ := json.Marshal(readSelfBaseInfo())
	return p
}
