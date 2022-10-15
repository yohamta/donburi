package main

import (
	"image"
	"log"
	"math/rand"
	"time"
	"unsafe"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/assets"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/component"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/helper"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/scripts"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/system"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	_ "net/http/pprof"
)

type System interface {
	Update(w donburi.World)
}

type Drawable interface {
	Draw(w donburi.World, screen *ebiten.Image)
}

type Game struct {
	ecs    *ecs.ECS
	bounds image.Rectangle
}

const (
	LayerBackground ecs.Layer = iota
	LayerBunnies
	LayerMetrics
)

func NewGame() *Game {
	g := &Game{
		bounds: image.Rectangle{},
		ecs:    createECS(),
	}

	metrics := system.NewMetrics(&g.bounds)

	g.ecs.AddUpdateSystems(
		system.NewSpawn(),
		metrics,
	)
	g.ecs.AddDrawSystem(LayerBackground, &system.Background{})
	g.ecs.AddDrawSystem(LayerMetrics, metrics)

	g.ecs.AddUpdateScript(scripts.NewBounce(&g.bounds),
		query.NewQuery(filter.Contains(
			component.Position,
			component.Velocity,
			component.Sprite,
		)))

	g.ecs.AddUpdateScript(&scripts.Velocity{},
		query.NewQuery(filter.Contains(component.Position, component.Velocity)))

	g.ecs.AddUpdateScript(scripts.Gravity,
		query.NewQuery(filter.Contains(component.Velocity, component.Gravity)))

	g.ecs.AddDrawScript(LayerBunnies, scripts.Render,
		query.NewQuery(filter.Contains(
			component.Position,
			component.Hue,
			component.Sprite,
		)))

	return g
}

func createECS() *ecs.ECS {
	world := createWorld()
	ecs := ecs.NewECS(world)
	return ecs
}

func createWorld() donburi.World {
	world := donburi.NewWorld()
	setting := world.Create(component.Settings)
	world.Entry(setting).SetComponent(component.Settings,
		unsafe.Pointer(&component.SettingsData{
			Ticker:   time.NewTicker(500 * time.Millisecond),
			Gpu:      helper.GpuInfo(),
			Tps:      helper.NewPlot(20, 60),
			Fps:      helper.NewPlot(20, 60),
			Objects:  helper.NewPlot(20, 60000),
			Sprite:   assets.LoadSprite(),
			Colorful: false,
			Amount:   1000,
		}))
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
