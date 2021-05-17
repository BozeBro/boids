package vector

import "math"

func Angle(x, y float64) float64 {
	rotation := math.Pi / 2
	theta := math.Atan2(-y, x)
	return -1*theta + rotation
}

type Vector2D struct {
	X float64
	Y float64
}

func (v *Vector2D) Add(v2 Vector2D) {
	v.X += v2.X
	v.Y += v2.Y
}

// Rotates points by an angle theta
func RotatePoints(theta float64, points ...*Vector2D) {
	sin, cos := math.Sincos(theta)
	for _, point := range points {
		nx, ny := cos*point.X+-1*sin*point.Y, sin*point.X+cos*point.Y
		point.X, point.Y = nx, ny
	}
}
