package ui

import (
	"gmsender/utils"
	"image"
	"image/color"

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

	nowSliderF       float64         // 当前滑动位置
	subRect          image.Rectangle // 遮罩范围
	sliderYRange     float64         // 滑动y的范围
	sliderHintCanvas *CanvasUi       // 滑动提示
	sliderDraw       func(screen *ebiten.Image)
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

	h.setSliderRect()
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
	if h.sliderDraw != nil {
		h.sliderDraw(screen)
		return
	}
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

// 滑动距离，正数是展示更多下面的内容
func (h *VerticalBox) Slider(f float64) {
	if h.Size().Y <= h.sliderYRange {
		// 能否展示完全，不能滑动
		h.offset = h.offset.SubY(-h.nowSliderF)
		h.nowSliderF = 0
		h.SetPosOffset(h.offset)
		return
	}

	h.nowSliderF += f

	h.offset = h.offset.SubY(f)
	h.SetPosOffset(h.offset)

	h.setSliderRect()
}

func (h *VerticalBox) setSliderRect() {
	llx, lly := h.pos.Add(h.offset).AddY(h.nowSliderF).BreakInt()
	h.subRect = image.Rect(llx, lly, llx+int(h.Size().X), lly+int(h.sliderYRange))
	if h.sliderHintCanvas != nil {
		h.sliderHintCanvas.SetPos(utils.NewPoint(h.subRect.Dx()/2+h.subRect.Min.X, h.subRect.Max.Y-20), utils.MM)
	}

}

// 检查检查点是否在渲染范围内
func (h *VerticalBox) CheckMouseInSlider(checkPos utils.Point) bool {
	return checkPos.IsRangeIn(utils.NewPoint(h.subRect.Min.X, h.subRect.Max.X), utils.NewPoint(h.subRect.Min.Y, h.subRect.Max.Y))
}

// 设置滑动窗口，超过尺寸则滑动
func (h *VerticalBox) SetSlider(sizeY float64, hintColor color.Color) {
	h.sliderYRange = sizeY
	screenCache := ebiten.NewImage(utils.LogicalSize.BreakInt())
	op := &ebiten.DrawImageOptions{}
	h.sliderHintCanvas = newRoundRectCanvasUi(utils.NewPoint(h.subRect.Dx()/2+h.subRect.Min.X, h.subRect.Max.Y-20), utils.MM, hintColor, 0).LockSize(utils.NewPoint(46, 22))
	h.sliderHintCanvas.AddKid(NewStaticTextUiAsKid("···", SmallSize, color.White))
	h.setSliderRect()
	h.sliderDraw = func(screen *ebiten.Image) {
		screenCache.Clear()
		for i := range h.kids {
			h.kids[i].Draw(screenCache)
		}
		op.GeoM.Reset()
		op.GeoM.Translate(float64(h.subRect.Min.X), float64(h.subRect.Min.Y))
		if h.Size().Y > h.sliderYRange {
			// 尺寸过大展示
			h.sliderHintCanvas.Draw(screenCache)
		}

		screen.DrawImage(screenCache.SubImage(h.subRect).(*ebiten.Image), op)
	}
}
