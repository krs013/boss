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
