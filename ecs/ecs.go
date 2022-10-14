package ecs

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

// ECS represents an entity-component-system.
type ECS struct {
	// World is the underlying world of the ECS.
	World donburi.World

	// UpdateCount is the number of updates
	UpdateCount int64

	// Time manages the time of the world.
	Time *Time

	// ScriptSystem manages the scripts of the world.
	ScriptSystem *ScriptSystem

	*innerECS
}

type SystemOpts struct {
	ImageToDraw *ebiten.Image
	Priority    int
}

type updaterEntry struct {
	Updater  Updater
	Priority int
}

type drawerEntry struct {
	Drawer      Drawer
	Priority    int
	ImageToDraw *ebiten.Image
}

type innerECS struct {
	updaters []*updaterEntry
	drawers  []*drawerEntry
}

// NewECS creates a new ECS with the specified world.
func NewECS(w donburi.World) *ECS {
	ecs := &ECS{
		World: w,
		Time:  NewTime(),
		innerECS: &innerECS{
			updaters: []*updaterEntry{},
			drawers:  []*drawerEntry{},
		},
	}

	ecs.ScriptSystem = NewScriptSystem()
	ecs.AddSystem(ecs.ScriptSystem, &SystemOpts{})

	return ecs
}

// AddSystem adds new system either Updater or Drawer
func (ecs *ECS) AddSystem(system interface{}, opts *SystemOpts) {
	if opts == nil {
		opts = &SystemOpts{}
	}
	flag := false
	if system, ok := system.(Updater); ok {
		ecs.addUpdater(&updaterEntry{
			Updater:  system,
			Priority: opts.Priority,
		})
		flag = true
	}
	if system, ok := system.(Drawer); ok {
		ecs.addDrawer(&drawerEntry{
			Drawer:      system,
			Priority:    opts.Priority,
			ImageToDraw: opts.ImageToDraw,
		})
		flag = true
	}
	if !flag {
		panic("ECS system should be either Updater or Drawer at least.")
	}
}

// AddUpdaterWithPriority adds an Updater to the ECS with the specified priority.
// Higher priority is executed first.
func (ecs *ECS) addUpdater(entry *updaterEntry) {
	ecs.updaters = append(ecs.updaters, entry)
	sortUpdaterEntries(ecs.updaters)
}

// AddDrawer adds a Drawer to the ECS.
func (ecs *ECS) addDrawer(entry *drawerEntry) {
	ecs.drawers = append(ecs.drawers, entry)
	sortDrawerEntries(ecs.drawers)
}

// Update calls Updater's Update() methods.
func (ecs *ECS) Update() {
	ecs.UpdateCount++
	ecs.Time.Update()
	for _, u := range ecs.updaters {
		u.Updater.Update(ecs)
	}
}

// Draw calls Drawer's Draw() methods.
func (ecs *ECS) Draw(screen *ebiten.Image) {
	for _, d := range ecs.drawers {
		if d.ImageToDraw != nil {
			d.Drawer.Draw(ecs, d.ImageToDraw)
			continue
		}
		d.Drawer.Draw(ecs, screen)
	}
}

func sortUpdaterEntries(entries []*updaterEntry) {
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Priority > entries[j].Priority
	})
}

func sortDrawerEntries(entries []*drawerEntry) {
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Priority > entries[j].Priority
	})
}
