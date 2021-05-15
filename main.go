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
	screenWidth  = 500
	screenHeight = 500
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
	option := ebiten.DrawImageOptions{}
	for _, boid := range sim.population {
		boid.Draw(screen, option)
	}
}

func (sim *Sim) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenHeight, screenHeight
}

func main() {
	sim := Sim{}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Boid Simulation")
	image := makeImage("images/arrowV2.png")
	w, h := image.Size()
	boid := b.Arrow{
		ImageWidth:  w,
		ImageHeight: h,
		SightDis:    w,
		Pos:         &v.Vector2D{X: screenWidth / 2, Y: screenHeight / 2},
		Vel:         &v.Vector2D{X: 3, Y: 1},
		Accel:       &v.Vector2D{X: 0, Y: 0},
	}
	sq := b.Square{
		Image:       image,
		ImageWidth:  boid.ImageWidth,
		ImageHeight: boid.ImageHeight,
		SightDis:    3,
		Pos:         boid.Pos,
		Vel:         boid.Vel,
		Accel:       boid.Accel,
	}
	sim.population = append(sim.population, &boid, &sq)
	if err := ebiten.RunGame(&sim); err != nil {
		log.Fatal(err)
	}
}
