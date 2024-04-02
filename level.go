package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"log/slog"
)

type Level struct {
	index        int
	ship         *Ship
	depthCharges []*DepthCharge
	submarines   []*Submarine

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

	for _, sub := range l.submarines {
		sub.Update()
	}

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

	opts := &ebiten.DrawImageOptions{}
	opts.ColorScale.ScaleAlpha(0.45)
	opts.GeoM.Translate(0, 140)
	screen.DrawImage(l.oceanImage, opts)
}

func (l *Level) requestCharge() {
	slog.Info("Fire requested")
	// Depth charge fired from where the ship
	l.depthCharges = append(l.depthCharges, NewDepthCharge(l.ship.X, l.ship.Y))
}

func (l *Level) requestTorpedo(submarine *Submarine) {

}
