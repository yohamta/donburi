package scripts

import (
	"github.com/hajimehoshi/ebiten/v2"
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

func (g *gravity) Draw(entry *donburi.Entry, screen *ebiten.Image) {}
