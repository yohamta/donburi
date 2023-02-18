package archetypes

import (
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/examples/platformer/components"
	"github.com/yohamta/donburi/examples/platformer/tags"
)

var (
	Platform = ecs.NewArchetype(
		tags.Platform,
		components.Object,
	)
	FloatingPlatform = ecs.NewArchetype(
		tags.FloatingPlatform,
		components.Object,
		components.Tween,
	)
	Player = ecs.NewArchetype(
		tags.Player,
		components.Player,
		components.Object,
	)
	Ramp = ecs.NewArchetype(
		tags.Ramp,
		components.Object,
	)
	Space = ecs.NewArchetype(
		components.Space,
	)
	Wall = ecs.NewArchetype(
		tags.Wall,
		components.Object,
	)
)
