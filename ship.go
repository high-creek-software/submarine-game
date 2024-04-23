package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"gitlab.com/high-creek-software/go2d/components/debug"
	"gitlab.com/high-creek-software/go2d/components/display"
	"golang.org/x/image/colornames"
	"slices"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Ship struct {
	*Entity
	shipSpeed     float64
	requestCharge func(isFront bool)

	sprite          *ebiten.Image
	hitboxComponent *debug.HitBoxDrawComponent

	maxHealth     float64
	health        float64
	idleAnimation *display.AnimateComponent

	lastFire time.Time
	cooldown time.Duration

	score int64
}

func NewShip() *Ship {
	ship := &Ship{
		Entity:    NewEntity(40, 105, 128, 40),
		shipSpeed: 3,
		maxHealth: 6,
		health:    6,
		cooldown:  800 * time.Millisecond,
	}

	//ship.sprite = assetLoader.MustLoadImage("assets/ship/0.png")

	idleShip := assetLoader.LoadConcurrentDirectory("assets/ship2", "png", 4)
	ship.idleAnimation = display.NewAnimateComponent(ship, 3, idleShip)

	ship.hitboxComponent = debug.NewHitBoxDrawComponent(ship, false, ship.Width, ship.Height)

	return ship
}

func (s *Ship) UpdateForLevel(requestCharge func(isFront bool)) {
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

func (s *Ship) internalFire(isFront bool) {
	now := time.Now()
	if now.Sub(s.lastFire) > s.cooldown {
		s.requestCharge(isFront)
		s.lastFire = now
	}
}

func (s *Ship) fire() {
	keys := inpututil.AppendJustPressedKeys(nil)
	if slices.Contains(keys, ebiten.KeySpace) {
		s.internalFire(false)
	} else if slices.Contains(keys, ebiten.KeyA) {
		s.internalFire(false)
	} else if slices.Contains(keys, ebiten.KeyD) {
		s.internalFire(true)
	}
}

func (s *Ship) Draw(screen *ebiten.Image) {
	// Ship
	//vector.DrawFilledRect(screen, s.X, s.Y, s.Width, s.Height, colornames.Orange, true)

	/*opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(s.X), float64(s.Y))
	screen.DrawImage(s.sprite, opts)*/
	s.idleAnimation.Draw(screen, 0, 0)
	//s.hitboxComponent.Draw(screen, 0, 0)

	maxWidth := float32(300.0)
	vector.DrawFilledRect(screen, 50, 15, maxWidth, 15, colornames.Red, true)
	healthRatio := s.health / s.maxHealth
	currentHealthWidth := maxWidth * float32(healthRatio)
	vector.DrawFilledRect(screen, 50, 15, currentHealthWidth, 15, colornames.Green, true)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", s.score), 400, 20)
}

func (s *Ship) incrementScore(delta int64) {
	s.score += delta
}

func (s *Ship) WasHit() {
	s.health -= 1
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
