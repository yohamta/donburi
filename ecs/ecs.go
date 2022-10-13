package ecs

import (
	"time"

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

	// DeltaTime is the time between the last update and the current update
	DeltaTime time.Duration

	// Sleep is the time to sleep
	Sleep time.Duration

	// Time manages the time of the world.
	Time *Time
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
