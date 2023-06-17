package main

import (
	"math"
)

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
