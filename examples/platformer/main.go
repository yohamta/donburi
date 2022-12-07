package main

import (
	"image"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi/examples/platformer/scenes"
)

type Scene interface {
	Update()
	Draw(screen *ebiten.Image)
}

type Game struct {
	bounds image.Rectangle
	scene  Scene
}

func NewGame() *Game {
	g := &Game{
		bounds: image.Rectangle{},
		scene:  scenes.NewPlatformerScene(),
	}

	return g
}

func (g *Game) Update() error {
	g.scene.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	g.scene.Draw(screen)
}

func (g *Game) Layout(width, height int) (int, int) {
	g.bounds = image.Rect(0, 0, width, height)
	return width, height
}

func main() {
	ebiten.SetWindowSize(640, 360)
	ebiten.SetWindowResizable(false)
	rand.Seed(time.Now().UTC().UnixNano())
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
