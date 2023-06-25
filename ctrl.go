package main

import "github.com/hajimehoshi/ebiten/v2"

type CtrlScene struct {
	RoomScene
}

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
		RoomScene{
			Room: room,
			Boss: boss,
		},
	}
}

func (c *CtrlScene) Update(g *Game) {
	c.Boss.Update(c.RoomScene)
	c.Room.Update(c.RoomScene)
}

func (c *CtrlScene) Draw(screen *ebiten.Image) {
	c.Room.Draw(screen)
	c.Boss.Draw(screen)
}
