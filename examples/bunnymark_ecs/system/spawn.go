package system

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/component"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/helper"
	"github.com/yohamta/donburi/examples/bunnymark_ecs/layers"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type Spawn struct {
	settings *component.SettingsData
}

func NewSpawn() *Spawn {
	return &Spawn{
		settings: nil,
	}
}

func (s *Spawn) Update(ecs *ecs.ECS) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s.addBunnies(ecs)
	}

	if ids := ebiten.AppendTouchIDs(nil); len(ids) > 0 {
		s.addBunnies(ecs) // not accurate, cause no input manager for this
	}

	if _, offset := ebiten.Wheel(); offset != 0 {
		s.settings.Amount += int(offset * 10)
		if s.settings.Amount < 0 {
			s.settings.Amount = 0
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		s.settings.Colorful = !s.settings.Colorful
	}
}

func (s *Spawn) addBunnies(ecs *ecs.ECS) {
	if s.settings == nil {
		query := query.NewQuery(filter.Contains(component.Settings))
		query.EachEntity(ecs.World, func(entry *donburi.Entry) {
			s.settings = component.GetSettings(entry)
		})
	}

	entities := ecs.CreateMany(
		layers.LayerBunnies,
		s.settings.Amount,
		component.Position,
		component.Velocity,
		component.Hue,
		component.Gravity,
		component.Sprite,
	)
	for i, entity := range entities {
		entry := ecs.World.Entry(entity)
		position := component.GetPosition(entry)
		*position = component.PositionData{
			X: float64(i % 2), // Alternate screen edges
		}
		donburi.SetValue(
			entry, component.Velocity, component.VelocityData{
				X: helper.RangeFloat(0, 0.005),
				Y: helper.RangeFloat(0.0025, 0.005),
			})
		donburi.SetValue(
			entry, component.Hue, component.HueData{
				Colorful: &s.settings.Colorful,
				Value:    helper.RangeFloat(0, 2*math.Pi),
			})
		donburi.SetValue(entry, component.Gravity,
			component.GravityData{Value: 0.00095})
		donburi.SetValue(entry, component.Sprite,
			component.SpriteData{Image: s.settings.Sprite})
	}
}
