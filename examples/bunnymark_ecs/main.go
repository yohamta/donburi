package main

import (
	"image"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/assets"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/component"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/helper"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/layers"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/system"

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

	metrics := system.NewMetrics(&g.bounds)

	g.ecs.
		AddSystem(system.NewSpawn().Update).
		AddSystem(metrics.Update).
		AddSystem(system.NewBounce(&g.bounds).Update).
		AddSystem(system.Velocity.Update).
		AddSystem(system.Gravity.Update).
		AddRenderer(layers.LayerBackground, system.DrawBackground).
		AddRenderer(layers.LayerMetrics, metrics.Draw).
		AddRenderer(layers.LayerBunnies, system.Render.Draw)

	return g
}

func createECS() *ecs.ECS {
	world := createWorld()
	ecs := ecs.NewECS(world)
	return ecs
}

func createWorld() donburi.World {
	world := donburi.NewWorld()
	settings := world.Entry(world.Create(component.Settings))
	donburi.SetValue(
		settings,
		component.Settings,
		component.SettingsData{
			Ticker:   time.NewTicker(500 * time.Millisecond),
			Gpu:      helper.GpuInfo(),
			Tps:      helper.NewPlot(20, 60),
			Fps:      helper.NewPlot(20, 60),
			Objects:  helper.NewPlot(20, 60000),
			Sprite:   assets.LoadSprite(),
			Colorful: false,
			Amount:   1000,
		},
	)
	return world
}

func (g *Game) Update() error {
	g.ecs.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	g.ecs.DrawLayer(layers.LayerBackground, screen)
	g.ecs.DrawLayer(layers.LayerBunnies, screen)
	g.ecs.DrawLayer(layers.LayerMetrics, screen)
}

func (g *Game) Layout(width, height int) (int, int) {
	g.bounds = image.Rect(0, 0, width, height)
	return width, height
}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowSizeLimits(300, 200, -1, -1)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
	ebiten.SetWindowResizable(true)
	rand.Seed(time.Now().UTC().UnixNano())
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
