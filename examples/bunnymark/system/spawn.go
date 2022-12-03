package system

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/examples/bunnymark/component"
	"github.com/yohamta/donburi/examples/bunnymark/helper"
)

type Spawn struct {
	settings *component.SettingsData
}

func NewSpawn() *Spawn {
	return &Spawn{
		settings: nil,
	}
}

func (s *Spawn) Update(w donburi.World) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s.addBunnies(w)
	}

	if ids := ebiten.AppendTouchIDs(nil); len(ids) > 0 {
		s.addBunnies(w) // not accurate, cause no input manager for this
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

func (s *Spawn) addBunnies(w donburi.World) {
	if s.settings == nil {
		if entry, ok := component.Settings.FirstEntity(w); ok {
			s.settings = component.Settings.Get(entry)
		} else {
			panic("no settings")
		}
	}

	entities := w.CreateMany(
		s.settings.Amount,
		component.Position,
		component.Velocity,
		component.Hue,
		component.Gravity,
		component.Sprite,
	)
	for i, entity := range entities {
		entry := w.Entry(entity)
		position := component.Position.Get(entry)
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
