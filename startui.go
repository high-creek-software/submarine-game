package submarinegame

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

type StartUI struct {
	gameStarted func()

	ui *ebitenui.UI
}

func NewStartUI(gameStarted func()) *StartUI {

	su := &StartUI{
		gameStarted: gameStarted,
	}

	rc := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}),
			widget.GridLayoutOpts.Spacing(0, 20),
		)),
	)

	rc.AddChild(su.header())

	su.ui = &ebitenui.UI{
		Container: rc,
	}

	return su
}

func (s *StartUI) Update() error {
	s.ui.Update()
	return nil
}

func (s *StartUI) Draw(screen *ebiten.Image) {
	s.ui.Draw(screen)
}

func (s *StartUI) header() widget.PreferredSizeLocateableWidget {
	c := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(15),
		)),
	)

	//ff :=
	//c.AddChild(widget.NewText(
	//	widget.TextOpts.Text("Submarine Game", mplusFaceSource, colornames.White),
	//))

	return c
}
