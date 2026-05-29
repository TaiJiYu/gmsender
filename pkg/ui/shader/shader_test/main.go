package main

import (
	_ "embed"
	"fmt"
	"gmsender/pkg/ui/shader"
	"gmsender/utils"
	"image/color"
	"os"

	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// shader测试
// go run pkg/ui/shader/shader_test/main.go

var testShaderFile = "E:\\GMStudioProject\\Swift\\pkg\\ui\\shader\\shaderfile\\rounded_rect_lerp.kage"

// var testShaderFile = "E:\\GMStudioProject\\Swift\\pkg\\ui\\shader\\shaderfile\\rounded_rect.kage"
// var testShaderFile = "E:\\GMStudioProject\\Swift\\pkg\\ui\\shader\\shaderfile\\core_rect.kage"

type Game struct {
	begin      time.Time
	shaderop   *ebiten.DrawRectShaderOptions
	testShader *ebiten.Shader
	width      int
}

var (
	gameClt  *Game
	gameOnce sync.Once
)

func main() {
	shader.InitShader()
	ebiten.SetTPS(120)
	ebiten.SetWindowSize(960, 540)
	ebiten.SetWindowTitle("SwiftUiShaderTest")
	g := &Game{
		begin: time.Now(),
		shaderop: &ebiten.DrawRectShaderOptions{
			Uniforms: map[string]any{},
		},
		width: 200,
	}
	// x, y := float64(utils.HalfLogicalSizeX-noise.Bounds().Dx()/2), float64(utils.HalfLogicalSizeY-noise.Bounds().Dy()/2)
	// fmt.Println(x, y)
	// g.shaderop.GeoM.Translate(x, y)
	g.TryRefshShader(testShaderFile)
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	g.TryRefshShader(testShaderFile)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		g.shaderop.GeoM.Translate(50, 0)
	} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.shaderop.GeoM.Translate(0, 50)
	}
	_, y := ebiten.Wheel()
	g.width = max(0, min(500, g.width+2*int(y)))

	g.shaderop.Uniforms["Time"] = time.Since(g.begin).Seconds()

	return nil
}

var lastfreshtime time.Time

func (g *Game) TryRefshShader(shaderfilename string) {
	if time.Since(lastfreshtime) < time.Second {
		return
	}
	lastfreshtime = time.Now()
	data, err := os.ReadFile(shaderfilename)
	if err != nil {
		fmt.Printf("TryRefshShader failed,err:%v\n", err)
		return
	}

	s, err := ebiten.NewShader(data)
	if err != nil {
		fmt.Printf("TryRefshShader failed,err:%v\n", err)
		return
	}
	g.testShader = s
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 143, G: 173, B: 200, A: 255})

	// screen.DrawRectShader(610, 610, g.testShader, g.shaderop)
	screen.DrawRectShader(g.width, 100, g.testShader, g.shaderop)

	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf("TPS: %0.2f\n FPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()),
	)

}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return utils.LogicalSizeX, utils.LogicalSizeY
}
