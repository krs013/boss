package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Button struct {
	AABB
	Pressed bool
}

// CtrlScene has the Boss in the control room with the level design buttons.
type CtrlScene struct {
	SceneData
	Buttons []Button
}

// NewCtrlScene creates a new control room scheme with boss at the door.
func NewCtrlScene(g *Game) *CtrlScene {
	room := &Room{
		Width:  ScreenWidth,
		Height: ScreenHeight,
		Obstacles: []AABB{
			{0, 0, ScreenWidth, 40},
			{0, ScreenHeight - 40, ScreenWidth/2 - 100, 40},
			{ScreenWidth/2 + 100, ScreenHeight - 40, ScreenWidth/2 - 100, 40},
			{0, 40, 40, ScreenHeight - 80},
			{ScreenWidth - 40, 40, 40, ScreenHeight - 80},
		},
		Floor:  Color7,
		Wall:   Color6,
		Button: Color3,
	}
	idle, left, right := BossAnimations()
	bossX, bossY := ScreenWidth/2-64., ScreenHeight-128.
	boss := &Boss{
		Mob: Mob{
			AABB:           AABB{bossX, bossY, 128, 128},
			IdleAnimation:  idle,
			LeftAnimation:  left,
			RightAnimation: right,
			Color:          Color0,
		},
	}
	room.Triggers = append(room.Triggers, Trigger{
		AABB: AABB{ScreenWidth/2 - 150, ScreenHeight - 1, 300, 1},
		Fn: func() {
			g.Scene = NewWaitScene(g)
		},
	})
	var buttons []Button
	for i := 0; i < 8; i++ {
		i := i
		button := Button{
			AABB: AABB{100 + 150*float64(i), 100, 32, 32},
		}
		buttons = append(buttons, button)
		room.Triggers = append(room.Triggers, Trigger{
			AABB: button.AABB,
			Fn: func() {
				buttons[i].Pressed = true
			},
		})
	}
	return &CtrlScene{
		SceneData{
			Room: room,
			Boss: boss,
		},
		buttons,
	}
}

// Update updates the boss and triggers in the scene.
func (c *CtrlScene) Update(g *Game) {
	c.Boss.Update(c.SceneData)
	c.Room.Update(c.SceneData)
}

// Draw draws the boss in the control room.
func (c *CtrlScene) Draw(screen *ebiten.Image) {
	c.Room.Draw(screen)
	for _, btn := range c.Buttons {
		if btn.Pressed {
			ebitenutil.DrawRect(screen, btn.X, btn.Y, btn.W, btn.H, Color6)
		}
	}
	c.Boss.Draw(screen)
}
