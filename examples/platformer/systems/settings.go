package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/examples/platformer/components"
)

func UpdateSettings(ecs *ecs.ECS) {
	settings := GetOrCreateSettings(ecs)

	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		settings.Debug = !settings.Debug
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		settings.ShowHelpText = !settings.ShowHelpText
	}
}

func GetOrCreateSettings(ecs *ecs.ECS) *components.SettingsData {
	if _, ok := components.Settings.First(ecs.World); !ok {
		ent := ecs.World.Entry(ecs.World.Create(components.Settings))
		components.Settings.SetValue(ent, components.SettingsData{
			ShowHelpText: true,
		})
	}

	ent, _ := components.Settings.First(ecs.World)
	return components.Settings.Get(ent)
}
