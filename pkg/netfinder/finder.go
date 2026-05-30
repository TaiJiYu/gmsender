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

	isInitDoneB   atomic.Bool  // 是否初始化完成
	multicastCoon *net.UDPConn // 组播监听链接

	finderErr error

	// nodes []baseInfo // 所有节点的信息，包含自身和master，所有节点与master同步
	fileLock sync.RWMutex      // 文件锁
	files    []File            // 公开的文件列表
	filesMap map[File]struct{} // 去重检查
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
			files:                    make([]File, 0),
			filesMap:                 make(map[File]struct{}),
		}
		finderCli.pTopTcpListener()
		finderCli.multicastListener()
		finderCli.askMaster()
	})
	return finderCli
}

// 读取自己的文件列表
func (f *finder) readFiles() []File {
	f.fileLock.RLock()
	files := make([]File, len(f.files))
	copy(files, f.files)
	f.fileLock.RUnlock()
	return files
}

// 重写自己的文件列表
func (f *finder) writeFiles(files []File) {
	f.fileLock.Lock()
	f.files = make([]File, len(files))
	copy(f.files, files)
	clear(f.filesMap)
	for i := range files {
		f.filesMap[files[i]] = struct{}{}
	}
	f.fileLock.Unlock()
}

// 删除files中的文件
func (f *finder) delFile(file File) {
	f.fileLock.Lock()
	defer f.fileLock.Unlock()
	_, isIn := f.filesMap[file]
	if isIn {
		delete(f.filesMap, file)
		for i := range f.files {
			if file == f.files[i] {
				f.files = append(f.files[:i], f.files[i+1:]...)
				return
			}
		}
	}
}

// 添加一个文件到列表中,返回全部文件，重复的不会添加，返回是否有变化
func (f *finder) appendFile(file File) ([]File, bool) {
	f.fileLock.Lock()
	defer f.fileLock.Unlock()
	_, isIn := f.filesMap[file]
	if isIn {
		files := make([]File, len(f.files))
		copy(files, f.files)
		return files, false
	} else {
		f.files = append(f.files, file)
		f.filesMap[file] = struct{}{}
		ret := make([]File, len(f.files))
		copy(ret, f.files)
		return ret, true
	}
}

// 返回是否初始化完成和是否有错误
func (f *finder) isInitDone() (bool, error) {
	return f.isInitDoneB.Load(), f.finderErr
}

// 自身的点对点通信监听
func (f *finder) pTopTcpListener() {
	// 端口写 0，系统自动分配
	lister, err := net.Listen("tcp", fmt.Sprintf("%s:0", getLocalIp()))
	if err != nil {
		f.finderErr = err
		return
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

// 组播监听
func (f *finder) multicastListener() {
	// 绑定本地端口接收组播
	conn, err := net.ListenUDP("udp", broadcastListenAddr)
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
				f.multicastCoon.WriteToUDP(askMasterBytes(), broadcastAddr)

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

func (f *finder) closeNetFinder() {
	close(f.ptopTcpOutChan)
	if f.multicastCoon != nil {
		f.multicastCoon.Close()
	}
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

// 删除公开文件
func (f *finder) delPublicFile(file File) {
	if f.selfIsMaster.Load() {
		// 自己是master
		f.delFile(file)
		go f.masterAnswerFiles()
	} else {
		// 自己是节点
		go f.multicastCoon.WriteToUDP(delPublicSelfFileBytes(file), f.masterAddr)
	}
}

// 公开文件,自己调用了公开本机文件
func (f *finder) publicFile(filename string) {
	fileS := File{
		Ip:       getLocalIp(),
		Port:     fmt.Sprintf("%v", getLocalTcpPort()),
		Id:       Id(),
		FileName: filename,
	}
	if f.selfIsMaster.Load() {
		// 自己master
		if files, change := f.appendFile(fileS); change {
			filesCallback(files)
			go f.masterAnswerFiles()
		}
	} else {
		// 自己是节点
		go f.multicastCoon.WriteToUDP(publicSelfFileBytes(fileS), f.masterAddr)
	}
}

// 作为一个节点监听信息
func (f *finder) beNode() {
	f.selfIsMaster.Store(false)
	f.isInitDoneB.Store(true)
	fmt.Println("node", Id())
	// 监听线程
	go func() {
		// 询问当前的公开文件,3次重试
		for i := 0; i < 3; i++ {
			if _, err := f.multicastCoon.WriteTo(askFilesBytes(), f.masterAddr); err != nil {
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
	f.isInitDoneB.Store(true)
	fmt.Println("master", Id())

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
