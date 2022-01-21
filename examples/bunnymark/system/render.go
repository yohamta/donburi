package system

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/examples/bunnymark/component"
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

func (r *Render) Draw(w donburi.World, screen *ebiten.Image) {
	r.query.EachEntity(w, func(entry *donburi.Entry) {
		var position *component.PositionData = (*component.PositionData)(entry.Component(component.Position))
		var hue *component.HueData = (*component.HueData)(entry.Component(component.Hue))
		var sprite *component.SpriteData = (*component.SpriteData)(entry.Component(component.Sprite))

		op := &ebiten.DrawImageOptions{}
		sw, sh := float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())
		op.GeoM.Translate(position.X*sw, position.Y*sh)
		if *hue.Colorful {
			op.ColorM.RotateHue(hue.Value)
		}
		screen.DrawImage(sprite.Image, op)
	})
}
