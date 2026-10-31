package goasteroids

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

const (
	baseMeteorVelocity  = 0.25
	meteorSpawnTime     = 100 * time.Millisecond
	meteorSpeedUpAmount = 0.1
	meteorSpeedUpTime   = 1000 * time.Millisecond
)

type GameScene struct {
	player           *Player
	baseVelocity     float64
	meteorCount      int
	meteorSpawnTimer *Timer
	meteoro          map[int]*Meteor
	meteorsForLevel  int
	velocityTimer    *Timer
	space            *resolv.Space
}

func NewGameScene() *GameScene {
	g := &GameScene{
		meteorSpawnTimer: NewTimer(meteorSpawnTime),
		baseVelocity:     baseMeteorVelocity,
		velocityTimer:    NewTimer(meteorSpeedUpTime),
		meteoro:          make(map[int]*Meteor),
		meteorCount:      0,
		meteorsForLevel:  2,
		space:            resolv.NewSpace(ScreenWidth, ScreenHeight, 16, 16),
	}
	g.player = NewPlayer(g)
	g.space.Add(g.player.playerObj)

	return g
}

func (g *GameScene) Update(state *State) error {
	g.player.Update()
	g.spawnMeteors()
	g.speedUpMeteors()

	for _, m := range g.meteoro {
		m.Update()
	}
	g.speedUpMeteors()
	g.isPlayerCollidingWithMeteor()

	return nil
}

func (g *GameScene) Draw(screen *ebiten.Image) {
	g.player.Draw(screen)

	for _, m := range g.meteoro {
		m.Draw(screen)
	}
}

func (g *GameScene) spawnMeteors() {
	g.meteorSpawnTimer.Update()

	if g.meteorSpawnTimer.IsReady() {
		g.meteorSpawnTimer.Reset()
		if len(g.meteoro) < g.meteorsForLevel && g.meteorCount < g.meteorsForLevel {
			m := NewMeteor(g.baseVelocity, g, len(g.meteoro)-1)
			g.space.Add(m.meteorObj)
			g.meteorCount++
			g.meteoro[g.meteorCount] = m
		}
	}
}

func (g *GameScene) speedUpMeteors() {
	g.velocityTimer.Update()
	if g.velocityTimer.IsReady() {
		g.velocityTimer.Reset()
		g.baseVelocity += meteorSpeedUpAmount
	}
}

func (g *GameScene) isPlayerCollidingWithMeteor() {
	for _, m := range g.meteoro {
		if m.meteorObj.IsIntersecting(g.player.playerObj) {
			data := m.meteorObj.Data().(*ObjectData)
			fmt.Println("Playr collided with meteor", data.index)
		}
	}
}
