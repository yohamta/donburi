package systems

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/examples/platformer/fonts"
	"golang.org/x/image/font"
)

func DrawHelp(ecs *ecs.ECS, screen *ebiten.Image) {
	drawText(screen, 16, 16,
		"~ Platformer Demo ~",
		"Move Player: Left, Right Arrow",
		"Jump: X Key",
		"Wallslide: Move into wall in air",
		"Walljump: Jump while wallsliding",
		"Fall through platforms: Down + X",
		"",
		"F1: Toggle Debug View",
		"F2: Show / Hide help text",
		"F4: Toggle fullscreen",
		"R: Restart world",
		"E: Next world",
		"Q: Previous world",
		fmt.Sprintf("%d FPS (frames per second)", int(ebiten.CurrentFPS())),
		fmt.Sprintf("%d TPS (ticks per second)", int(ebiten.CurrentTPS())),
	)
}

func drawText(screen *ebiten.Image, x, y int, textLines ...string) {
	f := fonts.Excel.Get()
	rectHeight := 10
	for _, txt := range textLines {
		w := float64(font.MeasureString(f, txt).Round())
		ebitenutil.DrawRect(screen, float64(x), float64(y-8), w, float64(rectHeight), color.RGBA{0, 0, 0, 192})

		text.Draw(screen, txt, f, x+1, y+1, color.RGBA{0, 0, 150, 255})
		text.Draw(screen, txt, f, x, y, color.RGBA{100, 150, 255, 255})
		y += rectHeight
	}
}
