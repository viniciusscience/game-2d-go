package goasteroids

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/viniciusscience/2d-games/assets"
)

type LifeIndicator struct {
	posisiton Vector
	rotation  float64
	sprite    *ebiten.Image
}

func NewLifeIndicator(pos Vector) *LifeIndicator {
	sprite := assets.LifeIndicator
	return &LifeIndicator{
		posisiton: pos,
		sprite:    sprite,
	}
}

func (l *LifeIndicator) Update() {
}

func (l *LifeIndicator) Draw(screen *ebiten.Image) {
	bounds := l.sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op := &colorm.DrawImageOptions{}
	op.GeoM.Translate(halfW, halfH)
	cm := colorm.ColorM{}
	cm.Scale(1.0, 1.0, 1.0, 0.2)

	op.GeoM.Translate(l.posisiton.X, l.posisiton.Y)

	colorm.DrawImage(screen, l.sprite, cm, op)
}
