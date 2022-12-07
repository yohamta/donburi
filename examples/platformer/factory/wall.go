package factory

import (
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/examples/platformer/archetypes"
	dresolv "github.com/yohamta/donburi/examples/platformer/resolv"
)

func CreateWall(ecs *ecs.ECS, obj *resolv.Object) *donburi.Entry {
	wall := archetypes.NewWall(ecs)
	dresolv.SetObject(wall, obj)
	return wall
}
