package main

import (
	"math"
)

type AABB struct {
	X, Y, W, H float64
}

func (a *AABB) Move(dx, dy float64) {
	a.X += dx
	a.Y += dy
}

func (a AABB) Collide(b AABB) bool {
	return b.X-a.W < a.X && a.X < b.X+b.W &&
		b.Y-a.H < a.Y && a.Y < b.Y+b.H
}

func (a *AABB) Detangle(b AABB) {
	x := Unclamp(b.X-a.W, a.X, b.X+b.W)
	y := Unclamp(b.Y-a.H, a.Y, b.Y+b.H)
	if math.Abs(a.X-x) < math.Abs(a.Y-y) {
		a.X = x
	} else {
		a.Y = y
	}
}

func (a *AABB) DetangleRoom(r *Room) {
	for _, o := range r.Obstacles {
		a.Detangle(o)
	}
	a.ClampToBound(r.Width, r.Height)
}

func (a *AABB) ClampToBound(width, height float64) {
	a.X = Clamp(0, a.X, width-a.W)
	a.Y = Clamp(0, a.Y, height-a.H)
}

func Clamp(a, x, b float64) float64 {
	return math.Max(a, math.Min(x, b))
}

func Unclamp(a, x, b float64) float64 {
	if x <= a || x >= b {
		return x
	}
	if x >= (a+b)/2 {
		return b
	}
	return a
}

func SliceContains[T comparable](slice []T, x T) bool {
	for _, y := range slice {
		if x == y {
			return true
		}
	}
	return false
}

func UpdateDelta(delta float64, inc, dec bool) float64 {
	const (
		DeltaDecay = .75
		DeltaStep  = 1
		DeltaMax   = 5
	)

	if inc == dec {
		return delta * DeltaDecay
	}
	if inc {
		delta += DeltaStep
	}
	if dec {
		delta -= DeltaStep
	}
	return Clamp(-DeltaMax, delta, DeltaMax)
}
