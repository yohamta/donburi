package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/component"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type Velocity struct {
	query *query.Query
}

func NewVelocity() *Velocity {
	return &Velocity{
		query: query.NewQuery(filter.Contains(component.Position, component.Velocity)),
	}
}

func (v *Velocity) Update(ecs *ecs.ECS) {
	v.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		position := component.GetPosition(entry)
		velocity := component.GetVelocity(entry)

		position.X += velocity.X
		position.Y += velocity.Y
	})
}
