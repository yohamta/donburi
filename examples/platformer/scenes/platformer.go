package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type PlatformerScene struct {
	ecs *ecs.ECS
}

func NewPlatformerScene() *PlatformerScene {
	return &PlatformerScene{
		ecs: createECS(),
	}
}

func (ps *PlatformerScene) Update() {
	ps.ecs.Update()
}

func (ps *PlatformerScene) Draw(screen *ebiten.Image) {
}

func createECS() *ecs.ECS {
	world := donburi.NewWorld()
	ecs := ecs.NewECS(world)
	return ecs
}
