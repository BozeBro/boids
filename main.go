package main

import (
	"image/color"
	_ "image/png"
	"log"

	b "github.com/BozeBro/boids/boid"
	v "github.com/BozeBro/boids/vector"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 1000
	screenHeight = 1000
)

func makeImage(path string) *ebiten.Image {
	loadedImage, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return loadedImage
}

type Sim struct {
	image      *ebiten.Image
	population []*b.Arrow
}

func (sim *Sim) Update() error {
	for _, object := range sim.population {
		//object.AvoidWalls(float64(screenWidth), float64(screenHeight))
		object.Pos.X += object.Vel.X
		object.Pos.Y += object.Vel.Y
		object.Vel.X += object.Accel.X
		object.Vel.Y += object.Accel.Y
		object.Pos.X = b.Teleport(object.Pos.X, screenWidth)
		object.Pos.Y = b.Teleport(object.Pos.Y, screenHeight)
	}
	return nil
}

func (sim *Sim) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	option := ebiten.DrawImageOptions{}
	for _, boid := range sim.population {
		theta := v.Angle(boid.Vel.X, boid.Vel.Y)
		option.GeoM.Reset()
		option.GeoM.Rotate(theta)
		option.GeoM.Translate(boid.Pos.X, boid.Pos.Y)
		screen.DrawImage(sim.image, &option)
	}
}

func (sim *Sim) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenHeight, screenHeight
}

func main() {
	sim := Sim{image: makeImage("images/arrow.png")}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Boid Simulation")
	w, h := sim.image.Size()
	boid := b.Arrow{
		ImageWidth:  w,
		ImageHeight: h,
		SightDis:    w,
		Pos:         &v.Vector2D{screenWidth / 2, screenHeight / 2},
		Vel:         &v.Vector2D{1, 3},
		Accel:       &v.Vector2D{0, 0},
	}
	sim.population = append(sim.population, &boid)
	if err := ebiten.RunGame(&sim); err != nil {
		log.Fatal(err)
	}
}
