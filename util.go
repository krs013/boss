package main

import (
	"math"
)

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
