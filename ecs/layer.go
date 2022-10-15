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
	scriptSystem *scriptSystem
}

func newLayer() *layer {
	return &layer{
		systems:      []*system{},
		scriptSystem: newScriptSystem(),
	}
}

func (l *layer) Update(e *ECS) {
	for _, u := range l.systems {
		u.System.Update(e)
	}
	l.scriptSystem.Update(e)
}

// Draw calls Drawer's Draw() methods.
func (l *layer) Draw(e *ECS, s *ebiten.Image) {
	for _, d := range l.systems {
		if d.Options.Image != nil {
			d.System.Draw(e, d.Options.Image)
			continue
		}
		d.System.Draw(e, s)
	}
	l.scriptSystem.Draw(e, s)
}

func (l *layer) addScript(s *script) {
	l.scriptSystem.AddScript(s)
}

func (l *layer) addSystem(s *system) {
	l.systems = append(l.systems, s)
}
