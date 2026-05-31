package ui

import (
	gametime "gmsender/pkg/game_time"
	statemachine "gmsender/pkg/state_machine"
	"gmsender/utils"
	"image/color"

	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type ButtonUi struct {
	buttonRender buttonRenderI
	// img          *asset.ImgAsset
	// op       *camera.CameraDrawImgOp
	pos utils.Point // 逻辑坐标
	// ali      utils.AlignmentType
	isUseing bool // 是否有效，只有有效时才会做检测判断

	checkKey            HotKeyI                     // 需要检查鼠标或者逻辑坐标是否在其上面
	isInRange           bool                        // 判定点是否在按钮范围内
	event               func(bu *ButtonUi)          // 判定满足时触发
	dropIngPos          func(*ButtonUi) utils.Point // 拖动控制
	dropAfrerDo         func(*ButtonUi)             // 单次拖动后的操作
	dropBeginLogicalPos utils.Point                 // 拖动开始时的逻辑坐标，由拖动开始时的dropIngPos决定
	dropKey             HotKeyI                     // 拖动按键
	allowDrop           bool                        // 允许拖动
	isDroping           bool                        // 是否正在拖动
	dropBeginPos        utils.Point                 // 开始拖动的位置
	enterEvent          func(*ButtonUi)             // 进入事件，当判定点进入时触发

	scale      float64 // 缩放
	scaleState statemachine.StateMachineI
}

const (
	buttonScaleDone = iota // 缩放完成
	buttonScaleUp          // 缩放增加
	buttonScaleDown        // 缩放减少

	scaleDuration = 200 * time.Millisecond
)

var (
	scaleDurationSec = scaleDuration.Seconds()
)

// 适配手柄的按钮控制
func newButton(render buttonRenderI, pos utils.Point, timer gametime.TimerType) *ButtonUi {
	b := &ButtonUi{
		isUseing:     true,
		pos:          pos,
		buttonRender: render,
		scale:        1,
		scaleState:   statemachine.NewStateMachine(timer),
	}
	b.scaleState.NewState(buttonScaleDone)

	scaleBegin := 0.0
	b.scaleState.NewState(buttonScaleUp).SetEnterFunc(func() {
		scaleBegin = b.scale
	}).BindU(func() {
		t := b.scaleState.ReadStateLastTimeSec() / scaleDurationSec
		b.scale = utils.Lerp(scaleBegin, 1.5, t)
	})
	b.scaleState.NewState(buttonScaleDown).SetEnterFunc(func() {
		scaleBegin = b.scale
	}).BindU(func() {
		t := b.scaleState.ReadStateLastTimeSec() / scaleDurationSec
		b.scale = utils.Lerp(scaleBegin, 1, t)
	})

	b.scaleState.SToSWithTimeLimit(buttonScaleUp, buttonScaleDone, scaleDuration)
	b.scaleState.SToSWithTimeLimit(buttonScaleDown, buttonScaleDone, scaleDuration)

	b.scaleState.Go(buttonScaleDone)

	return b
}
func (b *ButtonUi) SetFillColor(fillColor, inColor color.Color) {
	if r, ok := b.buttonRender.(*buttonCanvasRender); ok {
		r.setFillColor(fillColor, inColor)
	}
}

type HotKeyI interface {
	Check() bool
}

// func (b *ButtonUi) SetImg(img *ebiten.Image) {
// 	if render, ok := b.buttonRender.(*buttonImgRender); ok {
// 		render.setImg(img)
// 		render.drawPos(b.pos)
// 	}
// }

// 需要检查位置鼠标或者逻辑坐标是否在其上面
func (b *ButtonUi) SetCheckKey(key HotKeyI, event func(bu *ButtonUi)) *ButtonUi {
	b.checkKey = key
	b.event = event
	return b
}

// 设置拖动行为,拖动时会将坐标绑定到拖动，如果拖动发生，则触发event不生效,afterDo为拖动单次结束后要做的额外的事，例如根据按钮拖动后的真实坐标做纠正
func (b *ButtonUi) SetDropEvent(key HotKeyI, dropIngPos func(*ButtonUi) utils.Point, afterDo func(*ButtonUi)) {
	b.dropIngPos = dropIngPos
	b.dropKey = key
	b.dropAfrerDo = afterDo
	b.allowDrop = key != nil && dropIngPos != nil
}

func (b *ButtonUi) eventCheck() {
	if b.checkKey != nil && b.checkKey.Check() && b.event != nil {
		// 命中了
		b.event(b)
	}
}

func (b *ButtonUi) SetPos(p utils.Point) {
	b.pos = p
}

func (b *ButtonUi) Pos() utils.Point { return b.pos }

// // 返回当前的缩放值
// func (b *ButtonUi) Scale() float64 { return camera.Scale() * b.scale }

// 判定是否在范围内
func (b *ButtonUi) IsInRange() bool { return b.isInRange }

// 设置进入时的事件
func (b *ButtonUi) SetEnterEvent(enterEvent func(*ButtonUi)) {
	b.enterEvent = enterEvent
}

// checkPos是用于检查判定坐标的，可以提供鼠标坐标或者手柄逻辑坐标
func (b *ButtonUi) Update(checkPos utils.Point) {
	llpos := b.buttonRender.drawPos(b.pos)

	if b.isDroping {
		// 正在拖动不做范围检查
		if b.dropKey.Check() {
			//正在拖动且拖动按键触发
			b.pos = b.dropIngPos(b).Sub(b.dropBeginLogicalPos).Add(b.dropBeginPos)
			if b.dropAfrerDo != nil {
				b.dropAfrerDo(b)
			}
		} else {
			// 拖动结束
			b.isDroping = false
			// 如果小于5则认为没有拖动则可以触发命中检查
			if b.dropBeginPos.Sub(b.pos).IsLenLess(5) {
				b.eventCheck()
			}
		}
	} else {
		// 范围检查
		if b.buttonRender.check(llpos, checkPos) {
			// 范围内
			if !b.isInRange {
				// 进入范围
				b.scaleState.Go(buttonScaleUp)
				if b.enterEvent != nil {
					// 额外的进入事件
					b.enterEvent(b)
				}
			}
			b.isInRange = true
			if b.allowDrop && b.dropKey.Check() {
				// 允许拖动
				b.isDroping = true
				b.dropBeginPos = b.pos
				b.dropBeginLogicalPos = b.dropIngPos(b)
			} else {
				b.eventCheck()
				b.isDroping = false
			}
		} else {
			if b.isInRange {
				// 退出范围
				b.scaleState.Go(buttonScaleDown)
			}
			b.isInRange = false
		}
	}

	b.scaleState.Update()
	b.buttonRender.setScaleK(b.scale)
}

// func (b *ButtonUi) SetDrawColorScale(clr color.Color) {
// 	b.op.SetDrawColorScale(clr)
// }

// func (b *ButtonUi) ResetDrawColorScale() {
// 	b.op.ResetDrawColorScale()
// }

func (b *ButtonUi) Draw(screen *ebiten.Image) {
	b.buttonRender.draw(screen)
}
