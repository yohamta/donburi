package main

import (
	"image"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"

	_ "net/http/pprof"
)

type Game struct {
	ecs    *ecs.ECS
	bounds image.Rectangle
}

func NewGame() *Game {
	g := &Game{
		bounds: image.Rectangle{},
		ecs:    createECS(),
	}

	return g
}

func createECS() *ecs.ECS {
	world := createWorld()
	ecs := ecs.NewECS(world)
	return ecs
}

func createWorld() donburi.World {
	world := donburi.NewWorld()
	return world
}

func (g *Game) Update() error {
	g.ecs.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
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
