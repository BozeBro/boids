package boid

import (
	v "github.com/BozeBro/boids/vector"
	"github.com/hajimehoshi/ebiten/v2"
)

type Boid interface {
	Update(float64, float64, []Boid)
	Draw(*ebiten.Image)
	Coords() v.Vector2D
	Velocity() v.Vector2D
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

// Teleport places point on opposite end of the screen when offscreen.
func Teleport(pos, edge float64) float64 {
	if pos < 0 {
		pos = edge
	} else if pos > edge {
		pos = 0.
	}
	return pos
}
func (a *Arrow) Update(sx, sy float64, population []Boid) {
	a.Accel.Add(a.align(population))
	a.Accel.Add((a.cohesion(population)))
	a.Accel.Add((a.separation(population)))
	a.Vel.Add(*a.Accel)
	maxi := 5.
	if a.Vel.X > maxi {
		a.Vel.X = maxi
	} else if a.Vel.X < -maxi {
		a.Vel.X = -maxi
	}
	if a.Vel.Y > maxi {
		a.Vel.Y = maxi
	} else if a.Vel.Y < -maxi {
		a.Vel.Y = -maxi
	}
	a.Pos.Add(*a.Vel)

	a.Pos.X = Teleport(a.Pos.X, sx)
	a.Pos.Y = Teleport(a.Pos.Y, sy)
	a.Accel = &v.Vector2D{}
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
func (a *Arrow) align(population []Boid) v.Vector2D {
	maxspeed := 5.
	maxforce := 0.07
	var counter int
	steering := v.Vector2D{}
	for _, boid := range population {
		pos := boid.Coords()
		d := v.Distance(*a.Pos, pos)
		if boid != a && d <= float64(a.SightDis) {
			counter++
			steering.Add(boid.Velocity())
		}
	}
	if counter > 0 {
		steering.Divide(float64(counter))
		steering.SetMagnitude(maxspeed)
		steering.Subtract(*a.Vel)
		steering.Limit(maxforce)
	}
	return steering
}
func (a *Arrow) cohesion(population []Boid) v.Vector2D {
	maxspeed := 5.
	maxforce := 0.07
	var counter int
	steering := v.Vector2D{}
	for _, boid := range population {
		d := v.Distance(*a.Pos, boid.Coords())
		if boid != a && d <= float64(a.SightDis) {
			counter++
			steering.Add(boid.Coords())
		}
	}
	if counter > 0 {
		steering.Divide(float64(counter))
		steering.Subtract(*a.Pos)
		steering.SetMagnitude(maxspeed)
		steering.Subtract(*a.Vel)
		steering.SetMagnitude(maxforce * .8)
	}
	return steering
}

func (a *Arrow) separation(population []Boid) v.Vector2D {
	maxspeed := 5.
	maxforce := 0.1
	var counter int
	perception := 50.
	steering := v.Vector2D{}
	for _, boid := range population {
		d := v.Distance(*a.Pos, boid.Coords())
		if boid != a && d <= perception {
			counter++
			diff := *a.Pos
			diff.Subtract(boid.Coords())
			diff.Divide(d)
			steering.Add(diff)
		}
	}
	if counter > 0 {
		steering.Divide(float64(counter))
		steering.SetMagnitude(maxspeed)
		steering.Subtract(*a.Vel)
		steering.SetMagnitude(maxforce * 1.4)
	}
	return steering
}

func (a *Arrow) Coords() v.Vector2D {
	return *a.Pos
}

func (a *Arrow) Velocity() v.Vector2D {
	return *a.Vel
}
