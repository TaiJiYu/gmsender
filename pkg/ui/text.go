package ui

import (
	"bytes"
	_ "embed"
	"gmsender/pkg/asset"
	"gmsender/utils"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type textUiFace struct {
	smallFace, midFace, bigFace *text.MultiFace // 小中大3种字体，游戏中只需要3种
	options                     *text.DrawOptions
}
type FontSize float64

const (
	SmallSize FontSize = 15
	MidSize   FontSize = 32
	BigSize   FontSize = 64
)

var textUiFaces *textUiFace

var fontS *text.GoTextFaceSource

func initTextUi() {
	var err error
	// 其他
	fontS, err = text.NewGoTextFaceSource(bytes.NewReader(asset.FontData()))
	if err != nil {
		panic(err)
	}

	textUiFaces = &textUiFace{
		options: &text.DrawOptions{},
	}
	smallFace, err := text.NewMultiFace(&text.GoTextFace{
		Source: fontS,
		Size:   float64(SmallSize),
	})
	if err != nil {
		panic(err)
	}
	textUiFaces.smallFace = smallFace

	midFace, err := text.NewMultiFace(&text.GoTextFace{
		Source: fontS,
		Size:   float64(MidSize),
	})
	if err != nil {
		panic(err)
	}
	textUiFaces.midFace = midFace
	bigFace, err := text.NewMultiFace(&text.GoTextFace{
		Source: fontS,
		Size:   float64(BigSize),
	})
	if err != nil {
		panic(err)
	}
	textUiFaces.bigFace = bigFace
}

func (s FontSize) face() *text.MultiFace {
	switch s {
	case SmallSize:
		return textUiFaces.smallFace
	case MidSize:
		return textUiFaces.midFace
	case BigSize:
		return textUiFaces.bigFace
	}
	return textUiFaces.smallFace
}

type TextUi struct {
	offset   utils.Point
	isFormat bool // 是不是字符串格式化格式
	size     FontSize
	pos      utils.Point
	ali      utils.AlignmentType
	op       *ebiten.DrawImageOptions
	s        string // 原文本
	img      *ebiten.Image
	halfSize utils.Point
	color    color.Color
}

// func newTextUi(textKey asset.LocationKey, size FontSize, pos utils.Point, ali utils.AlignmentType, color color.Color) *TextUi {
// 	t := &TextUi{
// 		textKey: textKey,
// 		size:    size,
// 		pos:     pos,
// 		ali:     ali,
// 		op:      &ebiten.DrawImageOptions{},
// 		color:   color,
// 	}
// 	t.setText(textKey.Text())
// 	t.setOp(utils.Point{})
// 	return t
// }

// 静态文本，一般是数字之类的
func newStaticTextUi(text string, size FontSize, pos utils.Point, ali utils.AlignmentType, color color.Color) *TextUi {
	t := &TextUi{
		size:  size,
		pos:   pos,
		ali:   ali,
		op:    &ebiten.DrawImageOptions{},
		color: color,
	}
	t.setText(text)
	t.setOp(utils.Point{})
	return t
}

// func newTextUiByFormat(formats string, textKeys []asset.LocationKey, size FontSize, pos utils.Point, ali utils.AlignmentType, color color.Color) *TextUi {
// 	t := &TextUi{
// 		isFormat:   true,
// 		formats:    formats,
// 		formatKeys: textKeys,
// 		size:       size,
// 		pos:        pos,
// 		ali:        ali,
// 		op:         &ebiten.DrawImageOptions{},
// 		color:      color,
// 	}
// 	sl := []any{}
// 	for _, k := range textKeys {
// 		sl = append(sl, k.Text())
// 	}
// 	t.setText(fmt.Sprintf(formats, sl...))
// 	return t
// }

// // 设置新的formats,如果formats和已经设置的相同，则不会做任何处理，禁止每帧调用！！！
//
//	func (t *TextUi) SetByFormat(formats string) {
//		if t.formats == formats {
//			return
//		}
//		t.formats = formats
//		sl := []any{}
//		for _, k := range t.formatKeys {
//			sl = append(sl, k.Text())
//		}
//		t.setText(fmt.Sprintf(formats, sl...))
//	}
//
// 在两侧增加空格保证补齐到指定尺寸
func (t *TextUi) AddSpaceToSizeX(x float64) *TextUi {
	nowx := t.Size().X
	if nowx >= x {
		// 尺寸已经大于了就无效
		return t
	}
	w, _ := text.Measure(" ", t.size.face(), 5) // 字体宽高
	spaceLeftCount := int((x - nowx) / w / 2)   // 一侧所需的空格数量
	spaceS := strings.Repeat(" ", spaceLeftCount)
	t.SetText(spaceS + t.s + spaceS)
	return t
}

func (t *TextUi) SetText(s string) {
	t.setText(s)
}

func (t *TextUi) setText(s string) {
	if t.s == s {
		return
	}
	t.s = s
	w, h := text.Measure(s, t.size.face(), 5) // 字体宽高
	halfSize := utils.NewPoint(w, h).Divf1(2)
	wi, hi := int(w), int(h)
	hsxi, hsyi := halfSize.BreakInt()
	if t.img != nil && hsxi == wi && hsyi == hi {
		// 尺寸一致不用新建
		t.img.Clear()
	} else {
		//尺寸不一致，todo还可以优化，把newimg提前存一张大图，这里只取子图
		t.halfSize = halfSize
		t.setOp(t.offset)
		t.img = ebiten.NewImage(int(w), int(h))
	}

	textUiFaces.options.ColorScale.Reset()
	textUiFaces.options.ColorScale.ScaleWithColor(t.color)
	text.Draw(t.img, t.s, t.size.face(), textUiFaces.options)
}

// 设置坐标偏移
func (t *TextUi) SetPosOffset(offset utils.Point) {
	t.offset = offset
	t.setOp(offset)
}
func (t *TextUi) SetPos(pos utils.Point, ali utils.AlignmentType) {
	t.pos = pos
	t.ali = ali
	t.setOp(t.offset)
}

func (t *TextUi) setOp(offset utils.Point) {
	pos := t.ali.GetAlignmentPos(t.pos, t.halfSize).Add(offset) // 渲染坐标
	t.op.GeoM.Reset()
	t.op.GeoM.Translate(pos.Break())
}

func (t *TextUi) Draw(screen *ebiten.Image) {
	screen.DrawImage(t.img, t.op)
}

func (t *TextUi) Size() utils.Point {
	return t.halfSize.MulF1(2)
}
