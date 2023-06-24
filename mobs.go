package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Mob struct {
	AABB
	HasTarget        bool
	TargetX, TargetY float64
	DX, DY           float64

	color.Color
	IdleAnimation  *Animation
	LeftAnimation  *Animation
	RightAnimation *Animation
}

func (m Mob) Draw(dst *ebiten.Image) {
	ebitenutil.DrawRect(dst, m.X, m.Y, m.W, m.H, m.Color)
	if m.DX == 0 && m.DY == 0 && m.IdleAnimation != nil {
		m.IdleAnimation.Draw(dst, m.X, m.Y)
	} else if m.DX > 0 && m.RightAnimation != nil {
		m.RightAnimation.Draw(dst, m.X, m.Y)
	} else if m.LeftAnimation != nil {
		m.LeftAnimation.Draw(dst, m.X, m.Y)
	}
}

// UpdateDelta gets a desired component velocity subject to game constraints.
func UpdateDelta(delta, actual, target float64) float64 {
	const (
		DeltaDecay = .65
		DeltaStep  = 1
		DeltaMax   = 5
		DeltaMin   = 0.0001
	)

	err := target - actual
	if err == 0 {
		delta *= DeltaDecay
	} else {
		step := Clamp(-DeltaStep, err-delta*DeltaMax/DeltaStep, DeltaStep)
		delta = Clamp(-DeltaMax, delta+step, DeltaMax)
	}
	if math.Abs(delta) < DeltaMin {
		return 0
	}
	return delta
}

// TrackTarget updates the Mob DX, DY in the direction of the target.
func (m *Mob) TrackTarget() {
	if m.HasTarget {
		m.DX = UpdateDelta(m.DX, m.X, m.TargetX)
		m.DY = UpdateDelta(m.DY, m.Y, m.TargetY)
	} else {
		m.DX = UpdateDelta(m.DX, m.X, m.X)
		m.DY = UpdateDelta(m.DY, m.Y, m.Y)
	}
}

func (m *Mob) Move() {
	m.Translate(m.DX, m.DY)
}

type Boss struct {
	Mob
}

func (b *Boss) KeyTarget() (x, y float64, ok bool) {
	x, y = b.X, b.Y
	ok = false
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		x += ScreenWidth
		ok = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		x -= ScreenWidth
		ok = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		y += ScreenHeight
		ok = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		y -= ScreenHeight
		ok = true
	}
	return x, y, ok
}

func (b *Boss) MouseTarget() (x, y float64, ok bool) {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		cx, cy := ebiten.CursorPosition()
		return float64(cx) - b.W/2, float64(cy) - b.H/2, true
	}
	return b.X, b.Y, false
}

func (b *Boss) Update(w *Wait) {
	tx, ty, ok := b.KeyTarget()
	if ok {
		b.HasTarget = false
	} else if tx, ty, ok = b.MouseTarget(); ok {
		b.HasTarget = true
		b.TargetX, b.TargetY = tx, ty
	} else if b.HasTarget {
		tx, ty = b.TargetX, b.TargetY
	}
	b.DX = UpdateDelta(b.DX, b.X, tx)
	b.DY = UpdateDelta(b.DY, b.Y, ty)

	b.Move()
	b.DetangleRoom(w.Room)
	b.ClampToRoom(w.Room)

	// Boss pushes the hero, so hero does the detangle, not boss.
	w.Hero.Detangle(b.AABB)
	// After hero gets pushed, hero might need to avoid obstacles.
	w.Hero.DetangleRoom(w.Room)
	w.Hero.ClampToRoom(w.Room)
	// Hero might push back if they don't fit, so now boss detangles
	b.Detangle(w.Hero.AABB)
}

type Hero struct {
	Mob
}

func (h *Hero) Update(w *Wait) {
	// Set target so hero center moves towards the boss center.
	h.TargetX = w.Boss.X + w.Boss.W/2 - h.W/2
	h.TargetY = w.Boss.Y + w.Boss.H/2 - h.H/2
	h.TrackTarget()

	// Just do the move. Boss will handle collisions and pushing.
	h.Move()
}
