package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type layer struct {
	systems      []DrawSystem
	scriptSystem *scriptSystem
	image        *ebiten.Image
}

func newLayer() *layer {
	return &layer{
		systems:      []DrawSystem{},
		scriptSystem: newScriptSystem(),
	}
}

func (l *layer) Draw(e *ECS, i *ebiten.Image) {
	screen := i
	if l.image != nil {
		screen = l.image
	}
	for _, s := range l.systems {
		s(e, screen)
	}
	l.scriptSystem.Draw(e, screen)
}

func (l *layer) addDrawSystem(s DrawSystem) {
	l.systems = append(l.systems, s)
}
