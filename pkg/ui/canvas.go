package ui

import (
	"gmsender/pkg/ui/shader"
	"gmsender/utils"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// 背景是圆角矩形的画布ui
type CanvasUi struct {
	offset    utils.Point
	pos       utils.Point
	drawPosLL utils.Point
	ali       utils.AlignmentType
	kids      []UICompentI
	size      utils.Point
	lockSize  bool
	rectGeo   rectShdaerI
	boundSize float64
}

type rectShdaerI interface {
	Draw(screen *ebiten.Image)
	SetPos(pos utils.Point, ali utils.AlignmentType) utils.Point
	SetSize(size utils.Point)
	SetFillColor(color.Color)  // 设置填充色
	SetBoundColor(color.Color) // 设置边缘色
	SetLerp(t float64)         // 设置颜色插值
}

// boundSize为边缘填充,两边是圆的,支持颜色插值
func newRoundLerpRectCanvasUi(pos utils.Point, ali utils.AlignmentType, fillColor, backColor color.Color, boundSize float64) *CanvasUi {
	return &CanvasUi{
		pos:       pos,
		ali:       ali,
		kids:      make([]UICompentI, 0),
		rectGeo:   shader.NewRoundLerpRect(pos, ali, fillColor, backColor),
		boundSize: boundSize,
	}
}

// boundSize为边缘填充,两边是圆的
func newRoundRectCanvasUi(pos utils.Point, ali utils.AlignmentType, fillColor color.Color, boundSize float64) *CanvasUi {
	return &CanvasUi{
		pos:       pos,
		ali:       ali,
		kids:      make([]UICompentI, 0),
		rectGeo:   shader.NewRoundRect(pos, ali, fillColor),
		boundSize: boundSize,
	}
}

// boundSize为边缘填充，四角是圆的
func newCoreRectCanvasUi(pos utils.Point, ali utils.AlignmentType, fillColor color.Color, boundSize float64) *CanvasUi {
	return &CanvasUi{
		pos:       pos,
		ali:       ali,
		kids:      make([]UICompentI, 0),
		rectGeo:   shader.NewCoreRect(pos, ali, fillColor),
		boundSize: boundSize,
	}
}

// 空白画布，仅控制布局
func newEmptyCanvasUi(pos utils.Point, ali utils.AlignmentType, boundSize float64) *CanvasUi {
	return &CanvasUi{
		pos:       pos,
		ali:       ali,
		kids:      make([]UICompentI, 0),
		rectGeo:   shader.NewEmptyRect(pos, ali),
		boundSize: boundSize,
	}
}

// ui组件接口
type UICompentI interface {
	Size() utils.Point // 尺寸
	Draw(screen *ebiten.Image)
	SetPosOffset(offset utils.Point)
	SetPos(pos utils.Point, ali utils.AlignmentType)
}

// 设置填充色
func (c *CanvasUi) SetFillColor(fillColor color.Color) {
	c.rectGeo.SetFillColor(fillColor)
}

// 设置插值颜色，小于t的使用fillColor，大于t的使用backColor
func (c *CanvasUi) SetLerpColor(t float64) {
	c.rectGeo.SetLerp(t)
}

// 设置填充色
func (c *CanvasUi) BoundSize() float64 {
	return c.boundSize
}

// 锁定尺寸，尺寸固定，不受子对象影响
func (c *CanvasUi) LockSize(size utils.Point) *CanvasUi {
	c.size = size
	c.lockSize = true
	c.rectGeo.SetSize(c.size.AddFToXY(c.boundSize * 2))
	return c
}

// 返回子对象个数
func (c *CanvasUi) KidsCount() int {
	return len(c.kids)
}

// 从末尾移除指定个数的子对象,不做边界检查，请自行确认
func (c *CanvasUi) RemoveKidsCount(count int) {
	c.kids = c.kids[:len(c.kids)-count]
}

// 清空子类
func (c *CanvasUi) ClearKids() {
	c.kids = c.kids[:0]
	// 不考虑尺寸变化，因为一般用于锁定尺寸的画布
}

// 增加子ui组件,画布的子组件添加为填充
// 注意：请先给组件add完所有子组件后再调用其父组件add自身
func (c *CanvasUi) AddKid(u UICompentI) UICompentI {
	c.kids = append(c.kids, u)
	if !c.lockSize {
		c.size = u.Size().Max(c.size)
		c.rectGeo.SetSize(c.size.AddFToXY(c.boundSize * 2))
	}
	c.SetPos(c.pos, c.ali)
	return u
}

// // 设置x的归一化缩放，从0-1,0时是0,1时代表真实尺寸
// func (c *CanvasUi) SetNormalSizeX(f float64) {
// 	c.rectGeo.SetSize(c.size.AddFToXY(c.boundSize*2).MulF(f, 1))
// 	c.SetPos(c.pos, c.ali)
// }

func (c *CanvasUi) SetPos(pos utils.Point, ali utils.AlignmentType) {
	c.pos = pos
	c.ali = ali
	c.setOp(c.offset)
}

// 设置坐标偏移
func (c *CanvasUi) SetPosOffset(offset utils.Point) {
	c.offset = offset
	c.setOp(offset)
}

func (c *CanvasUi) setOp(offset utils.Point) {
	c.drawPosLL = c.rectGeo.SetPos(c.pos.Add(offset), c.ali).AddFToXY(c.boundSize)
	for i := range c.kids {
		c.kids[i].SetPosOffset(c.drawPosLL)
	}
}

// 渲染坐标，左上角的坐标
func (c *CanvasUi) drawPos() utils.Point {
	return c.drawPosLL
}

func (c *CanvasUi) Draw(screen *ebiten.Image) {
	c.rectGeo.Draw(screen)
	for i := range c.kids {
		c.kids[i].Draw(screen)
	}
}

// 尺寸
func (c *CanvasUi) Size() utils.Point {
	return c.size.AddFToXY(c.boundSize * 2)
}
