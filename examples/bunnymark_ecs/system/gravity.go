package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/component"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type Gravity struct {
	query *query.Query
}

func NewGravity() *Gravity {
	return &Gravity{
		query: query.NewQuery(filter.Contains(component.Velocity, component.Gravity)),
	}
}

func (g *Gravity) Update(ecs *ecs.ECS) {
	g.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		gravity := component.GetGravity(entry)
		velocity := component.GetVelocity(entry)

		velocity.Y += gravity.Value
	})
}
