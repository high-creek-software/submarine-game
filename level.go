package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"log/slog"
	"sync"
)

type Level struct {
	index        int
	ship         *Ship
	depthCharges []*DepthCharge
	submarines   []*Submarine
	torpedoes    []*Torpedo

	oceanImage *ebiten.Image
}

func NewLevel(index int, ship *Ship) *Level {
	oceanImg := assetLoader.MustLoadImage("assets/ocean.png")
	level := &Level{
		index:      index,
		ship:       ship,
		oceanImage: ebiten.NewImageFromImage(oceanImg),
	}
	ship.UpdateForLevel(level.requestCharge)

	for idx := 0; idx < index+1; idx++ {
		level.submarines = append(level.submarines, NewSubmarine(level.requestTorpedo))
	}

	return level
}

func (l *Level) Update() error {
	l.ship.Update()

	// This is the new array for active depth charges
	var keepDepthCharges []*DepthCharge
	for _, depthCharge := range l.depthCharges {
		depthCharge.Update()
		if depthCharge.IsActive {
			keepDepthCharges = append(keepDepthCharges, depthCharge)
		}
	}
	// Setting the active array back to the depthCharges
	l.depthCharges = keepDepthCharges

	var keepSubs []*Submarine
	for _, sub := range l.submarines {
		sub.Update()
		if sub.IsActive {
			keepSubs = append(keepSubs, sub)
		}
	}
	l.submarines = keepSubs

	var keepTorpedoes []*Torpedo
	for _, torpedo := range l.torpedoes {
		torpedo.Update()
		if torpedo.IsActive {
			keepTorpedoes = append(keepTorpedoes, torpedo)
		}
	}
	l.torpedoes = keepTorpedoes

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		// torpedo hits ship
		for _, torpedo := range l.torpedoes {
			if torpedo.Rect().AlignedCollides(l.ship.Rect()) {
				torpedo.IsActive = false
				l.ship.WasHit()
			}
		}
		wg.Done()
	}()

	go func() {
		// depth charge hits a submarine
		for _, dc := range l.depthCharges {
			for _, sub := range l.submarines {
				if dc.Rect().AlignedCollides(sub.Rect()) {
					dc.IsActive = false
					sub.IsActive = false
				}
			}
		}
		wg.Done()
	}()

	go func() {
		// depth charge hits a torpedo
		for _, dc := range l.depthCharges {
			for _, torpedo := range l.torpedoes {
				if dc.Rect().AlignedCollides(torpedo.Rect()) {
					torpedo.IsActive = false
					dc.IsActive = false
				}
			}
		}
		wg.Done()
	}()

	wg.Wait()
	return nil
}

func (l *Level) Draw(screen *ebiten.Image) {
	// Background
	vector.DrawFilledRect(screen, 0, 0, SCREEN_WIDTH, SCREEN_HEIGHT, colornames.Aliceblue, true)

	l.ship.Draw(screen)
	// Iteration type of loop
	slog.Info("Depth Charges", "length", len(l.depthCharges))
	for _, depthCharge := range l.depthCharges {
		depthCharge.Draw(screen)
	}

	for _, sub := range l.submarines {
		sub.Draw(screen)
	}

	for _, torpedo := range l.torpedoes {
		torpedo.Draw(screen)
	}

	opts := &ebiten.DrawImageOptions{}
	opts.ColorScale.ScaleAlpha(0.45)
	opts.GeoM.Translate(0, WATER_SURFACE)
	screen.DrawImage(l.oceanImage, opts)
}

func (l *Level) requestCharge(isFront bool) {
	slog.Info("Fire requested")
	x := l.ship.X
	if isFront {
		x += l.ship.Width
	}
	// Depth charge fired from where the ship
	l.depthCharges = append(l.depthCharges, NewDepthCharge(x, l.ship.Y+15))
}

func (l *Level) requestTorpedo(submarine *Submarine) {
	torpedo := NewTorpedo(submarine.X, submarine.Y, 2)
	l.torpedoes = append(l.torpedoes, torpedo)
}
