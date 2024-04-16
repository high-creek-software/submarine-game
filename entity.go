package main

import (
	"github.com/quartercastle/vector"
	"gitlab.com/high-creek-software/go2d/collision"
)

type Entity struct {
	X, Y          float64
	Width, Height float64
	IsActive      bool
	IsCentered    bool
}

func NewEntity(x, y, width, height float64) *Entity {
	return &Entity{
		X:        x,
		Y:        y,
		Width:    width,
		Height:   height,
		IsActive: true,
	}
}

func (e *Entity) SetCentered() *Entity {
	e.IsCentered = true
	return e
}

func (e *Entity) At() vector.MutableVector {
	return vector.MutableVector{float64(e.X), float64(e.Y)}
}

func (e *Entity) Rect() collision.Rect {
	x, y := e.X, e.Y
	if e.IsCentered {
		x = x - (e.Width / 2)
		y = y - (e.Height / 2)
	}

	return collision.NewRect(x, y, e.Width, e.Height)
}
