package gmsender

import (
	"container/list"
	"fmt"
	gametime "gmsender/pkg/game_time"
	"gmsender/pkg/input"
	"gmsender/pkg/netfinder"
	"gmsender/pkg/ui"
	"gmsender/utils"
	"image/color"
	"path/filepath"
	"time"
)

const (
	// 文件列表尺寸
	fileListX = utils.LogicalSizeX - 40*2
	// 单个文件高度（增加一行的高度）
	fileY = 70
)

// formatFileSize 格式化文件大小为合适的单位
func formatFileSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	} else if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
	} else {
		return fmt.Sprintf("%.2f GB", float64(bytes)/(1024*1024*1024))
	}
}

// formatSpeed 格式化下载速度
func formatSpeed(bytesPerSecond float64) string {
	if bytesPerSecond < 1024 {
		return fmt.Sprintf("%d B/s", int(bytesPerSecond))
	} else if bytesPerSecond < 1024*1024 {
		return fmt.Sprintf("%.1f KB/s", bytesPerSecond/1024)
	} else {
		return fmt.Sprintf("%.1f MB/s", bytesPerSecond/(1024*1024))
	}
}

// 公开文件列表
type fileList struct {
	files   *list.List
	fileMap map[netfinder.File]struct{}
}

var fileListCli *fileList

func InitFileList() {
	fileListCli = &fileList{
		files:   list.New(),
		fileMap: make(map[netfinder.File]struct{}),
	}
}
func (l *fileList) refreshFiles(files []netfinder.File) {
	i := 0
	for e := l.files.Front(); e != nil; e = e.Next() {
		f := e.Value.(*fileCmp)
		if i < len(files) {
			f.changeByFile(files[i])
		} else {
			l.files.Remove(e)
			delFileCmp(f.canvas)
		}
		i++
	}
	for i < len(files) {
		appendFileCmp(l.AddFile(files[i]))
		i++
	}
}

func (l *fileList) AddFile(file netfinder.File) *ui.CanvasUi {
	f := newFileCmp(file)
	l.files.PushBack(f)
	return f.canvas
}
func (l *fileList) Update(checkPos utils.Point) {
	for e := l.files.Front(); e != nil; e = e.Next() {
		f := e.Value.(*fileCmp)
		if f.isDel {
			l.files.Remove(e)
			delFileCmp(f.canvas)
			netfinder.DelPublicFile(f.file)
		} else {
			f.update(checkPos)
		}
	}
}

// 列表内的单个组件
type fileCmp struct {
	canvas, funcCanvas *ui.CanvasUi
	button             *ui.ButtonUi

	fileNameText, orginText, funcText, progressText *ui.TextUi //funcText是功能文本，progressText是进度文本

	file netfinder.File

	isDel bool

	// 下载状态
	downloading    bool
	downloaded     int64
	totalSize      int64
	lastDownloaded int64   // 上次记录的下载量，用于计算速度
	lastUpdateTime int64   // 上次更新时间戳
	currentSpeed   float64 // 当前下载速度（字节/秒）
}

func (f *fileCmp) changeByFile(file netfinder.File) {
	if file == f.file {
		return
	}
	f.file = file
	f.fileNameText.SetText(filepath.Base(file.FileName))
	f.orginText.SetText("-来自[" + file.Id + "]")
	if f.file.Id == netfinder.Id() {
		// 自己的
		f.funcText.SetText("关闭")
		f.funcText.AddSpaceToSizeX(100)
		f.button.SetFillColor(closeFileColor, closeFileColor)
	} else {
		// 别人的
		f.funcText.SetText("下载")
		f.funcText.AddSpaceToSizeX(100)
		f.button.SetFillColor(downloadColor, downloadColor)
	}

	f.isDel = false
}

func (f *fileCmp) buttonFunc(bu *ui.ButtonUi) {
	if f.file.Id == netfinder.Id() {
		// 自己的
		f.isDel = true
	} else {
		// 别人的
		if f.downloading {
			return // 下载中，不允许重复点击
		}
		folderName := utils.OpenWinFolder()
		if folderName == "" {
			return
		}

		f.downloading = true
		f.downloaded = 0
		f.totalSize = 0
		f.lastDownloaded = 0
		f.lastUpdateTime = time.Now().UnixNano()
		f.currentSpeed = 0
		f.button.SetFillColor(downloadingColor, downloadingColor)
		f.funcText.SetText("下载中")
		f.funcText.AddSpaceToSizeX(100)

		netfinder.DownLoadFile(folderName, f.file,
			func(downloaded int64, total int64) {
				f.downloaded = downloaded
				f.totalSize = total

				// 计算下载速度
				currentTime := time.Now().UnixNano()
				timeDiff := float64(currentTime-f.lastUpdateTime) / 1e9 // 转换为秒
				if timeDiff > 0.1 {                                     // 至少0.1秒才更新速度
					bytesDiff := downloaded - f.lastDownloaded
					f.currentSpeed = float64(bytesDiff) / timeDiff
					f.lastDownloaded = downloaded
					f.lastUpdateTime = currentTime
				}

				if total > 0 {
					percent := int(float64(downloaded) / float64(total) * 100)
					sizeStr := formatFileSize(total)
					f.progressText.SetText(fmt.Sprintf("%d%% %s / %s", percent, formatFileSize(downloaded), sizeStr))
				}
			},
			func(status string) {
				if status == "completed" {
					f.downloading = false
					f.funcText.SetText("下载")
					f.funcText.AddSpaceToSizeX(100)
					f.button.SetFillColor(downloadColor, downloadColor)
					// f.progressText.SetText("")
				} else if status == "failed" {
					f.downloading = false
					f.funcText.SetText("下载")
					f.funcText.AddSpaceToSizeX(100)
					f.button.SetFillColor(downloadColor, downloadColor)
					// f.progressText.SetText("下载失败")
				}
			},
		)
	}
}

// 新建一个文件组件,isSelf为是否是自身的
func newFileCmp(file netfinder.File) *fileCmp {
	f := &fileCmp{
		file: file,
	}

	f.canvas = ui.NewCoreRectCanvasUiAsKid(fileColor, 10).LockSize(utils.NewPoint(fileListX-20, fileY))

	hBox := ui.NewHorizontalBox(0)

	vbox := ui.NewVerticalBox(3).LockSize(utils.NewPoint(fileListX-20-100, fileY))

	// 第一行：文件名
	f.fileNameText = vbox.AddKid(ui.NewStaticTextUiAsKid(filepath.Base(file.FileName), ui.SmallSize, fileTextColor)).(*ui.TextUi)

	// 第二行：文件来源
	f.orginText = vbox.AddKid(ui.NewStaticTextUiAsKid("-来自["+file.Id+"]", ui.SmallSize, fileTextColor)).(*ui.TextUi)

	// 第三行：进度文本（默认显示 "-"）
	f.progressText = vbox.AddKid(ui.NewStaticTextUiAsKid("-", ui.SmallSize, fileProgressColor)).(*ui.TextUi)

	hBox.AddKid(vbox)
	var funcText string
	var funcColor color.Color

	// 接收或者关闭按钮
	if file.Id == netfinder.Id() {
		// 自己的文件，关闭按钮
		funcText = "关闭"
		funcColor = utils.ColorByColorI(closeFileColor)
	} else {
		// 别人的文件，下载按钮
		funcText = "下载"
		funcColor = utils.ColorByColorI(downloadColor)
	}

	f.funcCanvas = ui.NewRoundLerpRectCanvasUiAsKid(funcColor, color.White, 0).LockSize(utils.NewPoint(100, 50))
	hhbox := ui.NewHorizontalBox(0).LockSize(utils.NewPoint(100, 50))
	f.funcText = hhbox.AddKid(ui.NewStaticTextUiAsKid(funcText, ui.SmallSize, color.White).AddSpaceToSizeX(100)).(*ui.TextUi)
	f.funcCanvas.AddKid(hhbox)
	bCa := hBox.AddKid(f.funcCanvas)

	f.button = ui.NewButtonByCanvas(bCa.(*ui.CanvasUi), funcColor, funcColor, gametime.BigTimerType)
	f.button.SetCheckKey(input.GameMainReleasedAction, f.buttonFunc)
	f.canvas.AddKid(hBox)

	return f

}

func (f *fileCmp) update(checkPos utils.Point) {
	f.button.Update(checkPos)
}
