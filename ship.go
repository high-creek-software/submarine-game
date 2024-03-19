package main

import (
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

type Ship struct {
	Entity
	shipSpeed     float32
	requestCharge func()
}

func NewShip(requestCharge func()) *Ship {
	return &Ship{
		Entity:        *NewEntity(40, 100, 150, 75),
		shipSpeed:     5,
		requestCharge: requestCharge,
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
	moveX := s.X + (move * s.shipSpeed)

	// Bounds check
	if moveX < 0 { // going off screen left
		moveX = 0
	} else if moveX+s.Width > SCREEN_WIDTH { // going off screen right
		moveX = SCREEN_WIDTH - s.Width
	}

	s.X = moveX
}

func (s *Ship) fire() {
	keys := inpututil.AppendJustPressedKeys(nil)
	if slices.Contains(keys, ebiten.KeySpace) {
		s.requestCharge()
		// Deploy a depth charge
		// Pass methods to this struct to request fire
	}
}

func (s *Ship) Draw(screen *ebiten.Image) {
	// Ship
	vector.DrawFilledRect(screen, s.X, s.Y, s.Width, s.Height, colornames.Orange, true)
}
