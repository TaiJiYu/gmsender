package netfinder

import (
	"encoding/json"
	"fmt"
	"time"
)

type netOrder byte

const (
	askMasterNetOrder         netOrder = (iota + 1) << 5 // 询问master命令
	masterAnswerNetOrder                                 // master回复
	downLoadFileNetOrder                                 // 请求文件下载
	askFilesNetOrder                                     // 询问公开文件列表
	masterAnswerFilesNetOrder                            // master回复公开文件列表
	publicFileNetOrder                                   // 公开自己的文件
	delfPublicFileNetOrder                               // 删除自己的公开文件
)

const (
	afterMask    byte = 0x1f
	netOrderMask byte = ^afterMask
	// 协议魔数字符串，用于标识本软件的消息，避免误识别其他软件的广播
	protocolMagic = "GMStudio:GMSenderOrder:Broadcast"
)

// 建联基础信息
type baseInfo struct {
	Ip   string `json:"ip"`
	Port string `json:"port"`
	Id   string `json:"id"`
}

func (b baseInfo) addr() string {
	return fmt.Sprintf("%s:%s", b.Ip, b.Port)
}

// 解码与编码
func askMasterBytes() []byte {
	header := byte(askMasterNetOrder) | (byte(readWaitSec()) & afterMask)
	return append([]byte(protocolMagic), append([]byte{header}, readSelfBaseInfoBytes()...)...)
}

// master回应
func masterAnswerBytes() []byte {
	header := byte(masterAnswerNetOrder)
	return append([]byte(protocolMagic), append([]byte{header}, readSelfBaseInfoBytes()...)...)
}

func (f *finder) decode(data []byte) bool {
	magicLen := len(protocolMagic)
	if len(data) < magicLen+1 {
		return false
	}
	// 首先检查协议魔数
	if string(data[:magicLen]) != protocolMagic {
		return false
	}
	order := decodeOrder(data[magicLen])
	info := baseInfo{}
	err := json.Unmarshal(data[magicLen+1:], &info)
	if err != nil {
		return false
	}
	if info.Id == Id() {
		// 来自自己的消息，直接忽略
		return false
	}
	switch order {
	case askMasterNetOrder:
		//  询问指令
		if !CompetitiveMaster(info) {
			// 竞争失败主动退让
			waitSec := (time.Duration(data[magicLen]&afterMask) + 1) * time.Second
			f.askWait(waitSec)
		}
		return false
	case masterAnswerNetOrder:
		// master回复
		receivedMasterAnwer(info)
		return true
	}
	return false
}

// 竞争master，如果返回true说明竞争成功，有资格成为master，不退让
func CompetitiveMaster(otherInfo baseInfo) bool {
	if getLocalIp() > otherInfo.Ip {
		// 自己ip大于对方，则有资格
		return true
	}
	if getLocalIp() == otherInfo.Ip {
		// 与对方ip相等，再考虑id
		return Id() > otherInfo.Id
	}

	return false
}

// 询问现在的公开文件列表
func askFilesBytes() []byte {
	header := byte(askFilesNetOrder)
	return append([]byte(protocolMagic), header)

}
func publicSelfFileBytes(f File) []byte {
	header := byte(publicFileNetOrder)
	data, _ := json.Marshal(f)
	return append([]byte(protocolMagic), append([]byte{header}, data...)...)
}

// 删除自己的公开文件
func delPublicSelfFileBytes(f File) []byte {
	header := byte(delfPublicFileNetOrder)
	data, _ := json.Marshal(f)
	return append([]byte(protocolMagic), append([]byte{header}, data...)...)
}

// 某个节点解码组播消息
// 返回来源的基础信息和是否需要回复
func (f *finder) nodeDecode(data []byte) {
	magicLen := len(protocolMagic)
	if len(data) < magicLen+1 {
		return
	}
	// 首先检查协议魔数
	if string(data[:magicLen]) != protocolMagic {
		return
	}
	order := decodeOrder(data[magicLen])
	switch order {
	case masterAnswerFilesNetOrder:
		//  master回复的文件列表
		files := []File{}
		err := json.Unmarshal(data[magicLen+1:], &files)
		if err != nil {
			return
		}
		f.writeFiles(files)
		filesCallback(files)
	}
}

// master同步所有文件
func (f *finder) masterAnswerFiles() {
	files := f.readFiles()
	if dataFiles, err := json.Marshal(files); err != nil {
		// 文件列表序列化失败，一般不太可能
		return
	} else {
		filesOrdBytes := append([]byte(protocolMagic), append([]byte{byte(masterAnswerFilesNetOrder)}, dataFiles...)...)
		for i := 0; i < 3; i++ {
			if _, err := f.broadcastCoon.WriteToUDP(filesOrdBytes, broadcastAddr); err != nil {
				time.Sleep(time.Second)
				continue
			}
			break
		}
	}
}

// master解码组播消息
// 返回是否放弃了master身份
func (f *finder) masterDecode(data []byte) bool {
	magicLen := len(protocolMagic)
	if len(data) < magicLen+1 {
		return false
	}
	// 首先检查协议魔数
	if string(data[:magicLen]) != protocolMagic {
		return false
	}
	order := decodeOrder(data[magicLen])
	switch order {
	case masterAnswerNetOrder:
		// 收到了master回复
		info := baseInfo{}
		err := json.Unmarshal(data[magicLen+1:], &info)
		if err != nil {
			return false
		}
		if info.Id == Id() {
			// 来自自己的消息，直接忽略
			return false
		}
		if !CompetitiveMaster(info) {
			// 竞争失败，放弃master身份
			f.masterToNode(info)
			return true
		}

	case askMasterNetOrder:
		//  询问指令
		f.broadcastCoon.WriteToUDP(masterAnswerBytes(), broadcastAddr)
	case askFilesNetOrder:
		// 收到询问文件列表消息
		f.masterAnswerFiles()
	case publicFileNetOrder:
		// 收到了节点的公开请求
		file := File{}
		if err := json.Unmarshal(data[magicLen+1:], &file); err != nil {
			// 文件有问题
			return false
		}
		if newfiles, change := f.appendFile(file); change {
			// 保存后通知其他节点
			filesCallback(newfiles)
			f.masterAnswerFiles()
		}

	case delfPublicFileNetOrder:
		// 收到了删除公开文件的请求
		file := File{}
		if err := json.Unmarshal(data[magicLen+1:], &file); err != nil {
			// 文件有问题
			return false
		}
		f.delFile(file)
		filesCallback(f.readFiles())
		f.masterAnswerFiles()
	}

	return false
}

// 解码出指令类型
func decodeOrder(b byte) netOrder {
	return netOrder(b & netOrderMask)
}

// 请求文件下载的信息
type downloadRequest struct {
	FileName string `json:"file_name"`
	PubKey   string `json:"pub_key"`
}

// 请求文件下载（带公钥）
func downLoadFileBytes(fileName, pubKey string) []byte {
	req := downloadRequest{
		FileName: fileName,
		PubKey:   pubKey,
	}
	data, _ := json.Marshal(req)
	return append([]byte(protocolMagic), append([]byte{byte(downLoadFileNetOrder)}, data...)...)
}

// 解码下载请求（支持带公钥的新格式和旧的纯文件名格式）
func decodeDownloadFileInfo(data []byte) downloadRequest {
	magicLen := len(protocolMagic)
	if len(data) < magicLen+1 {
		return downloadRequest{}
	}

	// 首先检查协议魔数
	if string(data[:magicLen]) != protocolMagic {
		return downloadRequest{}
	}

	if decodeOrder(data[magicLen]) == downLoadFileNetOrder {
		req := downloadRequest{}
		if err := json.Unmarshal(data[magicLen+1:], &req); err != nil {
			// 尝试兼容旧的纯文件名格式
			return downloadRequest{FileName: string(data[magicLen+1:])}
		}
		return req
	}
	return downloadRequest{}
}
