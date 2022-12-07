package systems

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/examples/platformer/components"
	dresolv "github.com/yohamta/donburi/examples/platformer/resolv"
	"github.com/yohamta/donburi/examples/platformer/tags"
)

func DrawPlayer(ecs *ecs.ECS, screen *ebiten.Image) {
	tags.Player.EachEntity(ecs.World, func(e *donburi.Entry) {
		player := components.Player.Get(e)
		o := dresolv.GetObject(e)
		playerColor := color.RGBA{0, 255, 60, 255}
		if player.OnGround == nil {
			// We draw the player as a different color when jumping so we can visually see when he's in the air.
			playerColor = color.RGBA{200, 0, 200, 255}
		}
		ebitenutil.DrawRect(screen, o.X, o.Y, o.W, o.H, playerColor)
	})
}
