package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

// LayerID is used to specify a layer.
type LayerID int

// LayerDefault is the default layer.
const LayerDefault LayerID = 0

// ECS represents an entity-component-system.
type ECS struct {
	// World is the underlying world of the ECS.
	World donburi.World
	// Time manages the time of the world.
	Time *Time

	layers  []*Layer
	systems []System
}

// NewQuery creates a new query.
func NewQuery(l LayerID, f filter.LayoutFilter) *donburi.Query {
	layerFilter := filter.Contains(getLayer(l).tag)
	if f == nil {
		return donburi.NewQuery(layerFilter)
	}
	return donburi.NewQuery(filter.And(layerFilter, f))
}

// NewECS creates a new ECS with the specified world.
func NewECS(w donburi.World) *ECS {
	ecs := &ECS{
		World: w,
		Time:  NewTime(),

		systems: []System{},
		layers:  []*Layer{},
	}

	return ecs
}

// AddSystem adds a system.
func (ecs *ECS) AddSystem(s System) *ECS {
	ecs.systems = append(ecs.systems, s)
	return ecs
}

func (ecs *ECS) AddRenderer(l LayerID, r Renderer) *ECS {
	ecs.getLayer(l).addRenderer(r)
	return ecs
}

// Update runs systems
func (ecs *ECS) Update() {
	ecs.Time.Update()
	for _, sys := range ecs.systems {
		sys(ecs)
	}
}

// DrawLayer executes all draw systems of the specified layer.
func (ecs *ECS) DrawLayer(l LayerID, screen *ebiten.Image) {
	ecs.getLayer(l).draw(ecs, screen)
}

// Draw executes all draw systems.
func (ecs *ECS) Draw(screen *ebiten.Image) {
	for _, l := range ecs.layers {
		if l == nil {
			continue
		}
		l.draw(ecs, screen)
	}
}

// Create creates a new entity
func (ecs *ECS) Create(l LayerID, components ...donburi.IComponentType) donburi.Entity {
	entry := ecs.World.Entry(ecs.World.Create(components...))
	entry.AddComponent(ecs.getLayer(l).tag)
	return entry.Entity()
}

// Create creates a new entity
func (ecs *ECS) CreateMany(l LayerID, n int, components ...donburi.IComponentType) []donburi.Entity {
	comps := append(components, ecs.getLayer(l).tag)
	return ecs.World.CreateMany(n, comps...)
}

// Pause pauses the world.
func (ecs *ECS) Pause() {
	ecs.Time.Pause()
}

// Resume resumes the world.
func (ecs *ECS) Resume() {
	ecs.Time.Resume()
}

// IsPaused returns a boolean value indicating whether the world is paused.
func (ecs *ECS) IsPaused() bool {
	return ecs.Time.isPaused
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

func (ecs *ECS) addRenderer(l LayerID, r Renderer) *ECS {
	ecs.getLayer(l).addRenderer(r)
	return ecs
}
