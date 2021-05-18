package boid

import (
	"image/color"

	v "github.com/BozeBro/boids/vector"
	"github.com/hajimehoshi/ebiten/v2"
)

type Triangle struct {
	ImageWidth  int
	ImageHeight int
	SightDis    int
	Theta       float64
	Top         *v.Vector2D // Vertex of the initial top point
	Left        *v.Vector2D // Vertex of the initial left point
	Right       *v.Vector2D // Vertex of the initial right point
	Vel         *v.Vector2D
	Accel       *v.Vector2D
}

// Only change triangle's index when all of the points are off screen
func (t *Triangle) offscreen(sx, sy float64) {
	var counter int
	vertex := []*v.Vector2D{t.Top, t.Left, t.Right}
	// placeholder for new points
	points := make([]*v.Vector2D, 3)
	for index, vertex := range vertex {
		points[index] = &v.Vector2D{
			X: Teleport(vertex.X, sx),
			Y: Teleport(vertex.Y, sy),
		}
		if vertex.X < 0 || vertex.X > sx || vertex.Y < 0 || vertex.Y > sy {
			counter++
		}
	}
	if counter == 3 {
		for index, vertex := range vertex {
			vertex.X = points[index].X
			vertex.Y = points[index].Y
		}
	}
}
func (t *Triangle) Add(vector v.Vector2D, points ...*v.Vector2D) {
	for _, point := range points {
		point.Add(vector)
	}
}
func (t *Triangle) Update(sx, sy float64) {
	t.Add(*t.Accel, t.Vel)
	//_ = v.RotatePoints(0, t.Top.X, t.Top.Y, t.Top, t.Left, t.Right)
	//theta := v.Angle(t.Vel.X, t.Vel.Y)
	if theta := v.RegAngle(t.Vel.X, t.Vel.Y); theta != t.Theta {
		_ = v.RotatePoints(theta, t.Top.X, t.Top.Y, t.Top, t.Left, t.Right)
		t.Theta = theta
	}
	t.Add(*t.Vel, t.Top, t.Left, t.Right)
	t.offscreen(sx, sy)
}
func (t *Triangle) Draw(screen *ebiten.Image) {
	option := &ebiten.DrawTrianglesOptions{}
	triangleIm := ebiten.NewImage(t.ImageWidth, t.ImageWidth)
	triangleIm.Fill(color.RGBA{255, 255, 255, 1})
	vertex := makeVertex(*t.Top, *t.Left, *t.Right)
	// Draw the vertex onto the triangleIm.
	// Draw the triangleIm onto the screen
	screen.DrawTriangles(vertex, []uint16{0, 1, 2}, triangleIm, option)
}

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
		point := ebiten.Vertex{
			DstX:   float32(vector.X),
			DstY:   float32(vector.Y),
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: 254,
		}
		vertex = append(vertex, point)
	}
	return vertex
}
