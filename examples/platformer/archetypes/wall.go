package archetypes

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/examples/platformer/components"
	"github.com/yohamta/donburi/examples/platformer/layers"
	"github.com/yohamta/donburi/examples/platformer/tags"
)

func NewWall(ecs *ecs.ECS) *donburi.Entry {
	w := ecs.World

	entry := w.Entry(ecs.Create(
		layers.Default,
		tags.Wall,
		components.Object,
	))

	return entry
}
