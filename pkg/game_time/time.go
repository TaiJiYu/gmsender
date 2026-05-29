package gametime

import (
	"fmt"
	"gmsender/utils"
)

type gametimer struct {
	bigTick   int // 全局时间tick数
	smallTick int // 游玩时间的tick数
	// worldPlayTick int // 世界游玩时间

	bigTimeSec   float64 // 全局时间
	smallTimeSec float64 // 游玩当前时间
	// worldPlayTimeSec float64 // 世界游玩当前时间
}

// 计时器类型
type TimerType int

const (
	BigTimerType   TimerType = iota // 放大时的时间
	SmallTimerType                  // 缩小时的时间
)

var (
	gametimerCli *gametimer
)

// 初始化时间计时器
func initGameTimer() {
	gametimerCli = &gametimer{}
}

// // 仅重置游玩时间
// func (ti *gametimer) resetPlayTime() {
// 	ti.playTick = 0
// 	ti.playTimeSec = 0
// }

// // 仅重置游玩时间
// func (ti *gametimer) resetWorldPlayTime() {
// 	ti.worldPlayTick = 0
// 	ti.worldPlayTimeSec = 0
// }

// 游玩时间运行，请在update中调用,只有在调用run时，时间才会增加，否则不会改变值
func (ti *gametimer) smallTimeRun() {
	ti.smallTick++
	ti.smallTimeSec = float64(ti.smallTick) / utils.SmallTPS
}

// 全局时间运行，请在update中调用,只有在调用run时，时间才会增加，否则不会改变值,在该函数中会同步tps
func (ti *gametimer) bigTimeRun() {
	ti.bigTick++
	ti.bigTimeSec = float64(ti.bigTick) / utils.BigTPS
}

// 游玩时间，暂停不计时
func (ti *gametimer) smallNowSec() float64 { return ti.smallTimeSec }

// 游玩时间到t的时间差，暂停不计时
func (ti *gametimer) smallNowSecSince(t float64) float64 { return ti.smallNowSec() - t }

// 全局时间，暂停也计时
func (ti *gametimer) bigNowSec() float64 { return ti.bigTimeSec }

// 全局时间到t的时间差，暂停也计时
func (ti *gametimer) bigNowSecSince(t float64) float64 { return ti.bigNowSec() - t }

// 读取当前时间，单位秒
func (ti *gametimer) nowSec(timeType TimerType) float64 {
	switch timeType {
	case BigTimerType:
		return ti.bigNowSec()
	case SmallTimerType:
		return ti.smallNowSec()
	}
	panic(fmt.Errorf("未设置类型计时器：%v", timeType))
}

// 获取到tSec的时间
func (ti *gametimer) since(tSec float64, timeType TimerType) float64 {
	return ti.nowSec(timeType) - tSec
}
