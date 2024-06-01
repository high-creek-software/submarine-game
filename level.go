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
	MaxSubsOnScreen = 8
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
	keepSubmarines    []*Submarine
	torpedoes         []*Torpedo
	keepTorpedoes     []*Torpedo
	interactor        LevelInteractor
	particles         []*Particle
	keepParticles     []*Particle

	oceanImage *ebiten.Image
}

func NewLevel(interactor LevelInteractor, index int, ship *Ship) *Level {
	oceanImg := assetLoader.MustLoadImage("assets/ocean.png")
	level := &Level{
		interactor:     interactor,
		index:          index,
		ship:           ship,
		oceanImage:     ebiten.NewImageFromImage(oceanImg),
		submarines:     make([]*Submarine, 0, 50),
		keepSubmarines: make([]*Submarine, 0, 50),
		torpedoes:      make([]*Torpedo, 0, 50),
		keepTorpedoes:  make([]*Torpedo, 0, 50),
		particles:      make([]*Particle, 0, 50),
		keepParticles:  make([]*Particle, 0, 50),
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
	//start := time.Now()
	l.ship.Update()

	// This is the new array for active depth charges
	var keepDepthCharges []*DepthCharge
	for _, depthCharge := range l.depthCharges {
		depthCharge.Update()
		if depthCharge.IsActive {
			keepDepthCharges = append(keepDepthCharges, depthCharge)
		} else {
			depthChargePool.Return(depthCharge)
		}
	}
	// Setting the active array back to the depthCharges
	l.depthCharges = keepDepthCharges

	l.keepSubmarines = l.keepSubmarines[:0]
	for _, sub := range l.submarines {
		sub.Update()
		if sub.IsActive {
			l.keepSubmarines = append(l.keepSubmarines, sub)
		} else {
			l.destroyedSubCount += 1

			if l.index > MaxSubsOnScreen {
				if l.destroyedSubCount+MaxSubsOnScreen <= l.index {
					l.keepSubmarines = append(l.keepSubmarines, SpawnSub(l.requestTorpedo))
					l.interactor.PlaySoundDelay(SOUNDS_PING, 1500*time.Millisecond)
				}
			}
		}
	}
	l.submarines = l.keepSubmarines

	if len(l.submarines) == 0 && len(l.particles) == 0 {
		// Finish the level
		l.interactor.LevelComplete()
		return nil
	}

	l.keepTorpedoes = l.keepTorpedoes[:0]
	for _, torpedo := range l.torpedoes {
		torpedo.Update()
		if torpedo.IsActive {
			l.keepTorpedoes = append(l.keepTorpedoes, torpedo)
		} else {
			torpedoPool.Return(torpedo)
		}
	}
	l.torpedoes = l.keepTorpedoes

	l.keepParticles = l.keepParticles[:0]
	for _, particle := range l.particles {
		particle.Update()
		if particle.IsActive {
			l.keepParticles = append(l.keepParticles, particle)
		}
	}
	l.particles = l.keepParticles

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
				if isMobile {
					if l.ship.health > 0 {
						ebiten.Vibrate(&ebiten.VibrateOptions{Duration: 50 * time.Millisecond, Magnitude: 0.45})
					} else {
						ebiten.Vibrate(&ebiten.VibrateOptions{Duration: 500 * time.Millisecond, Magnitude: 0.75})
					}
				}
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

	//slog.Info("Update Duration", "duration", time.Now().Sub(start))

	return nil
}

func (l *Level) Draw(screen *ebiten.Image) {
	// Background
	screen.Fill(colornames.Aliceblue)

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
	depthCharge := depthChargePool.Borrow()
	depthCharge.X = x
	depthCharge.Y = l.ship.Y + 15
	depthCharge.IsActive = true
	// Depth charge fired from where the ship
	l.depthCharges = append(l.depthCharges, depthCharge)
	l.interactor.PlaySound(SOUNDS_SPLASH)
}

func (l *Level) requestTorpedo(submarine *Submarine) {
	//torpedo := NewTorpedo(submarine.X, submarine.Y, 2)
	torpedo := torpedoPool.Borrow()
	torpedo.X = submarine.X
	torpedo.Y = submarine.Y
	torpedo.speed = 2
	torpedo.IsActive = true
	l.torpedoes = append(l.torpedoes, torpedo)
	l.interactor.PlaySound(SOUNDS_MISSLE)
}
