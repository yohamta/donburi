package system

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/component"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type Render struct {
	query *query.Query
}

func NewRender() *Render {
	return &Render{
		query: query.NewQuery(filter.Contains(
			component.Position,
			component.Hue,
			component.Sprite,
		))}
}

func (r *Render) Update(ecs *ecs.ECS) {}

func (r *Render) Draw(ecs *ecs.ECS, screen *ebiten.Image) {
	r.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		position := component.GetPosition(entry)
		hue := component.GetHue(entry)
		sprite := component.GetSprite(entry)

		op := &ebiten.DrawImageOptions{}
		sw, sh := float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())
		op.GeoM.Translate(position.X*sw, position.Y*sh)
		if *hue.Colorful {
			op.ColorM.RotateHue(hue.Value)
		}
		screen.DrawImage(sprite.Image, op)
	})
}
