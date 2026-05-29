package statemachine

import (
	"gmsender/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

// 状态结构
type stateWithDraw struct {
	enterFunc, exitFunc func()              // 进入和退出函数
	updateFunc          func()              // 更新函数
	drawFunc            func(*ebiten.Image) // 渲染函数
}

func newStateWithDraw() *stateWithDraw {
	return &stateWithDraw{
		enterFunc:  utils.Empty,
		exitFunc:   utils.Empty,
		updateFunc: utils.Empty,
		drawFunc:   utils.EmptyDraw,
	}
}

func (s *stateWithDraw) BindUD(uFunc func(), dFunc func(*ebiten.Image)) {
	s.updateFunc = uFunc
	s.drawFunc = dFunc
}

// SetEnterFunc 进入时执行
func (s *stateWithDraw) SetEnterFunc(e func()) StateWithDrawI {
	s.enterFunc = e
	return s
}

// SetExitFunc 退出时执行
func (s *stateWithDraw) SetExitFunc(e func()) StateWithDrawI {
	s.exitFunc = e
	return s
}

// Enter 进入时执行
func (s *stateWithDraw) Enter() {
	s.enterFunc()
}

// Exit 退出时执行
func (s *stateWithDraw) Exit() {
	s.exitFunc()
}

func (s *stateWithDraw) Update() {
	s.updateFunc()
}
func (s *stateWithDraw) Draw(screen *ebiten.Image) {
	s.drawFunc(screen)
}

// 状态结构
type state struct {
	enterFunc, exitFunc func() // 进入和退出函数
	updateFunc          func() // 更新函数
}

func newState() *state {
	return &state{
		enterFunc:  utils.Empty,
		exitFunc:   utils.Empty,
		updateFunc: utils.Empty,
	}
}

func (s *state) BindU(uFunc func()) {
	s.updateFunc = uFunc
}

// SetEnterFunc 进入时执行
func (s *state) SetEnterFunc(e func()) StateI {
	s.enterFunc = e
	return s
}

// SetExitFunc 退出时执行
func (s *state) SetExitFunc(e func()) StateI {
	s.exitFunc = e
	return s
}

// Enter 进入时执行
func (s *state) Enter() {
	s.enterFunc()
}

// Exit 退出时执行
func (s *state) Exit() {
	s.exitFunc()
}

func (s *state) Update() {
	s.updateFunc()
}
