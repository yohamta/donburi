package systems

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	dresolv "github.com/yohamta/donburi/examples/platformer/resolv"
	"github.com/yohamta/donburi/examples/platformer/tags"
)

func DrawWall(ecs *ecs.ECS, screen *ebiten.Image) {
	tags.Wall.EachEntity(ecs.World, func(e *donburi.Entry) {
		o := dresolv.GetObject(e)
		drawColor := color.RGBA{60, 60, 60, 255}
		ebitenutil.DrawRect(screen, o.X, o.Y, o.W, o.H, drawColor)
	})
}
