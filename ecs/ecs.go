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

	systems      []UpdateSystem
	layers       []*layer
	scriptSystem *scriptSystem
}

// NewECS creates a new ECS with the specified world.
func NewECS(w donburi.World) *ECS {
	ecs := &ECS{
		World: w,
		Time:  NewTime(),

		systems:      []UpdateSystem{},
		layers:       []*layer{},
		scriptSystem: newScriptSystem(),
	}

	return ecs
}

// AddUpdateSystems adds new update systems
func (ecs *ECS) AddUpdateSystems(systems ...UpdateSystem) {
	ecs.systems = append(ecs.systems, systems...)
}

// AddUpdateSystem adds new update system
func (ecs *ECS) AddUpdateSystem(s UpdateSystem) {
	ecs.addSystem(s)
}

// AddDrawSystem adds new draw system
func (ecs *ECS) AddDrawSystem(l Layer, s DrawSystem) {
	ecs.getLayer(l).addDrawSystem(s)
}

// AddUpdateScript adds a script to the entities matched by the query.
func (ecs *ECS) AddUpdateScript(s UpdateScript, q *query.Query) {
	ecs.scriptSystem.AddUpdateScript(s, q)
}

// AddDrawScript adds a script to the entities matched by the query.
func (ecs *ECS) AddDrawScript(l Layer, s DrawScript, q *query.Query) {
	ecs.getLayer(l).scriptSystem.AddDrawScript(s, q)
}

// Update calls Updater's Update() methods.
func (ecs *ECS) Update() {
	ecs.Time.Update()
	for _, sys := range ecs.systems {
		sys(ecs)
	}
	ecs.scriptSystem.Update(ecs)
}

// Draw calls Drawer's Draw() methods.
func (ecs *ECS) Draw(l Layer, screen *ebiten.Image) {
	ecs.getLayer(l).Draw(ecs, screen)
}

func (ecs *ECS) addSystem(s UpdateSystem) {
	ecs.systems = append(ecs.systems, s)
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
