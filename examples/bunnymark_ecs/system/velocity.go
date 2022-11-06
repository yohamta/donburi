package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/component"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type velocity struct {
	query *query.Query
}

var Velocity = &velocity{
	query: query.NewQuery(
		filter.Contains(
			component.Position,
			component.Velocity,
		)),
}

func (v *velocity) Update(ecs *ecs.ECS) {
	v.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		position := component.Position.Get(entry)
		velocity := component.Velocity.Get(entry)

		position.X += velocity.X
		position.Y += velocity.Y
	})
}
