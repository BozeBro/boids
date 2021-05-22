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
	SightDis    float64
	SightAngle  float64
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
	a.Pos.X = Teleport(a.Pos.X, sx)
	a.Pos.Y = Teleport(a.Pos.Y, sy)
	align, cohesion, separation := a.rules(population)
	a.Accel.Add(*align)
	a.Accel.Add(*cohesion)
	a.Accel.Add(*separation)
	/* align.Divide(1)
	cohesion.Divide(1)
	separation.Divide(1) */
	maxi := 2.
	a.Pos.Add(*a.Vel)
	a.Vel.Mini(maxi)
	a.Vel.Add(*a.Accel)
	//a.Vel.Limit(3.)

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
	maxforce := 0.1
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
	maxspeed := 6.
	maxforce := 0.03
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

func (a *Arrow) rules(population []Boid) (steeringA, steeringC, steeringS *v.Vector2D) {
	steeringA, steeringC, steeringS = &v.Vector2D{}, &v.Vector2D{}, &v.Vector2D{}
	maxspeedA, maxspeedC, maxspeedS := 4., 4., 4.
	maxforceA, maxforceC, maxforceS := 1., 0.9, 1.2
	perceptionA, perceptionC, perceptionS := 75., 100., 50.
	var counterA, counterC, counterS int
	for _, boid := range population {
		pos := boid.Coords()
		//align
		if boid != a {
			d := v.Distance(*a.Pos, pos)
			if d <= perceptionA {
				counterA++
				steeringA.Add(boid.Velocity())
			}
			if d <= perceptionC {
				counterC++
				steeringC.Add(boid.Coords())
			}
			if d <= perceptionS {
				counterS++
				diff := *a.Pos
				diff.Subtract(boid.Coords())
				diff.Divide(d)
				steeringS.Add(diff)
			}
		}
	}
	if counterA > 0 {
		steeringA.Divide(float64(counterA))
		steeringA.SetMagnitude(maxspeedA)
		steeringA.Subtract(*a.Vel)
		steeringA.Limit(maxforceA)
	}
	if counterC > 0 {
		steeringC.Divide(float64(counterC))
		steeringC.Subtract(*a.Pos)
		steeringC.SetMagnitude(maxspeedC)
		steeringC.Subtract(*a.Vel)
		steeringC.Limit(maxforceC)
	}
	if counterS > 0 {
		steeringS.Divide(float64(counterS))
		steeringS.SetMagnitude(maxspeedS)
		steeringS.Subtract(*a.Vel)
		steeringS.Limit(maxforceS)
	}
	return steeringA, steeringC, steeringS
}

func (a *Arrow) Coords() v.Vector2D {
	return *a.Pos
}

func (a *Arrow) Velocity() v.Vector2D {
	return *a.Vel
}
