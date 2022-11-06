package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/examples/bunnymark/component"
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
		gravity := component.Gravity.Get(entry)
		velocity := component.Velocity.Get(entry)

		velocity.Y += gravity.Value
	})
}
