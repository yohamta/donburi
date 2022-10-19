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
	"github.com/yohamta/donburi/examples/bunnymark_ecs/system"

	_ "net/http/pprof"
)

type Game struct {
	ecs    *ecs.ECS
	bounds image.Rectangle
}

const (
	LayerBackground ecs.LayerID = iota
	LayerBunnies
	LayerMetrics
)

func NewGame() *Game {
	g := &Game{
		bounds: image.Rectangle{},
		ecs:    createECS(),
	}

	metrics := system.NewMetrics(&g.bounds)

	g.ecs.AddSystems(
		// Systems are executed in the order they are added.
		ecs.System{
			Update: system.NewSpawn().Update,
		},
		ecs.System{
			Layer: LayerBackground,
			Draw:  system.DrawBackground,
		},
		ecs.System{
			Layer:  LayerMetrics,
			Update: metrics.Update,
			Draw:   metrics.Draw,
		},
		ecs.System{
			Update: system.NewBounce(&g.bounds).Update,
		},
		ecs.System{
			Update: system.Velocity.Update,
		},
		ecs.System{
			Update: system.Gravity.Update,
		},
		ecs.System{
			Layer: LayerBunnies,
			Draw:  system.Render.Draw,
		},
	)

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
	g.ecs.Draw(LayerBackground, screen)
	g.ecs.Draw(LayerBunnies, screen)
	g.ecs.Draw(LayerMetrics, screen)
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
