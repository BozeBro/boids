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
	population []b.Boid
}

func (sim *Sim) Update() error {
	for _, object := range sim.population {
		object.Update(screenWidth, screenHeight)
	}
	return nil
}

func (sim *Sim) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	for _, boid := range sim.population {
		boid.Draw(screen)
	}
}

func (sim *Sim) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenHeight, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Boid Simulation")
	image := makeImage("images/arrow.png")
	w, h := image.Size()
	boid := &b.Arrow{
		Image:       image,
		ImageWidth:  w,
		ImageHeight: h,
		SightDis:    w,
		Pos:         &v.Vector2D{X: screenWidth / 2, Y: screenHeight / 2},
		Vel:         &v.Vector2D{X: 1, Y: 0},
		Accel:       &v.Vector2D{X: 0, Y: 0},
	}
	sq := &b.Square{
		ImageWidth:  boid.ImageWidth,
		ImageHeight: boid.ImageHeight,
		SightDis:    3,
		Pos:         boid.Pos,
		Vel:         boid.Vel,
		Accel:       boid.Accel,
	}
	tri := &b.Triangle{
		ImageWidth:  50,
		ImageHeight: 25,
		SightDis:    3,
		Top:         &v.Vector2D{X: screenWidth / 2, Y: screenHeight / 2},
		Vel:         &v.Vector2D{X: 10, Y: 1},
		Accel:       &v.Vector2D{X: -1, Y: 1},
	}
	tri.Left = &v.Vector2D{
		X: tri.Top.X - float64(tri.ImageWidth),
		Y: tri.Top.Y - float64(tri.ImageHeight)/2,
	}
	tri.Right = &v.Vector2D{
		X: tri.Top.X - float64(tri.ImageWidth),
		Y: tri.Top.Y + float64(tri.ImageHeight)/2,
	}
	sim := &Sim{}
	sq.SightDis = 3
	sim.population = append(sim.population, tri)
	if err := ebiten.RunGame(sim); err != nil {
		log.Fatal(err)
	}
}
