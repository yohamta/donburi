package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// UpdateSystem is a system that updates the world.
type UpdateSystem func(ecs *ECS)

// DrawSystem is a system that draws the world.
type DrawSystem func(ecs *ECS, screen *ebiten.Image)

// System represents a system.
type System struct {
	Update    UpdateSystem
	DrawLayer LayerID
	Draw      DrawSystem
}
