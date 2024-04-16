package main

import (
	"embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"gitlab.com/high-creek-software/go2d/loader"
	_ "image/png"
)

const (
	SCREEN_WIDTH  = 1280
	SCREEN_HEIGHT = 720
	WATER_SURFACE = 140
)

type Game struct {
	ship               *Ship
	levelIndex         int
	level              *Level
	pointerX, pointerY int
}

func (g *Game) Update() error {
	// slog.Info("Update")
	g.level.Update()

	g.pointerX, g.pointerY = ebiten.CursorPosition()

	// Array of torpedos

	// Array of submarines

	// Check for collisions

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.level.Draw(screen)
	// Sentinel loop with a slice
	// for i := 0; i < len(g.depthCharges); i++ {
	// 	g.depthCharges[i].Draw(screen)
	// }

	// Sentinel type of loop
	// for i := 0; i < 10; i++ {
	// 	slog.Info("Sentinel", "i", i)
	// }

	// While loops cousin
	// for true {

	// }

	ebitenutil.DebugPrint(screen, fmt.Sprintf("%d, %d", g.pointerX, g.pointerY))
}

func (g *Game) Layout(ow, oh int) (w, h int) {

	return SCREEN_WIDTH, SCREEN_HEIGHT
}

//go:embed assets
var assets embed.FS

var assetLoader *loader.AssetLoader

func main() {
	assetLoader = loader.NewAssetLoader(assets)

	subGame := &Game{}
	subGame.ship = NewShip()

	subGame.level = NewLevel(5, subGame.ship)

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Submarine Game")

	ebiten.RunGame(subGame)
}
