package statemachine

import (
	gametime "gmsender/pkg/game_time"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// StateMachineI 状态机接口
type StateMachineI interface {
	NewState(state int) StateI

	// ReGo 重新进入当前状态，并且执行新状态的进入函数，但不会先执行当前状态的退出，当前状态的退出应该在上次生命周期结束时执行
	ReGo()
	// Go 进入某个状态，并退出当前状态,状态机第一次进入某个状态，只会执行这个状态的enter不会执行当前状态的exit
	Go(state int)

	// GoIfNot 如果当前状态不等于state则进入
	GoIfNot(state int)

	// ReadStateLastTimeSec 读取状态持续当前时间，单位：秒
	ReadStateLastTimeSec() float64

	// SToS 从from到to的转换条件，多次调用只保留最后一次与SToSWithTimeLimit会相互覆盖
	SToS(frome, to int, transCondition func() bool)

	// SToSWithTimeLimit 从from到to的转换条件，多次调用只保留最后一次与SToS会相互覆盖
	SToSWithTimeLimit(frome, to int, limit time.Duration)

	Update()
}

// StateMachineWithDrawI 状态机接口
type StateMachineWithDrawI interface {
	NewState(state int) StateWithDrawI

	// ReGo 重新进入当前状态，并且执行新状态的进入函数，但不会先执行当前状态的退出，当前状态的退出应该在上次生命周期结束时执行
	ReGo()
	// Go 进入某个状态，并退出当前状态,状态机第一次进入某个状态，只会执行这个状态的enter不会执行当前状态的exit
	Go(state int)

	// GoIfNot 如果当前状态不等于state则进入
	GoIfNot(state int)

	// ReadStateLastTimeSec 读取状态持续当前时间，单位：秒
	ReadStateLastTimeSec() float64

	// SToS 从from到to的转换条件
	SToS(frome, to int, transCondition func() bool)

	// SToSWithTimeLimit 从from到to的转换条件
	SToSWithTimeLimit(frome, to int, limit time.Duration)

	Update()
	Draw(screen *ebiten.Image)
}

// StateWithDrawI 状态接口，包含渲染函数
type StateWithDrawI interface {
	// BindUD 绑定更新与绘制函数
	BindUD(uFunc func(), dFunc func(*ebiten.Image))

	// SetEnterFunc 进入时执行
	SetEnterFunc(func()) StateWithDrawI

	// SetExitFunc 退出时执行
	SetExitFunc(func()) StateWithDrawI

	// Enter 进入时执行
	Enter()
	// Exit 退出时执行
	Exit()

	Update()
	Draw(scteen *ebiten.Image)
}

// StateI 状态接口,不包含渲染函数
type StateI interface {
	// BindUD 绑定更新与绘制函数
	BindU(uFunc func())

	// SetEnterFunc 进入时执行
	SetEnterFunc(func()) StateI

	// SetExitFunc 退出时执行
	SetExitFunc(func()) StateI

	// Enter 进入时执行
	Enter()
	// Exit 退出时执行
	Exit()

	Update()
}

// NewStateMachine 新建状态机
func NewStateMachine(timer gametime.TimerType) StateMachineI {
	return newStateMachine(timer)
}

// NewState 新建状态
func NewState() StateI {
	return newState()
}

// NewStateMachineWithDraw 新建状态机，timer为计时器获取，单位：秒
// Deprecated:draw函数似乎有些多余
func NewStateMachineWithDraw(timer gametime.TimerType) StateMachineWithDrawI {
	return newStateMachineWithDraw(timer)
}

// NewStateWithDraw 新建状态
// Deprecated:draw函数似乎有些多余
func NewStateWithDraw() StateWithDrawI {
	return newStateWithDraw()
}
