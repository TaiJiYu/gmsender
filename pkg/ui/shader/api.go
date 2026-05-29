package shader

import (
	_ "embed"
	"gmsender/utils"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	//go:embed shaderfile/rounded_rect.kage
	rounded_rect_kage        []byte
	rounded_rect_kage_shader *ebiten.Shader // 两侧圆角矩形

	//go:embed shaderfile/rounded_rect_lerp.kage
	rounded_rect_lerp_kage        []byte
	rounded_rect_lerp_kage_shader *ebiten.Shader // 两侧圆角矩形,支持插值染色

	//go:embed shaderfile/core_rect.kage
	core_rect_kage        []byte
	core_rect_kage_shader *ebiten.Shader // 四角圆角矩形
)

func InitShader() {
	var err error
	rounded_rect_kage_shader, err = ebiten.NewShader(rounded_rect_kage)
	if err != nil {
		panic(err)
	}
	rounded_rect_lerp_kage_shader, err = ebiten.NewShader(rounded_rect_lerp_kage)
	if err != nil {
		panic(err)
	}
	core_rect_kage_shader, err = ebiten.NewShader(core_rect_kage)
	if err != nil {
		panic(err)
	}
}

const (
	TimeUniformsKey = "Time"
)

var (

// 所有shader所需的key

)

// 圆角矩形,两侧是圆的
type RoundRect struct {
	width   int
	height  int
	options *ebiten.DrawRectShaderOptions
}

// 新建圆角矩形,两侧是圆的
func NewRoundRect(pos utils.Point, ali utils.AlignmentType, fillColor color.Color) *RoundRect {
	r := &RoundRect{
		options: &ebiten.DrawRectShaderOptions{
			Uniforms: map[string]any{
				"FillColor":  utils.ColorToList32(fillColor),
				"BoundColor": utils.ColorToList32(color.Black),
			},
			Images: [4]*ebiten.Image{},
		},
	}
	r.SetPos(pos, ali)
	return r
}
func (r *RoundRect) SetLerp(float64) {}

// 设置填充色
func (r *RoundRect) SetFillColor(fillColor color.Color) {
	r.options.Uniforms["FillColor"] = utils.ColorToList32(fillColor)
}

// 设置边缘颜色
func (r *RoundRect) SetBoundColor(boundColor color.Color) {
	r.options.Uniforms["BoundColor"] = utils.ColorToList32(boundColor)
}

// 需要的时候调用,返回真实渲染的左上角坐标
func (r *RoundRect) SetPos(pos utils.Point, ali utils.AlignmentType) utils.Point {
	pos = ali.GetAlignmentPos(pos, utils.NewPoint(r.width, r.height).Divf1(2))
	r.options.GeoM.Reset()
	r.options.GeoM.Translate(pos.Break())
	return pos
}

// 需要的时候调用
func (r *RoundRect) SetSize(size utils.Point) {
	r.width, r.height = size.BreakInt()
}

func (r *RoundRect) Draw(screen *ebiten.Image) {
	screen.DrawRectShader(r.width, r.height, rounded_rect_kage_shader, r.options)
}

// 圆角矩形,两侧是圆的
type RoundLerpRect struct {
	width   int
	height  int
	options *ebiten.DrawRectShaderOptions
}

// 新建圆角矩形,两侧是圆的,支持插值渲染
func NewRoundLerpRect(pos utils.Point, ali utils.AlignmentType, fillColor, backColor color.Color) *RoundLerpRect {
	r := &RoundLerpRect{
		options: &ebiten.DrawRectShaderOptions{
			Uniforms: map[string]any{
				"FillColor":  utils.ColorToList32(fillColor),
				"BackColor":  utils.ColorToList32(backColor),
				"BoundColor": utils.ColorToList32(color.Black),
				"T":          1.0,
			},
			Images: [4]*ebiten.Image{},
		},
	}
	r.SetPos(pos, ali)
	return r
}

// 设置填充色
func (r *RoundLerpRect) SetFillColor(fillColor color.Color) {
	r.options.Uniforms["FillColor"] = utils.ColorToList32(fillColor)
}

// 设置填充色插值
func (r *RoundLerpRect) SetLerp(t float64) {
	r.options.Uniforms["T"] = t
}

// 设置边缘颜色
func (r *RoundLerpRect) SetBoundColor(boundColor color.Color) {
	r.options.Uniforms["BoundColor"] = utils.ColorToList32(boundColor)
}

// 需要的时候调用,返回真实渲染的左上角坐标
func (r *RoundLerpRect) SetPos(pos utils.Point, ali utils.AlignmentType) utils.Point {
	pos = ali.GetAlignmentPos(pos, utils.NewPoint(r.width, r.height).Divf1(2))
	r.options.GeoM.Reset()
	r.options.GeoM.Translate(pos.Break())
	return pos
}

// 需要的时候调用
func (r *RoundLerpRect) SetSize(size utils.Point) {
	r.width, r.height = size.BreakInt()
}

func (r *RoundLerpRect) Draw(screen *ebiten.Image) {
	screen.DrawRectShader(r.width, r.height, rounded_rect_lerp_kage_shader, r.options)
}

// 圆角矩形
type CoreRect struct {
	width   int
	height  int
	options *ebiten.DrawRectShaderOptions
}

// 新建圆角矩形,四角是圆的
func NewCoreRect(pos utils.Point, ali utils.AlignmentType, fillColor color.Color) *CoreRect {
	r := &CoreRect{
		options: &ebiten.DrawRectShaderOptions{
			Uniforms: map[string]any{
				"FillColor":  utils.ColorToList32(fillColor),
				"BoundColor": utils.ColorToList32(color.Black),
			},
			Images: [4]*ebiten.Image{},
		},
	}
	r.SetPos(pos, ali)
	return r
}
func (r *CoreRect) SetLerp(float64) {}

// 设置填充色
func (r *CoreRect) SetFillColor(fillColor color.Color) {
	r.options.Uniforms["FillColor"] = utils.ColorToList32(fillColor)
}

// 设置边缘颜色
func (r *CoreRect) SetBoundColor(boundColor color.Color) {
	r.options.Uniforms["BoundColor"] = utils.ColorToList32(boundColor)
}

// 需要的时候调用,返回真实渲染的左上角坐标
func (r *CoreRect) SetPos(pos utils.Point, ali utils.AlignmentType) utils.Point {
	pos = ali.GetAlignmentPos(pos, utils.NewPoint(r.width, r.height).Divf1(2))
	r.options.GeoM.Reset()
	r.options.GeoM.Translate(pos.Break())
	return pos
}

// 需要的时候调用
func (r *CoreRect) SetSize(size utils.Point) {
	r.width, r.height = size.BreakInt()
}

func (r *CoreRect) Draw(screen *ebiten.Image) {
	screen.DrawRectShader(r.width, r.height, core_rect_kage_shader, r.options)
}

// 空矩形
type EmptyRect struct {
	width  int
	height int
}

// 新建空矩形，但包含尺寸和位置计算
func NewEmptyRect(pos utils.Point, ali utils.AlignmentType) *EmptyRect {
	r := &EmptyRect{}
	r.SetPos(pos, ali)
	return r
}

// 需要的时候调用,返回真实渲染的左上角坐标
func (r *EmptyRect) SetPos(pos utils.Point, ali utils.AlignmentType) utils.Point {
	pos = ali.GetAlignmentPos(pos, utils.NewPoint(r.width, r.height).Divf1(2))
	return pos
}

// 需要的时候调用
func (r *EmptyRect) SetSize(size utils.Point) {
	r.width, r.height = size.BreakInt()
}
func (r *EmptyRect) SetLerp(float64) {}

// 设置填充色
func (r *EmptyRect) SetFillColor(color.Color) {
}

// 设置边缘颜色
func (r *EmptyRect) SetBoundColor(color.Color) {
}

func (r *EmptyRect) Draw(screen *ebiten.Image) {}

// 通用shader
//kage:unit pixels

/*
func UV(tex vec2) vec2 {
	return (tex - imageDstOrigin()) / Size
}

func Lerp(l, r vec4, t float) vec4 {
	return (1.0-t)*l + t*r
}

func LerpT(l, r, t float) float {
	return (1-t)*l + t*r
}

// 最里面是0，最外层是1
func CircleS(uv vec2, radius float) float {
	uv = uv*2 - 1
	uv.y *= -1
	dist := length(uv)
	return step(dist, radius) * dist / radius
}


func CircleLine(uv vec2, radius float, borderWidth float) float {
	uv = uv*2 - 1
	uv.y *= -1
	dist := length(uv)
	return step(dist, radius) * step(radius-borderWidth, dist)
}

// 旋转的0-1
func CircleRound(uv vec2) float {
	uv -= 0.5
	down := step(uv.y, 0)
	up := 1 - down
	zero_v := vec2(1.0, 0.0)
	uv = normalize(uv)
	cos_a := dot(uv, zero_v)
	c := acos(cos_a) / 3.14159265
	return (c*up + down*(2.0-c)) / 2.0
}


func Circle(uv vec2,center vec2, radius float) float{
	dist := distance(uv, center)
	return step(dist,radius)
}


func norColor(r, g, b float) vec3 {
	return vec3(r/255.0, g/255.0, b/255.0)
}


func MoveCircleLine(uv vec2, radius float, speed float, borderWidth float, color vec3, times float, offsetT float, timeNow float) vec4 {
	a := CircleRound(uv)
	max_height := radius // 最大振幅
	height := LerpT(-max_height, max_height, abs(fract(timeNow*speed)*2-1))
	c := CircleLine(uv, height*sin((a*times*2+offsetT)*3.14159265)+1.0-max_height, borderWidth)
	return vec4(color*c, c)
}

// x==v返回1，否则为0
func equl(x,v float)float{return 1.0-sign(abs(x-v))}


func noise(uv vec2) float {
	return fract(sin(dot(uv,vec2(12.9898,78.233)))*43758.5453)
}



// 返回0-1的随机数
func RandomF(seed vec2)float{
    return fract(sin(dot(seed, vec2(12.9898, 4.1414)))*43758.5453)
}



// 返回0-1的随机数
func RandomFF(s1,s2 float)float{
    return RandomF(vec2(s1,s2))
}


// 绘制一条线
func DrawLine(uv,p0, p1 vec2,width,border float)float{
	i := uv - p0
	k := p1 - p0
	return 1-clamp((distance(i-clamp(dot(i, k)/dot(k,k),0,1)*k,vec2(0,0))-width)/(width*border),0,1)
}


// 绕(0,0)中心旋转
func Rotate(uv vec2, angle float)vec2{
    c := cos(angle)
    s := sin(angle)
    return vec2(uv.x*c-uv.y*s,uv.x*s+uv.y*c)
}

// 绕中心旋转图片的srcPos坐标系
func RotateByMid(drawPos vec2,srcPos vec2,angle float)vec2{
    srcOrigin := imageSrc0Origin()
    halfSize := imageDstSize()/2-drawPos

    uv := srcPos-halfSize-srcOrigin
    dist := length(uv)
    angle = atan2(uv.y,uv.x)+angle

    newUV := vec2(cos(angle),sin(angle))*dist
    newUV += srcOrigin+halfSize
    return newUV
}


*/
