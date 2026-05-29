package gmsender

import (
	"container/list"
	"fmt"
	gametime "gmsender/pkg/game_time"
	"gmsender/pkg/input"
	"gmsender/pkg/ui"
	"gmsender/utils"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	// 文件列表尺寸
	fileListX = utils.LogicalSizeX - 40*2
	// 单个文件高度
	fileY = 50
)

// 公开文件列表
type fileList struct {
	files *list.List
}

var fileListCli *fileList

func InitFileList() {
	fileListCli = &fileList{
		files: list.New(),
	}
}

func (l *fileList) AddFile(isSelf bool, fileName, orgin string) *ui.CanvasUi {
	f := newFileCmp(isSelf, fileName, orgin)
	l.files.PushBack(f)
	return f.canvas
}
func (l *fileList) Update(checkPos utils.Point) {
	for e := l.files.Front(); e != nil; e = e.Next() {
		f := e.Value.(*fileCmp)
		if f.isDel {
			l.files.Remove(e)
			delFileCmp(f.canvas)
		} else {
			f.update(checkPos)
		}
	}
}
func (l *fileList) Draw(screen *ebiten.Image) {
	for e := l.files.Front(); e != nil; e = e.Next() {
		e.Value.(*fileCmp).draw(screen)
	}
}

// 列表内的单个组件
type fileCmp struct {
	fileName, orgin string
	canvas          *ui.CanvasUi
	button          *ui.ButtonUi

	isDel bool
}

// 新建一个文件组件,isSelf为是否是自身的
func newFileCmp(isSelf bool, fileName, orgin string) *fileCmp {
	f := &fileCmp{
		fileName: fileName,
		orgin:    orgin,
	}

	f.canvas = ui.NewCoreRectCanvasUiAsKid(fileColor, 10).LockSize(utils.NewPoint(fileListX-20, fileY))

	hBox := ui.NewHorizontalBox(0)

	vbox := ui.NewVerticalBox(5).LockSize(utils.NewPoint(fileListX-20-100, fileY))
	vbox.AddKid(ui.NewStaticTextUiAsKid(fileName, ui.SmallSize, fileTextColor))
	vbox.AddKid(ui.NewStaticTextUiAsKid("-来自["+orgin+"]", ui.SmallSize, fileTextColor))

	hBox.AddKid(vbox)
	// 接收或者关闭按钮
	if isSelf {
		// 自己的文件，关闭按钮
		closeCanvas := ui.NewRoundLerpRectCanvasUiAsKid(closeFileColor, utils.ColorRGBByOx(0xD93125), 0).LockSize(utils.NewPoint(100, 50))
		hhbox := ui.NewHorizontalBox(0).LockSize(utils.NewPoint(100, 50))
		hhbox.AddKid(ui.NewStaticTextUiAsKid("关闭", ui.SmallSize, color.White).AddSpaceToSizeX(100))
		closeCanvas.AddKid(hhbox)
		bCa := hBox.AddKid(closeCanvas)
		f.button = ui.NewButtonByCanvas(bCa.(*ui.CanvasUi), closeFileColor, utils.ColorRGBByOx(0xD93125), gametime.BigTimerType)
		f.button.SetCheckKey(input.GameMainReleasedAction, func(bu *ui.ButtonUi) {
			f.isDel = true
		})
	} else {
		// 别人的文件，下载按钮
		downloadCanvas := ui.NewRoundLerpRectCanvasUiAsKid(downloadColor, color.White, 0).LockSize(utils.NewPoint(100, 50))
		hhbox := ui.NewHorizontalBox(0).LockSize(utils.NewPoint(100, 50))
		hhbox.AddKid(ui.NewStaticTextUiAsKid("下载", ui.SmallSize, color.White).AddSpaceToSizeX(100))
		downloadCanvas.AddKid(hhbox)
		bCa := hBox.AddKid(downloadCanvas)
		f.button = ui.NewButtonByCanvas(bCa.(*ui.CanvasUi), downloadColor, color.White, gametime.BigTimerType)
		f.button.SetCheckKey(input.GameMainReleasedAction, func(bu *ui.ButtonUi) {
			fmt.Println("下载")
		})
	}

	f.canvas.AddKid(hBox)

	return f

}

func (f *fileCmp) update(checkPos utils.Point) {
	f.button.Update(checkPos)
}
func (f *fileCmp) draw(screen *ebiten.Image) {
	f.canvas.Draw(screen)
}
