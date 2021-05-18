package boid

import (
	v "github.com/BozeBro/boids/vector"
	"github.com/hajimehoshi/ebiten/v2"
)

type Boid interface {
	// screen sizes
	Update(float64, float64)
	Draw(*ebiten.Image)
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
	a.Pos.Add(*a.Vel)
	a.Vel.Add(*a.Accel)
	a.Pos.X = Teleport(a.Pos.X, sx)
	a.Pos.Y = Teleport(a.Pos.Y, sy)
}
func (a *Arrow) Draw(screen *ebiten.Image) {
	option := &ebiten.DrawImageOptions{}
	// Do this translation so PosX, PosY is near the center of the arrow.
	option.GeoM.Translate(-1*float64(a.ImageWidth)/2, -1*float64(a.ImageHeight)/2)
	// Don't rotate if vectors are nil
	if a.Vel.X != 0 || a.Vel.Y != 0 {
		theta := v.Angle(*a.Vel)
		option.GeoM.Rotate(theta)
	}
	option.GeoM.Translate(a.Pos.X, a.Pos.Y)
	screen.DrawImage(a.Image, option)
}
