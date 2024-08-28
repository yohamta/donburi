package systems

import (
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/examples/platformer/components"
	dresolv "github.com/yohamta/donburi/examples/platformer/resolv"
)

func UpdateObjects(ecs *ecs.ECS) {
	for e := range components.Object.Iter(ecs.World) {
		obj := dresolv.GetObject(e)
		obj.Update()
	}
}
