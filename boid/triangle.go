package boid

import (
	"image/color"

	v "github.com/BozeBro/boids/vector"
	"github.com/hajimehoshi/ebiten/v2"
)

// Triangle is an image object for the game screen.
// It satisfies the boid interface.
type Triangle struct {
	ImageWidth  int
	ImageHeight int
	SightAngle  float64
	SightDis    int         // Distance that boid can see in front of it
	Theta       float64     // Angle that the Velocity vectors create.
	VelTheta    float64     // Angle that the Acceleration vectors create.
	Top         *v.Vector2D // Vertex of the initial top point
	Left        *v.Vector2D // Vertex of the initial left point
	Right       *v.Vector2D // Vertex of the initial right point
	Vel         *v.Vector2D
	Accel       *v.Vector2D
}

// TrianglePointsMake uses (x, y) as coordinates of the top.
// Finds the Right and Left Point using ImageWidth, ImageHeight, and (x, y).
func (t *Triangle) TrianglePointsMake(x, y float64) {
	t.Top = &v.Vector2D{X: x, Y: y}
	t.Left = &v.Vector2D{
		X: x - float64(t.ImageWidth),
		Y: y - float64(t.ImageHeight)/2,
	}
	t.Right = &v.Vector2D{
		X: x - float64(t.ImageWidth),
		Y: y + float64(t.ImageHeight)/2,
	}
}

// Only change triangle's location when all of the points are off screen
func (t *Triangle) offscreen(sx, sy float64) {
	var (
		counter  = 0
		vertices = [3]*v.Vector2D{t.Top, t.Left, t.Right}
		// placeholder for new point positions.
		// Will only use the new position for t.Top.
		points = [3]*v.Vector2D{}
	)
	for index, vertex := range vertices {
		if !(vertex.X < 0 || vertex.X > sx || vertex.Y < 0 || vertex.Y > sy) {
			break
		}
		points[index] = &v.Vector2D{
			X: Teleport(vertex.X, sx),
			Y: Teleport(vertex.Y, sy),
		}
		counter++
	}
	if counter != 3 {
		return
	}
	// Place t.Top coords on the other side of the screen.
	// t.Right and t.Right are behind t.Top and still offscreen.
	place := points[0]
	diff := v.Vector2D{
		X: t.Top.X - place.X,
		Y: t.Top.Y - place.Y,
	}
	t.Top = place
	t.Left.Subtract(diff)
	t.Right.Subtract(diff)
}
func (t *Triangle) Add(vector v.Vector2D, points ...*v.Vector2D) {
	for _, point := range points {
		point.Add(vector)
	}
}
func (t *Triangle) Apply(*v.Vector2D, *v.Vector2D) {

}

// Update gives new values to the vertices and velocity vectors.
func (t *Triangle) Update(sx, sy float64, population []Boid, index int, info chan *Data) {
	t.Accel = &v.Vector2D{}
	t.Accel.Add(t.align(population))
	maxi := 7.
	if t.Vel.X > maxi {
		t.Vel.X = maxi
	} else if t.Vel.X < -maxi {
		t.Vel.X = -maxi
	}
	if t.Vel.Y > maxi {
		t.Vel.Y = maxi
	} else if t.Vel.Y < -maxi {
		t.Vel.Y = -maxi
	}
	// Add the vectors
	t.Add(*t.Vel, t.Top, t.Left, t.Right)
	t.Add(*t.Accel, t.Vel)
	// rotate points
	velTheta := v.AngleReg(*t.Accel)
	v.RotatePoints(velTheta-t.VelTheta, *t.Vel, t.Vel)
	t.VelTheta = velTheta

	theta := v.AngleReg(*t.Vel)
	v.RotatePoints(theta-t.Theta, *t.Top, t.Left, t.Right)
	t.Theta = theta
	t.offscreen(sx, sy)
}
func (t *Triangle) Draw(screen *ebiten.Image) {
	option := &ebiten.DrawTrianglesOptions{}
	triangleIm := ebiten.NewImage(t.ImageWidth, t.ImageWidth)
	// Make alpha 1, so the colors will go over it.
	// Colors and Alphas are defined in makeVertex()
	triangleIm.Fill(color.RGBA{255, 255, 255, 1})
	vertex := makeVertex(*t.Top, *t.Left, *t.Right)
	// Draw the vertex onto the triangleIm.
	// Draw the triangleIm onto the screen
	screen.DrawTriangles(vertex, []uint16{0, 1, 2}, triangleIm, option)
}

// makeVertex converts v.Vector2D to ebiten.Vertex
// It defines colors coming from each vertex and its alphas.
func makeVertex(vectors ...v.Vector2D) (vertex []ebiten.Vertex) {
	for index, vector := range vectors {
		var r, g, b float32
		// make rainbow-ish triangle
		switch index {
		case 0:
			r++
		case 1:
			g++
		case 2:
			b++
		}
		// Standard rgba range from 0-255. Give Alpha maximum value
		point := ebiten.Vertex{
			DstX:   float32(vector.X),
			DstY:   float32(vector.Y),
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: 255,
		}
		vertex = append(vertex, point)
	}
	return vertex
}
func (t *Triangle) align(population []Boid) v.Vector2D {
	var counter int
	steering := v.Vector2D{X: 0, Y: 0}
	for _, boid := range population {
		pos := boid.Coords()
		d := v.Distance(*t.Top, pos)
		if boid != t && d <= float64(t.SightDis) {
			counter++
			steering.Add(boid.Velocity())
		}
	}
	if counter > 0 {
		steering.Divide(float64(counter))
		steering.Subtract(*t.Vel)
	}
	return steering
}
func (t *Triangle) Coords() v.Vector2D {
	return *t.Top
}
func (t *Triangle) Velocity() v.Vector2D {
	return *t.Vel
}
