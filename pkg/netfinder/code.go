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
)

const (
	afterMask byte = 0x1f
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
	switch order {
	case askMasterNetOrder:
		//  询问指令
		waitSec := (time.Duration(data[0]&afterMask) + 1) * time.Second
		askWait(waitSec)
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
		files := []file{}
		err := json.Unmarshal(data[1:], &files)
		if err != nil {
			return
		}

		f.files = f.files[:0]
		f.files = append(f.files, files...)
	}
}

// master解码组播消息
// 返回来源的基础信息和是否需要回复
func (f *finder) masterDecode(data []byte) {
	if len(data) < 1 {
		return
	}
	order := decodeOrder(data[0])
	switch order {
	case askMasterNetOrder:
		//  询问指令
		f.multicastCoon.Write(masterAnswerBytes())
	case askFilesNetOrder:
		// 收到询问文件列表消息
		if dataFiles, err := json.Marshal(f.files); err != nil {
			// 文件列表序列化失败，一般不太可能
			return
		} else {
			filesOrdBytes := append([]byte{byte(masterAnswerFilesNetOrder)}, dataFiles...)
			for i := 0; i < 3; i++ {
				if _, err := f.multicastCoon.Write(filesOrdBytes); err != nil {
					time.Sleep(time.Second)
					continue
				}
				break
			}
		}

	}
}

// 解码出指令类型
func decodeOrder(b byte) netOrder {
	return netOrder(b >> 5)
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
