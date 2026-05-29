package utils

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	maxUint32 float32 = 0xffff
)

func ColorToList32(col color.Color) []float32 {
	r, g, b, a := col.RGBA()
	return []float32{float32(r) / maxUint32, float32(g) / maxUint32, float32(b) / maxUint32, float32(a) / maxUint32}
}

// 将彩色图片转换为灰度图,保留透明通道,会返回新的图
func ConvertToGrayscale(img *ebiten.Image) *ebiten.Image {
	bounds := img.Bounds()
	ret := ebiten.NewImage(bounds.Max.X, bounds.Max.Y)
	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			// 转换为灰度
			orginColor := img.At(x, y)
			_, _, _, a := orginColor.RGBA()
			r, g, b, _ := color.GrayModel.Convert(orginColor).RGBA()
			ret.Set(x, y, color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)})
		}
	}
	return ret
}

// 仅rgb，a为255
func ColorRGBByOx(c uint32) color.Color {
	return color.RGBA{R: uint8(c >> 16), G: uint8(c >> 8 & 0xff), B: uint8(c & 0xff), A: 255}
}

// 颜色插值
func ColorRGBLerp(col0, col1 color.Color, t float64) color.Color {
	r0, g0, b0, _ := col0.RGBA()
	r1, g1, b1, _ := col1.RGBA()
	return color.RGBA{
		R: uint8(float64(r0>>8)*(1-t) + float64(r1>>8)*t),
		G: uint8(float64(g0>>8)*(1-t) + float64(g1>>8)*t),
		B: uint8(float64(b0>>8)*(1-t) + float64(b1>>8)*t),
		A: 255,
	}
}
