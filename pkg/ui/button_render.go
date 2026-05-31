package ui

import (
	"gmsender/utils"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// 按钮渲染器
type buttonRenderI interface {
	draw(screen *ebiten.Image)
	drawPos(pos utils.Point) utils.Point // 提供按钮的逻辑坐标，返回实际渲染的左上角坐标
	size() utils.Point
	check(llPos, checkPos utils.Point) bool
	setScaleK(s float64)
}

type buttonImgRender struct {
	img      *ebiten.Image
	op       *ebiten.DrawImageOptions
	xysize   utils.Point
	scaleK   float64
	ali      utils.AlignmentType
	halfSize utils.Point
}

// 使用图片资产渲染
func newImgRender(img *ebiten.Image, pos utils.Point, ali utils.AlignmentType) *buttonImgRender {
	b := &buttonImgRender{
		img:    img,
		xysize: utils.NewPoint(img.Bounds().Dx(), img.Bounds().Dy()),
		scaleK: 1,
		ali:    ali,
	}
	b.halfSize = b.xysize.Divf1(2)
	b.op = &ebiten.DrawImageOptions{}
	b.drawPos(pos)
	// op, _ := camera.NewCameraUiDrawImgOp(pos, 0, ali, b.size().Divf1(2), 1)
	// b.op = op
	return b
}

func (b *buttonImgRender) setImg(img *ebiten.Image) {
	b.img = img
	b.xysize = utils.NewPoint(img.Bounds().Dx(), img.Bounds().Dy())
}
func (b *buttonImgRender) setScaleK(s float64) {
	b.scaleK = s
}

func (b *buttonImgRender) drawPos(pos utils.Point) utils.Point {
	llpos := b.ali.GetAlignmentPos(pos, b.halfSize)
	halfSize := b.halfSize.MulF1(b.scaleK)
	pos = b.ali.GetAlignmentPos(pos, halfSize)
	x, y := halfSize.Break()
	b.op.GeoM.Reset()
	b.op.GeoM.Scale(b.scaleK, b.scaleK)
	b.op.GeoM.Translate(-x, -y)
	b.op.GeoM.Translate(pos.X+x, pos.Y+y)
	return llpos
}

func (b *buttonImgRender) draw(screen *ebiten.Image) {
	screen.DrawImage(b.img, b.op)
	// b.op.Draw(screen, b.img)
}

func (b *buttonImgRender) size() utils.Point {
	return b.xysize
}

// 判定检查
func (b *buttonImgRender) check(llPos, checkPos utils.Point) bool {
	return checkPos.Sub(llPos).IsRangeInFloat(0, b.xysize.X, 0, b.xysize.Y)
}

type buttonCanvasRender struct {
	canvas                       *CanvasUi
	xysize                       utils.Point
	staticFillColor, inFillColor color.Color
}

// 画布渲染
func newCanvasRender(canvas *CanvasUi, staticFillColor, inFillColor color.Color) *buttonCanvasRender {
	return &buttonCanvasRender{
		staticFillColor: staticFillColor,
		inFillColor:     inFillColor,
		canvas:          canvas,
		xysize:          canvas.Size(),
	}
}

func (b *buttonCanvasRender) setFillColor(fillColor, inColor color.Color) {
	b.staticFillColor = fillColor
	b.inFillColor = inColor
	b.canvas.SetFillColor(fillColor)
}

// 对于画布渲染器，调整缩放系数时，直接修改边缘填充颜色
func (b *buttonCanvasRender) setScaleK(s float64) {
	s = (s - 1) * 2
	b.canvas.rectGeo.SetBoundColor(color.RGBA{R: uint8(255 * s), G: uint8(255 * s), B: uint8(255 * s), A: 255})
	b.canvas.SetFillColor(utils.ColorRGBLerp(b.staticFillColor, b.inFillColor, s))
}

func (b *buttonCanvasRender) drawPos(utils.Point) utils.Point {
	return b.canvas.drawPos()
}

// 渲染由画布自行渲染
func (b *buttonCanvasRender) draw(*ebiten.Image) {
	// b.op.Draw(screen, b.img)
}

func (b *buttonCanvasRender) size() utils.Point {
	return b.xysize
}

// 判定检查
func (b *buttonCanvasRender) check(llPos, checkPos utils.Point) bool {
	return checkPos.Sub(llPos).IsRangeInFloat(0, b.xysize.X, 0, b.xysize.Y)
}
