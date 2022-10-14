package scripts

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/component"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type velocity struct{}

var Velocity = ecs.NewScript(
	query.NewQuery(filter.Contains(component.Position, component.Velocity)),
	&velocity{},
	nil,
)

func (v *velocity) Update(entry *donburi.Entry) {
	position := component.GetPosition(entry)
	velocity := component.GetVelocity(entry)

	position.X += velocity.X
	position.Y += velocity.Y
}

func (v *velocity) Draw(entry *donburi.Entry, image *ebiten.Image) {}
