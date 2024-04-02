package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

type DepthCharge struct {
	*Entity
	gravity float64
}

func NewDepthCharge(x, y float64) *DepthCharge {
	return &DepthCharge{
		Entity:  NewEntity(x, y, 20, 20),
		gravity: 3,
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
	vector.DrawFilledRect(screen, float32(dc.X), float32(dc.Y), float32(dc.Width), float32(dc.Height), colornames.Purple, true)
}
