package systems

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi/ecs"
	dresolv "github.com/yohamta/donburi/examples/platformer/resolv"
	"github.com/yohamta/donburi/examples/platformer/tags"
)

func DrawRamp(ecs *ecs.ECS, screen *ebiten.Image) {
	for e := range tags.Ramp.Iter(ecs.World) {
		o := dresolv.GetObject(e)
		drawColor := color.RGBA{255, 50, 100, 255}
		tri := o.Shape.(*resolv.ConvexPolygon)
		drawPolygon(screen, tri, drawColor)
	}
}

func drawPolygon(screen *ebiten.Image, polygon *resolv.ConvexPolygon, color color.Color) {
	for _, line := range polygon.Lines() {
		ebitenutil.DrawLine(screen, line.Start.X(), line.Start.Y(), line.End.X(), line.End.Y(), color)
	}
}
