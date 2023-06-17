package main

import (
	"math"
)

func Clamp(a, x, b float64) float64 {
	return math.Max(a, math.Min(x, b))
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
