package system

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/component"
	"github.com/yohamta/donburi/filter"
)

type render struct {
	query        *donburi.Query
	orderedQuery *donburi.OrderedQuery[component.PositionData]
}

var Render = &render{
	query: donburi.NewQuery(
		filter.Contains(
			component.Position,
			component.Hue,
			component.Sprite,
		)),
	orderedQuery: donburi.NewOrderedQuery[component.PositionData](
		filter.Contains(
			component.Position,
			component.Hue,
			component.Sprite,
		)),
}

func (r *render) Draw(ecs *ecs.ECS, screen *ebiten.Image) {
	if !UsePositionOrdering {
		for entry := range r.query.Iter(ecs.World) {
			position := component.Position.Get(entry)
			hue := component.Hue.Get(entry)
			sprite := component.Sprite.Get(entry)

			op := &ebiten.DrawImageOptions{}
			sw, sh := float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())
			op.GeoM.Translate(position.X*sw, position.Y*sh)
			if *hue.Colorful {
				op.ColorM.RotateHue(hue.Value)
			}
			screen.DrawImage(sprite.Image, op)
		}
	} else {
		for entry := range r.orderedQuery.IterOrdered(ecs.World, component.Position) {
			position := component.Position.Get(entry)
			hue := component.Hue.Get(entry)
			sprite := component.Sprite.Get(entry)

			op := &ebiten.DrawImageOptions{}
			sw, sh := float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())
			op.GeoM.Translate(position.X*sw, position.Y*sh)
			if *hue.Colorful {
				op.ColorM.RotateHue(hue.Value)
			}
			screen.DrawImage(sprite.Image, op)
		}
	}
}
