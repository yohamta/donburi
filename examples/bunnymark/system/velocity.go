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
		var position *component.PositionData = (*component.PositionData)(entry.Component(component.Position))
		var velocity *component.VelocityData = (*component.VelocityData)(entry.Component(component.Velocity))

		position.X += velocity.X
		position.Y += velocity.Y
	})
}
