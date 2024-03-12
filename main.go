package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

const (
	SCREEN_WIDTH  = 1280
	SCREEN_HEIGHT = 720
)

type Game struct {
	ship *Ship
}

func (g *Game) Update() error {
	// slog.Info("Update")
	g.ship.Update()

	// Array of depth charges

	// Array of torpedos

	// Array of submarines

	// Check for collisions
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// Background
	vector.DrawFilledRect(screen, 0, 0, SCREEN_WIDTH, SCREEN_HEIGHT, colornames.Aliceblue, true)

	g.ship.Draw(screen)
}

func (g *Game) Layout(ow, oh int) (w, h int) {

	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func main() {

	subGame := &Game{
		ship: NewShip(),
	}

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Submarine Game")

	ebiten.RunGame(subGame)
}
