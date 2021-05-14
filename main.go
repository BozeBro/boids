package main

import (
	"image/color"
	_ "image/png"
	"log"

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
	image *ebiten.Image
}

func (sim *Sim) Update() error {
	return nil
}

func (sim *Sim) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	option := ebiten.DrawImageOptions{}
	option.GeoM.Translate(screenWidth/2, screenHeight/2)
	screen.DrawImage(sim.image, &option)
}

func (sim *Sim) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenHeight, screenHeight
}

func main() {
	sim := Sim{image: makeImage("images/boid.png")}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Boid Simulation")
	if err := ebiten.RunGame(&sim); err != nil {
		log.Fatal(err)
	}
}
