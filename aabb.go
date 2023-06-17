package main

import "math"

type AABB struct {
	X, Y, W, H float64
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
