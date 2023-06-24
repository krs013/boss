package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

func (r Room) Update(w *WaitScene) {
	for _, t := range r.Triggers {
		if w.Boss.Collide(t.AABB) {
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
