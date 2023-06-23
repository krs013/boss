package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenWidth  = 1280
	ScreenHeight = 800 // Or 720 for 16:9 aspect ratio
)

type Trigger struct {
	AABB
	Fn func()
}

type Room struct {
	Width, Height float64
	Obstacles     []AABB
	Triggers      []Trigger

	Floor  color.Color
	Wall   color.Color
	Button color.Color
}

func (r Room) Update(g *Game) {
	for _, t := range r.Triggers {
		if g.Boss.Collide(t.AABB) {
			if t.Fn != nil {
				t.Fn()
			}
		}
	}
}

func (r Room) Draw(dst *ebiten.Image) {
	ebitenutil.DrawRect(dst, 0, 0, r.Width, r.Height, r.Floor)
	for _, t := range r.Triggers {
		ebitenutil.DrawRect(dst, t.X, t.Y, t.W, t.H, r.Button)
	}
	for _, o := range r.Obstacles {
		ebitenutil.DrawRect(dst, o.X, o.Y, o.W, o.H, r.Wall)
	}
}

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
