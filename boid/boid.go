package boid

import (
	v "github.com/BozeBro/boids/vector"
	"github.com/hajimehoshi/ebiten/v2"
)

type Boid interface {
	// screen sizes
	Update(float64, float64)
	Draw(*ebiten.Image, ebiten.DrawImageOptions)
}
type Arrow struct {
	ImageWidth  int
	ImageHeight int
	SightDis    int
	Pos         *v.Vector2D
	Vel         *v.Vector2D
	Accel       *v.Vector2D
	Image       *ebiten.Image
}

func Teleport(pos, edge float64) float64 {
	if pos < 0 {
		pos = edge
	} else if pos > edge {
		pos = 0.
	}
	return pos
}
func (a *Arrow) Update(sx, sy float64) {
	a.Pos.X += a.Vel.X
	a.Pos.Y += a.Vel.Y
	a.Vel.X += a.Accel.X
	a.Vel.Y += a.Accel.Y
	a.Pos.X = Teleport(a.Pos.X, sx)
	a.Pos.Y = Teleport(a.Pos.Y, sy)
}
func (a *Arrow) SupplyImage() *ebiten.Image {
	return nil
}

func (a *Arrow) Draw(screen *ebiten.Image, option ebiten.DrawImageOptions) {
	theta := v.Angle(a.Vel.X, a.Vel.Y)
	option.GeoM.Rotate(theta)
	option.GeoM.Translate(a.Pos.X, a.Pos.Y)
	screen.DrawImage(a.Image, &option)
}
