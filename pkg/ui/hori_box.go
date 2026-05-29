package ui

import (
	"gmsender/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

// 水平框,纯粹的布局容器
type HorizontalBox struct {
	pos      utils.Point
	ali      utils.AlignmentType
	offset   utils.Point
	size     utils.Point
	lockSize bool
	otherX   float64 // 多余的x，包含所有的kidSpace
	kidSpace float64
	kids     []UICompentI
}

// kidSpace为水平元素的间隔
func newHorizontalBox(kidSpace float64) *HorizontalBox {
	return &HorizontalBox{
		kidSpace: kidSpace,
		kids:     make([]UICompentI, 0),
	}
}

func (h *HorizontalBox) Pos() utils.Point {
	return h.offset
}

func (h *HorizontalBox) LockSize(size utils.Point) *HorizontalBox {
	h.size = size
	h.lockSize = true
	return h
}

// 返回子对象个数
func (h *HorizontalBox) KidsCount() int {
	return len(h.kids)
}

// 从末尾移除指定个数的子对象,不做边界检查，请自行确认
func (h *HorizontalBox) RemoveKidsCount(count int) {
	h.kids = h.kids[:len(h.kids)-count]
}

// 添加子项,xAli为x的对齐方式，只有当锁定尺寸后xAli才会生效，默认左对齐，只支持让最后一个右对齐
// 注意：请先给组件add完所有子组件后再调用其父组件add自身
func (h *HorizontalBox) AddKid(u UICompentI, xAli ...utils.SingleAlignmentType) UICompentI {
	// 子组件在水平框中，x居左，y居中
	h.kids = append(h.kids, u)
	if !h.lockSize {
		h.size.X += u.Size().X
		if len(h.kids) > 1 {
			h.size.X += h.kidSpace
			h.otherX = float64(len(h.kids)-1) * h.kidSpace
		}
		h.size.Y = max(h.size.Y, u.Size().Y)
	}
	x := 0.0
	for i := range h.kids {
		sx, sy := h.kids[i].Size().Break()
		h.kids[i].SetPos(utils.Point{X: x, Y: (h.size.Y - sy) / 2}, utils.LL)
		h.kids[i].SetPosOffset(h.offset)
		x += sx + h.kidSpace
	}
	if len(xAli) > 0 {
		switch xAli[0] {
		case utils.R:
			u.SetPos(utils.Point{X: h.Size().X, Y: (h.size.Y - u.Size().Y) / 2}, utils.RL)
		}
	}

	return u
}

func (h *HorizontalBox) SetPos(pos utils.Point, ali utils.AlignmentType) {
	h.pos = pos
	h.ali = ali
	h.SetPosOffset(h.offset)
}

func (h *HorizontalBox) SetPosOffset(offset utils.Point) {
	h.offset = offset
	for i := range h.kids {
		h.kids[i].SetPosOffset(offset.Add(h.pos))
	}
}

func (h *HorizontalBox) Size() utils.Point {
	return h.size.AddX(h.otherX)
}

func (h *HorizontalBox) Draw(screen *ebiten.Image) {
	for i := range h.kids {
		h.kids[i].Draw(screen)
	}
}
