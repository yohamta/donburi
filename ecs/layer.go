package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type system struct {
	System DrawSystem
}

type layer struct {
	systems      []*system
	scriptSystem *scriptSystem
	image        *ebiten.Image
}

func newLayer() *layer {
	return &layer{
		systems:      []*system{},
		scriptSystem: newScriptSystem(),
	}
}

func (l *layer) Update(e *ECS) {
	l.scriptSystem.Update(e)
}

func (l *layer) Draw(e *ECS, i *ebiten.Image) {
	screen := i
	if l.image != nil {
		screen = l.image
	}
	for _, d := range l.systems {
		d.System.Draw(e, screen)
	}
	l.scriptSystem.Draw(e, screen)
}

func (l *layer) addSystem(s *system) {
	l.systems = append(l.systems, s)
}
