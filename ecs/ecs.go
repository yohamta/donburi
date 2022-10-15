package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/query"
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
func (ecs *ECS) AddSystem(l Layer, s System) {
	ecs.getLayer(l).addSystem(&system{System: s})
}

// AddScript adds a script to the entities matched by the query.
func (ecs *ECS) AddScript(l Layer, s Script, q *query.Query) {
	ecs.addScript(l, newScript(s, q))

}

// ConfigLayer sets the layer configuration.
func (ecs *ECS) ConfigLayer(l Layer, cfg *LayerConfig) {
	ecs.getLayer(l).Config(cfg)
}

// Update calls Updater's Update() methods.
func (ecs *ECS) Update() {
	ecs.Time.Update()
	for _, l := range ecs.layers {
		l.Update(ecs)
	}
}

// Draw calls Drawer's Draw() methods.
func (ecs *ECS) Draw(l Layer, screen *ebiten.Image) {
	ecs.getLayer(l).Draw(ecs, screen)
}

// AddScript adds a script to the entities matched by the query.
func (ecs *ECS) addScript(l Layer, s *script) {
	ecs.getLayer(l).addScript(s)
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
