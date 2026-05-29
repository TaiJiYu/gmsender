package input

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// 所有按键控制

// 系统层面
var (
	// 设置
	// GameExitAction = NewAction(PressedType).AddKey(ebiten.KeyEscape).AddPad(ebiten.StandardGamepadButtonCenterRight)

	// 鼠标左键释放
	GameMainReleasedAction = NewAction(ReleasedType).AddMouse(ebiten.MouseButtonLeft)
	// // 鼠标左键与手柄的A按住
	// GameMainPressHoldAction = NewAction(PressHold).AddMouse(ebiten.MouseButtonLeft)
	// // 键盘esc与手柄的B松开
	// GameReturnAction = NewAction(ReleasedType).AddKey(ebiten.KeyEscape)
)
