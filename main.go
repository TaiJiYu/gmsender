package main

import (
	"gmsender/gmsender"
	"gmsender/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(utils.LogicalSizeX, utils.LogicalSizeY)
	ebiten.SetTPS(utils.BigTPS)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetWindowDecorated(false)
	ebiten.RunGameWithOptions(gmsender.NewGMSender(), &ebiten.RunGameOptions{
		ScreenTransparent: true,
	})
	// netfinder.Close()
}
