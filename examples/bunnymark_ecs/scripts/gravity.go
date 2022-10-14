package scripts

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/component"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type gravity struct{}

var Gravity = ecs.NewScript(
	query.NewQuery(filter.Contains(component.Velocity, component.Gravity)),
	&gravity{},
	nil,
)

func (g *gravity) Update(entry *donburi.Entry) {
	gravity := component.GetGravity(entry)
	velocity := component.GetVelocity(entry)

	velocity.Y += gravity.Value
}

func (g *gravity) Draw(entry *donburi.Entry, screen *ebiten.Image) {}
