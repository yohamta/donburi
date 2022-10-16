package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// UpdateSystem is a system that updates the world.
type UpdateSystem interface {
	Update(ecs *ECS)
}

// DrawSystem is a system that draws the world.
type DrawSystem interface {
	Draw(ecs *ECS, screen *ebiten.Image)
}
