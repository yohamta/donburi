package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/example/bunnymark/component"
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

func (g *Gravity) Update(w donburi.World) {
	g.query.EachEntity(w, func(entry *donburi.Entry) {
		var gravity *component.GravityData = (*component.GravityData)(entry.Component(component.Gravity))
		var velocity *component.VelocityData = (*component.VelocityData)(entry.Component(component.Velocity))

		velocity.Y += gravity.Value
	})
}
