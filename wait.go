package main

import "github.com/hajimehoshi/ebiten/v2"

type Wait struct {
	Room *Room
	Boss *Boss
	Hero *Hero
}

func NewWait() *Wait {
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
	return &Wait{
		Room: room,
		Boss: boss,
		Hero: hero,
	}
}

func (w *Wait) Update(g *Game) error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.Scene = NewMenu()
	}
	// NB: Order matters here! Only the Boss resolves hero-boss pushing
	// interaction, so boss must go after hero has done naive moves.
	w.Hero.Update(w)
	w.Boss.Update(w)
	// Once boss has moved, see if any triggers were tripped.
	w.Room.Update(w)

	return nil
}

func (w *Wait) Draw(screen *ebiten.Image) {
	// NB: Order matters here! Later stuff draws over earlier stuff.
	w.Room.Draw(screen)
	w.Hero.Draw(screen)
	w.Boss.Draw(screen)
}
