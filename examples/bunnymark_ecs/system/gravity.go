package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/component"
	"github.com/yohamta/donburi/filter"
)

type gravity struct {
	query *donburi.Query
}

var Gravity *gravity = &gravity{
	query: donburi.NewQuery(
		filter.Contains(
			component.Velocity,
			component.Gravity,
		)),
}

func (g *gravity) Update(ecs *ecs.ECS) {
	for entry := range g.query.Iter(ecs.World) {
		gravity := component.Gravity.Get(entry)
		velocity := component.Velocity.Get(entry)

		velocity.Y += gravity.Value
	}
}
