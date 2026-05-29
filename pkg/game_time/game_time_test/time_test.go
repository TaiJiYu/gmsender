package gametime_test

import (
	"fmt"
	gametime "gmsender/pkg/game_time"
	"testing"
	"time"
)

// 误差测试
func TestAcc(t *testing.T) {
	gametime.InitGameTimer()
	// allTime 测试时间,tick 测试帧率
	f := func(allTime time.Duration, tick int) {
		c := int(allTime.Seconds()) * tick
		for range c {
			gametime.BigTimeRun()
			gametime.SmallTimeRun()
		}
		glNow := gametime.BigNowSec()
		now := gametime.SmallNowSec()
		acc := now - allTime.Seconds()
		fmt.Printf("测试帧率：%v,累加次数：%v,测试时长：%v(=%vs),游玩时长：%v(=%vs),全局时长：%v(=%vs),游玩误差：%vs,游玩误差小于1毫秒:%v\n", tick, c, allTime, allTime.Seconds(), time.Duration(now*1e9), now, time.Duration(glNow*1e9), glNow, acc, acc < time.Millisecond.Seconds())
		// gametime.ResetPlayTime()
	}
	for i := range 24 {
		f(time.Duration(i+1)*time.Hour, 120)
	}
}
