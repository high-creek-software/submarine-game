package submarinegame

import (
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"gitlab.com/high-creek-software/go2d/collision"
	"gitlab.com/high-creek-software/go2d/components/debug"
	"gitlab.com/high-creek-software/go2d/components/display"
	"golang.org/x/image/colornames"
	"slices"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ShipType int

const (
	ShipTypeOne ShipType = iota
	ShipTypeTwo
)

type touch struct {
	x, y int
}

type Ship struct {
	*Entity
	shipSpeed     float64
	requestCharge func(isFront bool)

	sprite          *ebiten.Image
	hitboxComponent *debug.HitBoxDrawComponent

	maxHealth       float64
	health          float64
	idleAnimation   *display.AnimateComponent
	brokenAnimation *display.AnimateComponent

	lastFire time.Time
	cooldown time.Duration

	score int64

	leftControl     *Control
	rightControl    *Control
	fireAftControl  *Control
	fireForeControl *Control
	touchIDS        []ebiten.TouchID
	touches         map[ebiten.TouchID]*touch
}

func NewShip(shipType ShipType) *Ship {
	ship := &Ship{
		Entity:  NewEntity(40, 105, 128, 40),
		touches: make(map[ebiten.TouchID]*touch),
	}

	var idleShip []*ebiten.Image
	var brokenShip []*ebiten.Image
	if shipType == ShipTypeOne {
		ship.shipSpeed = 3
		ship.maxHealth = 6
		ship.health = 6
		ship.cooldown = 800 * time.Millisecond
		idleShip = assetLoader.LoadConcurrentDirectory("assets/ship/idle", "png", 4)
		brokenShip = assetLoader.LoadConcurrentDirectory("assets/ship/broken", "png", 4)
	} else {
		ship.shipSpeed = 4
		ship.maxHealth = 6
		ship.health = 6
		ship.cooldown = 1000 * time.Millisecond
		idleShip = assetLoader.LoadConcurrentDirectory("assets/ship2/idle", "png", 4)
		brokenShip = assetLoader.LoadConcurrentDirectory("assets/ship2/broken", "png", 4)
	}

	ship.idleAnimation = display.NewAnimateComponent(ship, 3, idleShip)
	ship.brokenAnimation = display.NewAnimateComponent(ship, 3, brokenShip)
	ship.hitboxComponent = debug.NewHitBoxDrawComponent(ship, false, ship.Width, ship.Height)

	if isMobile {
		ship.leftControl = NewControl(ControlTypeLeft, 75, SCREEN_HEIGHT-75)
		ship.rightControl = NewControl(ControlTypeRight, 75+100+40, SCREEN_HEIGHT-75)
		ship.fireForeControl = NewControl(ControlTypeFireFore, SCREEN_WIDTH-75, SCREEN_HEIGHT-75)
		ship.fireAftControl = NewControl(ControlTypeFireAft, SCREEN_WIDTH-75-100-40, SCREEN_HEIGHT-75)

	}

	return ship
}

func (s *Ship) UpdateForLevel(requestCharge func(isFront bool)) {
	s.requestCharge = requestCharge
}

func (s *Ship) Update() error {
	if isMobile {
		// Checking any touches that have just been released
		for id, _ := range s.touches {
			if inpututil.IsTouchJustReleased(id) {
				delete(s.touches, id) // Removing touches if the touchID has just been released
			}
		}
		// Looking for any new touches in this frame
		s.touchIDS = inpututil.AppendJustPressedTouchIDs(s.touchIDS[:0])
		for _, id := range s.touchIDS {
			x, y := ebiten.TouchPosition(id)
			// Creating a new touch at this position and adding it to a touches map
			s.touches[id] = &touch{
				x: x,
				y: y,
			}
		}
		// Updating the location of a touch (ie the finger was moved).
		s.touchIDS = ebiten.AppendTouchIDs(s.touchIDS[:0])
		for _, id := range s.touchIDS {
			if t, ok := s.touches[id]; ok {
				t.x, t.y = ebiten.TouchPosition(id)
			}
		}
	}
	if s.health > 0 {
		s.move()
		s.fire()
		s.idleAnimation.Update()
	} else {
		s.brokenAnimation.Update()
		s.Y += 1
		if s.Y > SCREEN_HEIGHT {
			s.IsActive = false
		}
	}

	return nil
}

func (s *Ship) move() {
	var move float64

	// Movement input
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyJ) {
		move = -1
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyL) {
		move = 1
	}

	if s.leftControl != nil {
		for _, t := range s.touches {
			x, y := t.x, t.y
			touchCircle := collision.NewCircle(float64(x), float64(y), 10)
			if touchCircle.Collides(s.leftControl.HitCircle()) {
				move = -1
			} else if touchCircle.Collides(s.rightControl.HitCircle()) {
				move = 1
			}
		}
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

	if s.leftControl != nil {
		for _, t := range s.touches {
			x, y := t.x, t.y
			touchCircle := collision.NewCircle(float64(x), float64(y), 10)
			if touchCircle.Collides(s.fireAftControl.HitCircle()) {
				s.internalFire(false)
			} else if touchCircle.Collides(s.fireForeControl.HitCircle()) {
				s.internalFire(true)
			}
		}
	}
}

func (s *Ship) Draw(screen *ebiten.Image) {

	if s.health > 0 {
		s.idleAnimation.Draw(screen, 0, 0)
	} else {
		s.brokenAnimation.Draw(screen, 0, 0)
	}
	//s.hitboxComponent.Draw(screen, 0, 0)

	maxWidth := float32(300.0)
	vector.DrawFilledRect(screen, 50, 15, maxWidth, 15, colornames.Red, true)
	healthRatio := s.health / s.maxHealth
	currentHealthWidth := maxWidth * float32(healthRatio)
	vector.DrawFilledRect(screen, 50, 15, currentHealthWidth, 15, colornames.Green, true)

	//ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", s.score), 400, 20)
	txtOpts := &text.DrawOptions{}
	txtOpts.GeoM.Translate(400, 10)
	txtOpts.ColorScale.ScaleWithColor(colornames.Black)
	text.Draw(screen,
		gamePrinter.Sprintf("%d", s.score),
		&text.GoTextFace{Source: mplusFaceSource, Size: 18},
		txtOpts,
	)

	if s.leftControl != nil {
		s.leftControl.Draw(screen)
		s.rightControl.Draw(screen)
		s.fireAftControl.Draw(screen)
		s.fireForeControl.Draw(screen)
	}
}

func (s *Ship) incrementScore(delta int64) {
	s.score += delta
}

func (s *Ship) WasHit() {
	if s.health == 0 {
		return
	}
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
