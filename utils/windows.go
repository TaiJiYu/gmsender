package utils

import (
	"bytes"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/text/transform"

	"github.com/lxn/win"
	"golang.org/x/text/encoding/simplifiedchinese"
)

var defaultFolder string

// decodeGBK 转换Windows下的GBK编码到UTF-8
func decodeGBK(s []byte) (string, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// cleanPath 清理路径字符串
func cleanPath(path string) string {
	path = strings.TrimSpace(path)
	// 处理macOS的路径前缀
	if runtime.GOOS == "darwin" && strings.HasPrefix(path, "alias:") {
		path = path[6:]
	}
	// 处理可能的换行符和引号
	path = strings.Trim(path, "\n\r\"'")
	return path
}

var muf sync.Mutex
var folderDir string

func chooseFolder() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// 使用PowerShell的FolderBrowserDialog
		script := `Add-Type -AssemblyName System.Windows.Forms
$folder = New-Object System.Windows.Forms.FolderBrowserDialog
if($folder.ShowDialog() -eq "OK") { $folder.SelectedPath }`
		cmd = exec.Command("powershell", "-Command", script)
		// case "darwin":
		// 	// macOS使用osascript
		// 	cmd = exec.Command("osascript", "-e", `choose folder with prompt "选择文件夹"`)
		// default:
		// 	// Linux使用zenity
		// 	cmd = exec.Command("zenity", "--file-selection", "--directory")
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err == nil {
		path := out.String()
		if runtime.GOOS == "windows" {
			if decoder, err := decodeGBK(out.Bytes()); err == nil {
				path = decoder
			} else {
				path = ""
			}
		}
		path = cleanPath(path)
		// if runtime.GOOS == "darwin" {
		// 	// macOS返回的路径格式需要处理
		// 	path = path[7 : len(path)-1] // 去除前缀"alias:"和换行
		// }

		muf.Lock()
		folderDir = path
		muf.Unlock()
	} else {
		folderDir = ""
	}
}

// 打开windows文件夹选择
func ChooseFolder() string {
	chooseFolder()
	return folderDir
}

// 打开windows任意文件选择器
func OpenWinChooseFile() string {
	var ofn win.OPENFILENAME
	fileName := make([]uint16, win.MAX_PATH)
	ofn.LStructSize = uint32(unsafe.Sizeof(ofn))
	ofn.LpstrFile = &fileName[0]
	ofn.NMaxFile = uint32(len(fileName))
	var err error
	ofn.LpstrFilter, err = syscall.UTF16PtrFromString("file") //\000*.*\000Text Files\000*.txt\000
	if err != nil {
		panic(err)
	}
	ofn.NFilterIndex = 1
	ofn.Flags = win.OFN_FILEMUSTEXIST | win.OFN_PATHMUSTEXIST

	if win.GetOpenFileName(&ofn) {
		selectedFile := syscall.UTF16ToString(fileName)
		return selectedFile
	} else {
		return ""
	}
}

// 打开windows保存文件名
func OpenWinSaveFileName() string {
	var ofn win.OPENFILENAME
	fileName := make([]uint16, win.MAX_PATH)
	ofn.LStructSize = uint32(unsafe.Sizeof(ofn))
	ofn.LpstrFile = &fileName[0]
	ofn.NMaxFile = uint32(len(fileName))
	// r := utf16.Encode([]rune("所有文件(*.*)"))
	// r = append(r, 0)
	ofn.LpstrFilter = nil
	// ofn.NFilterIndex = 1
	ofn.Flags = win.OFN_FILEMUSTEXIST | win.OFN_PATHMUSTEXIST

	if win.GetSaveFileName(&ofn) {
		return syscall.UTF16ToString(fileName)
	} else {
		return ""
	}
}
