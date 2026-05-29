package ui

import (
	"gmsender/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

// 尺寸框
type SizeBox struct {
	size utils.Point // 尺寸
}

func newSizeBox(size utils.Point) *SizeBox {
	return &SizeBox{
		size: size,
	}
}
func (s *SizeBox) Size() utils.Point                               { return s.size }
func (s *SizeBox) Draw(screen *ebiten.Image)                       {}
func (s *SizeBox) SetPosOffset(offset utils.Point)                 {}
func (s *SizeBox) SetPos(pos utils.Point, ali utils.AlignmentType) {}
