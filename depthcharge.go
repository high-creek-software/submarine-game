package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type DepthCharge struct {
	*Entity
	gravity float64

	sprite *ebiten.Image
}

func NewDepthCharge(x, y float64) *DepthCharge {

	sprite := assetLoader.MustLoadImage("assets/depthcharge/idle/0.png")
	return &DepthCharge{
		Entity:  NewEntity(x, y, 20, 20).SetCentered(),
		gravity: 1.5,
		sprite:  sprite,
	}
}

func (dc *DepthCharge) Update() error {

	moveY := dc.Y + dc.gravity

	dc.Y = moveY

	if dc.Y > SCREEN_HEIGHT {
		dc.IsActive = false
	}

	return nil
}

func (dc *DepthCharge) Draw(screen *ebiten.Image) {
	//vector.DrawFilledRect(screen, float32(dc.X), float32(dc.Y), float32(dc.Width), float32(dc.Height), colornames.Purple, true)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(-(dc.Width / 2), -(dc.Height / 2))

	opts.GeoM.Translate(dc.X, dc.Y)
	screen.DrawImage(dc.sprite, opts)
}
