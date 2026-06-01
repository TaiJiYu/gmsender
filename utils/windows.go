package utils

import (
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

// 回调
func browseCallback(hwnd win.HWND, uMsg uint, lParam, lpData uintptr) uintptr {
	return 0
}

// 打开windows文件夹选择器
func OpenWinFolder() string {
	var bi win.BROWSEINFO
	pathBuffer := make([]uint16, win.MAX_PATH)
	bi.PszDisplayName = &pathBuffer[0]                      // 接收显示名称的缓冲区
	bi.LpszTitle, _ = syscall.UTF16PtrFromString("选择下载位置：") // 对话框标题
	bi.Lpfn = syscall.NewCallback(browseCallback)           // 回调函数指针

	// 2. 调用 Windows API 显示对话框
	pidl := win.SHBrowseForFolder(&bi)
	if pidl == 0 {
		return ""
	}

	defer win.CoTaskMemFree(pidl) // 记得释放 PIDL 内存

	// 3. 将用户选择的路径 (PIDL) 转换为可读的路径字符串
	if !win.SHGetPathFromIDList(pidl, &pathBuffer[0]) {
		return ""
	}

	// 4. 将 UTF-16 字节数组转换为 Go 字符串
	selectedPath := syscall.UTF16ToString(pathBuffer[:])
	return selectedPath
}

// 打开windows任意文件选择器
func OpenWinChooseFile() string {
	var ofn win.OPENFILENAME
	fileName := make([]uint16, win.MAX_PATH)
	ofn.LStructSize = uint32(unsafe.Sizeof(ofn))
	ofn.LpstrFile = &fileName[0]
	ofn.NMaxFile = uint32(len(fileName))

	ofn.LpstrFilter = nil
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
	ofn.LpstrFilter = nil
	ofn.Flags = win.OFN_FILEMUSTEXIST | win.OFN_PATHMUSTEXIST

	if win.GetSaveFileName(&ofn) {
		return syscall.UTF16ToString(fileName)
	} else {
		return ""
	}
}
