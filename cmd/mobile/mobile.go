package mobile

import (
	"github.com/hajimehoshi/ebiten/v2/mobile"
	submarinegame "github.com/high-creek-software/submarine-game"
)

func init() {
	mobile.SetGame(submarinegame.NewGame())
}

func Dummy() {}
