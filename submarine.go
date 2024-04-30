package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand/v2"
	"time"
)

type RedOctober struct {
	*Submarine
}

func (r *RedOctober) Update() error {
	r.Submarine.Update()

	// Random call
	// Crazy Ivan
	crazy := rand.IntN(100)%55 == 0
	// change direction
	if crazy {
		r.horizontalFlipped = true
	}

	return nil
}

type Submarine struct {
	*Entity
	subSpeed          float64
	requestTorpedo    func(sub *Submarine)
	horizontalFlipped bool

	lastFire time.Time
	cooldown time.Duration

	sprite *ebiten.Image
}

func NewSubmarine(requestTorpedo func(sub *Submarine)) *Submarine {
	// Dimensions
	// Location

	// Height groups: 200-320, 320-500, 500-720

	y := pickLayer()

	lane := rand.IntN(3)

	var x float64
	switch lane {
	case 0:
		x = 64.0
	case 1:
		x = SCREEN_WIDTH/2 - 128/2
	case 2:
		x = SCREEN_WIDTH - 128/2
	}

	// Picking the milliseconds of cooldown
	cooldownDuration := (1500 + rand.IntN(4500))
	// Parsing the duration, and interval of time
	cooldown, _ := time.ParseDuration(fmt.Sprintf("%dms", cooldownDuration))

	img := assetLoader.MustLoadImage("assets/sub1/0.png")
	sub := &Submarine{
		Entity:         NewEntity(x, float64(y), 128, 25).SetCentered(),
		subSpeed:       2,
		sprite:         img,
		lastFire:       time.Now(),
		cooldown:       cooldown,
		requestTorpedo: requestTorpedo,
	}

	return sub
}

func SpawnSub(requestTorpedo func(sub *Submarine)) *Submarine {
	y := pickLayer()
	img := assetLoader.MustLoadImage("assets/sub1/0.png")
	// Picking the milliseconds of cooldown
	cooldownDuration := (1500 + rand.IntN(2500))
	// Parsing the duration, and interval of time
	cooldown, _ := time.ParseDuration(fmt.Sprintf("%dms", cooldownDuration))

	horizontalFlipped := false
	x := -128

	if rand.IntN(100)%2 == 0 {
		horizontalFlipped = true
		x = SCREEN_WIDTH + 128
	}

	sub := &Submarine{
		Entity:            NewEntity(float64(x), float64(y), 128, 25).SetCentered(),
		subSpeed:          3,
		sprite:            img,
		lastFire:          time.Now(),
		cooldown:          cooldown,
		requestTorpedo:    requestTorpedo,
		horizontalFlipped: horizontalFlipped,
	}

	return sub
}

func pickLayer() float64 {
	layer := rand.IntN(3)
	var y float64
	switch layer {
	case 0:
		y = float64(200 + rand.IntN(120))
	case 1:
		y = float64(320 + rand.IntN(180))
	case 2:
		y = float64(500 + rand.IntN(220))
	}

	return y
}

func (s *Submarine) Update() error {
	speedModifier := 1
	if s.horizontalFlipped {
		speedModifier = -1
	}
	s.X = s.X + s.subSpeed*float64(speedModifier)

	if s.X > SCREEN_WIDTH {
		s.horizontalFlipped = true
	} else if s.X < 0 {
		s.horizontalFlipped = false
	}

	now := time.Now()
	if now.Sub(s.lastFire) > s.cooldown {
		s.requestTorpedo(s)
		s.lastFire = now
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
