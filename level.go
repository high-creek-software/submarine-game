package submarinegame

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/colornames"
	"log/slog"
	"sync"
	"time"
)

const (
	MaxSubsOnScreen = 5
)

/*type Attacker interface {
	Update() error
	Draw(screen *ebiten.Image)
}*/

type LevelInteractor interface {
	GameOver()
	PlaySound(soundType SoundType)
	PlaySoundDelay(soundType SoundType, delay time.Duration)
	LevelComplete()
}

type Level struct {
	index             int
	destroyedSubCount int
	ship              *Ship
	depthCharges      []*DepthCharge
	submarines        []*Submarine
	torpedoes         []*Torpedo
	interactor        LevelInteractor
	particles         []*Particle

	oceanImage *ebiten.Image
}

func NewLevel(interactor LevelInteractor, index int, ship *Ship) *Level {
	oceanImg := assetLoader.MustLoadImage("assets/big_ocean.png")
	level := &Level{
		interactor: interactor,
		index:      index,
		ship:       ship,
		oceanImage: ebiten.NewImageFromImage(oceanImg),
	}
	ship.UpdateForLevel(level.requestCharge)

	startSubs := index
	if index > MaxSubsOnScreen {
		startSubs = MaxSubsOnScreen
	}
	for idx := 0; idx < startSubs; idx++ {
		level.submarines = append(level.submarines, NewSubmarine(level.requestTorpedo))
	}

	return level
}

func (l *Level) Update() error {
	start := time.Now()
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
		} else {
			l.destroyedSubCount += 1

			if l.index > MaxSubsOnScreen {
				if l.destroyedSubCount+MaxSubsOnScreen <= l.index {
					keepSubs = append(keepSubs, SpawnSub(l.requestTorpedo))
					l.interactor.PlaySoundDelay(SOUNDS_PING, 1500*time.Millisecond)
				}
			}
		}
	}
	l.submarines = keepSubs

	if len(l.submarines) == 0 && len(l.particles) == 0 {
		// Finish the level
		l.interactor.LevelComplete()
		return nil
	}

	var keepTorpedoes []*Torpedo
	for _, torpedo := range l.torpedoes {
		torpedo.Update()
		if torpedo.IsActive {
			keepTorpedoes = append(keepTorpedoes, torpedo)
		}
	}
	l.torpedoes = keepTorpedoes

	var keepParticles []*Particle
	for _, particle := range l.particles {
		particle.Update()
		if particle.IsActive {
			keepParticles = append(keepParticles, particle)
		}
	}
	l.particles = keepParticles

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		// torpedo hits ship
		for _, torpedo := range l.torpedoes {
			if torpedo.Rect().AlignedCollides(l.ship.Rect()) {
				torpedo.IsActive = false
				l.ship.WasHit()
				l.particles = append(l.particles, NewParticle(ParticleShipHit, torpedo.X, torpedo.Y-torpedo.Height/2))
				l.interactor.PlaySound(SOUNDS_HIT)
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
					l.particles = append(l.particles, NewParticle(ParticleSubExplosion, sub.X, sub.Y))
					l.interactor.PlaySound(SOUNDS_EXP_UNDERWATER)
					l.ship.incrementScore(250)
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
					l.ship.incrementScore(125)
					l.particles = append(l.particles, NewParticle(ParticleDepthChargeTorpedoExplosion, torpedo.X, torpedo.Y-torpedo.Height/2))
				}
			}
		}
		wg.Done()
	}()

	wg.Wait()

	// Are we still alive
	if !l.ship.IsActive {
		l.interactor.GameOver()
	}

	slog.Info("Update Duration", "duration", time.Now().Sub(start))

	return nil
}

func (l *Level) Draw(screen *ebiten.Image) {
	// Background
	screen.Fill(colornames.Aliceblue)

	// 006680ff
	//clr := color.NRGBA{52, 91, 235, 128}
	//vector.DrawFilledRect(screen, 0, float32(WATER_SURFACE), float32(SCREEN_WIDTH), float32(SCREEN_HEIGHT-WATER_SURFACE), clr, true)

	l.ship.Draw(screen)
	// Iteration type of loop
	for _, depthCharge := range l.depthCharges {
		depthCharge.Draw(screen)
	}

	for _, sub := range l.submarines {
		sub.Draw(screen)
	}

	for _, torpedo := range l.torpedoes {
		torpedo.Draw(screen)
	}

	for _, particle := range l.particles {
		particle.Draw(screen)
	}

	// Drawing ocean
	opts := &ebiten.DrawImageOptions{}
	opts.ColorScale.ScaleAlpha(0.45)
	opts.GeoM.Translate(0, WATER_SURFACE)
	screen.DrawImage(l.oceanImage, opts)

	// Drawing level screen
	levelStr := gamePrinter.Sprintf("Level: %d", l.index)
	face := &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   18,
	}
	width, _ := text.Measure(levelStr, face, 1)
	txtOpts := &text.DrawOptions{}
	txtOpts.GeoM.Translate(SCREEN_WIDTH-width-50, 10)
	txtOpts.ColorScale.ScaleWithColor(colornames.Black)
	//text.Measure()
	text.Draw(
		screen,
		levelStr,
		face,
		txtOpts,
	)

	if l.ship.health <= 0 {

	}
}

func (l *Level) requestCharge(isFront bool) {
	slog.Info("Fire requested")
	x := l.ship.X
	if isFront {
		x += l.ship.Width
	}
	// Depth charge fired from where the ship
	l.depthCharges = append(l.depthCharges, NewDepthCharge(x, l.ship.Y+15))
	l.interactor.PlaySound(SOUNDS_SPLASH)
}

func (l *Level) requestTorpedo(submarine *Submarine) {
	torpedo := NewTorpedo(submarine.X, submarine.Y, 2)
	l.torpedoes = append(l.torpedoes, torpedo)
	l.interactor.PlaySound(SOUNDS_MISSLE)
}
