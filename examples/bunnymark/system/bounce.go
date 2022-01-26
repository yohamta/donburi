package system

import (
	"image"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/examples/bunnymark/component"
	"github.com/yohamta/donburi/examples/bunnymark/helper"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type Bounce struct {
	bounds *image.Rectangle
	query  *query.Query
}

func NewBounce(bounds *image.Rectangle) *Bounce {
	return &Bounce{
		bounds: bounds,
		query: query.NewQuery(filter.Contains(
			component.Position,
			component.Velocity,
			component.Sprite,
		))}
}

func (b *Bounce) Update(w donburi.World) {
	b.query.EachEntity(w, func(entry *donburi.Entry) {
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
	})
}
