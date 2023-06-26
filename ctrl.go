package main

import "github.com/hajimehoshi/ebiten/v2"

// CtrlScene has the Boss in the control room with the level design buttons.
type CtrlScene struct {
	SceneData
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
	return &CtrlScene{
		SceneData{
			Room: room,
			Boss: boss,
		},
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
	c.Boss.Draw(screen)
}
