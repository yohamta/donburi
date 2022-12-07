package systems

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/examples/platformer/components"
	dresolv "github.com/yohamta/donburi/examples/platformer/resolv"
)

func UpdateObjects(ecs *ecs.ECS) {
	components.Object.EachEntity(ecs.World, func(e *donburi.Entry) {
		obj := dresolv.GetObject(e)
		obj.Update()
	})
}
