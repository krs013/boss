package main

import "github.com/hajimehoshi/ebiten/v2"

// WaitScene has the boss waiting in the throne room for something to happen.
type WaitScene struct {
	SceneData
}

// NewWaitScene creates the boss in their throne/waiting room. How boring.
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
		SceneData{
			Room: room,
			Boss: boss,
			Hero: hero,
		},
	}
}

// Update updates the boss, and runs any transition triggers in the scene.
func (w *WaitScene) Update(g *Game) {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.Scene = NewMainMenu()
	}
	// NB: Order matters here! Only the Boss resolves hero-boss pushing
	// interaction, so boss must go after hero has done naive moves.
	w.Hero.Update(w.SceneData)
	w.Boss.Update(w.SceneData)
	// Once boss has moved, see if any triggers were tripped.
	w.Room.Update(w.SceneData)
}

// Draw draws the boss in the waiting room.
func (w *WaitScene) Draw(screen *ebiten.Image) {
	// NB: Order matters here! Later stuff draws over earlier stuff.
	w.Room.Draw(screen)
	w.Hero.Draw(screen)
	w.Boss.Draw(screen)
}
