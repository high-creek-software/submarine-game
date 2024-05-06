package submarinegame

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"slices"
)

type GameOverScreen struct {
	restart func()
}

func NewGameOverScreen(restart func()) *GameOverScreen {
	return &GameOverScreen{
		restart: restart,
	}
}

func (g *GameOverScreen) Update() error {
	keys := inpututil.AppendJustPressedKeys(nil)
	if slices.Contains(keys, ebiten.KeyEnter) {
		g.restart()
	}
	return nil
}

func (g *GameOverScreen) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0.0, 0.0, float32(SCREEN_WIDTH), float32(SCREEN_HEIGHT), colornames.Red, true)
	ebitenutil.DebugPrintAt(screen, "You Died: Press Enter To Start", int(SCREEN_WIDTH/2-100), int(SCREEN_HEIGHT/2-10))
}
