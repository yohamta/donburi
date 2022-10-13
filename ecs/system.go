package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Updater is a system that updates the world.
type Updater interface {
	Update(ecs *ECS)
}

// Drawer is a system that draws the world.
type Drawer interface {
	Draw(ecs *ECS, screen *ebiten.Image)
}
