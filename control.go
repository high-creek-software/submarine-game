package submarinegame

import (
	"github.com/hajimehoshi/ebiten/v2"
	"gitlab.com/high-creek-software/go2d/collision"
	"gitlab.com/high-creek-software/go2d/components/debug"
)

type ControlType int

const (
	ControlTypeLeft ControlType = iota
	ControlTypeRight
	ControlTypeFireAft
	ControlTypeFireFore
)

type Control struct {
	*Entity

	sprite      *ebiten.Image
	controlType ControlType
	radius      float64
	isHFlipped  bool

	originComponent *debug.OriginDrawComponent
}

func NewControl(ct ControlType, x, y float64) *Control {

	c := &Control{Entity: NewEntity(x, y, 256, 256)}
	c.radius = 64

	switch ct {
	case ControlTypeLeft:
		c.isHFlipped = true
		c.sprite = assetLoader.MustLoadImage("assets/control_direction.png")
	case ControlTypeRight:
		c.sprite = assetLoader.MustLoadImage("assets/control_direction.png")
	case ControlTypeFireAft:
		c.isHFlipped = true
		c.sprite = assetLoader.MustLoadImage("assets/control_fire.png")
	case ControlTypeFireFore:
		c.sprite = assetLoader.MustLoadImage("assets/control_fire.png")
	}

	c.originComponent = debug.NewOriginDrawComponent(c)

	return c
}

func (c *Control) Update() error {

	return nil
}

func (c *Control) HitCircle() collision.Circle {
	return collision.NewCircle(c.X, c.Y, c.radius)
}

func (c *Control) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(0.35, 0.35)
	opts.GeoM.Translate(-c.radius, -c.radius)
	if c.isHFlipped {
		opts.GeoM.Scale(-1, 1)
	}
	opts.GeoM.Translate(c.X, c.Y)

	screen.DrawImage(c.sprite, opts)

	//c.originComponent.Draw(screen, 0, 0)
}
