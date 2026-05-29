package utils

import "math"

const (
	LogicalSizeX = 480
	LogicalSizeY = 600

	SmallLogicalSize = 64 // 缩小后的尺寸
	BigTPS           = 24 // 运行帧率
	SmallTPS         = 12 // 缩小帧率

	DoubleClickTime = 500 * SmallTPS / 1000 // 这个帧差以内可看做双击
)

var (
	LogicalSize = Point{X: LogicalSizeX, Y: LogicalSizeY}
)

const (
	Angle5   float64 = 5 * math.Pi / 180   // 弧度制5度
	Angle10  float64 = 10 * math.Pi / 180  // 弧度制10度
	Angle15  float64 = 15 * math.Pi / 180  // 弧度制15度
	Angle20  float64 = 20 * math.Pi / 180  // 弧度制20度
	Angle25  float64 = 25 * math.Pi / 180  // 弧度制25度
	Angle30  float64 = 30 * math.Pi / 180  // 弧度制30度
	Angle35  float64 = 35 * math.Pi / 180  // 弧度制35度
	Angle40  float64 = 40 * math.Pi / 180  // 弧度制40度
	Angle45  float64 = 45 * math.Pi / 180  // 弧度制45度
	Angle50  float64 = 50 * math.Pi / 180  // 弧度制50度
	Angle55  float64 = 55 * math.Pi / 180  // 弧度制55度
	Angle60  float64 = 60 * math.Pi / 180  // 弧度制60度
	Angle65  float64 = 65 * math.Pi / 180  // 弧度制65度
	Angle70  float64 = 70 * math.Pi / 180  // 弧度制70度
	Angle75  float64 = 75 * math.Pi / 180  // 弧度制75度
	Angle80  float64 = 80 * math.Pi / 180  // 弧度制80度
	Angle85  float64 = 85 * math.Pi / 180  // 弧度制85度
	Angle90  float64 = 90 * math.Pi / 180  // 弧度制90度
	Angle95  float64 = 95 * math.Pi / 180  // 弧度制95度
	Angle100 float64 = 100 * math.Pi / 180 // 弧度制100度
	Angle105 float64 = 105 * math.Pi / 180 // 弧度制105度
	Angle110 float64 = 110 * math.Pi / 180 // 弧度制110度
	Angle115 float64 = 115 * math.Pi / 180 // 弧度制115度
	Angle120 float64 = 120 * math.Pi / 180 // 弧度制120度
	Angle125 float64 = 125 * math.Pi / 180 // 弧度制125度
	Angle130 float64 = 130 * math.Pi / 180 // 弧度制130度
	Angle135 float64 = 135 * math.Pi / 180 // 弧度制135度
	Angle140 float64 = 140 * math.Pi / 180 // 弧度制140度
	Angle145 float64 = 145 * math.Pi / 180 // 弧度制145度
	Angle150 float64 = 150 * math.Pi / 180 // 弧度制150度
	Angle155 float64 = 155 * math.Pi / 180 // 弧度制155度
	Angle160 float64 = 160 * math.Pi / 180 // 弧度制160度
	Angle165 float64 = 165 * math.Pi / 180 // 弧度制165度
	Angle170 float64 = 170 * math.Pi / 180 // 弧度制170度
	Angle175 float64 = 175 * math.Pi / 180 // 弧度制175度
	Angle180 float64 = 180 * math.Pi / 180 // 弧度制180度
	Angle185 float64 = 185 * math.Pi / 180 // 弧度制185度
	Angle190 float64 = 190 * math.Pi / 180 // 弧度制190度
	Angle195 float64 = 195 * math.Pi / 180 // 弧度制195度
	Angle200 float64 = 200 * math.Pi / 180 // 弧度制200度
	Angle205 float64 = 205 * math.Pi / 180 // 弧度制205度
	Angle210 float64 = 210 * math.Pi / 180 // 弧度制210度
	Angle215 float64 = 215 * math.Pi / 180 // 弧度制215度
	Angle220 float64 = 220 * math.Pi / 180 // 弧度制220度
	Angle225 float64 = 225 * math.Pi / 180 // 弧度制225度
	Angle230 float64 = 230 * math.Pi / 180 // 弧度制230度
	Angle235 float64 = 235 * math.Pi / 180 // 弧度制235度
	Angle240 float64 = 240 * math.Pi / 180 // 弧度制240度
	Angle245 float64 = 245 * math.Pi / 180 // 弧度制245度
	Angle250 float64 = 250 * math.Pi / 180 // 弧度制250度
	Angle255 float64 = 255 * math.Pi / 180 // 弧度制255度
	Angle260 float64 = 260 * math.Pi / 180 // 弧度制260度
	Angle265 float64 = 265 * math.Pi / 180 // 弧度制265度
	Angle270 float64 = 270 * math.Pi / 180 // 弧度制270度
	Angle275 float64 = 275 * math.Pi / 180 // 弧度制275度
	Angle280 float64 = 280 * math.Pi / 180 // 弧度制280度
	Angle285 float64 = 285 * math.Pi / 180 // 弧度制285度
	Angle290 float64 = 290 * math.Pi / 180 // 弧度制290度
	Angle295 float64 = 295 * math.Pi / 180 // 弧度制295度
	Angle300 float64 = 300 * math.Pi / 180 // 弧度制300度
	Angle305 float64 = 305 * math.Pi / 180 // 弧度制305度
	Angle310 float64 = 310 * math.Pi / 180 // 弧度制310度
	Angle315 float64 = 315 * math.Pi / 180 // 弧度制315度
	Angle320 float64 = 320 * math.Pi / 180 // 弧度制320度
	Angle325 float64 = 325 * math.Pi / 180 // 弧度制325度
	Angle330 float64 = 330 * math.Pi / 180 // 弧度制330度
	Angle335 float64 = 335 * math.Pi / 180 // 弧度制335度
	Angle340 float64 = 340 * math.Pi / 180 // 弧度制340度
	Angle345 float64 = 345 * math.Pi / 180 // 弧度制345度
	Angle350 float64 = 350 * math.Pi / 180 // 弧度制350度
	Angle355 float64 = 355 * math.Pi / 180 // 弧度制355度
	Angle360 float64 = 360 * math.Pi / 180 // 弧度制360度
)
