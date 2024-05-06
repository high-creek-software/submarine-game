package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	submarinegame "github.com/high-creek-software/submarine-game"
)

func main() {

	subGame := submarinegame.NewGame()

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Submarine Game")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	ebiten.RunGame(subGame)
}
