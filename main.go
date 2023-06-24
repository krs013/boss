package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 1280
	ScreenHeight = 800 // Or 720 for 16:9 aspect ratio
)

type Game struct {
	Room *Room
	Boss *Boss
	Hero *Hero
}

func NewGame() *Game {
	room := &Room{
		Width:  ScreenWidth,
		Height: ScreenHeight,
		Obstacles: []AABB{
			{10, 10, 300, 40},
			{10, 40, 40, 200},
		},
		Floor:  Color7,
		Wall:   Color6,
		Button: Color3,
	}
	idle, left, right := BossAnimations()
	boss := &Boss{
		Mob: Mob{
			AABB:           AABB{50, 50, 128, 128},
			IdleAnimation:  idle,
			LeftAnimation:  left,
			RightAnimation: right,
			Color:          Color0,
		},
	}
	hero := &Hero{
		Mob: Mob{
			AABB:  AABB{250, 400, 64, 64},
			Color: Color1,
		},
	}
	room.Triggers = append(room.Triggers, Trigger{
		AABB: AABB{500, 500, 64, 64},
		Fn: func() {
			hero.X = ScreenWidth / 2
			hero.Y = ScreenHeight / 2
		},
	})
	return &Game{room, boss, hero}
}

func (g *Game) Update() error {
	// NB: Order matters here! Only the Boss resolves hero-boss pushing
	// interaction, so boss must go after hero has done naive moves.
	g.Hero.Update(g)
	g.Boss.Update(g)
	// Once boss has moved, see if any triggers were tripped.
	g.Room.Update(g)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// NB: Order matters here! Later stuff draws over earlier stuff.
	g.Room.Draw(screen)
	g.Hero.Draw(screen)
	g.Boss.Draw(screen)
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
