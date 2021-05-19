package boid

import (
	"image/color"

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
}

func (s *Square) Update(sx, sy float64) {
	s.Pos.Add(*s.Vel)
	s.Vel.Add(*s.Accel)
	s.Pos.X = Teleport(s.Pos.X, sx)
	s.Pos.Y = Teleport(s.Pos.Y, sy)
}
func (s *Square) Draw(screen *ebiten.Image) {
	option := &ebiten.DrawImageOptions{}
	// Don't rotate if vectors are nil
	if s.Vel.X != 0 || s.Vel.Y != 0 {
		theta := v.Angle(*s.Vel)
		option.GeoM.Rotate(theta)
	}
	option.GeoM.Translate(s.Pos.X, s.Pos.Y)
	sq := ebiten.NewImage(2, 2)
	sq.Fill(color.RGBA{255, 0, 0, 254})
	screen.DrawImage(sq, option)
}

/* func (s *Square) Trig() (x, y float64) {
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
} */
