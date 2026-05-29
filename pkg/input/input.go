package input

import (
	gametime "gmsender/pkg/game_time"
	"gmsender/utils"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type InputCli struct {
	keyBuf []ebiten.Key

	keyPressedBuf   map[ebiten.Key]bool
	keyReleasedBuf  map[ebiten.Key]bool
	keyPressHoldBuf map[ebiten.Key]bool

	mousePressedBuf   map[ebiten.MouseButton]bool
	mouseReleasedBuf  map[ebiten.MouseButton]bool
	mousePressHoldBuf map[ebiten.MouseButton]bool
	wheelMoveBuf      map[keyPressedType]bool

	hasAnyKeyPressed, hasAnyKeyReleased bool // 是否有按键被按下或松开

	isPause bool // 是否暂停
}
type keyPressedType int

const (
	PressedType     keyPressedType = iota // 按下
	ReleasedType                          // 松开
	PressHold                             // 按住，只支持键盘,鼠标和手柄，不支持滑轮
	AnyPressedType                        // 任意按键按下
	AnyReleasedType                       // 任意按键松开
	MouseWheelDown                        // 鼠标滑轮向下运动
	MouseWheelUp                          // 鼠标滑轮向上运动
)

// 行为计时器，如果设置时间限制可以使用
type actionTimer struct {
	trueTimeSec  float64            // 上次触发的时间,单位秒
	timeLimitSec float64            // 多久内不能再次触发，单位秒
	timer        gametime.TimerType // 计时器类型
}

type actionKeyInfo struct {
	keyType    keyPressedType
	keys       []inputKeyI
	timerLimit *actionTimer // 时间限制
}

// 时间限制，一定时间内不能多次触发
func (a *actionKeyInfo) AddTimeLimit(t time.Duration, timer gametime.TimerType) *actionKeyInfo {
	a.timerLimit = &actionTimer{
		timeLimitSec: t.Seconds(),
		timer:        timer,
	}
	return a
}

// 时间限制，一定时间内不能多次触发
func (a *actionKeyInfo) SetTimeLimit(t time.Duration) *actionKeyInfo {
	a.timerLimit.timeLimitSec = t.Seconds()
	return a
}

// 添加键盘按键
func (a *actionKeyInfo) AddKey(k ebiten.Key) *actionKeyInfo {
	a.keys = append(a.keys, keyboardKey{
		key: k,
	})
	return a
}

// 添加鼠标按键
func (a *actionKeyInfo) AddMouse(m ebiten.MouseButton) *actionKeyInfo {
	a.keys = append(a.keys, mouseKey{
		key: m,
	})
	return a
}

// 各种按键的接口
type inputKeyI interface {
	pressedCheck(*InputCli) bool   // 按下检查
	releasedCheck(*InputCli) bool  // 松开检查
	pressHoldCheck(*InputCli) bool // 按住检查
	pressFloat(*InputCli) float64  // 读取按下的力度参数从-1到1
}

// 键盘按键
type keyboardKey struct {
	key ebiten.Key
}

func (k keyboardKey) pressedCheck(i *InputCli) bool {
	return i.checkKeysPressed(k.key)
}
func (k keyboardKey) releasedCheck(i *InputCli) bool {
	return i.checkKeysReleased(k.key)
}
func (k keyboardKey) pressHoldCheck(i *InputCli) bool {
	return i.checkKeyPressHold(k.key)
}
func (k keyboardKey) pressFloat(i *InputCli) float64 {
	return i.keyPressedFloat(k.key)
}

// 鼠标按键
type mouseKey struct {
	key ebiten.MouseButton
}

func (m mouseKey) pressedCheck(i *InputCli) bool {
	return i.checkMousePressed(m.key)
}
func (m mouseKey) releasedCheck(i *InputCli) bool {
	return i.checkMouseReleased(m.key)
}
func (m mouseKey) pressHoldCheck(i *InputCli) bool {
	return i.checkMousePressHold(m.key)
}
func (m mouseKey) pressFloat(_ *InputCli) float64 {
	return 0 // 鼠标暂不需要按住力度检查
}

// 键盘按下键
func (i *InputCli) setKeyPressed() {
	i.hasAnyKeyPressed = i.hasAnyKeyPressed || len(i.keyBuf) > 0
	for _, key := range i.keyBuf {
		i.keyPressedBuf[key] = true
	}
}

// 检查后清空
func (i *InputCli) checkKeysPressed(key ebiten.Key) bool {
	ret := i.keyPressedBuf[key]
	i.keyPressedBuf[key] = false
	return ret
}

// 键盘松开键
func (i *InputCli) setKeyReleased() {
	i.hasAnyKeyReleased = i.hasAnyKeyReleased || len(i.keyBuf) > 0
	for _, key := range i.keyBuf {
		i.keyReleasedBuf[key] = true
	}
}

// 检查后清空
func (i *InputCli) checkKeysReleased(key ebiten.Key) bool {
	ret := i.keyReleasedBuf[key]
	i.keyReleasedBuf[key] = false
	return ret
}

// 键盘按住
func (i *InputCli) setKeyPressHold() {
	for _, key := range i.keyBuf {
		i.keyPressHoldBuf[key] = true
	}
}

// 检查键盘按住
func (i *InputCli) checkKeyPressHold(key ebiten.Key) bool {
	ret := i.keyPressHoldBuf[key]
	i.keyPressHoldBuf[key] = false
	return ret
}

// 键盘按下的力度，键盘只有0和1
func (i *InputCli) keyPressedFloat(key ebiten.Key) float64 {
	if i.keyPressHoldBuf[key] {
		i.keyPressHoldBuf[key] = false
		return 1
	} else {
		return 0
	}
}

// 鼠标键
func (i *InputCli) setMouse() {
	for key := ebiten.MouseButton(0); key <= ebiten.MouseButtonMax; key++ {
		i.mousePressedBuf[key] = inpututil.IsMouseButtonJustPressed(key)
		i.hasAnyKeyPressed = i.hasAnyKeyPressed || i.mousePressedBuf[key]
		i.mouseReleasedBuf[key] = inpututil.IsMouseButtonJustReleased(key)
		i.hasAnyKeyReleased = i.hasAnyKeyReleased || i.mouseReleasedBuf[key]
		i.mousePressHoldBuf[key] = inpututil.MouseButtonPressDuration(key) > 0 // 可以灵活调整
	}
}

// 检查后清空
func (i *InputCli) checkMousePressed(b ebiten.MouseButton) (ret bool) {
	ret = i.mousePressedBuf[b]
	i.mousePressedBuf[b] = false
	return
}

// 检查后清空
func (i *InputCli) checkMouseReleased(b ebiten.MouseButton) (ret bool) {
	ret = i.mouseReleasedBuf[b]
	i.mouseReleasedBuf[b] = false
	return
}

// 检查后清空
func (i *InputCli) checkMousePressHold(b ebiten.MouseButton) (ret bool) {
	ret = i.mousePressHoldBuf[b]
	i.mousePressHoldBuf[b] = false
	return
}

// 鼠标滑轮滑动
func (i *InputCli) setMouseWheel() {
	_, wheelMoveY := ebiten.Wheel()
	i.wheelMoveBuf[MouseWheelUp] = wheelMoveY > 0
	i.wheelMoveBuf[MouseWheelDown] = wheelMoveY < 0

}

// 检查后清空
func (i *InputCli) checkMouseWheel(moveType keyPressedType) (ret bool) {
	ret = i.wheelMoveBuf[moveType]
	i.wheelMoveBuf[moveType] = false
	return
}

func (i *InputCli) clearActionBuf() {
	i.hasAnyKeyPressed = false
	i.hasAnyKeyReleased = false
	// 键盘
	clear(i.keyPressedBuf)
	clear(i.keyReleasedBuf)
	clear(i.keyPressHoldBuf)
	// 鼠标
	clear(i.mousePressedBuf)
	clear(i.mouseReleasedBuf)
	clear(i.mousePressHoldBuf)
	// 鼠标滑轮
	clear(i.wheelMoveBuf)

}

func (i *InputCli) input() {
	if i.isPause {
		return
	}

	i.clearActionBuf()
	// 键盘
	i.keyBuf = inpututil.AppendJustPressedKeys(i.keyBuf[:0])
	i.setKeyPressed()

	i.keyBuf = inpututil.AppendJustReleasedKeys(i.keyBuf[:0])
	i.setKeyReleased()

	i.keyBuf = inpututil.AppendPressedKeys(i.keyBuf[:0])
	i.setKeyPressHold()

	// 鼠标
	i.setMouse()
	// 鼠标滑轮
	i.setMouseWheel()

}

func (i *InputCli) checkKey(key *actionKeyInfo) bool {
	setTimeF := utils.Empty
	if key.timerLimit != nil {
		// 先检查时间限制
		if key.timerLimit.timer.Since(key.timerLimit.trueTimeSec) > key.timerLimit.timeLimitSec {
			// 超过了时间可以触发
			setTimeF = func() {
				key.timerLimit.trueTimeSec = key.timerLimit.timer.NowSec()
			}
		} else {
			// 没超过，直接退出
			return false
		}
	}

	switch keytype := key.keyType; keytype {
	case MouseWheelUp, MouseWheelDown:
		ret := i.checkMouseWheel(keytype)
		if ret {
			setTimeF()
		}
		return ret
	case PressedType:
		ret := false
		for index := 0; index < len(key.keys); index++ {
			ret = ret || key.keys[index].pressedCheck(i)
		}
		if ret {
			setTimeF()
		}
		return ret
	case ReleasedType:
		ret := false
		for index := 0; index < len(key.keys); index++ {
			ret = ret || key.keys[index].releasedCheck(i)
		}
		if ret {
			setTimeF()
		}
		return ret
	case PressHold:
		ret := false
		for index := 0; index < len(key.keys); index++ {
			ret = ret || key.keys[index].pressHoldCheck(i)
		}
		if ret {
			setTimeF()
		}
		return ret
	case AnyPressedType:
		// 任意键按下检查不支持时间限制
		return i.hasAnyKeyPressed
	case AnyReleasedType:
		// 任意键松开检查不支持时间限制
		return i.hasAnyKeyReleased
	}
	return false
}

func (i *InputCli) checkAllMouseReleasedAction() bool {
	return len(i.mouseReleasedBuf) != 0
}

func (i *InputCli) checkAllMousePerssedAction() bool {
	return len(i.mousePressedBuf) != 0
}

func (i *InputCli) pauseInput()    { i.clearActionBuf(); i.isPause = true }
func (i *InputCli) continueInput() { i.isPause = false }

var (
	inputCli  *InputCli
	inputOnce sync.Once
)

func nowInputCli() *InputCli {
	inputOnce.Do(func() {
		inputCli = &InputCli{
			isPause:           true, // 默认禁止输入
			keyBuf:            make([]ebiten.Key, 0),
			keyPressedBuf:     make(map[ebiten.Key]bool),
			keyReleasedBuf:    make(map[ebiten.Key]bool),
			keyPressHoldBuf:   make(map[ebiten.Key]bool),
			mousePressedBuf:   make(map[ebiten.MouseButton]bool),
			mouseReleasedBuf:  make(map[ebiten.MouseButton]bool),
			mousePressHoldBuf: make(map[ebiten.MouseButton]bool),
			wheelMoveBuf:      make(map[keyPressedType]bool),
		}
	})
	return inputCli
}
