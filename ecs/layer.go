package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type system struct {
	System  System
	Options *SystemOpts
}

type layer struct {
	systems      []*system
	scriptSystem *ScriptSystem
}

func newLayer() *layer {
	return &layer{
		systems:      []*system{},
		scriptSystem: NewScriptSystem(),
	}
}

func (l *layer) Update(ecs *ECS) {
	for _, u := range l.systems {
		u.System.Update(ecs)
	}
	l.scriptSystem.Update(ecs)
}

// Draw calls Drawer's Draw() methods.
func (l *layer) Draw(ecs *ECS, screen *ebiten.Image) {
	for _, d := range l.systems {
		if d.Options.Image != nil {
			d.System.Draw(ecs, d.Options.Image)
			continue
		}
		d.System.Draw(ecs, screen)
	}
	l.scriptSystem.Draw(ecs, screen)
}

func (l *layer) addScript(script *Script) {
	l.scriptSystem.AddScript(script)
}

func (l *layer) addSystem(sys *system) {
	l.systems = append(l.systems, sys)
}
