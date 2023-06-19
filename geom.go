package main

import (
	"math"
)

// AABB is an axis-aligned bounding box. It's the only shape we use.
type AABB struct {
	X, Y, W, H float64
}

// Translate moves the AABB by (dx, dy).
func (a *AABB) Translate(dx, dy float64) {
	a.X += dx
	a.Y += dy
}

// Collide returns true if this AABB overlaps with the other AABB.
func (a AABB) Collide(b AABB) bool {
	return b.X-a.W < a.X && a.X < b.X+b.W &&
		b.Y-a.H < a.Y && a.Y < b.Y+b.H
}

// Detangle moves this AABB so that it no longer overlaps with the other AABB.
// It has no effect of this AABB does not overlap the other.
func (a *AABB) Detangle(b AABB) {
	// Move a.X so that [a.X, a.X+a.W] does not overlap with [b.X, b.x+b.W].
	// For that to happen, either the left side of a must be to the right of b,
	// meaning that b.X < a.X+a.W (equivalently, b.X-a.W < a.X), or the right
	// side of a must be to the left of b, meaning that a.X > b.X+b.W. We do
	// this with Unclamp, making sure that a.X is outside [b.X-a.W,b.X+b.W].
	// We just figure out what x should be, we don't apply the update yet.
	x := Unclamp(b.X-a.W, a.X, b.X+b.W)
	// Do the same thing to y that we did to x.
	y := Unclamp(b.Y-a.H, a.Y, b.Y+b.H)
	// Now that we know how to move either x or y to detangle a and b, actually
	// do the move. We move in either x or y, whichever is smaller. If a and b
	// just collided after a small movement, this adjustment will likely not be
	// visible to the player, even if it does mean that a moves contrary to its
	// desired velocity.
	if math.Abs(a.X-x) < math.Abs(a.Y-y) {
		a.X = x
	} else {
		a.Y = y
	}
}

// DetangleRoom detangles this AABB with every obstactle in the room.
func (a *AABB) DetangleRoom(r *Room) {
	for _, o := range r.Obstacles {
		a.Detangle(o)
	}
}

// ClampToBound moves the AABB so that it is within the Room bounds.
func (a *AABB) ClampToRoom(r *Room) {
	a.X = Clamp(0, a.X, r.Width-a.W)
	a.Y = Clamp(0, a.Y, r.Height-a.H)
}

// Clamp ensures that x is in the range [a, b]. It returns x if x is in [a,
// b], a if x < a, and b if x > b. Clamp does NOT check that a < b, and the
// behavior is undefined in this case.
func Clamp(a, x, b float64) float64 {
	return math.Max(a, math.Min(x, b))
}

// Unclamp ensures that x is outside the range (a, b). It returns x if x is
// outside (a, b), and either a or b, whichever is closest to x, if x is inside
// (a, b). Unclamp does NOT check that a < b, and the behavior is undefined in
// this case.
func Unclamp(a, x, b float64) float64 {
	if x <= a || x >= b {
		return x
	}
	if x >= (a+b)/2 {
		return b
	}
	return a
}

// SliceContains returns true if the slice contains x.
func SliceContains[T comparable](slice []T, x T) bool {
	for _, y := range slice {
		if x == y {
			return true
		}
	}
	return false
}

// UpdateDelta gets a desired component velocity subject to game constraints.
func UpdateDelta(delta, actual, target float64) float64 {
	const (
		DeltaDecay = .65
		DeltaStep  = 1
		DeltaMax   = 5
	)

	err := target - actual
	if err == 0 {
		return delta * DeltaDecay
	}
	step := Clamp(-DeltaStep, err-delta*DeltaMax/DeltaStep, DeltaStep)
	return Clamp(-DeltaMax, delta+step, DeltaMax)
}
