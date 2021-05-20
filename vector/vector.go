package vector

import "math"

// Finds angle relative to the y-axis.
// Clockwise is positive and Counterclockwise is negative.
func Angle(v Vector2D) float64 {
	x, y := Components(v)
	rotation := math.Pi / 2
	theta := math.Atan2(-y, x)
	return -1*theta + rotation
}

// Finds standard angle
func AngleReg(v Vector2D) (theta float64) {
	x, y := Components(v)
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

// RotatePoints Rotates points by an angle theta about an origin point
// Rotates in-place
func RotatePoints(theta float64, origin Vector2D, points ...*Vector2D) {
	sin, cos := math.Sincos(theta)
	originX, originY := Components(origin)
	for _, point := range points {
		tildeX, tildeY := point.X-originX, point.Y-originY
		point.X = cos*tildeX + -1*sin*tildeY + originX
		point.Y = sin*tildeX + cos*tildeY + originY
	}
}

// Ccomponents returns the components of the vector
func Components(v Vector2D) (x, y float64) {
	x, y = v.X, v.Y
	return x, y
}

// IsIntersect detects if two linesegments intersect.
// Intersect if 0 <= t <= 1 and 0 <= u <= 1
// https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection
// Points a and b is 1 segment
// Points c and d is the 2 segment
func IsIntersect(a, b, c, d Vector2D) (t, u float64, intersected bool) {
	// Grabbing notation to look like the formula
	var (
		x1, y1 = a.X, a.Y
		x2, y2 = b.X, b.Y
		x3, y3 = c.X, c.Y
		x4, y4 = d.X, d.Y
	)
	z := 1 / ((x1-x2)*(y3-y4) - (y1-y2)*(x3-x4))
	t = ((x1-x3)*(y3-y4) - (y1-y3)*(x3-x4)) * z
	u = ((x2-x1)*(y1-y3) - (y2-y1)*(x1-x3)) * z
	if t >= 0. && t <= 1. && u >= 0. && u <= 1. {
		intersected = true
	}
	return t, u, intersected
}

// IntersectionPoints finds the actual intersection points.
// Will return non-valid numbers if t vector was invalid from IsIntersect
func IntersectionPoint(a, b Vector2D, t float64) (x, y float64) {
	return a.X + t*(b.X-a.X), a.Y + t*(b.Y-a.Y)
}
