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
	netOrderMask byte = 0xff - afterMask
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
	return append([]byte{header}, readSelfBaseInfoBytes()...)
}

// master回应
func masterAnswerBytes() []byte {
	header := byte(masterAnswerNetOrder)
	return append([]byte{header}, readSelfBaseInfoBytes()...)
}

func decode(data []byte) bool {
	if len(data) < 1 {
		return false
	}
	order := decodeOrder(data[0])
	info := baseInfo{}
	err := json.Unmarshal(data[1:], &info)
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
		if info.Ip > getLocalIp() || info.Id > Id() {
			// 如果对方ip高于自己或者id大于自己，则放弃选主，直接等待
			waitSec := (time.Duration(data[0]&afterMask) + 1) * time.Second
			askWait(waitSec)
		} else {
			// 否则不做处理继续询问
		}
		return false
	case masterAnswerNetOrder:
		// master回复
		receivedMasterAnwer(info)
		return true
	}
	return false
}

// 询问现在的公开文件列表
func askFilesBytes() []byte {
	header := byte(askFilesNetOrder)
	return []byte{header}

}
func publicSelfFileBytes(f File) []byte {
	header := byte(publicFileNetOrder)
	data, _ := json.Marshal(f)
	return append([]byte{header}, data...)
}

// 删除自己的公开文件
func delPublicSelfFileBytes(f File) []byte {
	header := byte(delfPublicFileNetOrder)
	data, _ := json.Marshal(f)
	return append([]byte{header}, data...)
}

// 某个节点解码组播消息
// 返回来源的基础信息和是否需要回复
func (f *finder) nodeDecode(data []byte) {
	if len(data) < 1 {
		return
	}
	order := decodeOrder(data[0])
	switch order {
	case masterAnswerFilesNetOrder:
		//  master回复的文件列表
		files := []File{}
		err := json.Unmarshal(data[1:], &files)
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
		filesOrdBytes := append([]byte{byte(masterAnswerFilesNetOrder)}, dataFiles...)
		for i := 0; i < 3; i++ {
			if _, err := f.multicastCoon.WriteToUDP(filesOrdBytes, broadcastAddr); err != nil {
				time.Sleep(time.Second)
				continue
			}
			break
		}
	}
}

// master解码组播消息
// 返回来源的基础信息和是否需要回复
func (f *finder) masterDecode(data []byte) {
	if len(data) < 1 {
		return
	}
	order := decodeOrder(data[0])
	fmt.Println("收到指令", order)
	switch order {
	case askMasterNetOrder:
		//  询问指令
		f.multicastCoon.WriteToUDP(masterAnswerBytes(), broadcastAddr)
	case askFilesNetOrder:
		// 收到询问文件列表消息
		f.masterAnswerFiles()
	case publicFileNetOrder:
		// 收到了节点的公开请求
		fmt.Println("收到公布文件开始")
		file := File{}
		if err := json.Unmarshal(data[1:], &file); err != nil {
			// 文件有问题
			return
		}
		fmt.Println("收到公布文件解码完成")
		if newfiles, change := f.appendFile(file); change {
			// 保存后通知其他节点
			filesCallback(newfiles)
			f.masterAnswerFiles()
		}

	case delfPublicFileNetOrder:
		// 收到了删除公开文件的请求
		file := File{}
		if err := json.Unmarshal(data[1:], &file); err != nil {
			// 文件有问题
			return
		}
		f.delFile(file)
		filesCallback(f.readFiles())
		f.masterAnswerFiles()

	}
}

// 解码出指令类型
func decodeOrder(b byte) netOrder {
	return netOrder(b & netOrderMask)
}

// 请求文件下载的信息
func downLoadFileBytes(fileName string) []byte {
	return append([]byte{byte(downLoadFileNetOrder)}, []byte(fileName)...)
}

// 解码下载信息
func decodeDownloadFileInfo(data []byte) string {
	if len(data) < 1 {
		return ""
	}
	if decodeOrder(data[0]) == downLoadFileNetOrder {
		return string(data[1:])
	}
	return ""
}
