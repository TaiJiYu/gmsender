package gametime

// 初始化时间系统
func InitGameTimer() { initGameTimer() }

// 全局时间计时
func BigTimeRun() {
	gametimerCli.bigTimeRun()
}

// 游玩时间计时
func SmallTimeRun() {
	gametimerCli.smallTimeRun()
}

// 游玩时间，暂停不计时
func SmallNowSec() float64 { return gametimerCli.smallNowSec() }

// 游玩时间到t的时间差，暂停不计时
func SmallNowSecSince(tSec float64) float64 { return gametimerCli.smallNowSecSince(tSec) }

// 全局时间，暂停也计时
func BigNowSec() float64 { return gametimerCli.bigNowSec() }

// 全局时间到t的时间差，暂停也计时
func BigNowSecSince(tSec float64) float64 { return gametimerCli.bigNowSecSince(tSec) }

// 根据计时器类型获取时间
func (t TimerType) NowSec() float64 { return gametimerCli.nowSec(t) }

// 根据计时器类型获取时间差,t的单位为秒
func (t TimerType) Since(tSec float64) float64 {
	return gametimerCli.since(tSec, t)
}
