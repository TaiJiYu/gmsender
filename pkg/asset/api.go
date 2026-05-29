package asset

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	//go:embed NotoSansSC-Regular.ttf
	fontBytes []byte
	//go:embed close.png
	closeBytes []byte
	closeImg   *ebiten.Image
	//go:embed small.png
	smallBytes []byte
	smallImg   *ebiten.Image

	//go:embed smallball.png
	smallBallBytes []byte
	smallBallImg   *ebiten.Image
)

func init() {
	if img, _, err := image.Decode(bytes.NewReader(closeBytes)); err != nil {
		panic(err)
	} else {
		closeImg = ebiten.NewImageFromImage(img)
	}

	if img, _, err := image.Decode(bytes.NewReader(smallBytes)); err != nil {
		panic(err)
	} else {
		smallImg = ebiten.NewImageFromImage(img)
	}

	if img, _, err := image.Decode(bytes.NewReader(smallBallBytes)); err != nil {
		panic(err)
	} else {
		smallBallImg = ebiten.NewImageFromImage(img)
	}
}

func FontData() []byte {
	return fontBytes
}

func CloseImg() *ebiten.Image {
	return closeImg
}

func SmallImg() *ebiten.Image {
	return smallImg
}
func SmallBallImg() *ebiten.Image {
	return smallBallImg
}
