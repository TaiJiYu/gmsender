package statemachine

import (
	gametime "gmsender/pkg/game_time"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// 用于状态机的状态控制器
type stateController struct {
	s               StateI
	transConditions map[int]func() bool // 转移到某个状态的条件
}

type statemachine struct {
	states   map[int]*stateController
	nowState *stateController

	enterTime        float64            // 进入当前状态的时间
	nowStateTimeLast float64            // 当前状态的持续时间
	timer            gametime.TimerType // 计时器类型，单位：秒
	isFirstGoDone    bool               // 是否第一次完成
}

func newStateMachine(timer gametime.TimerType) *statemachine {
	return &statemachine{
		states: make(map[int]*stateController),
		timer:  timer,
	}
}

func (s *statemachine) NewState(state int) StateI {
	s.states[state] = &stateController{
		s:               NewState(),
		transConditions: make(map[int]func() bool),
	}
	return s.states[state].s
}

// ReGo 重新进入当前状态，并且执行新状态的进入函数，但不会先执行当前状态的退出，当前状态的退出应该在上次生命周期结束时执行
func (s *statemachine) ReGo() {
	s.nowState.s.Enter() // 在enter后再重新计时，enter中可选择使用旧状态的时间，也可用0作为当前时间
	s.enterTime = s.timer.NowSec()
	s.nowStateTimeLast = 0
}

// GoIfNot 如果当前状态不等于state则进入
func (s *statemachine) GoIfNot(state int) {
	if s.nowState == s.states[state] {
		return
	}
	s.Go(state)
}

// Go 直接进入某个状态
func (s *statemachine) Go(state int) {
	if s.isFirstGoDone {
		s.nowState.s.Exit()
	} else {
		s.isFirstGoDone = true
	}
	s.nowState = s.states[state]
	s.nowState.s.Enter() // 在enter后再重新计时，enter中可选择使用旧状态的时间，也可用0作为当前时间
	s.enterTime = s.timer.NowSec()
	s.nowStateTimeLast = 0

}

// SToS 从from到to的转换条件
func (s *statemachine) SToS(from, to int, transCondition func() bool) {
	s.states[from].transConditions[to] = transCondition
}

func (s *statemachine) SToSWithTimeLimit(from, to int, limit time.Duration) {
	sec := limit.Seconds()
	s.states[from].transConditions[to] = func() bool {
		return s.nowStateTimeLast >= sec
	}
}

func (s *statemachine) ReadStateLastTimeSec() float64 {
	return s.nowStateTimeLast
}

func (s *statemachine) Update() {
	// 更新帧状态时间
	s.nowStateTimeLast = s.timer.Since(s.enterTime)
	for nextS, condition := range s.nowState.transConditions {
		if condition() {
			s.Go(nextS)
			break
		}
	}
	s.nowState.s.Update()
}

// 用于状态机的状态控制器
type stateControllerWithDraw struct {
	s               StateWithDrawI
	transConditions map[int]func() bool // 转移到某个状态的条件
}

type statemachineWithDraw struct {
	states   map[int]*stateControllerWithDraw
	nowState *stateControllerWithDraw

	enterTime        float64            // 进入当前状态的时间
	nowStateTimeLast float64            // 当前状态的持续时间
	timer            gametime.TimerType // 计时器类型，单位：秒
	isFirstGoDone    bool               // 是否第一次完成
}

func newStateMachineWithDraw(timer gametime.TimerType) *statemachineWithDraw {
	return &statemachineWithDraw{
		states: make(map[int]*stateControllerWithDraw),
		timer:  timer,
	}
}

func (s *statemachineWithDraw) NewState(state int) StateWithDrawI {
	s.states[state] = &stateControllerWithDraw{
		s:               NewStateWithDraw(),
		transConditions: make(map[int]func() bool),
	}
	return s.states[state].s
}

// ReGo 重新进入当前状态，并且执行新状态的进入函数，但不会先执行当前状态的退出，当前状态的退出应该在上次生命周期结束时执行
func (s *statemachineWithDraw) ReGo() {
	s.nowState.s.Enter() // 在enter后再重新计时，enter中可选择使用旧状态的时间，也可用0作为当前时间
	s.enterTime = s.timer.NowSec()
	s.nowStateTimeLast = 0
}

// GoIfNot 如果当前状态不等于state则进入
func (s *statemachineWithDraw) GoIfNot(state int) {
	if s.nowState == s.states[state] {
		return
	}
	s.Go(state)
}

// Go 直接进入某个状态
func (s *statemachineWithDraw) Go(state int) {
	if s.isFirstGoDone {
		s.nowState.s.Exit()
	} else {
		s.isFirstGoDone = true
	}
	s.nowState = s.states[state]
	s.nowState.s.Enter() // 在enter后再重新计时，enter中可选择使用旧状态的时间，也可用0作为当前时间
	s.enterTime = s.timer.NowSec()
	s.nowStateTimeLast = 0

}

// SToS 从from到to的转换条件
func (s *statemachineWithDraw) SToS(from, to int, transCondition func() bool) {
	s.states[from].transConditions[to] = transCondition
}

func (s *statemachineWithDraw) SToSWithTimeLimit(from, to int, limit time.Duration) {
	sec := limit.Seconds()
	s.states[from].transConditions[to] = func() bool {
		return s.nowStateTimeLast >= sec
	}
}

func (s *statemachineWithDraw) ReadStateLastTimeSec() float64 {
	return s.nowStateTimeLast
}

func (s *statemachineWithDraw) Update() {
	// 更新帧状态时间
	s.nowStateTimeLast = s.timer.Since(s.enterTime)
	for nextS, condition := range s.nowState.transConditions {
		if condition() {
			s.Go(nextS)
			break
		}
	}
	s.nowState.s.Update()
}

func (s *statemachineWithDraw) Draw(screen *ebiten.Image) {
	s.nowState.s.Draw(screen)
}
