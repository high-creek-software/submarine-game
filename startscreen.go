package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"slices"
)

type StartScreen struct {
	gameStarted func()
}

func NewStartScreen(gameStarted func()) *StartScreen {
	return &StartScreen{
		gameStarted: gameStarted,
	}
}

func (s *StartScreen) Update() error {
	keys := inpututil.AppendJustPressedKeys(nil)
	if slices.Contains(keys, ebiten.KeyEnter) {
		s.gameStarted()
	}
	return nil
}

func (s *StartScreen) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0.0, 0.0, SCREEN_WIDTH, SCREEN_HEIGHT, colornames.Lightblue, true)
	ebitenutil.DebugPrintAt(screen, "Press Enter To Start", SCREEN_WIDTH/2-100, SCREEN_HEIGHT/2-10)
}
