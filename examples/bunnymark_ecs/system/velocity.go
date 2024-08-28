package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/component"
	"github.com/yohamta/donburi/filter"
)

type velocity struct {
	query *donburi.Query
}

var Velocity = &velocity{
	query: donburi.NewQuery(
		filter.Contains(
			component.Position,
			component.Velocity,
		)),
}

func (v *velocity) Update(ecs *ecs.ECS) {
	for entry := range v.query.Iter(ecs.World) {
		position := component.Position.Get(entry)
		velocity := component.Velocity.Get(entry)

		position.X += velocity.X
		position.Y += velocity.Y
	}
}
