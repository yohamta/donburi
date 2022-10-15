package scripts

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/component"
)

type Velocity struct{}

func (v *Velocity) Update(entry *donburi.Entry) {
	position := component.GetPosition(entry)
	velocity := component.GetVelocity(entry)

	position.X += velocity.X
	position.Y += velocity.Y
}
