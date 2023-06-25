package main

import "github.com/hajimehoshi/ebiten/v2"

type RoomScene struct {
	Room *Room
	Boss *Boss
	Hero *Hero
}

type WaitScene struct {
	RoomScene
}

func NewWaitScene(g *Game) *WaitScene {
	room := &Room{
		Width:  ScreenWidth,
		Height: ScreenHeight,
		Obstacles: []AABB{
			{0, 0, ScreenWidth/2 - 100, 40},
			{ScreenWidth/2 + 100, 0, ScreenWidth/2 - 150, 40},
			{0, ScreenHeight - 40, ScreenWidth, 40},
			{0, 40, 40, ScreenHeight - 80},
			{ScreenWidth - 40, 40, 40, ScreenHeight - 80},
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
	// Trigger to move to control room.
	room.Triggers = append(room.Triggers, Trigger{
		AABB: AABB{ScreenWidth/2 - 100, 0, 200, 1},
		Fn: func() {
			g.Scene = NewCtrlScene(g)
		},
	})
	// Trigger to open the control room door.
	room.Triggers = append(room.Triggers, Trigger{
		AABB: AABB{ScreenWidth/2 - 150, ScreenHeight - 115, 300, 75},
		Fn:   nil,
	})
	// Trigger Lt. telling Boss to stay put.
	room.Triggers = append(room.Triggers, Trigger{
		AABB: AABB{ScreenWidth/2 - 150, 40, 300, 150},
		Fn:   nil,
	})
	return &WaitScene{
		RoomScene{
			Room: room,
			Boss: boss,
			Hero: hero,
		},
	}
}

func (w *WaitScene) Update(g *Game) {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.Scene = NewMenu()
	}
	// NB: Order matters here! Only the Boss resolves hero-boss pushing
	// interaction, so boss must go after hero has done naive moves.
	w.Hero.Update(w.RoomScene)
	w.Boss.Update(w.RoomScene)
	// Once boss has moved, see if any triggers were tripped.
	w.Room.Update(w.RoomScene)
}

func (w *WaitScene) Draw(screen *ebiten.Image) {
	// NB: Order matters here! Later stuff draws over earlier stuff.
	w.Room.Draw(screen)
	w.Hero.Draw(screen)
	w.Boss.Draw(screen)
}
