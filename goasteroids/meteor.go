package goasteroids

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/viniciusscience/2d-games/assets"
)

const (
	rotationSpeedMin                    = -0.02
	rotationSpeedMax                    = 0.02
	numberOfSmallMeteorsFromLargeMeteor = 4
)

type Meteor struct {
	game          *GameScene
	position      Vector
	rotation      float64
	movement      Vector
	angle         float64
	RotationSpeed float64
	sprint        *ebiten.Image
	meteorObj     *resolv.Circle
}

func NewMeteor(baseVelocity float64, g *GameScene, index int) *Meteor {
	target := Vector{
		X: ScreenWidth / 2,
		Y: ScreenHeight / 2,
	}

	angle := rand.Float64() * 2 * math.Pi

	r := ScreenWidth/2.0 + 500
	pos := Vector{
		X: target.X + math.Cos(angle)*r,
		Y: target.Y + math.Sin(angle)*r,
	}

	velocity := baseVelocity + rand.Float64()*1.5
	direction := Vector{
		X: target.X - pos.X,
		Y: target.Y - pos.Y,
	}

	normalizeDirection := direction.Normalize()

	moviment := Vector{
		X: normalizeDirection.X * velocity,
		Y: normalizeDirection.Y * velocity,
	}

	sprinte := assets.MeteorosSprits[rand.Intn(len(assets.MeteorosSprits))]

	// create the collision object
	meteorObj := resolv.NewCircle(pos.X, pos.Y, float64(sprinte.Bounds().Dx()/2))

	m := &Meteor{
		game:          g,
		position:      pos,
		movement:      moviment,
		RotationSpeed: rotationSpeedMin + rand.Float64()*(rotationSpeedMax-rotationSpeedMin),
		sprint:        sprinte,
		angle:         angle,
		meteorObj:     meteorObj,
	}

	m.meteorObj.SetPosition(pos.X, pos.Y)
	m.meteorObj.Tags().Set(TagMeteor | TagLarge)
	m.meteorObj.SetData(&ObjectData{index: index})

	return m
}

func NewSmallMeteor(baseVelocity float64, g *GameScene, index int) *Meteor {
	target := Vector{
		X: ScreenWidth / 2,
		Y: ScreenHeight / 2,
	}

	angle := rand.Float64() * 2 * math.Pi

	r := ScreenWidth/2.0 + 500
	pos := Vector{
		X: target.X + math.Cos(angle)*r,
		Y: target.Y + math.Sin(angle)*r,
	}

	velocity := baseVelocity + rand.Float64()*1.5
	direction := Vector{
		X: target.X - pos.X,
		Y: target.Y - pos.Y,
	}

	normalizeDirection := direction.Normalize()

	moviment := Vector{
		X: normalizeDirection.X * velocity,
		Y: normalizeDirection.Y * velocity,
	}

	sprinte := assets.MeteorosSprits[rand.Intn(len(assets.MeteorosSpritsSmall))]

	// create the collision object
	meteorObj := resolv.NewCircle(pos.X, pos.Y, float64(sprinte.Bounds().Dx()/2))

	m := &Meteor{
		game:          g,
		position:      pos,
		movement:      moviment,
		RotationSpeed: rotationSpeedMin + rand.Float64()*(rotationSpeedMax-rotationSpeedMin),
		sprint:        sprinte,
		angle:         angle,
		meteorObj:     meteorObj,
	}

	m.meteorObj.SetPosition(pos.X, pos.Y)
	m.meteorObj.Tags().Set(TagMeteor | TagSmall)
	m.meteorObj.SetData(&ObjectData{index: index})

	return m
}

func (m *Meteor) Update() {
	dx := m.movement.X
	dy := m.movement.Y

	m.position.X += dx
	m.position.Y += dy
	m.rotation += m.RotationSpeed

	m.keepOnScreen()
	m.meteorObj.SetPosition(m.position.X, m.position.Y)

}

func (m *Meteor) Draw(screen *ebiten.Image) {
	bounds := m.sprint.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(m.rotation)
	op.GeoM.Translate(halfW, halfH)

	op.GeoM.Translate(m.position.X, m.position.Y)

	screen.DrawImage(m.sprint, op)
}

func (m *Meteor) keepOnScreen() {
	if m.position.X >= float64(ScreenWidth) {
		m.position.X = 0
		m.meteorObj.SetPosition(0, m.position.Y)
	}
	if m.position.X < 0 {
		m.position.X = ScreenWidth
		m.meteorObj.SetPosition(ScreenWidth, m.position.Y)
	}
	if m.position.Y >= float64(ScreenHeight) {
		m.position.Y = 0
		m.meteorObj.SetPosition(m.position.X, 0)
	}
	if m.position.Y < 0 {
		m.position.Y = ScreenHeight
		m.meteorObj.SetPosition(m.position.X, ScreenHeight)
	}
}
