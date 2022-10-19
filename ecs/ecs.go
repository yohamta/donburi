package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

// LayerID is used to specify a layer.
type LayerID int

// ECS represents an entity-component-system.
type ECS struct {
	// World is the underlying world of the ECS.
	World donburi.World
	// Time manages the time of the world.
	Time *Time
	// UpdateCount is the number of times Update is called.
	UpdateCount int64

	layers         []*Layer
	systems        []UpdateSystem
	startupSystems []UpdateSystem
}

// NewQuery creates a new query.
func NewQuery(l LayerID, f filter.LayoutFilter) *query.Query {
	layerFilter := filter.Contains(getLayer(l).tag)
	if f == nil {
		return query.NewQuery(layerFilter)
	}
	return query.NewQuery(filter.And(layerFilter, f))
}

// NewECS creates a new ECS with the specified world.
func NewECS(w donburi.World) *ECS {
	ecs := &ECS{
		World: w,
		Time:  NewTime(),

		systems:        []UpdateSystem{},
		layers:         []*Layer{},
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
		ecs.addDrawSystem(s.Layer, s.Draw)
	}
	return ecs
}

// Update runs systems
func (ecs *ECS) Update() {
	ecs.Time.Update()
	for _, sys := range ecs.systems {
		sys(ecs)
	}
}

// Draw calls draw
func (ecs *ECS) Draw(l LayerID, screen *ebiten.Image) {
	ecs.getLayer(l).draw(ecs, screen)
}

// Create creates a new entity
func (ecs *ECS) Create(l LayerID, components ...*donburi.ComponentType) *donburi.Entry {
	entry := ecs.World.Entry(ecs.World.Create(components...))
	entry.AddComponent(ecs.getLayer(l).tag)
	return entry
}

func (ecs *ECS) getLayer(layerID LayerID) *Layer {
	if int(layerID) >= len(ecs.layers) {
		ecs.layers = append(ecs.layers, make([]*Layer, int(layerID)-len(ecs.layers)+1)...)
	}
	if ecs.layers[layerID] == nil {
		ecs.layers[layerID] = newLayer(getLayer(layerID))
	}
	return ecs.layers[layerID]
}

func (ecs *ECS) addUpdateSystem(systems ...UpdateSystem) *ECS {
	for _, s := range systems {
		ecs.systems = append(ecs.systems, s)
	}
	return ecs
}

func (ecs *ECS) addDrawSystem(l LayerID, s DrawSystem) *ECS {
	ecs.getLayer(l).addDrawSystem(s)
	return ecs
}
