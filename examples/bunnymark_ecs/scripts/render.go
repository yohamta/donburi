package scripts

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/component"
)

type render struct{}

var Render = &render{}

func (r *render) Draw(entry *donburi.Entry, screen *ebiten.Image) {
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
}
