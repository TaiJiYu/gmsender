package ui

import (
	"gmsender/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

// 垂直框,纯粹的布局容器
type VerticalBox struct {
	pos       utils.Point
	ali       utils.AlignmentType
	offset    utils.Point
	drawOrder DrawOrder
	size      utils.Point
	otherY    float64 // 包含所有的kidSpace
	kidSpace  float64
	kids      []UICompentI
	lockSize  bool
}

// kidSpace为垂直元素的间隔
func newVerticalBox(kidSpace float64) *VerticalBox {
	return &VerticalBox{
		kidSpace: kidSpace,
		kids:     make([]UICompentI, 0),
	}
}
func (h *VerticalBox) LockSize(size utils.Point) *VerticalBox {
	h.size = size
	h.lockSize = true
	return h
}

// 删除子组件
func (h *VerticalBox) DelKid(u UICompentI) {
	if !h.lockSize {
		h.size.Y -= u.Size().Y
		h.size.Y = max(0, h.size.Y-h.kidSpace)
		h.otherY = max(0, float64(len(h.kids)-2)*h.kidSpace)
	}
	y := 0.0
	delI := -1
	for i := range h.kids {
		if h.kids[i] == u {
			delI = i
			continue
		}

		sy := h.kids[i].Size().Y
		h.kids[i].SetPos(utils.Point{X: 0, Y: y}, utils.LL)
		h.kids[i].SetPosOffset(h.offset)
		y = max(y, y+sy+h.kidSpace)
	}
	if delI >= 0 {
		h.kids = append(h.kids[:delI], h.kids[delI+1:]...)
	}
}

// 添加子项
// 注意：请先给组件add完所有子组件后再调用其父组件add自身
func (h *VerticalBox) AddKid(u UICompentI) UICompentI {
	// 子组件在水平框中，x居左，y居中
	h.kids = append(h.kids, u)
	if !h.lockSize {
		h.size.Y += u.Size().Y
		if len(h.kids) > 1 {
			h.size.Y += h.kidSpace
			h.otherY = float64(len(h.kids)-1) * h.kidSpace
		}
		h.size.X = max(h.size.X, u.Size().X)
	}
	y := 0.0
	for i := range h.kids {
		sy := h.kids[i].Size().Y
		h.kids[i].SetPos(utils.Point{X: 0, Y: y}, utils.LL)
		h.kids[i].SetPosOffset(h.offset)
		y = max(y, y+sy+h.kidSpace)
	}
	return u
}

// 设置垂直间隔，为了性能，后续不会重新计算尺寸
func (h *VerticalBox) SetKidSpace(kidSpace float64) {
	h.kidSpace = kidSpace

	y := 0.0
	for i := range h.kids {
		sy := h.kids[i].Size().Y
		h.kids[i].SetPos(utils.Point{X: 0, Y: y}, utils.LL)
		h.kids[i].SetPosOffset(h.offset)
		y += sy + h.kidSpace
	}
}

func (h *VerticalBox) SetPos(pos utils.Point, ali utils.AlignmentType) {
	h.pos = pos
	h.ali = ali
	h.SetPosOffset(h.offset)
}

func (h *VerticalBox) SetPosOffset(offset utils.Point) {
	h.offset = offset
	for i := range h.kids {
		h.kids[i].SetPosOffset(offset.Add(h.pos))
	}
}

func (h *VerticalBox) Size() utils.Point {
	return h.size.AddY(h.otherY).MaxY(0)
}

func (h *VerticalBox) Draw(screen *ebiten.Image) {
	switch h.drawOrder {
	case ADrawOrder:
		for i := range h.kids {
			h.kids[i].Draw(screen)
		}
	case ReDrawOrder:
		for i := len(h.kids) - 1; i >= 0; i-- {
			h.kids[i].Draw(screen)
		}
	}
}

// 设置渲染顺序,o只影响渲染顺序
func (h *VerticalBox) SetDrawOrder(o DrawOrder) {
	h.drawOrder = o
}
