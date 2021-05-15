package vector

import "math"

type Vector2D struct {
	X float64
	Y float64
}

func Angle(x, y float64) float64 {
	rotation := math.Pi / 2
	theta := math.Atan2(-y, x)
	return -1*theta + rotation
}
