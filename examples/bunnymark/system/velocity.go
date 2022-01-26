package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/examples/bunnymark/component"
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

func (v *Velocity) Update(w donburi.World) {
	v.query.EachEntity(w, func(entry *donburi.Entry) {
		position := component.GetPosition(entry)
		velocity := component.GetVelocity(entry)

		position.X += velocity.X
		position.Y += velocity.Y
	})
}
