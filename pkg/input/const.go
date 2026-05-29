package input

const (
	axisLowest = 0.2            // 手柄死区，不得低于这个值，低于直接当做0
	axisRange  = 1 - axisLowest // 手柄真实活动范围
)
