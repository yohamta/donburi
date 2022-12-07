package factory

import (
	"github.com/solarlune/resolv"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/examples/platformer/archetypes"
	"github.com/yohamta/donburi/examples/platformer/components"
	"github.com/yohamta/donburi/examples/platformer/config"
	dresolv "github.com/yohamta/donburi/examples/platformer/resolv"
)

func CreatePlatform(ecs *ecs.ECS, object *resolv.Object) *donburi.Entry {
	platform := archetypes.NewPlatform(ecs)
	dresolv.SetObject(platform, object)

	return platform
}

func CreateFloatingPlatform(ecs *ecs.ECS, object *resolv.Object) *donburi.Entry {
	platform := archetypes.NewFloatingPlatform(ecs)

	gh := float64(config.C.Height)
	dresolv.SetObject(platform, resolv.NewObject(128, gh-32, 128, 8))
	obj := components.Object.Get(platform)

	// The floating platform moves using a *gween.Sequence sequence of tweens, moving it back and forth.
	tw := gween.NewSequence()
	tw.Add(
		gween.New(float32(obj.Y), float32(obj.Y-128), 2, ease.Linear),
		gween.New(float32(obj.Y-128), float32(obj.Y), 2, ease.Linear),
	)
	components.Tween.Set(platform, tw)

	return platform
}
