package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// System is a system that updates the world.
type System interface {
	Update(ecs *ECS)
	Draw(ecs *ECS, screen *ebiten.Image)
}

// SystemOpts represents options for systems.
type SystemOpts struct {
	// Image is the image to draw the system.
	Image *ebiten.Image
}
