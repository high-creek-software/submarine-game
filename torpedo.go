package submarinegame

import (
	"github.com/hajimehoshi/ebiten/v2"
	"gitlab.com/high-creek-software/go2d/components/debug"
)

type Torpedo struct {
	*Entity
	speed float64

	hitboxComponent *debug.HitBoxDrawComponent
	sprite          *ebiten.Image
}

func NewTorpedo() *Torpedo {

	sprite := assetLoader.MustLoadImage("assets/torpedo/idle/0.png")
	torpedo := &Torpedo{
		Entity: NewEntity(0, 0, 15, 30).SetCentered(),
		//speed:  speed,
		sprite: sprite,
	}
	//torpedo.hitboxComponent = debug.NewHitBoxDrawComponent(torpedo, true, torpedo.Width, torpedo.Height)

	return torpedo
}

func (t *Torpedo) Update() error {
	t.Y -= t.speed
	if t.Y < WATER_SURFACE {
		t.IsActive = false
	}
	return nil
}

func (t *Torpedo) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(-(t.Width / 2), -(t.Height / 2))

	opts.GeoM.Translate(t.X, t.Y)
	screen.DrawImage(t.sprite, opts)

	//t.hitboxComponent.Draw(screen, 0, 0)
}
