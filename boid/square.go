package boid

import (
	"math"

	v "github.com/BozeBro/boids/vector"
	"github.com/hajimehoshi/ebiten/v2"
)

type Square struct {
	ImageWidth  int
	ImageHeight int
	SightDis    int
	Pos         *v.Vector2D
	Vel         *v.Vector2D
	Accel       *v.Vector2D
	Image       *ebiten.Image
}

func (s *Square) Update(sx, sy float64) {
	s.Pos.X += s.Vel.X
	s.Pos.Y += s.Vel.Y
	s.Vel.X += s.Accel.X
	s.Vel.Y += s.Accel.Y
	s.Pos.X = Teleport(s.Pos.X, sx)
	s.Pos.Y = Teleport(s.Pos.Y, sy)
}

func (s *Square) Trig() (x, y float64) {
	const (
		partW   = 2
		partH   = 2
		offsetx = 0
		offsety = 0
	)
	unitVecX := -1.
	unitVecY := 1.
	if s.Vel.X != 0 {
		unitVecX = -1 * s.Vel.X / math.Abs(s.Vel.X)
	}
	if s.Vel.Y != 0 {
		unitVecY = -1 * s.Vel.Y / math.Abs(s.Vel.Y)
	}
	x = s.Pos.X + unitVecX*float64(s.ImageWidth)/partW - offsetx
	y = s.Pos.Y + unitVecY*float64(s.ImageHeight)/partH - offsety
	return x, y
}

func (s *Square) Draw(screen *ebiten.Image, option ebiten.DrawImageOptions) {
	offsetx, offsety := s.Trig()
	theta := v.Angle(s.Vel.X, s.Vel.Y)
	option.GeoM.Rotate(theta)
	option.GeoM.Translate(s.Pos.X+offsetx, s.Pos.Y+offsety)
	screen.DrawImage(s.Image, &option)
}
