package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/query"
)

type DrawLayer int

// ECS represents an entity-component-system.
type ECS struct {
	// World is the underlying world of the ECS.
	World donburi.World
	// Time manages the time of the world.
	Time *Time
	// UpdateCount is the number of times Update is called.
	UpdateCount int64

	systems        []UpdateSystem
	layers         []*layer
	scriptSystem   *scriptSystem
	startupSystems []UpdateSystem
}

// NewECS creates a new ECS with the specified world.
func NewECS(w donburi.World) *ECS {
	ecs := &ECS{
		World: w,
		Time:  NewTime(),

		systems:        []UpdateSystem{},
		layers:         []*layer{},
		scriptSystem:   newScriptSystem(),
		startupSystems: []UpdateSystem{},
	}

	return ecs
}

// AddSystems adds systems.
func (ecs *ECS) AddSystems(s ...System) *ECS {
	for _, ss := range s {
		ecs.AddSystem(ss)
	}
	return ecs
}

// AddSystem adds a system.
func (ecs *ECS) AddSystem(s System) *ECS {
	if s.Update != nil {
		ecs.addUpdateSystem(s.Update)
	}
	if s.Draw != nil {
		ecs.addDrawSystem(s.DrawLayer, s.Draw)
	}
	return ecs
}

// AddScripts adds scripts.
func (ecs *ECS) AddScripts(scripts ...Script) *ECS {
	for _, s := range scripts {
		ecs.AddScript(s)
	}
	return ecs
}

// AddScript adds a script to the entities matched by the query.
func (ecs *ECS) AddScript(s Script) *ECS {
	if s.Query == nil {
		panic("query must not be nil")
	}
	if s.Update != nil {
		ecs.scriptSystem.AddUpdateScript(s.Update, s.Query)
	}
	if s.Draw != nil {
		ecs.getLayer(s.DrawLayer).scriptSystem.AddDrawScript(s.Draw, s.Query)
	}
	return ecs
}

// AddUpdateSystem adds new update system
func (ecs *ECS) addUpdateSystem(systems ...UpdateSystem) *ECS {
	for _, s := range systems {
		ecs.systems = append(ecs.systems, s)
	}
	return ecs
}

// AddDrawSystem adds new draw system
func (ecs *ECS) addDrawSystem(l DrawLayer, s DrawSystem) *ECS {
	ecs.getLayer(l).addDrawSystem(s)
	return ecs
}

func (ecs *ECS) addUpdateScript(s UpdateScript, q *query.Query) *ECS {
	ecs.scriptSystem.AddUpdateScript(s, q)
	return ecs
}

func (ecs *ECS) addDrawScript(l DrawLayer, s DrawScript, q *query.Query) *ECS {
	ecs.getLayer(l).scriptSystem.AddDrawScript(s, q)
	return ecs
}

// Update runs systems
func (ecs *ECS) Update() {
	ecs.Time.Update()
	for _, sys := range ecs.systems {
		sys(ecs)
	}
	ecs.scriptSystem.Update(ecs)
}

// Draw calls draw
func (ecs *ECS) Draw(l DrawLayer, screen *ebiten.Image) {
	ecs.getLayer(l).Draw(ecs, screen)
}

func (ecs *ECS) getLayer(l DrawLayer) *layer {
	if int(l) >= len(ecs.layers) {
		// expand layers slice
		ecs.layers = append(ecs.layers, make([]*layer, int(l)-len(ecs.layers)+1)...)
	}
	if ecs.layers[l] == nil {
		ecs.layers[l] = newLayer()
	}
	return ecs.layers[l]
}
