package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 1280
	ScreenHeight = 800 // Or 720 for 16:9 aspect ratio
)

type Scene interface {
	Update(g *Game) error
	Draw(*ebiten.Image)
}

type Game struct {
	CurrenScene Scene
}

func NewGame() *Game {
	return nil
	//return &Game{room, boss, hero}
}

func (g *Game) Update() error {
	return g.CurrenScene.Update(g)
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.CurrenScene.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Don't get fancy here - just let Ebitengine handle scaling stuff for us.
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("You're the Boss!")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
