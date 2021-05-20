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

// loadImage loads an image from a filepath specified
func loadImage(path string) *ebiten.Image {
	loadedImage, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return loadedImage
}

// Satisfies the ebiten.Game interface.
// requires  Update() error,  Draw(screen *ebiten.Image), Layout(outsideWidth, outsideHeight int) (int, int)
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
	image := loadImage("images/arrow.png")
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
	// Used to visualize the center of the boid image.
	sq := &b.Square{
		ImageWidth:  boid.ImageWidth,
		ImageHeight: boid.ImageHeight,
		SightDis:    3,
		Pos:         boid.Pos,
		Vel:         boid.Vel,
		Accel:       boid.Accel,
	}
	sq.SightDis = 3
	/*
			The triangle is positioned sideways like this.
			Angle of 0 points in the same direction as the top of the triangle.
			|\
			| \
			|  \
			|   \
			|   /
			|  /
			| /
		    |/
	*/
	tri := b.Triangle{
		ImageWidth:  50,
		ImageHeight: 25,
		SightDis:    3,
		Vel:         &v.Vector2D{X: 1, Y: 1},
		Accel:       &v.Vector2D{X: 0, Y: 0},
	}
	tri.TrianglePointsMake(screenWidth/2, screenHeight/2)
	tri2 := b.Triangle{
		ImageWidth:  75,
		ImageHeight: 15,
		SightDis:    3,
		Vel:         &v.Vector2D{X: -3 / 2, Y: -1.75},
		Accel:       &v.Vector2D{X: 0.5, Y: .1},
	}
	tri2.TrianglePointsMake(screenWidth/2, screenHeight/2)
	sim := &Sim{
		population: []b.Boid{
			&tri,
			&tri2,
		},
	}
	if err := ebiten.RunGame(sim); err != nil {
		log.Fatal(err)
	}
}
