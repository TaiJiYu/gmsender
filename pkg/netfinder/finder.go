package netfinder

import (
	"encoding/json"
	"fmt"
	"gmsender/pkg/asset"
	"io"
	"math/rand/v2"
	"net"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

// 服务发现
type finder struct {
	isExWating               atomic.Bool // 是否在延长等待期间
	askMasterChan            chan time.Duration
	askMasterOutChan         chan struct{} // 询问线程退出信号
	multicastListenerOutChan chan struct{}
	ptopTcpOutChan           chan struct{} // tcp链接管道
	isClose                  atomic.Bool

	lastWaitSec time.Duration // 剩余询问等待时间

	masterInfo          baseInfo
	masterbroadcastAddr *net.UDPAddr // 应该是master的广播通信地址，而非master的tcp通信地址

	selfIsMaster atomic.Bool // 自己是否是master

	isInitDoneB   atomic.Bool  // 是否初始化完成
	broadcastCoon *net.UDPConn // 组播监听链接

	finderErr error

	// nodes []baseInfo // 所有节点的信息，包含自身和master，所有节点与master同步
	fileLock sync.RWMutex      // 文件锁
	files    []File            // 公开的文件列表
	filesMap map[File]struct{} // 去重检查

	// 加密相关
	keyManager *keyManager // ECDH密钥对管理器
}

var (
	finderCli *finder
	finderMut sync.Mutex
)

func defaultFinder() *finder {
	finderMut.Lock()
	if finderCli == nil {
		// 生成ECDH密钥对用于端到端加密
		keyMgr, err := generateKeyPair()
		if err != nil {
			fmt.Printf("生成加密密钥对失败: %v\n", err)
		}

		finderCli = &finder{
			isExWating:               atomic.Bool{},
			askMasterChan:            make(chan time.Duration, 1),
			askMasterOutChan:         make(chan struct{}, 1),
			multicastListenerOutChan: make(chan struct{}, 1),
			ptopTcpOutChan:           make(chan struct{}, 1),
			files:                    make([]File, 0),
			filesMap:                 make(map[File]struct{}),
			keyManager:               keyMgr,
		}

		finderCli.lastWaitSec = 3 * time.Second
		finderCli.pTopTcpListener()
		finderCli.broadcastListener()
		finderCli.askMaster()
	}
	r := finderCli
	finderMut.Unlock()
	return r
}

// 刷新
func (f *finder) reIn() {
	// 之前没有初始化，不能执行
	if f.isInitDoneB.Load() {
		// 已经初始化过了
		if f.selfIsMaster.Load() {
			// master只需要问一下谁是master
			f.broadcastCoon.WriteToUDP(askMasterBytes(), broadcastAddr)
			initDoneCallBack()
		} else {
			// 如果是节点需要先关闭广播请求重新问一下
			f.broadcastCoon.Close()
			f.broadcastCoon = nil
			f.isClose.Store(false)
			f.isInitDoneB.Store(false)
			f.askMasterOutChan = make(chan struct{}, 1)
			f.multicastListenerOutChan = make(chan struct{}, 1)
			finderCli.broadcastListener()
			finderCli.askMaster()
		}
	}
}

// 停止询问环节
func (f *finder) closeAsk() {
	if f.isClose.CompareAndSwap(false, true) {
		close(f.askMasterOutChan)         // 关闭询问
		close(f.multicastListenerOutChan) // 关闭询问监听
	}
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
func (f *finder) broadcastListener() {
	// 绑定本地端口接收组播
	if f.broadcastCoon == nil {
		conn, err := net.ListenUDP("udp", broadcastListenAddr)
		if err != nil {
			panic(err)
		}
		f.broadcastCoon = conn
	}

	// 监听线程
	go func() {
		buf := make([]byte, 32*1024) // 32kB缓存
		{
		listenerLoop:
			for {
				n, err := f.broadcastCoon.Read(buf)
				if err != nil {
					continue
				}
				if f.decode(buf[:n]) {
					// 收到master回复并退出监听，自己不是master，但保留通信组播监听，用于接受之后的组播消息
					f.closeAsk()
					f.beNode()
					break
				}

				select {
				case <-f.multicastListenerOutChan:
					// 没收到master回复但是收到了退出信号，自己成为master
					f.closeAsk()
					break listenerLoop
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
				f.broadcastCoon.WriteToUDP(askMasterBytes(), broadcastAddr)
				time.Sleep(time.Second)
				f.lastWaitSec -= time.Second
				select {
				case <-f.askMasterOutChan:
					// 收到退出信号
					// 通道已经关闭，被监听线程关闭
					break loop
				case waitTime := <-f.askMasterChan:
					time.Sleep(waitTime)
					f.isExWating.Store(false)
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
	if f.broadcastCoon != nil {
		f.broadcastCoon.Close()
	}
}

// 收到其他用户的下载请求
func (f *finder) handlerDownLoad(conn net.Conn) {

	buf := make([]byte, 32*1024)

	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
	}

	// 解码下载请求（包含文件名和客户端公钥）
	req := decodeDownloadFileInfo(buf[:n])
	if req.FileName == "" {
		conn.Close()
		return
	}

	// 打开文件
	fileS, err := os.Open(req.FileName)
	if err != nil {
		// 文件错误
		conn.Close()
		return
	}
	defer fileS.Close()

	// 获取文件大小
	fileInfo, err := fileS.Stat()
	if err != nil {
		fmt.Printf("获取文件信息失败: %v\n", err)
		conn.Close()
		return
	}
	fileSize := fileInfo.Size()

	// 如果客户端提供了公钥，则使用加密传输
	if req.PubKey != "" && f.keyManager != nil {
		// 计算共享密钥
		sharedSecret, err := f.keyManager.computeSharedSecret(req.PubKey)
		if err != nil {
			fmt.Printf("计算共享密钥失败: %v\n", err)
			conn.Close()
			return
		}

		// 创建加密连接
		encryptedConn, err := newEncryptedConn(conn, sharedSecret)
		if err != nil {
			fmt.Printf("创建加密连接失败: %v\n", err)
			conn.Close()
			return
		}
		defer encryptedConn.Close()

		// 先发送nonce给客户端
		if _, err := conn.Write(encryptedConn.getNonce()); err != nil {
			fmt.Printf("发送nonce失败: %v\n", err)
			return
		}

		// 先发送文件大小（8字节）
		sizeBuf := make([]byte, 8)
		sizeBuf[0] = byte(fileSize)
		sizeBuf[1] = byte(fileSize >> 8)
		sizeBuf[2] = byte(fileSize >> 16)
		sizeBuf[3] = byte(fileSize >> 24)
		sizeBuf[4] = byte(fileSize >> 32)
		sizeBuf[5] = byte(fileSize >> 40)
		sizeBuf[6] = byte(fileSize >> 48)
		sizeBuf[7] = byte(fileSize >> 56)
		if _, err := encryptedConn.Write(sizeBuf); err != nil {
			fmt.Printf("发送文件大小失败: %v\n", err)
			return
		}

		// 流式加密发送
		io.Copy(encryptedConn, fileS)
	} else {
		// 兼容旧版本：不使用加密
		// 先发送文件大小（8字节）
		sizeBuf := make([]byte, 8)
		sizeBuf[0] = byte(fileSize)
		sizeBuf[1] = byte(fileSize >> 8)
		sizeBuf[2] = byte(fileSize >> 16)
		sizeBuf[3] = byte(fileSize >> 24)
		sizeBuf[4] = byte(fileSize >> 32)
		sizeBuf[5] = byte(fileSize >> 40)
		sizeBuf[6] = byte(fileSize >> 48)
		sizeBuf[7] = byte(fileSize >> 56)
		if _, err := conn.Write(sizeBuf); err != nil {
			fmt.Printf("发送文件大小失败: %v\n", err)
			return
		}

		io.Copy(conn, fileS)
	}

	conn.Close()
}

func (f *finder) saveMasterInfo(info baseInfo) {
	f.masterInfo = info
	f.masterbroadcastAddr, _ = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", info.Ip, broadcastPort)) // 广播公用地址
}

// 删除公开文件
func (f *finder) delPublicFile(file File) {
	if f.selfIsMaster.Load() {
		// 自己是master
		f.delFile(file)
		go f.masterAnswerFiles()
	} else {
		// 自己是节点
		go f.broadcastCoon.WriteToUDP(delPublicSelfFileBytes(file), f.masterbroadcastAddr)
	}
}

const (
	retryTimesMax = 5 // 最多重试5次
)

// 请求下载别人的文件
func (f *finder) downloadFile(saveToFloderName string, info File, progressCallback func(downloaded int64, total int64), statusCallback func(status string)) {
	go func() {
		if statusCallback != nil {
			statusCallback("started")
		}

		for i := 0; i < retryTimesMax; i++ {
			// 点对点链接
			conn, err := net.Dial("tcp", info.addr())
			if err != nil {
				// 链接出错，等待1-3秒后重试
				time.Sleep(time.Duration(rand.IntN(3)+1) * time.Second)
				continue
			}

			// 生成临时密钥对用于本次下载
			privKey, err := generateKeyPair()
			if err != nil {
				fmt.Printf("生成临时密钥对失败: %v\n", err)
				conn.Close()
				break
			}

			// 请求下载文件（携带公钥）
			message := downLoadFileBytes(info.FileName, privKey.publicKeyBase64())
			_, err = conn.Write(message)
			if err != nil {
				// 接受失败，等会重试
				time.Sleep(time.Duration(rand.IntN(3)+1) * time.Second)
				conn.Close()
				continue
			}

			saveFileName := filepath.Join(saveToFloderName, info.FileNamePure())

			// 接收文件
			file, err := os.Create(saveFileName)
			if err != nil {
				conn.Close()
				break
			}

			var reader io.Reader = conn
			var totalSize int64 = 0

			// 如果服务端提供了公钥，则使用加密接收
			if info.PubKey != "" && f.keyManager != nil {
				// 计算共享密钥
				sharedSecret, err := privKey.computeSharedSecret(info.PubKey)
				if err != nil {
					fmt.Printf("计算共享密钥失败: %v\n", err)
					file.Close()
					conn.Close()
					break
				}

				// 先读取服务端发送的nonce
				nonce := make([]byte, 24) // XChaCha20-Poly1305需要24字节nonce
				n, err := io.ReadFull(conn, nonce)
				if err != nil || n != 24 {
					fmt.Printf("读取nonce失败: %v\n", err)
					file.Close()
					conn.Close()
					break
				}

				// 使用服务端的nonce创建加密连接
				encryptedConn, err := newEncryptedConnWithNonce(conn, sharedSecret, nonce)
				if err != nil {
					fmt.Printf("创建加密连接失败: %v\n", err)
					file.Close()
					conn.Close()
					break
				}
				defer encryptedConn.Close()
				reader = encryptedConn
			}

			// 先读取文件大小（8字节）
			sizeBuf := make([]byte, 8)
			_, err = io.ReadFull(reader, sizeBuf)
			if err != nil {
				fmt.Printf("读取文件大小失败: %v\n", err)
				file.Close()
				if info.PubKey != "" && f.keyManager != nil {
					conn.Close()
				}
				break
			}
			totalSize = int64(sizeBuf[0]) | int64(sizeBuf[1])<<8 | int64(sizeBuf[2])<<16 | int64(sizeBuf[3])<<24 |
				int64(sizeBuf[4])<<32 | int64(sizeBuf[5])<<40 | int64(sizeBuf[6])<<48 | int64(sizeBuf[7])<<56

			if statusCallback != nil {
				statusCallback("downloading")
			}

			// 分块读取并报告进度
			buffer := make([]byte, 32*1024) // 32KB 缓冲区
			var downloaded int64

			for {
				n, err := reader.Read(buffer)
				if n > 0 {
					_, writeErr := file.Write(buffer[:n])
					if writeErr != nil {
						fmt.Printf("写入文件失败: %v\n", writeErr)
						break
					}
					downloaded += int64(n)
					if progressCallback != nil && totalSize > 0 {
						progressCallback(downloaded, totalSize)
					}
				}
				if err != nil {
					if err == io.EOF {
						if statusCallback != nil {
							statusCallback("completed")
						}
						asset.PlayDoneMusic()
					} else {
						fmt.Println(err)
						if statusCallback != nil {
							statusCallback("failed")
						}
						asset.PlayFailMusic()
					}
					break
				}
			}

			file.Close()
			if info.PubKey != "" && f.keyManager != nil {
				conn.Close()
			} else {
				conn.Close()
			}
			break
		}
	}()
}

// 公开文件,自己调用了公开本机文件
func (f *finder) publicFile(filename string) {
	fileS := File{
		Ip:       getLocalIp(),
		Port:     fmt.Sprintf("%v", getLocalTcpPort()),
		Id:       Id(),
		FileName: filename,
		PubKey:   f.keyManager.publicKeyBase64(), // 添加ECDH公钥用于端到端加密
	}
	if f.selfIsMaster.Load() {
		// 自己master
		if files, change := f.appendFile(fileS); change {
			filesCallback(files)
			go f.masterAnswerFiles()
		}
	} else {
		// 自己是节点
		go f.broadcastCoon.WriteToUDP(publicSelfFileBytes(fileS), f.masterbroadcastAddr)
	}
}

// 放弃master身份转换为node
func (f *finder) masterToNode(master baseInfo) {
	// 先更新前端
	f.selfIsMaster.Store(false)
	f.saveMasterInfo(master)
	typeChangeCallback(false)
	go func() {
		time.Sleep(time.Second) // 等待一秒缓冲
		f.nodeLoop()
	}()
}

// 节点循环职责
func (f *finder) nodeLoop() {
	// 询问当前的公开文件,3次重试
	for i := 0; i < 3; i++ {
		if _, err := f.broadcastCoon.WriteTo(askFilesBytes(), f.masterbroadcastAddr); err != nil {
			time.Sleep(time.Second)
			continue
		}
		break
	}

	buf := make([]byte, 32*1024) // 32kB缓存

	for {
		n, err := f.broadcastCoon.Read(buf)
		if err != nil {
			return
		}
		f.nodeDecode(buf[:n])
	}
}

// 作为一个节点监听信息
func (f *finder) beNode() {
	if f.isInitDoneB.Load() {
		// 已经交换过了
		return
	}
	f.selfIsMaster.Store(false)
	f.isInitDoneB.Store(true)
	typeChangeCallback(false)
	initDoneCallBack()
	fmt.Println("node", Id())
	// 监听线程
	go f.nodeLoop()
}

// 成为master
func (f *finder) beMaster() {
	if f.isInitDoneB.Load() {
		// 已经交换过了
		return
	}

	f.saveMasterInfo(readSelfBaseInfo())
	f.selfIsMaster.Store(true)
	f.isInitDoneB.Store(true)
	typeChangeCallback(true)
	initDoneCallBack()
	fmt.Println("master", Id())

	// 成为master后需要继续监听组播消息，给其他人回复master信息
	// 监听线程
	go func() {
		buf := make([]byte, 32*1024) // 32kB缓存

		for {
			n, err := f.broadcastCoon.Read(buf)
			if err != nil {
				continue
			}
			if f.masterDecode(buf[:n]) {
				// 放弃了master身份
				return
			}
		}
	}()
}

// 内部接口
// 读取剩余等待秒数
func readWaitSec() int {
	return max(0, int(defaultFinder().lastWaitSec.Seconds()))
}

// 要求等待线程延长等待时间
func (f *finder) askWait(t time.Duration) {
	if f.isExWating.Load() {
		// 延长等待期间
		return
	}
	f.isExWating.Store(true)
	defaultFinder().askMasterChan <- t
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
