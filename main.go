package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/viniciusscience/2d-games/goasteroids"
)

func main() {

	ebiten.SetWindowTitle("Go Asteroids")
	ebiten.SetWindowSize(goasteroids.ScreenWidth, goasteroids.ScreenHeight)

	err := ebiten.RunGame(&goasteroids.Game{})
	if err != nil {
		panic(err)
	}

}
