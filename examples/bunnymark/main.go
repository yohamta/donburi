package main

import (
	"image"
	"log"
	"math/rand"
	"time"
	"unsafe"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/example/bunnymark/assets"
	"github.com/yohamta/donburi/example/bunnymark/component"
	"github.com/yohamta/donburi/example/bunnymark/helper"
	"github.com/yohamta/donburi/example/bunnymark/system"

	_ "net/http/pprof"
)

type System interface {
	Update(w donburi.World)
}

type Drawable interface {
	Draw(w donburi.World, screen *ebiten.Image)
}

type Game struct {
	bounds    image.Rectangle
	world     donburi.World
	systems   []System
	drawables []Drawable
}

func NewGame() *Game {
	g := &Game{
		bounds: image.Rectangle{},
	}
	g.createWorld()
	return g
}

func (g *Game) createWorld() {
	g.world = donburi.NewWorld()
	setting := g.world.Create(component.Settings)
	g.world.Entry(setting).SetComponent(component.Settings,
		unsafe.Pointer(&component.SettingsData{
			Ticker:   time.NewTicker(500 * time.Millisecond),
			Gpu:      helper.GpuInfo(),
			Tps:      helper.NewPlot(20, 60),
			Fps:      helper.NewPlot(20, 60),
			Objects:  helper.NewPlot(20, 60000),
			Sprite:   assets.LoadSprite(),
			Colorful: false,
			Amount:   100,
		}))

	metrics := system.NewMetrics(&g.bounds)
	g.systems = []System{
		system.NewSpawn(),
		system.NewGravity(),
		system.NewVelocity(),
		system.NewBounce(&g.bounds),
		metrics,
	}
	g.drawables = []Drawable{
		&system.Background{},
		system.NewRender(),
		metrics,
	}
}

func (g *Game) Update() error {
	for _, s := range g.systems {
		s.Update(g.world)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	for _, s := range g.drawables {
		s.Draw(g.world, screen)
	}
}

func (g *Game) Layout(width, height int) (int, int) {
	g.bounds = image.Rect(0, 0, width, height)
	return width, height
}

func main() {
	ebiten.SetWindowSize(300, 200)
	ebiten.SetWindowSizeLimits(300, 200, -1, -1)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
	ebiten.SetWindowResizable(true)
	rand.Seed(time.Now().UTC().UnixNano())
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
