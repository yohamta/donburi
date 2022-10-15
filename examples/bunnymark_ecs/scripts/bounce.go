package scripts

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/component"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/helper"
)

type bounce struct {
	bounds *image.Rectangle
}

func NewBounce(bounds *image.Rectangle) *bounce {
	return &bounce{bounds}
}

func (b *bounce) Update(entry *donburi.Entry) {
	position := component.GetPosition(entry)
	velocity := component.GetVelocity(entry)
	sprite := component.GetSprite(entry)

	sw, sh := float64(b.bounds.Dx()), float64(b.bounds.Dy())
	iw, ih := float64(sprite.Image.Bounds().Dx()), float64(sprite.Image.Bounds().Dy())
	relW, relH := iw/sw, ih/sh
	if position.X+relW > 1 {
		velocity.X *= -1
		position.X = 1 - relW
	}
	if position.X < 0 {
		velocity.X *= -1
		position.X = 0
	}
	if position.Y+relH > 1 {
		velocity.Y *= -0.85
		position.Y = 1 - relH
		if helper.Chance(0.5) {
			velocity.Y -= helper.RangeFloat(0, 0.009)
		}
	}
	if position.Y < 0 {
		velocity.Y = 0
		position.Y = 0
	}
}

func (b *bounce) Draw(entry *donburi.Entry, screen *ebiten.Image) {}
