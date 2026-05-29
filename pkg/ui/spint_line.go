package ui

import (
	"gmsender/utils"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	lineImg *ebiten.Image
)

func initSpitLint() {
	lineImg = ebiten.NewImage(utils.LogicalSizeX, 2)
	lineImg.Fill(color.Black)
}

// 分割线，纯粹的图形
type SplitLine struct {
	offset   utils.Point
	img      *ebiten.Image
	pos      utils.Point
	ali      utils.AlignmentType
	halfSize utils.Point
	op       *ebiten.DrawImageOptions
}

// kidSpace为垂直元素的间隔
func newSplitLine(pos utils.Point, lenMax int, ali utils.AlignmentType) *SplitLine {
	l := &SplitLine{
		pos:      pos,
		img:      lineImg.SubImage(image.Rect(0, 0, lenMax, 2)).(*ebiten.Image),
		halfSize: utils.NewPoint(float64(lenMax)/2, float64(lineImg.Bounds().Dy())/2),
		op:       &ebiten.DrawImageOptions{},
		ali:      ali,
	}
	l.setOp(utils.ZeroPoint)
	return l
}

func (l *SplitLine) SetPos(pos utils.Point, ali utils.AlignmentType) {
	l.pos = pos
	l.ali = ali
	l.setOp(l.offset)
}

func (l *SplitLine) SetPosOffset(offset utils.Point) {
	l.offset = offset
	l.setOp(offset)
}

func (l *SplitLine) setOp(offset utils.Point) {
	pos := l.ali.GetAlignmentPos(l.pos, l.halfSize).Add(offset) // 渲染坐标
	l.op.GeoM.Reset()
	l.op.GeoM.Translate(pos.Break())
}

// 分割线的size中高度始终看做0
func (l *SplitLine) Size() utils.Point {
	return l.halfSize.MulF(2, 0)
}

func (l *SplitLine) Draw(screen *ebiten.Image) {
	screen.DrawImage(l.img, l.op)
}
