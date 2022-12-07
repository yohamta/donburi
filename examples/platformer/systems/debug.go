package systems

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/examples/platformer/components"
)

func DrawDebug(ecs *ecs.ECS, screen *ebiten.Image) {
	settings := GetOrCreateSettings(ecs)
	if !settings.Debug {
		return
	}
	spaceEntry, ok := components.Space.FirstEntity(ecs.World)
	if !ok {
		return
	}
	space := components.Space.Get(spaceEntry)

	for y := 0; y < space.Height(); y++ {

		for x := 0; x < space.Width(); x++ {

			cell := space.Cell(x, y)

			cw := float64(space.CellWidth)
			ch := float64(space.CellHeight)
			cx := float64(cell.X) * cw
			cy := float64(cell.Y) * ch

			drawColor := color.RGBA{20, 20, 20, 255}

			if cell.Occupied() {
				drawColor = color.RGBA{255, 255, 0, 255}
			}

			ebitenutil.DrawLine(screen, cx, cy, cx+cw, cy, drawColor)

			ebitenutil.DrawLine(screen, cx+cw, cy, cx+cw, cy+ch, drawColor)

			ebitenutil.DrawLine(screen, cx+cw, cy+ch, cx, cy+ch, drawColor)

			ebitenutil.DrawLine(screen, cx, cy+ch, cx, cy, drawColor)
		}

	}
}
