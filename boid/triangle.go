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
	Top         *v.Vector2D // Vertex of the initial top point
	Left        *v.Vector2D // Vertex of the initial left point
	Right       *v.Vector2D // Vertex of the initial right point
	Vel         *v.Vector2D
	Accel       *v.Vector2D
}

func (t *Triangle) Add(vector v.Vector2D, points ...*v.Vector2D) {
	for _, point := range points {
		point.X += vector.X
		point.Y += vector.Y
		//point.X = Teleport(point.X, vector.X)
		//point.Y = Teleport(point.Y, vector.Y)
	}
}
func (t *Triangle) Update(sx, sy float64) {
	t.Add(*t.Vel, t.Top, t.Left, t.Right)
	t.Add(*t.Accel, t.Vel)
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
func makeVertex(vectors ...v.Vector2D) []ebiten.Vertex {
	var vertex []ebiten.Vertex
	for index, vector := range vectors {
		var r, g, b float32
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
			ColorA: 255,
		}
		vertex = append(vertex, point)
	}
	return vertex
}
