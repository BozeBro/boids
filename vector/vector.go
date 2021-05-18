package vector

import "math"

func Angle(x, y float64) float64 {
	rotation := math.Pi / 2
	theta := math.Atan2(-y, x)
	return -1*theta + rotation
}

func RegAngle(x, y float64) (theta float64) {
	theta = math.Atan2(y, x)
	return theta
}

type Vector2D struct {
	X float64
	Y float64
}

func (v *Vector2D) Add(v2 Vector2D) {
	v.X += v2.X
	v.Y += v2.Y
}
func (v *Vector2D) Subtract(v2 Vector2D) {
	v.X -= v2.X
	v.Y -= v2.Y
}

// Rotates points by an angle theta
func RotatePoints(theta, originX, originY float64, points ...*Vector2D) []*Vector2D {
	sin, cos := math.Sincos(theta)
	for _, point := range points {
		tildeX, tildeY := point.X-originX, point.Y-originY
		point.X = cos*tildeX + -1*sin*tildeY + originX
		point.Y = sin*tildeX + cos*tildeY + originY
	}
	return points
}
