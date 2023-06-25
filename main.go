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
	Update(g *Game)
	Draw(*ebiten.Image)
}

type Game struct {
	Scene
}

func NewGame() *Game {
	return &Game{
		Scene: NewMenu(),
	}
}

func (g *Game) Update() error {
	g.Scene.Update(g)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Scene.Draw(screen)
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
