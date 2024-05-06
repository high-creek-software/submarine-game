package submarinegame

import (
	"github.com/hajimehoshi/ebiten/v2"
	"gitlab.com/high-creek-software/go2d/components/display"
)

type ParticleType int

const (
	ParticleSubExplosion ParticleType = iota
	ParticleDepthChargeTorpedoExplosion
	ParticleShipHit
	ParticleShipExplosion
)

type Particle struct {
	*Entity

	animation    *display.AnimateComponent
	particleType ParticleType
}

func (p *Particle) Scale() (float64, float64) {
	return 1, 1
}

func (p *Particle) IsHorizontalFlipped() bool {
	return false
}

func (p *Particle) IsVerticalFlipped() bool {
	return false
}

func (p *Particle) DrawOffset() (float64, float64) {
	return 0, 0
}

func NewParticle(particleType ParticleType, x, y float64) *Particle {
	p := &Particle{
		particleType: particleType,
	}

	var path string
	var count int
	var width float64
	var height float64
	switch particleType {
	case ParticleSubExplosion:
		path = "assets/particles/subexplosion"
		count = 7
		width = 128
		height = 45
	case ParticleShipHit:
		path = "assets/particles/hit"
		count = 4
		width = 32
		height = 32
	case ParticleDepthChargeTorpedoExplosion:
		path = "assets/particles/depthchargetorpedo"
		count = 8
		width = 32
		height = 32
	default:
	}

	p.Entity = NewEntity(x, y, width, height)

	frames := assetLoader.LoadConcurrentDirectory(path, "png", count)
	p.animation = display.NewAnimateComponent(p, 6, frames).CenterImage().DisableLoop(p.AnimationComplete)

	return p
}

func (p *Particle) Update() error {
	p.animation.Update()
	return nil
}

func (p *Particle) Draw(screen *ebiten.Image) {
	p.animation.Draw(screen, 0, 0)
}

func (p *Particle) AnimationComplete() {
	p.IsActive = false
}
