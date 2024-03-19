package main

type Entity struct {
	X, Y          float32
	Width, Height float32
	IsActive      bool
}

func NewEntity(x, y, width, height float32) *Entity {
	return &Entity{
		X:        x,
		Y:        y,
		Width:    width,
		Height:   height,
		IsActive: true,
	}
}
