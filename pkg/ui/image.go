package ui

import (
	"gmsender/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

// 图片ui组件
type ImgUi struct {
	offset   utils.Point
	img      *ebiten.Image
	pos      utils.Point
	ali      utils.AlignmentType
	op       *ebiten.DrawImageOptions
	halfSize utils.Point

	drawPos utils.Point
}

func newImgUI(img *ebiten.Image, pos utils.Point, ali utils.AlignmentType) *ImgUi {
	i := &ImgUi{
		img:      img,
		pos:      pos,
		ali:      ali,
		halfSize: utils.NewPoint(float64(img.Bounds().Dx())/2, float64(img.Bounds().Dy())/2),
		op:       &ebiten.DrawImageOptions{},
	}
	i.setOp(utils.ZeroPoint)
	return i
}

func (i *ImgUi) SetImg(img *ebiten.Image) {
	i.img = img
	i.halfSize = utils.NewPoint(float64(img.Bounds().Dx())/2, float64(img.Bounds().Dy())/2)
	i.setOp(i.offset)
}

// 设置坐标偏移
func (i *ImgUi) SetPosOffset(offset utils.Point) {
	i.offset = offset
	i.setOp(offset)
}

func (i *ImgUi) SetPos(pos utils.Point, ali utils.AlignmentType) {
	i.pos = pos
	i.ali = ali
	i.setOp(i.offset)
}

func (i *ImgUi) setOp(offset utils.Point) {
	pos := i.ali.GetAlignmentPos(i.pos, i.halfSize).Add(offset) // 渲染坐标
	i.drawPos = pos
	i.op.GeoM.Reset()
	i.op.GeoM.Translate(pos.Break())
}

func (i *ImgUi) Draw(screen *ebiten.Image) {
	screen.DrawImage(i.img, i.op)
}

func (i *ImgUi) Size() utils.Point {
	return i.halfSize.MulF1(2)
}
