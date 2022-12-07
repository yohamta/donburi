package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/examples/platformer/config"
	"github.com/yohamta/donburi/examples/platformer/factory"
	"github.com/yohamta/donburi/examples/platformer/layers"
	dresolv "github.com/yohamta/donburi/examples/platformer/resolv"
	"github.com/yohamta/donburi/examples/platformer/systems"
)

type PlatformerScene struct {
	ecs *ecs.ECS
}

func NewPlatformerScene() *PlatformerScene {
	ecs := ecs.NewECS(donburi.NewWorld())

	ecs.AddSystem(systems.UpdateFloatingPlatform)
	ecs.AddSystem(systems.UpdatePlayer)
	ecs.AddSystem(systems.UpdateObjects)
	ecs.AddSystem(systems.UpdateSettings)

	ecs.AddRenderer(layers.Default, systems.DrawWall)
	ecs.AddRenderer(layers.Default, systems.DrawPlatform)
	ecs.AddRenderer(layers.Default, systems.DrawRamp)
	ecs.AddRenderer(layers.Default, systems.DrawFloatingPlatform)
	ecs.AddRenderer(layers.Default, systems.DrawPlayer)
	ecs.AddRenderer(layers.Default, systems.DrawDebug)
	ecs.AddRenderer(layers.Default, systems.DrawHelp)

	ps := &PlatformerScene{ecs: ecs}
	ps.init()

	return ps
}

func (ps *PlatformerScene) Update() {
	ps.ecs.Update()
}

func (ps *PlatformerScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 40, 255})
	ps.ecs.Draw(screen)
}

func (ps *PlatformerScene) init() {
	gw, gh := float64(config.C.Width), float64(config.C.Height)

	// Define the world's Space. Here, a Space is essentially a grid (the game's width and height, or 640x360), made up of 16x16 cells. Each cell can have 0 or more Objects within it,
	// and collisions can be found by checking the Space to see if the Cells at specific positions contain (or would contain) Objects. This is a broad, simplified approach to collision
	// detection.
	space := factory.CreateSpace(ps.ecs)

	dresolv.Add(space,
		// Construct the solid level geometry. Note that the simple approach of checking cells in a Space for collision works simply when the geometry is aligned with the cells,
		// as it all is in this platformer example.
		factory.CreateWall(ps.ecs, resolv.NewObject(0, 0, 16, gh, "solid")),
		factory.CreateWall(ps.ecs, resolv.NewObject(gw-16, 0, 16, gh, "solid")),
		factory.CreateWall(ps.ecs, resolv.NewObject(0, 0, gw, 16, "solid")),
		factory.CreateWall(ps.ecs, resolv.NewObject(0, gh-24, gw, 32, "solid")),
		factory.CreateWall(ps.ecs, resolv.NewObject(160, gh-56, 160, 32, "solid")),
		factory.CreateWall(ps.ecs, resolv.NewObject(320, 64, 32, 160, "solid")),
		factory.CreateWall(ps.ecs, resolv.NewObject(64, 128, 16, 160, "solid")),
		factory.CreateWall(ps.ecs, resolv.NewObject(gw-128, 64, 128, 16, "solid")),
		factory.CreateWall(ps.ecs, resolv.NewObject(gw-128, gh-88, 128, 16, "solid")),
		// Create the Player. NewPlayer adds it to the world's Space.
		factory.CreatePlayer(ps.ecs),
		// Non-moving floating Platforms.
		factory.CreatePlatform(ps.ecs, resolv.NewObject(352, 64, 48, 8, "platform")),
		factory.CreatePlatform(ps.ecs, resolv.NewObject(352, 64+64, 48, 8, "platform")),
		factory.CreatePlatform(ps.ecs, resolv.NewObject(352, 64+128, 48, 8, "platform")),
		factory.CreatePlatform(ps.ecs, resolv.NewObject(352, 64+192, 48, 8, "platform")),
		// Create the floating platform.
		factory.CreateFloatingPlatform(ps.ecs, resolv.NewObject(128, gh-32, 128, 8, "platform")),
		// A ramp, which is unique as it has a non-rectangular shape. For this, we will specify a different shape for collision testing.
		factory.CreateRamp(ps.ecs, resolv.NewObject(320, gh-56, 64, 32, "ramp")),
	)

}
