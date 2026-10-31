package goasteroids

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/viniciusscience/2d-games/assets"
)

const (
	rotationPerSecond = math.Pi
	maxAcceleration   = 8.0
	shootCoolDown     = time.Millisecond * 150
	burstCoolDown     = time.Microsecond * 500
	laserSpawnOffset  = 50.0
	maxShotsPerBurst  = 3
)

var shortsFired = 0

type Player struct {
	game            *GameScene
	sprite          *ebiten.Image
	rotation        float64
	position        Vector
	playerVelocity  float64
	curAcceleration float64
	playerObj       *resolv.Circle
	shootCoolDown   *Timer
	burstCoolDown   *Timer
}

func NewPlayer(game *GameScene) *Player {
	sprite := assets.PlayerSpring

	//Center
	bounds := sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	pos := Vector{
		X: ScreenWidth/2 - halfW,
		Y: ScreenHeight/2 - halfH,
	}

	//create a resolv object
	playerObj := resolv.NewCircle(pos.X, pos.Y, float64(sprite.Bounds().Dx()/2))

	p := &Player{
		sprite:        sprite,
		game:          game,
		position:      pos,
		playerObj:     playerObj,
		shootCoolDown: NewTimer(shootCoolDown),
		burstCoolDown: NewTimer(burstCoolDown),
	}

	p.playerObj.SetPosition(pos.X, pos.Y)
	p.playerObj.Tags().Set(TagPlayer)

	return p
}

func (p *Player) Draw(screen *ebiten.Image) {
	bounds := p.sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(p.rotation)
	op.GeoM.Translate(halfW, halfH)

	op.GeoM.Translate(p.position.X, p.position.Y)

	screen.DrawImage(p.sprite, op)
}

func (p *Player) Update() {
	speed := rotationPerSecond / float64(ebiten.TPS())

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.rotation -= speed
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.rotation += speed
	}

	p.accelerate()

	p.playerObj.SetPosition(p.position.X, p.position.Y)

	p.burstCoolDown.Update()

	p.shootCoolDown.Update()

	p.fireLasers()
}

func (p *Player) fireLasers() {
	if p.burstCoolDown.IsReady() {
		if p.shootCoolDown.IsReady() && ebiten.IsKeyPressed(ebiten.KeySpace) {
			p.shootCoolDown.Reset()
			shortsFired++
			if shortsFired <= maxShotsPerBurst {
				bounds := p.sprite.Bounds()
				halfW := float64(bounds.Dx()) / 2
				halfH := float64(bounds.Dy()) / 2

				spawnPos := Vector{
					p.position.X + halfW + math.Sin(p.rotation)*laserSpawnOffset,
					p.position.Y + halfH + math.Cos(p.rotation)*-laserSpawnOffset,
				}

				p.game.laserCount++
				laser := NewLaser(spawnPos, p.rotation, p.game.laserCount, p.game)
				p.game.lasers[p.game.laserCount] = laser
				p.game.space.Add(laser.laserObj)
			} else {
				p.burstCoolDown.Reset()
				shortsFired = 0
			}
		}
	}
}

func (p *Player) accelerate() {
	p.keepOnScreen()
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		if p.curAcceleration < maxAcceleration {
			p.curAcceleration = p.playerVelocity + 4
		}

		if p.curAcceleration >= 8 {
			p.curAcceleration = 8
		}

		p.playerVelocity = p.curAcceleration

		dx := math.Sin(p.rotation) * p.curAcceleration
		dy := math.Cos(p.rotation) * -p.curAcceleration

		p.position.X += dx
		p.position.Y += dy
	}
}

func (p *Player) keepOnScreen() {
	if p.position.X >= float64(ScreenWidth) {
		p.position.X = 0
		p.playerObj.SetPosition(0, p.position.Y)
	}
	if p.position.X < 0 {
		p.position.X = ScreenWidth
		p.playerObj.SetPosition(ScreenWidth, p.position.Y)

	}
	if p.position.Y >= float64(ScreenHeight) {
		p.position.Y = 0
		p.playerObj.SetPosition(p.position.X, 0)

	}
	if p.position.Y < 0 {
		p.position.Y = ScreenHeight
		p.playerObj.SetPosition(p.position.X, ScreenHeight)

	}
}
