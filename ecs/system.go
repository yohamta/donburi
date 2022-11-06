package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// UpdateSystem is a system that updates the world.
type System func(ecs *ECS)

// DrawSystem is a system that draws the world.
type Renderer func(ecs *ECS, screen *ebiten.Image)
