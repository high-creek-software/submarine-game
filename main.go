package main

import (
	"embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"gitlab.com/high-creek-software/go2d/loader"
	"golang.org/x/image/colornames"
	_ "image/png"
	"log/slog"
)

const (
	SCREEN_WIDTH  = 1280
	SCREEN_HEIGHT = 720
)

type Game struct {
	ship         *Ship
	depthCharges []*DepthCharge

	oceanImage *ebiten.Image

	pointerX, pointerY int
}

func (g *Game) requestCharge() {
	slog.Info("Fire requested")
	// Depth charge fired from where the ship
	g.depthCharges = append(g.depthCharges, NewDepthCharge(g.ship.X, g.ship.Y))
}

func (g *Game) Update() error {
	// slog.Info("Update")
	g.ship.Update()

	// This is the new array for active depth charges
	var keepDepthCharges []*DepthCharge
	for _, depthCharge := range g.depthCharges {
		depthCharge.Update()
		if depthCharge.IsActive {
			keepDepthCharges = append(keepDepthCharges, depthCharge)
		}
	}
	// Setting the active array back to the depthCharges
	g.depthCharges = keepDepthCharges

	g.pointerX, g.pointerY = ebiten.CursorPosition()

	// Array of torpedos

	// Array of submarines

	// Check for collisions

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// Background
	vector.DrawFilledRect(screen, 0, 0, SCREEN_WIDTH, SCREEN_HEIGHT, colornames.Aliceblue, true)

	g.ship.Draw(screen)
	// Iteration type of loop
	slog.Info("Depth Charges", "length", len(g.depthCharges))
	for _, depthCharge := range g.depthCharges {
		depthCharge.Draw(screen)
	}

	opts := &ebiten.DrawImageOptions{}
	opts.ColorScale.ScaleAlpha(0.45)
	opts.GeoM.Translate(0, 140)
	screen.DrawImage(g.oceanImage, opts)

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
	subGame.ship = NewShip(subGame.requestCharge)

	oceanImg := assetLoader.MustLoadImage("assets/ocean.png")

	subGame.oceanImage = ebiten.NewImageFromImage(oceanImg)
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Submarine Game")

	ebiten.RunGame(subGame)
}
