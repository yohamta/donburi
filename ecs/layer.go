package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// LayerConfig represents options for layers.
type LayerConfig struct {
	Image *ebiten.Image
}

type system struct {
	System System
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
	for _, u := range l.systems {
		u.System.Update(e)
	}
	l.scriptSystem.Update(e)
}

// Draw calls Drawer's Draw() methods.
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

func (l *layer) Config(c *LayerConfig) {
	l.image = c.Image
}

func (l *layer) addScript(s *script) {
	l.scriptSystem.AddScript(s)
}

func (l *layer) addSystem(s *system) {
	l.systems = append(l.systems, s)
}
