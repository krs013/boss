package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// The logical size of the screen. In theory changing these (within reason)
// should just work, but I can't say I've been too careful to always use these
// constants when computing various object sizes, so be careful.
const (
	ScreenWidth  = 1280
	ScreenHeight = 800 // Or 720 for 16:9 aspect ratio
)

// Scene represents a particular screen in the game.
type Scene interface {
	Update(g *Game)
	Draw(*ebiten.Image)
}

// SceneData contains all the room and mob info needed for a Scene. Depending
// on the Scene, some of the mobs may be nil, indicating they aren't in the
// scene. Consequently, mobs *must* check for nil in their updates!
type SceneData struct {
	Room *Room
	Boss *Boss
	Hero *Hero
	Mate *Mate
}

// Game is contains all game info, including the current Scene.
type Game struct {
	Scene
}

// NewGame creates a brand new Game, starting with the main menu.
func NewGame() *Game {
	return &Game{
		Scene: NewMainMenu(),
	}
}

// Update simply calls the Update of the current Scene.
func (g *Game) Update() error {
	g.Scene.Update(g)
	return nil
}

// Draw simply calls the Draw of the current Scene.
func (g *Game) Draw(screen *ebiten.Image) {
	g.Scene.Draw(screen)
}

// Layout gets the ScreenWidth and ScreenHeight.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Don't get fancy here - just let Ebitengine handle scaling stuff for us.
	return ScreenWidth, ScreenHeight
}

// main is were all the magic happens!
func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("You're the Boss!")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
