package ecs

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/query"
)

// ECS represents an entity-component-system.
type ECS struct {
	// World is the underlying world of the ECS.
	World donburi.World
	// Time manages the time of the world.
	Time *Time
	// ScriptSystem manages the scripts of the world.
	ScriptSystem *ScriptSystem

	systems []*system
}

type system struct {
	System  System
	Options *SystemOpts
}

// NewECS creates a new ECS with the specified world.
func NewECS(w donburi.World) *ECS {
	ecs := &ECS{
		World:   w,
		Time:    NewTime(),
		systems: []*system{},
	}

	ecs.ScriptSystem = NewScriptSystem()
	ecs.AddSystem(ecs.ScriptSystem, &SystemOpts{})

	return ecs
}

// AddSystems adds new systems either Updater or Drawer
func (ecs *ECS) AddSystems(systems ...System) {
	for _, system := range systems {
		ecs.AddSystem(system, nil)
	}
}

// AddSystem adds new system
func (ecs *ECS) AddSystem(sys System, opts *SystemOpts) {
	if opts == nil {
		opts = &SystemOpts{}
	}
	ecs.addSystem(&system{
		System:  sys,
		Options: opts,
	})
}

// AddScript adds a script to the entities matched by the query.
func (ecs *ECS) AddScript(q *query.Query, script Script, opts *ScriptOpts) {
	ecs.ScriptSystem.AddScript(q, script, opts)
}

// Update calls Updater's Update() methods.
func (ecs *ECS) Update() {
	ecs.Time.Update()
	for _, u := range ecs.systems {
		u.System.Update(ecs)
	}
}

// Draw calls Drawer's Draw() methods.
func (ecs *ECS) Draw(screen *ebiten.Image) {
	for _, d := range ecs.systems {
		if d.Options.Image != nil {
			d.System.Draw(ecs, d.Options.Image)
			continue
		}
		d.System.Draw(ecs, screen)
	}
}

func (ecs *ECS) addSystem(sys *system) {
	ecs.systems = append(ecs.systems, sys)
	sortSystems(ecs.systems)
}

func sortSystems(systems []*system) {
	sort.SliceStable(systems, func(i, j int) bool {
		return systems[i].Options.Priority > systems[j].Options.Priority
	})
}
