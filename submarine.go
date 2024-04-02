package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
)

type Submarine struct {
	*Entity
	subSpeed          float64
	requestTorpedo    func(sub *Submarine)
	horizontalFlipped bool

	sprite *ebiten.Image
}

func NewSubmarine(requestTorpedo func(sub *Submarine)) *Submarine {
	// Dimensions
	// Location

	// Height groups: 200-320, 320-500, 500-720

	layer := rand.Intn(3)
	var y float64
	switch layer {
	case 0:
		y = float64(200 + rand.Intn(120))
	case 1:
		y = float64(320 + rand.Intn(180))
	case 2:
		y = float64(500 + rand.Intn(220))
	}

	lane := rand.Intn(3)

	var x float64
	switch lane {
	case 0:
		x = 64.0
	case 1:
		x = SCREEN_WIDTH/2 - 128/2
	case 2:
		x = SCREEN_WIDTH - 128/2
	}

	img := assetLoader.MustLoadImage("assets/sub1/0.png")
	sub := &Submarine{
		Entity:   NewEntity(x, float64(y), 128, 25),
		subSpeed: 2,
		sprite:   img,
	}

	return sub
}

func (s *Submarine) Update() error {
	speedModifier := 1
	if s.horizontalFlipped {
		speedModifier = -1
	}
	s.X = s.X + s.subSpeed*float64(speedModifier)

	if s.X+s.Width/2 > SCREEN_WIDTH {
		s.horizontalFlipped = true
	} else if s.X-s.Width/2 < 0 {
		s.horizontalFlipped = false
	}

	return nil
}

func (s *Submarine) Draw(screen *ebiten.Image) {

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(-(s.Width / 2), -(s.Height / 2))
	if s.horizontalFlipped {
		opts.GeoM.Scale(-1, 1)
	}

	opts.GeoM.Translate(s.X, s.Y)

	screen.DrawImage(s.sprite, opts)
}
