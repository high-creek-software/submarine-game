package main

import (
	"log/slog"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

type Ship struct {
	shipX, shipY          float32
	shipSpeed             float32
	shipWidth, shipHeight float32
}

func NewShip() *Ship {
	return &Ship{
		shipX:      20,
		shipY:      40,
		shipSpeed:  5,
		shipWidth:  150,
		shipHeight: 50,
	}
}

func (s *Ship) Update() error {
	s.move()
	s.fire()
	return nil
}

func (s *Ship) move() {
	var move float32

	// Movement input
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		move = -1
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		move = 1
	}

	// New position
	moveX := s.shipX + (move * s.shipSpeed)

	// Bounds check
	if moveX < 0 { // going off screen left
		moveX = 0
	} else if moveX+s.shipWidth > SCREEN_WIDTH { // going off screen right
		moveX = SCREEN_WIDTH - s.shipWidth
	}

	s.shipX = moveX
}

func (s *Ship) fire() {
	keys := inpututil.AppendJustPressedKeys(nil)
	if slices.Contains(keys, ebiten.KeySpace) {
		slog.Info("Fire requested")
		// Deploy a depth charge
		// Pass methods to this struct to request fire
	}
}

func (s *Ship) Draw(screen *ebiten.Image) {
	// Ship
	vector.DrawFilledRect(screen, s.shipX, s.shipY, s.shipWidth, s.shipHeight, colornames.Orange, true)
}
