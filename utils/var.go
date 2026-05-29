package utils

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// 拷贝绘制
	CopyDrawImageOp = &ebiten.DrawImageOptions{
		Blend: ebiten.BlendCopy,
	}

	// 直接绘制
	BaseDrawImageOp = &ebiten.DrawImageOptions{}
)

var (
	weeks = []string{"日", "一", "二", "三", "四", "五", "六"}
)

// 当前周几的文本
func WeekStr() string {
	return "周" + weeks[time.Now().Weekday()]
}
