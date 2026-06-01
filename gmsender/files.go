package gmsender

import (
	"container/list"
	gametime "gmsender/pkg/game_time"
	"gmsender/pkg/input"
	"gmsender/pkg/netfinder"
	"gmsender/pkg/ui"
	"gmsender/utils"
	"image/color"
	"path/filepath"
)

const (
	// 文件列表尺寸
	fileListX = utils.LogicalSizeX - 40*2
	// 单个文件高度
	fileY = 50
)

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

	fileNameText, orginText, funcText *ui.TextUi //funcText是功能文本

	file netfinder.File

	isDel bool
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
		netfinder.DownLoadFile(utils.OpenWinSaveFileName(), f.file)
	}
}

// 新建一个文件组件,isSelf为是否是自身的
func newFileCmp(file netfinder.File) *fileCmp {
	f := &fileCmp{
		file: file,
	}

	f.canvas = ui.NewCoreRectCanvasUiAsKid(fileColor, 10).LockSize(utils.NewPoint(fileListX-20, fileY))

	hBox := ui.NewHorizontalBox(0)

	vbox := ui.NewVerticalBox(5).LockSize(utils.NewPoint(fileListX-20-100, fileY))
	f.fileNameText = vbox.AddKid(ui.NewStaticTextUiAsKid(filepath.Base(file.FileName), ui.SmallSize, fileTextColor)).(*ui.TextUi)
	f.orginText = vbox.AddKid(ui.NewStaticTextUiAsKid("-来自["+file.Id+"]", ui.SmallSize, fileTextColor)).(*ui.TextUi)

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
