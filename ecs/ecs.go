package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type Layer int

// ECS represents an entity-component-system.
type ECS struct {
	// World is the underlying world of the ECS.
	World donburi.World
	// Time manages the time of the world.
	Time *Time

	layers []*layer
}

// NewECS creates a new ECS with the specified world.
func NewECS(w donburi.World) *ECS {
	ecs := &ECS{
		World: w,
		Time:  NewTime(),

		layers: []*layer{},
	}

	return ecs
}

// AddSystem adds new system
func (ecs *ECS) AddSystem(layer Layer, sys System, opts *SystemOpts) {
	if opts == nil {
		opts = &SystemOpts{}
	}
	ecs.getLayer(layer).addSystem(&system{
		System:  sys,
		Options: opts,
	})
}

// AddScript adds a script to the entities matched by the query.
func (ecs *ECS) AddScript(layer Layer, script *Script) {
	ecs.getLayer(layer).addScript(script)
}

// Update calls Updater's Update() methods.
func (ecs *ECS) Update() {
	ecs.Time.Update()
	for _, l := range ecs.layers {
		l.Update(ecs)
	}
}

// Draw calls Drawer's Draw() methods.
func (ecs *ECS) Draw(screen *ebiten.Image) {
	for _, l := range ecs.layers {
		l.Draw(ecs, screen)
	}
}

func (ecs *ECS) getLayer(l Layer) *layer {
	if int(l) >= len(ecs.layers) {
		// expand layers slice
		ecs.layers = append(ecs.layers, make([]*layer, int(l)-len(ecs.layers)+1)...)
	}
	if ecs.layers[l] == nil {
		ecs.layers[l] = newLayer()
	}
	return ecs.layers[l]
}
