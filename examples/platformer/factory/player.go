package factory

import (
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/examples/platformer/archetypes"
	"github.com/yohamta/donburi/examples/platformer/components"
	dresolv "github.com/yohamta/donburi/examples/platformer/resolv"
)

func CreatePlayer(ecs *ecs.ECS) *donburi.Entry {
	player := archetypes.NewPlayer(ecs)

	obj := resolv.NewObject(32, 128, 16, 24)
	dresolv.SetObject(player, obj)
	components.Player.SetValue(player, components.PlayerData{
		FacingRight: true,
	})

	obj.SetShape(resolv.NewRectangle(0, 0, 16, 24))

	return player
}
