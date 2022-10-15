package scripts

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/component"
)

type gravity struct{}

var Gravity = &gravity{}

func (g *gravity) Update(entry *donburi.Entry) {
	gravity := component.GetGravity(entry)
	velocity := component.GetVelocity(entry)

	velocity.Y += gravity.Value
}
