package main

import (
	"gmsender/gmsender"
	"gmsender/pkg/asset"
	"gmsender/pkg/netfinder"
	"gmsender/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowIcon(asset.Ico())
	ebiten.SetWindowSize(utils.LogicalSizeX, utils.LogicalSizeY)
	ebiten.SetTPS(utils.BigTPS)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetWindowDecorated(false)
	ebiten.RunGameWithOptions(gmsender.NewGMSender(), &ebiten.RunGameOptions{
		ScreenTransparent: true,
	})
	netfinder.Close()
}
