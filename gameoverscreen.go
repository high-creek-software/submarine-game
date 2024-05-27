package submarinegame

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"gitlab.com/high-creek-software/go2d/collision"
	"golang.org/x/image/colornames"
	"log/slog"
)

type GameOverScreen struct {
	restart func()

	restartButton *BasicButton
}

func NewGameOverScreen(restart func()) *GameOverScreen {
	return &GameOverScreen{
		restart:       restart,
		restartButton: NewBasicButton(0, 0, 100, 50, "Restart"),
	}
}

func (g *GameOverScreen) Update() error {
	g.restartButton.Update()

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
		x, y := ebiten.CursorPosition()
		g.handleTouch(x, y)
	}

	touchIDS := inpututil.AppendJustReleasedTouchIDs(nil)
	for _, tid := range touchIDS {
		x, y := inpututil.TouchPositionInPreviousTick(tid)
		slog.Info("Location", "x", x, "y", y)
		g.handleTouch(x, y)
	}

	return nil
}

func (g *GameOverScreen) handleTouch(x, y int) {
	rect := collision.NewRect(float64(x), float64(y), 10, 10)
	if rect.AlignedCollides(g.restartButton.Rect()) {
		g.restart()
	}
}

func (g *GameOverScreen) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0.0, 0.0, float32(SCREEN_WIDTH), float32(SCREEN_HEIGHT), colornames.Red, true)

	g.restartButton.Draw(screen)
}
