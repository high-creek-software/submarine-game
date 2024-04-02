package main

import (
	"gitlab.com/high-creek-software/go2d/components/display"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Ship struct {
	*Entity
	shipSpeed     float64
	requestCharge func()

	sprite *ebiten.Image

	idleAnimation *display.AnimateComponent
}

func NewShip() *Ship {
	ship := &Ship{
		Entity:    NewEntity(40, 80, 150, 75),
		shipSpeed: 5,
	}

	//ship.sprite = assetLoader.MustLoadImage("assets/ship/0.png")

	idleShip := assetLoader.LoadConcurrentDirectory("assets/ship", "png", 4)
	ship.idleAnimation = display.NewAnimateComponent(ship, 3, idleShip)

	return ship
}

func (s *Ship) UpdateForLevel(requestCharge func()) {
	s.requestCharge = requestCharge
}

func (s *Ship) Update() error {
	s.move()
	s.fire()
	s.idleAnimation.Update()
	return nil
}

func (s *Ship) move() {
	var move float64

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
	//vector.DrawFilledRect(screen, s.X, s.Y, s.Width, s.Height, colornames.Orange, true)

	/*opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(s.X), float64(s.Y))
	screen.DrawImage(s.sprite, opts)*/
	s.idleAnimation.Draw(screen, 0, 0)
}

func (s *Ship) Scale() (float64, float64) {
	return 1, 1
}

func (s *Ship) IsHorizontalFlipped() bool {
	return false
}

func (s *Ship) IsVerticalFlipped() bool {
	return false
}

func (s *Ship) DrawOffset() (float64, float64) {
	return 0, 0
}
