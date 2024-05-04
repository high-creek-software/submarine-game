package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"slices"
)

type GameOverScreen struct {
	gameStarted func()
}

func NewGameOverScreen(gameStarted func()) *GameOverScreen {
	return &GameOverScreen{
		gameStarted: gameStarted,
	}
}

func (g *GameOverScreen) Update() error {
	keys := inpututil.AppendJustPressedKeys(nil)
	if slices.Contains(keys, ebiten.KeyEnter) {
		g.gameStarted()
	}
	return nil
}

func (g *GameOverScreen) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0.0, 0.0, float32(SCREEN_WIDTH), float32(SCREEN_HEIGHT), colornames.Red, true)
	ebitenutil.DebugPrintAt(screen, "You Died: Press Enter To Start", int(SCREEN_WIDTH/2-100), int(SCREEN_HEIGHT/2-10))
}
