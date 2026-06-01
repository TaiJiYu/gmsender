package asset

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"io"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
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

	//go:embed rein.png
	reInBytes []byte
	reInImg   *ebiten.Image

	//go:embed done.ogg
	doneMusicBytes []byte
	doneMusic      *audio.Player
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

	if img, _, err := image.Decode(bytes.NewReader(reInBytes)); err != nil {
		panic(err)
	} else {
		reInImg = ebiten.NewImageFromImage(img)
	}
	doneMusic = ogg(doneMusicBytes)
}

func ogg(oggBytes []byte) *audio.Player {
	s, err := vorbis.DecodeF32(bytes.NewReader(oggBytes))
	if err != nil {
		panic(err)
	}
	if data, err := io.ReadAll(s); err != nil {
		panic(err)
	} else {
		if audio.CurrentContext() == nil {
			audio.NewContext(44100)
		}
		return audio.CurrentContext().NewPlayerF32FromBytes(data)
	}
}

// 播放完成音效
func PlayDoneMusic() {
	doneMusic.Rewind()
	doneMusic.Play()
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

func ReinImg() *ebiten.Image {
	return reInImg
}

func SmallBallImg() *ebiten.Image {
	return smallBallImg
}
