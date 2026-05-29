package utils

import (
	"math"
	"time"
)

// 所有数字类型
type Number interface {
	~float32 | ~float64 | ~int | ~int32 | ~int64
}

// 浮点数类型
type Float interface {
	~float32 | ~float64
}

// 插值
func Lerp[T Number](minT, maxT, t T) T {
	return (1-t)*minT + t*maxT
}

func TimeLerp(minT, maxT time.Duration, t float64) time.Duration {
	return time.Duration((1-t)*float64(minT) + t*float64(maxT))
}

// 先快后慢
func QuickSlowLerp[T Float](now, end T) T {
	return 0.8*now + 0.2*end
}

// LogBase 计算以 base 为底，x 的对数 (log_base(x))
// 参数:
//
//	base: 对数的底数，必须大于 0 且不等于 1
//	x:   真数，必须大于 0
//
// 返回值:
//
//	对数值；如果输入不合法，则返回 NaN (Not a Number)
func LogBase(base, x float64) float64 {
	// 检查输入是否合法：底数必须 >0 且 ≠1，真数必须 >0
	if base <= 0 || base == 1 || x <= 0 {
		return math.NaN() // 返回"非数字"表示非法输入
	}
	// 使用换底公式：log_base(x) = ln(x) / ln(base)
	return math.Log(x) / math.Log(base)
}

// 限制范围
func Limit[T Number](f, minf, maxf T) T {
	return max(minf, min(f, maxf))
}

// 弧度制转角度制
func RadianToDegrees(r float64) float64 { return 180 * r / math.Pi }

// 角度制转弧度制
func DegreesToRadian(a float64) float64 { return a * math.Pi / 180 }

// 平方
func Square[T Number](x T) T { return x * x }

// 无符号数
type Uint interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64
}

// a-b，确保a减去以后不会从0溢出到最大值
func UintSub[T Uint](a, b T) T {
	return a - min(a, b)
}

// a+b，确保a不会正溢出到0
func Uint8Add(a, b uint8) uint8 {
	return a + min(math.MaxUint8-a, b)
}
