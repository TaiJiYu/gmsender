package utils

import (
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

var ZeroPoint = Point{}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// 泛型变量
func NewPoint[T, W ~int | ~float64 | ~float32](x T, y W) Point {
	return Point{X: float64(x), Y: float64(y)}
}

// 泛型变量
func NewPointT1[T ~int | ~float64 | ~float32](f T) Point {
	return Point{X: float64(f), Y: float64(f)}
}

func (p Point) Max(p1 Point) Point {
	p.X = max(p.X, p1.X)
	p.Y = max(p.Y, p1.Y)
	return p
}

func (p Point) MaxY(y float64) Point {
	p.Y = max(p.Y, y)
	return p
}

func (p Point) Min(p1 Point) Point {
	p.X = min(p.X, p1.X)
	p.Y = min(p.Y, p1.Y)
	return p
}

func (p Point) BreakToList(l [2]float64) [2]float64 {
	l[0] = p.X
	l[1] = p.Y
	return l
}

// 分开多个维度
func (p Point) Break() (float64, float64) {
	return p.X, p.Y
}

// 分开多个维度
func (p Point) BreakInt() (int, int) {
	return int(p.X), int(p.Y)
}

// 分开多个维度为float32
func (p Point) Break32() (float32, float32) {
	return float32(p.X), float32(p.Y)
}

// 是否相等
func (p Point) Eqal(p0 Point) bool {
	return p.X == p0.X && p.Y == p0.Y
}

func (p Point) IsNan() bool {
	return math.IsNaN(p.X) || math.IsNaN(p.Y)
}

func (p Point) IsZero() bool {
	return p.X == 0 && p.Y == 0
}

// x和y都要加f
func (p Point) AddFToXY(f float64) Point {
	p.X += f
	p.Y += f
	return p
}

func (p Point) AddXY(x, y float64) Point {
	p.X += x
	p.Y += y
	return p
}
func (p Point) AddX(x float64) Point {
	p.X += x
	return p
}

func (p Point) AddY(y float64) Point {
	p.Y += y
	return p
}

func (p Point) Add(p0 Point) Point {
	p.X += p0.X
	p.Y += p0.Y
	return p
}

func (p Point) Sub(p0 Point) Point {
	p.X -= p0.X
	p.Y -= p0.Y
	return p
}

func (p Point) SubY(y float64) Point {
	p.Y -= y
	return p
}

func (p Point) SubXY(f float64) Point {
	p.X -= f
	p.Y -= f
	return p
}

func (p Point) SubX(x float64) Point {
	p.X -= x
	return p
}

func (p Point) Divf1(f float64) Point {
	p.X /= f
	p.Y /= f
	return p
}

func (p Point) Mul(p0 Point) Point {
	p.X *= p0.X
	p.Y *= p0.Y
	return p
}
func (p Point) MulF(x, y float64) Point {
	p.X *= x
	p.Y *= y
	return p
}

func (p Point) MulF1(f float64) Point {
	p.X *= f
	p.Y *= f
	return p
}

// x+y
func (p Point) Sum() float64 {
	return p.X + p.Y
}

// 返回p的长度
func (p Point) Len() float64 {
	return math.Hypot(p.X, p.Y)
}

// 检查p的长度是否大于x
func (p Point) IsLenBig(x float64) bool {
	if x <= 0 {
		return false
	}
	return p.X*p.X+p.Y*p.Y > x*x
}

// 检查p的长度是否小于x
func (p Point) IsLenLess(x float64) bool {
	if x <= 0 {
		return false
	}
	return p.X*p.X+p.Y*p.Y < x*x
}

// AB线段是否长度超过x
func (A Point) IsLineLenBig(B Point, x float64) bool {
	return A.Sub(B).IsLenBig(x)
}

// 归一化,p为0则不变
func (p Point) Normal() Point {
	if p.IsZero() {
		return p
	}
	return p.Divf1(p.Len())
}

// 在p到obj间插值
func (p Point) Lerp(obj Point, t float64) Point {
	p.X = Lerp(p.X, obj.X, t)
	p.Y = Lerp(p.Y, obj.Y, t)
	return p
}

// 向量转角度弧度制，(0,-1)代表0度，顺时针为正,p为(0,0时)返回0
func (p Point) ToAngle() float64 {
	if p.IsZero() {
		return 0
	}
	return math.Atan2(p.Y, p.X) + Angle90
}

// 弧度制的角度转方向归一化向量，该函数是ToAngle的逆函数
func AngleToDir(a float64) Point {
	a -= Angle90
	return Point{
		X: math.Cos(a),
		Y: math.Sin(a),
	}
}

// p0在p上的投影
func (p Point) Shadow(p0 Point) Point {
	return p.MulF1(p.Dot(p0))
}

// 点积
func (p Point) Dot(p0 Point) float64 {
	return p.Mul(p0).Sum()
}

// 叉积,返回标量
func (p Point) CrossF(p0 Point) float64 {
	return p.X*p0.Y - p.Y*p0.X
}

// 距离约束，约束p0到p的距离为dis
func (p Point) DistanceConstraint(p0 Point, dis float64) Point {
	p0 = p0.Sub(p)
	return p0.MulF1(dis / p0.Len()).Add(p)
}

// p绕p0顺时针旋转弧度制的角度
func (p Point) RotateAngleWithP0(p0 Point, a float64) Point {
	sinA := math.Sin(a)
	cosA := math.Cos(a)
	dis := p.Sub(p0)
	p.X = dis.MulF(cosA, -sinA).Sum()
	p.Y = dis.MulF(sinA, cosA).Sum()
	return p.Add(p0)
}

// p绕0,0旋转弧度制的角度，顺时针为正方向
func (p Point) RotateAngle(a float64) Point {
	sinA := math.Sin(a)
	cosA := math.Cos(a)
	v := p
	p.X = v.MulF(cosA, -sinA).Sum()
	p.Y = v.MulF(sinA, cosA).Sum()
	return p
}

// 计算p的垂直向量
func (p Point) Vertical() Point {
	p.X, p.Y = -p.Y, p.X
	return p
}

// 获取p与p0的角平分单位向量
func (p Point) Split(p0 Point) Point {
	return p.Normal().Add(p0.Normal()).Normal()
}

// 计算角ABC的归一化角平分向量
func (A Point) SplitAngle(B, C Point) Point {
	return A.Sub(B).Split(C.Sub(B))
}

// 计算v1在v的哪一侧，顺时针方向为-1，否则为1
func (v Point) Dir(v1 Point) float64 {
	return math.Copysign(1, v.AngleV(v1.Vertical())-Angle90)
}

// 计算v和v1的弧度制夹角
func (v Point) AngleV(v1 Point) float64 {
	return math.Acos(Limit(v.Normal().Dot(v1.Normal()), -1, 1))
}

// 计算角ABC的弧度制角度
func (A Point) Angle(B, C Point) float64 {
	return math.Acos(Limit(A.Sub(B).Normal().Dot(C.Sub(B).Normal()), -1, 1))
}

// 计算角ABC的弧度制角度，正方向为正数
func (A Point) AngleS(B, C Point) float64 {
	return A.Angle(B, C) * C.PointDir(B, A)
}

// 计算向量v和v1的弧度制角度，正方向为正数
func (v Point) AngleVS(v1 Point) float64 {
	return v.AngleV(v1) * v1.Dir(v)
}

// 把向量BA与BC的角度拉开到至少弧度制anlgeMin，BC远离BA，BA不变，并确保C在向量BA的cDir测
func (A Point) FarAwayPKeepCWithDir(B, C Point, anlgeMin, cDir float64) Point {
	v := A.Sub(B)
	v1 := C.Sub(B)
	dir := v.Dir(v1)
	if dir != cDir {
		// 如果不等于应该先纠正到anlgeMin度的位置
		// return v.Normal().RotateAngle(cDir * anlgeMin).MulF1(v1.Len()).Add(B)

		// 如果不等于，先纠正到180度
		return v.Normal().MulF1(v1.Len()).MulF1(-1).Add(B)
	}
	return v.FarAwayV(v1, anlgeMin).Add(B)
}

// 计算C在向量BA的哪一侧，顺时针方向为-1，否则为1
func (A Point) PointDir(B, C Point) float64 {
	return A.Sub(B).Dir(C.Sub(B))
}

// 把向量BA与BC的角度拉开到至少弧度制anlgeMin，BC远离BA，BA不变
func (A Point) FarAwayP(B, C Point, anlgeMin float64) Point {
	return A.Sub(B).FarAwayV(C.Sub(B), anlgeMin).Add(B)
}

// 把向量v1与v的角度拉开到至少弧度制anlgeMin，v1远离v，v不变
func (v Point) FarAwayV(v1 Point, anlgeMin float64) Point {
	return v1.RotateAngle(v.Dir(v1) * max((anlgeMin-v.AngleV(v1)), 0))
}

// p是否在范围内
func (p Point) IsRangeIn(rangeX Point, rangeY Point) bool {
	return p.X >= rangeX.X && p.Y >= rangeY.X && p.X <= rangeX.Y && p.Y <= rangeY.Y
}

// p是否在范围内
func (p Point) IsRangeInFloat(xmin, xmax, ymin, ymax float64) bool {
	return p.X >= xmin && p.Y >= ymin && p.X <= xmax && p.Y <= ymax
}

// 获取p关于p0-p1直线的对称点
func (p Point) Symmetric(p0, p1 Point) Point {
	v := p1.Sub(p0) // 直线方向向量
	w := p.Sub(p0)

	// 投影点
	t := w.Dot(v) / v.Dot(v)
	proj := p0.Add(v.MulF1(t))
	return proj.MulF1(2).Sub(p)
}

// 计算points的所有控制点并添加到数组中，需要在points中预留控制点的位置，points为闭合曲线，个数一定为3的倍数
func BezierContralPoints(points []Point) {
	size := len(points) / 3
	for i := 0; i < size; i++ {
		p0Index := (i - 1 + size) % size * 3
		p0 := points[p0Index]
		p1 := points[i*3]
		p2 := points[(i+1)%size*3]
		v := p0.SplitAngle(p1, p2).Vertical()
		p1_per := p1.Sub(v.Shadow(p1.Sub(p0)).MulF1(0.3))
		p1_after := p1.Add(v.Shadow(p2.Sub(p1)).MulF1(0.3))
		points[p0Index+2] = p1_per
		points[i*3+1] = p1_after
	}
}

// 在范围内生成随机点
func RandPoint(rangeX Point, rangeY Point) Point {
	return Point{
		X: Lerp(rangeX.X, rangeX.Y, rand.Float64()),
		Y: Lerp(rangeY.X, rangeY.Y, rand.Float64()),
	}
}

// 极坐标转直角坐标,p的x代表角度，y代表半径
func (p Point) PolarToCartesian() Point {
	cosa := Limit(math.Cos(p.X), -1, 1)
	sina := Limit(math.Sin(p.X), -1, 1)
	p.X = cosa * p.Y
	p.Y = sina * p.Y
	return p
}

// 获取鼠标位置
func MousePos() Point {
	x, y := ebiten.CursorPosition()
	return Point{float64(x), float64(y)}
}

// 返回x和y的绝对值
func (p Point) Abs() Point {
	p.X = math.Abs(p.X)
	p.Y = math.Abs(p.Y)
	return p
}

// x是否等于y
func (p Point) XEqulY() bool {
	return p.X == p.Y
}

// x或者y为0
func (p Point) XOrYIsZero() bool {
	return p.X == 0 || p.Y == 0
}

// 任意归一化向量射向与单位圆的外接正方形的距离
func (p Point) DistanceToSquare() float64 {
	if p.IsZero() {
		return 0
	}
	if p.XOrYIsZero() {
		return 1
	}
	p = p.Abs()
	if p.XEqulY() {
		return math.Sqrt2
	}
	if p.X < p.Y {
		return 1 / p.Y
	}
	return 1 / p.X
}

// GenerateFibonacciSpiralPoints 使用斐波那契螺旋算法生成更均匀松散的点
// 这种方法生成的点在视觉上更均匀,count为点个数,r为点半径,innerR是内部半径，不会在这里面刷新
func (p Point) GenerateFibonacciSpiralPoints(count int, r float64, innerR float64, randFloatFunc func() float64) []Point {
	if count <= 0 {
		return []Point{}
	}

	points := make([]Point, count)

	// 黄金角（弧度）
	goldenAngle := math.Pi * (3 - math.Sqrt(5))

	// maxRadius := r * math.Sqrt(float64(count))

	for i := 0; i < count; i++ {
		// 使用平方根分布，让点从内向外自然分布
		radius := r*math.Sqrt(float64(i)+0.5) + innerR

		// 角度使用黄金角，确保点不会重叠
		angle := float64(i) * goldenAngle

		// 添加少量随机扰动，让分布更自然
		radiusJitter := radius + (randFloatFunc()-0.5)*r*0.2
		angleJitter := angle + (randFloatFunc()-0.5)*0.2

		x := radiusJitter * math.Cos(angleJitter)
		y := radiusJitter * math.Sin(angleJitter)

		points[i] = Point{X: x, Y: y}.Add(p)
	}

	return points
}

// 约束在屏幕范围内,r为半径
func (p Point) LimitInScreen(r float64) Point {
	p.X = Limit(p.X, r, LogicalSizeX-r)
	p.Y = Limit(p.Y, r, LogicalSizeY-r)
	return p
}
