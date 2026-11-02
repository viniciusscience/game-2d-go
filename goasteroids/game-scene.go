package goasteroids

import (
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/solarlune/resolv"
	"github.com/viniciusscience/2d-games/assets"
)

const (
	baseMeteorVelocity   = 0.25
	meteorSpawnTime      = 100 * time.Millisecond
	meteorSpeedUpAmount  = 0.1
	meteorSpeedUpTime    = 1000 * time.Millisecond
	cleanUpExplosionTime = 200 * time.Millisecond
	baseBeatWaitTime     = 1600
)

var audioContext = audio.NewContext(assets.SampleRate)

type GameScene struct {
	player               *Player
	baseVelocity         float64
	meteorCount          int
	meteorSpawnTimer     *Timer
	meteoro              map[int]*Meteor
	meteorsForLevel      int
	velocityTimer        *Timer
	space                *resolv.Space
	lasers               map[int]*Laser
	laserCount           int
	score                int
	explosionSmallSprite *ebiten.Image
	explosionSprite      *ebiten.Image
	explosionFrames      []*ebiten.Image
	cleanUpTimer         *Timer
	playerIsDead         bool
	audioContext         *audio.Context
	thrustPlayer         *audio.Player
	exhaust              *Exhaust
	laserOnePlayer       *audio.Player
	laserTwoPlayer       *audio.Player
	laserThreePlayer     *audio.Player
	explosionPlayer      *audio.Player
	beatOnePlayer        *audio.Player
	beatTwoPlayer        *audio.Player
	beatTimer            *Timer
	beatWaitTime         int
	playBeatOne          bool
}

func NewGameScene() *GameScene {
	g := &GameScene{
		meteorSpawnTimer:     NewTimer(meteorSpawnTime),
		baseVelocity:         baseMeteorVelocity,
		velocityTimer:        NewTimer(meteorSpeedUpTime),
		meteoro:              make(map[int]*Meteor),
		meteorCount:          0,
		meteorsForLevel:      2,
		space:                resolv.NewSpace(ScreenWidth, ScreenHeight, 16, 16),
		lasers:               make(map[int]*Laser),
		laserCount:           0,
		explosionSprite:      assets.ExplosionSprite,
		explosionSmallSprite: assets.ExplosionSmallSprite,
		cleanUpTimer:         NewTimer(cleanUpExplosionTime),
		beatTimer:            NewTimer(2 * time.Second),
		beatWaitTime:         baseBeatWaitTime,
	}
	g.player = NewPlayer(g)
	g.space.Add(g.player.playerObj)

	g.explosionFrames = assets.Explosion

	// Load Audio.
	g.audioContext = audioContext
	thrustPlayer, _ := g.audioContext.NewPlayer(assets.ThrustSound)
	g.thrustPlayer = thrustPlayer

	laserOnePlayer, _ := g.audioContext.NewPlayer(assets.LaserOneSound)
	g.laserOnePlayer = laserOnePlayer

	laserTwoPlayer, _ := g.audioContext.NewPlayer(assets.LaserTwoSound)
	g.laserTwoPlayer = laserTwoPlayer

	laserThreePlayer, _ := g.audioContext.NewPlayer(assets.LaserThreeSound)
	g.laserThreePlayer = laserThreePlayer

	explosionPlayer, _ := g.audioContext.NewPlayer(assets.ExplosionSound)
	g.explosionPlayer = explosionPlayer

	beatOnePlayer, _ := g.audioContext.NewPlayer(assets.BeatOneSound)
	beatTwoPlayer, _ := g.audioContext.NewPlayer(assets.BeatTwoSound)
	g.beatOnePlayer = beatOnePlayer
	g.beatTwoPlayer = beatTwoPlayer

	return g
}

func (g *GameScene) Update(state *State) error {
	g.player.Update()
	g.updateExhaust()
	g.isPlayerDying()
	g.isPlayerDead(state)

	g.spawnMeteors()
	g.speedUpMeteors()

	for _, m := range g.meteoro {
		m.Update()
	}

	for _, l := range g.lasers {
		l.Update()
	}

	g.speedUpMeteors()
	g.isPlayerCollidingWithMeteor()
	g.isMeteorHitByPlayerLaser()
	g.cleanUpMeteorsAndAliens()

	return nil
}

func (g *GameScene) updateExhaust() {
	if g.exhaust != nil {
		g.exhaust.Update()
	}
}

func (g *GameScene) beatSound() {
	g.beatTimer.Update()
	if g.beatTimer.IsReady() {
		if g.playBeatOne {
			_ = g.beatOnePlayer.Rewind()
			g.beatOnePlayer.Play()
			g.beatTimer.Reset()
		} else {
			_ = g.beatTwoPlayer.Rewind()
			g.beatTwoPlayer.Play()
			g.beatTimer.Reset()
		}
		g.playBeatOne = !g.playBeatOne

		// speed up the timer
		if g.beatWaitTime > 400 {
			g.beatWaitTime = g.beatWaitTime - 25
			g.beatTimer = NewTimer(time.Millisecond * time.Duration(g.beatWaitTime))

		}

	}
}

func (g *GameScene) isPlayerDying() {
	if g.player.isDying {
		g.player.dyingTimer.Update()
		if g.player.dyingTimer.IsReady() {
			g.player.dyingTimer.Reset()
			g.player.DyingCounter++
			if g.player.DyingCounter == 12 {
				g.player.isDying = false
				g.player.isDead = true
			} else if g.player.DyingCounter < 12 {
				g.player.sprite = g.explosionFrames[g.player.DyingCounter]
			} else {
				// Do nothing.
			}
		}
	}
}

func (g *GameScene) isPlayerDead(state *State) {
	if g.playerIsDead {
		g.player.livesRemaining--
		if g.player.livesRemaining == 0 {
			g.Reset()
			state.SceneManager.GoToScene(g)
		}
	}

}

func (g *GameScene) Draw(screen *ebiten.Image) {
	g.player.Draw(screen)

	if g.exhaust != nil {
		g.exhaust.Draw(screen)
	}

	for _, m := range g.meteoro {
		m.Draw(screen)
	}

	for _, l := range g.lasers {
		l.Draw(screen)
	}
}

func (g *GameScene) Reset() {
	g.player = NewPlayer(g)
	g.meteoro = make(map[int]*Meteor)
	g.meteorCount = 0
	g.lasers = make(map[int]*Laser)
	g.laserCount = 0
	g.score = 0
	g.meteorSpawnTimer.Reset()
	g.baseVelocity = baseMeteorVelocity
	g.velocityTimer.Reset()
	g.playerIsDead = false
	g.exhaust = nil
	g.space.RemoveAll()
	g.space.Add(g.player.playerObj)
}

func (g *GameScene) isMeteorHitByPlayerLaser() {
	for _, m := range g.meteoro {
		for _, l := range g.lasers {
			if m.meteorObj.IsIntersecting(l.laserObj) {
				if m.meteorObj.Tags().Has(TagSmall) {
					//Small meteor
					m.sprint = g.explosionSmallSprite
					g.score++

					if !g.explosionPlayer.IsPlaying() {
						_ = g.explosionPlayer.Rewind()
						g.explosionPlayer.Play()
					}
				} else {
					oldPos := m.position
					m.sprint = g.explosionSprite

					g.score++

					if !g.explosionPlayer.IsPlaying() {
						_ = g.explosionPlayer.Rewind()
						g.explosionPlayer.Play()
					}
					numToSpawn := rand.Intn(numberOfSmallMeteorsFromLargeMeteor)

					for i := 0; i < numToSpawn; i++ {
						meteor := NewSmallMeteor(baseMeteorVelocity, g, len(m.game.meteoro)-1)
						meteor.position = Vector{oldPos.X + float64(rand.Intn(100-50)) + 50, oldPos.Y + float64(rand.Intn(100-50))}
						meteor.meteorObj.SetPosition(meteor.position.X, meteor.position.Y)
						g.space.Add(meteor.meteorObj)
						g.meteorCount++
						g.meteoro[m.game.meteorCount] = meteor
					}
				}
			}
		}
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
			if !g.player.isShielded {
				m.game.player.isDying = true

				if !g.explosionPlayer.IsPlaying() {
					_ = g.explosionPlayer.Rewind()
					g.explosionPlayer.Play()
				}

				break
			} else {
				// bounce the meteor
			}

		}
	}
}

func (g *GameScene) cleanUpMeteorsAndAliens() {
	g.cleanUpTimer.Update()
	if g.cleanUpTimer.IsReady() {
		for i, m := range g.meteoro {
			if m.sprint == g.explosionSprite || m.sprint == g.explosionSmallSprite {
				delete(g.meteoro, i)
				g.space.Remove(m.meteorObj)
			}
		}
		g.cleanUpTimer.Reset()
	}

}
