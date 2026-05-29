package ui

import (
	gametime "gmsender/pkg/game_time"
	"gmsender/pkg/ui/shader"
	"gmsender/utils"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

func Init() {
	shader.InitShader()
	initTextUi()
	initSpitLint()
}

func NewImgUI(img *ebiten.Image, pos utils.Point, ali utils.AlignmentType) *ImgUi {
	return newImgUI(img, pos, ali)
}

func NewImgUIAsKid(img *ebiten.Image) *ImgUi {
	return newImgUI(img, utils.ZeroPoint, utils.LL)
}

// 新建尺寸框
func NewSizeBox(size utils.Point) *SizeBox {
	return newSizeBox(size)
}

// boundSize为边缘填充,两边是圆角，支持插值颜色
func NewRoundLerpRectCanvasUi(pos utils.Point, ali utils.AlignmentType, fillColor color.Color, backColor color.Color, boundSize float64) *CanvasUi {
	return newRoundLerpRectCanvasUi(pos, ali, fillColor, backColor, boundSize)
}

// boundSize为边缘填充,两边是圆角，支持插值颜色
func NewRoundLerpRectCanvasUiAsKid(fillColor color.Color, backColor color.Color, boundSize float64) *CanvasUi {
	return newRoundLerpRectCanvasUi(utils.ZeroPoint, utils.LL, fillColor, backColor, boundSize)
}

// boundSize为边缘填充,两边是圆角
func NewRoundRectCanvasUi(pos utils.Point, ali utils.AlignmentType, fillColor color.Color, boundSize float64) *CanvasUi {
	return newRoundRectCanvasUi(pos, ali, fillColor, boundSize)
}

// boundSize为边缘填充,两边是圆角
func NewRoundRectCanvasUiAsKid(fillColor color.Color, boundSize float64) *CanvasUi {
	return newRoundRectCanvasUi(utils.ZeroPoint, utils.LL, fillColor, boundSize)
}

// boundSize为边缘填充，四角是圆的
func NewCoreRectCanvasUi(pos utils.Point, ali utils.AlignmentType, fillColor color.Color, boundSize float64) *CanvasUi {
	return newCoreRectCanvasUi(pos, ali, fillColor, boundSize)
}

// boundSize为边缘填充，四角是圆的
func NewCoreRectCanvasUiAsKid(fillColor color.Color, boundSize float64) *CanvasUi {
	return newCoreRectCanvasUi(utils.ZeroPoint, utils.LL, fillColor, boundSize)
}

// 新建空白画布，仅包含尺寸和位置计算
func NewEmptyCanvasUi(pos utils.Point, ali utils.AlignmentType, boundSize float64) *CanvasUi {
	return newEmptyCanvasUi(pos, ali, boundSize)
}

// 新建空白画布，仅包含尺寸和位置计算
func NewEmptyCanvasUiAsKid(boundSize float64) *CanvasUi {
	return newEmptyCanvasUi(utils.ZeroPoint, utils.LL, boundSize)
}

// 新建水平框，kidSpace是子组件间隔
func NewHorizontalBox(kidSpace float64) *HorizontalBox {
	return newHorizontalBox(kidSpace)
}

// 新建垂直框，kidSpace是子组件间隔
func NewVerticalBox(kidSpace float64) *VerticalBox {
	return newVerticalBox(kidSpace)
}

// 新建分割线
func NewSplitLine(pos utils.Point, lenMax int, ali utils.AlignmentType) *SplitLine {
	return newSplitLine(pos, lenMax, ali)
}

// 新建分割线
func NewSplitLineAsKid(lenMax int) *SplitLine {
	return newSplitLine(utils.ZeroPoint, lenMax, utils.LL)
}

//	func NewTextUi(textKey asset.LocationKey, size FontSize, pos utils.Point, ali utils.AlignmentType, col color.Color) *TextUi {
//		return newTextUi(textKey, size, pos, ali, col)
//	}
//
//	func NewTextUiAsKid(textKey asset.LocationKey, size FontSize, col color.Color) *TextUi {
//		return newTextUi(textKey, size, utils.ZeroPoint, utils.LL, col)
//	}
func NewStaticTextUi(text string, size FontSize, pos utils.Point, ali utils.AlignmentType, color color.Color) *TextUi {
	return newStaticTextUi(text, size, pos, ali, color)
}
func NewStaticTextUiAsKid(text string, size FontSize, color color.Color) *TextUi {
	return newStaticTextUi(text, size, utils.ZeroPoint, utils.LL, color)
}

// func NewTextUiByFormat(formats string, textKeys []asset.LocationKey, size FontSize, pos utils.Point, ali utils.AlignmentType, color color.Color) *TextUi {
// 	return newTextUiByFormat(formats, textKeys, size, pos, ali, color)
// }

// func NewTextUiByFormatAsKid(formats string, textKeys []asset.LocationKey, size FontSize, color color.Color) *TextUi {
// 	return newTextUiByFormat(formats, textKeys, size, utils.ZeroPoint, utils.LL, color)
// }

func NewButton(img *ebiten.Image, pos utils.Point, ali utils.AlignmentType, timer gametime.TimerType) *ButtonUi {
	return newButton(newImgRender(img, pos, ali), pos, timer)
}

// 使用画布作为渲染器,staticFillColor,inFillColor分别为鼠标不在上面与在上面的填充色，会覆盖画布本身的填充色
func NewButtonByCanvas(canvas *CanvasUi, staticFillColor, inFillColor color.Color, timer gametime.TimerType) *ButtonUi {
	canvas.SetFillColor(staticFillColor)
	return newButton(newCanvasRender(canvas, staticFillColor, inFillColor), canvas.drawPos(), timer)
}
