package main

import "github.com/quartercastle/vector"

type Entity struct {
	X, Y          float64
	Width, Height float64
	IsActive      bool
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

func (e *Entity) At() vector.MutableVector {
	return vector.MutableVector{float64(e.X), float64(e.Y)}
}
