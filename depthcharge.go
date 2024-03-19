package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

type DepthCharge struct {
	Entity
	gravity float32
}

func NewDepthCharge(x, y float32) *DepthCharge {
	return &DepthCharge{
		Entity:  *NewEntity(x, y, 20, 20),
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
	vector.DrawFilledRect(screen, dc.X, dc.Y, dc.Width, dc.Height, colornames.Purple, true)
}
