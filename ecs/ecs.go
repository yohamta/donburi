package ecs

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

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

type DrawerOpts struct {
	ImageToDraw *ebiten.Image
}

// NewECS creates a new ECS with the specified world.
func NewECS(w donburi.World) *ECS {
	ecs := &ECS{
		World: w,
		innerECS: &innerECS{
			updaters: []updaterEntry{},
			drawers:  []drawerEntry{},
		},
	}

	ecs.ScriptSystem = NewScriptSystem()
	ecs.AddUpdater(ecs.ScriptSystem)
	ecs.AddDrawer(ecs.ScriptSystem, nil)

	return ecs
}

// AddUpdater adds an Updater to the ECS.
func (ecs *ECS) AddUpdater(u Updater) {
	ecs.updaters = append(ecs.updaters, updaterEntry{Updater: u})
}

// AddDrawer adds a Drawer to the ECS.
func (ecs *ECS) AddDrawer(d Drawer, opts *DrawerOpts) {
	entry := drawerEntry{Drawer: d, Options: opts}
	if entry.Options == nil {
		entry.Options = &DrawerOpts{}
	}
	ecs.drawers = append(ecs.drawers, entry)
}

// AddUpdaterWithPriority adds an Updater to the ECS with the specified priority.
// Higher priority is executed first.
func (ecs *ECS) AddUpdaterWithPriority(u Updater, priority int) {
	ecs.updaters = append(ecs.updaters, updaterEntry{Updater: u, Priority: priority})
	sortUpdaterEntries(ecs.updaters)
}

// AddDrawerWithPriority adds a Drawer to the ECS with the specified priority.
func (ecs *ECS) AddDrawerWithPriority(d Drawer, priority int, opts *DrawerOpts) {
	ecs.drawers = append(ecs.drawers, drawerEntry{Drawer: d, Priority: priority, Options: opts})
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
		if d.Options.ImageToDraw != nil {
			d.Drawer.Draw(ecs, d.Options.ImageToDraw)
			continue
		}
		d.Drawer.Draw(ecs, screen)
	}
}

func sortUpdaterEntries(entries []updaterEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Priority > entries[j].Priority
	})
}

func sortDrawerEntries(entries []drawerEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Priority > entries[j].Priority
	})
}

type updaterEntry struct {
	Updater  Updater
	Priority int
}

type drawerEntry struct {
	Drawer   Drawer
	Priority int
	Options  *DrawerOpts
}

type innerECS struct {
	updaters []updaterEntry
	drawers  []drawerEntry
}
