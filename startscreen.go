package submarinegame

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"gitlab.com/high-creek-software/go2d/collision"
	"gitlab.com/high-creek-software/go2d/components/debug"
	"gitlab.com/high-creek-software/go2d/components/display"
	"golang.org/x/image/colornames"
	"slices"
)

const (
	GAME_TITLE = "Submarine Game"
)

type StartScreen struct {
	gameStarted func(shipType ShipType)

	shipOne *ShipSelection
	shipTwo *ShipSelection

	startButton *BasicButton
}

func NewStartScreen(gameStarted func(shipType ShipType)) *StartScreen {

	halfWidth := SCREEN_WIDTH / 2

	shipOne := NewShipSelection(halfWidth-halfWidth/2-128/2, 200, "assets/ship/idle")
	shipOne.selected = true
	shipTwo := NewShipSelection(SCREEN_WIDTH/2+halfWidth/2+128/2, 200, "assets/ship2/idle")

	return &StartScreen{
		gameStarted: gameStarted,
		shipOne:     shipOne,
		shipTwo:     shipTwo,
		startButton: NewBasicButton(0, 0, 100, 50),
	}
}

func (s *StartScreen) Update() error {
	keys := inpututil.AppendJustPressedKeys(nil)
	if slices.Contains(keys, ebiten.KeyEnter) {

		st := ShipTypeOne

		if s.shipTwo.selected {
			st = ShipTypeTwo
		}

		s.gameStarted(st)
	}

	s.shipOne.Update()
	s.shipTwo.Update()
	s.startButton.Update()

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
		x, y := ebiten.CursorPosition()

		rect := collision.NewRect(float64(x), float64(y), 1, 1)

		if rect.AlignedCollides(s.shipOne.Rect()) {
			s.shipOne.selected = true
			s.shipTwo.selected = false
		} else if rect.AlignedCollides(s.shipTwo.Rect()) {
			s.shipTwo.selected = true
			s.shipOne.selected = false
		} else if rect.AlignedCollides(s.startButton.Rect()) {
			st := ShipTypeOne

			if s.shipTwo.selected {
				st = ShipTypeTwo
			}

			s.gameStarted(st)
		}
	}

	return nil
}

func (s *StartScreen) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0.0, 0.0, float32(SCREEN_WIDTH), float32(SCREEN_HEIGHT), colornames.Lightblue, true)
	//ebitenutil.DebugPrintAt(screen, "Press Enter To Start", int(SCREEN_WIDTH/2-100), int(SCREEN_HEIGHT/2-10))

	s.shipOne.Draw(screen)
	s.shipTwo.Draw(screen)
	s.startButton.Draw(screen)

	face := &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   22,
	}
	width, _ := text.Measure(GAME_TITLE, face, 1)
	titleOpts := &text.DrawOptions{}
	titleOpts.GeoM.Translate(SCREEN_WIDTH/2-width/2, 50)
	titleOpts.ColorScale.ScaleWithColor(colornames.Black)
	text.Draw(
		screen,
		GAME_TITLE,
		face,
		titleOpts,
	)
}

type ShipSelection struct {
	*Entity

	animation       *display.AnimateComponent
	hitBoxComponent *debug.HitBoxDrawComponent
	selected        bool
}

func (s *ShipSelection) Scale() (float64, float64) {
	return 1, 1
}

func (s *ShipSelection) IsHorizontalFlipped() bool {
	return false
}

func (s *ShipSelection) IsVerticalFlipped() bool {
	return false
}

func (s *ShipSelection) DrawOffset() (float64, float64) {
	return 0, 0
}

func (s *ShipSelection) Update() error {
	s.animation.Update()
	return nil
}

func (s *ShipSelection) Draw(screen *ebiten.Image) {
	s.animation.Draw(screen, 0, 0)

	if s.selected {
		s.hitBoxComponent.Draw(screen, 0, 0)
	}
}

func NewShipSelection(x, y float64, root string) *ShipSelection {

	ship := &ShipSelection{
		Entity: NewEntity(x, y, 128, 45),
	}

	sprites := assetLoader.LoadConcurrentDirectory(root, "png", 4)
	animation := display.NewAnimateComponent(ship, 3, sprites).CenterImage()

	ship.animation = animation
	ship.hitBoxComponent = debug.NewHitBoxDrawComponent(ship, true, 128, 45)

	return ship
}

type BasicButton struct {
	*Entity
}

func (bb *BasicButton) Update() error {

	x := SCREEN_WIDTH/2 - bb.Width/2
	y := SCREEN_HEIGHT - bb.Height - 50

	bb.X = x
	bb.Y = y

	return nil
}

func (bb *BasicButton) Draw(screen *ebiten.Image) {

	vector.DrawFilledRect(screen, float32(bb.X), float32(bb.Y), float32(bb.Width), float32(bb.Height), colornames.White, true)

	lbl := "Start"
	face := &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   16,
	}
	width, height := text.Measure(lbl, face, 1)
	titleOpts := &text.DrawOptions{}
	titleOpts.GeoM.Translate(bb.X+bb.Width/2-width/2, bb.Y+bb.Height/2-height/2)
	titleOpts.ColorScale.ScaleWithColor(colornames.Black)
	text.Draw(
		screen,
		lbl,
		face,
		titleOpts,
	)

}

func NewBasicButton(x, y, width, height float64) *BasicButton {
	return &BasicButton{
		Entity: NewEntity(x, y, width, height),
	}
}
