package boid

import (
	"math"

	v "github.com/BozeBro/boids/vector"
	"github.com/hajimehoshi/ebiten/v2"
)

type Boid interface {
	Update(float64, float64, []Boid, int, chan *Data)
	Draw(*ebiten.Image)
	Coords() v.Vector2D
	Velocity() v.Vector2D
	Apply(*v.Vector2D, *v.Vector2D)
}

type Data struct {
	Index  int
	NewPos *v.Vector2D
	NewVel *v.Vector2D
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

func (a *Arrow) avoidWalls(sx, sy float64) *v.Vector2D {
	steering := &v.Vector2D{0, 0}
	theta := v.AngleReg(*a.Vel)
	heading := v.Vector2D{a.Pos.X + a.SightDis*math.Cos(theta), a.Pos.Y + a.SightDis*math.Sin(theta)}
	topL := v.Vector2D{0, 0}
	topR := v.Vector2D{sx, 0}
	botL := v.Vector2D{0, sy}
	botR := v.Vector2D{sx, sy}
	var (
		lt, _, lbool = v.IsIntersect(*a.Pos, heading, topL, v.Vector2D{0, sy - 1})
		ut, _, ubool = v.IsIntersect(*a.Pos, heading, v.Vector2D{1, 0}, topR)
		rt, _, rbool = v.IsIntersect(*a.Pos, heading, v.Vector2D{sx - 1, 0}, botR)
		dt, _, dbool = v.IsIntersect(*a.Pos, heading, botL, v.Vector2D{sx - 1, sy})
	)
	//var x, y float64
	if lbool {
		steerAway := *a.Vel
		x, y := v.IntersectionPoint(*a.Pos, heading, lt)
		d := v.Distance(*a.Pos, v.Vector2D{x, y})
		//wall := v.Vector2D{x, y}
		if a.Vel.Y < 0 {
			v.RotatePoints(math.Pi/2, v.Vector2D{}, &steerAway)
		} else {
			v.RotatePoints(-math.Pi/2, v.Vector2D{}, &steerAway)
		}
		steerAway.Multiply(d)
		steering.Add(steerAway)
	}
	if ubool {
		steerAway := *a.Vel
		x, y := v.IntersectionPoint(*a.Pos, heading, ut)
		d := v.Distance(*a.Pos, v.Vector2D{x, y})
		//wall = v.Vector2D{x, y}
		if a.Vel.X < 0 {
			v.RotatePoints(-math.Pi/2, v.Vector2D{}, &steerAway)
		} else {
			v.RotatePoints(math.Pi/2, v.Vector2D{}, &steerAway)
		}
		steerAway.Multiply(d)
		steering.Add(steerAway)
	}
	if rbool {
		steerAway := *a.Vel
		x, y := v.IntersectionPoint(*a.Pos, heading, rt)
		d := v.Distance(*a.Pos, v.Vector2D{x, y})
		//wall = v.Vector2D{x, y}
		if a.Vel.Y < 0 {
			v.RotatePoints(-math.Pi/2, v.Vector2D{}, &steerAway)
		} else {
			v.RotatePoints(math.Pi/2, v.Vector2D{}, &steerAway)
		}
		steerAway.Multiply(d)
		steering.Add(steerAway)
	}
	if dbool {
		steerAway := *a.Vel
		x, y := v.IntersectionPoint(*a.Pos, heading, dt)
		d := v.Distance(*a.Pos, v.Vector2D{x, y})
		//wall = v.Vector2D{x, y}
		if a.Vel.X < 0 {
			v.RotatePoints(math.Pi/2, v.Vector2D{}, &steerAway)
		} else {
			v.RotatePoints(-math.Pi/2, v.Vector2D{}, &steerAway)
		}
		steerAway.Multiply(d)
		steering.Add(steerAway)
	}
	//steering.SetMagnitude(4)
	//steering.Limit(1.)
	return steering
}
func (a *Arrow) Update(sx, sy float64, population []Boid, index int, info chan *Data) {
	a.Pos.X = Teleport(a.Pos.X, sx)
	a.Pos.Y = Teleport(a.Pos.Y, sy)
	align, cohesion, separation := a.rules(population)
	a.Accel.Add(*align)
	a.Accel.Add(*cohesion)
	a.Accel.Add(*separation)
	//a.Accel.Limit(0.3)
	a.Accel.SetMagnitude(1.)
	avoid := a.avoidWalls(sx, sy)
	a.Accel = &v.Vector2D{
		X: a.Accel.X + avoid.X,
		Y: a.Accel.Y + avoid.Y,
	}
	newPos := &v.Vector2D{a.Pos.X + a.Vel.X, a.Pos.Y + a.Vel.Y}
	newVel := &v.Vector2D{a.Vel.X + a.Accel.X, a.Vel.Y + a.Accel.Y}
	a.Accel = &v.Vector2D{0, 0}
	newVel.Limit(4)
	data := &Data{
		Index:  index,
		NewPos: newPos,
		NewVel: newVel,
	}
	info <- data

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
func (a *Arrow) Apply(newPos, newVel *v.Vector2D) {
	a.Pos = newPos
	a.Vel = newVel
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
	maxforceA, maxforceC, maxforceS := 1., 0.9, 1.5
	perceptionA, perceptionC, perceptionS := 75., 70., 50.
	var counterA, counterC, counterS int
	cosView := math.Cos(a.SightAngle)
	newPos := &v.Vector2D{a.Pos.X + a.Vel.X, a.Pos.Y + a.Vel.Y}
	for _, boid := range population {
		pos := boid.Coords()
		//align
		angle := v.IsSeen(a.Pos, newPos, &pos)
		if boid != a && angle > cosView {
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
				//diff.Divide(d)
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

/* func (a *Arrow) avoidWalls(sx, sy float64) {
	//theta := math.Atan2(a.Pos.X, a.Pos.Y)
} */

func (a *Arrow) Coords() v.Vector2D {
	return *a.Pos
}

func (a *Arrow) Velocity() v.Vector2D {
	return *a.Vel
}
