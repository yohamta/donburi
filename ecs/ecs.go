package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type ECS struct {
	// World is the underlying world of the ECS.
	World donburi.World

	// Updaters is a list of systems that update the world.
	Updaters []Updater

	// Drawers is a list of systems that draw the world.
	Drawers []Drawer
}

// NewECS creates a new ECS with the specified world.
func NewECS(w donburi.World) *ECS {
	return &ECS{
		World: w,
	}
}

// Update calls Updater's Update() methods.
func (ecs *ECS) Update() {
	for _, u := range ecs.Updaters {
		u.Update(ecs)
	}
}

// Draw calls Drawer's Draw() methods.
func (ecs *ECS) Draw(screen *ebiten.Image) {
	for _, d := range ecs.Drawers {
		d.Draw(ecs, screen)
	}
}
