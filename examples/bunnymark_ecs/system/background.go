package system

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi/ecs"
)

type Background struct{}

func (b *Background) Draw(ecs *ecs.ECS, screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 41, G: 44, B: 45, A: 255})
}
