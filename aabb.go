package main

type AABB struct {
	X, Y float64
	W, H float64
}

func Collide(a, b AABB) bool {
	return !(a.X+a.W < b.X ||
		a.X > b.X+b.W ||
		a.Y+a.H < b.Y ||
		a.Y > b.Y+b.H)
}
