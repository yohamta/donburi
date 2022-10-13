package ecs

import (
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

	*innerECS
}

type DrawerOpts struct {
	ImageToDraw *ebiten.Image
}

type updaterEntry struct {
	Updater Updater
}

type drawerEntry struct {
	Drawer  Drawer
	Options *DrawerOpts
}

type innerECS struct {
	updaters []updaterEntry
	drawers  []drawerEntry
}

// NewECS creates a new ECS with the specified world.
func NewECS(w donburi.World) *ECS {
	return &ECS{
		World: w,
		innerECS: &innerECS{
			updaters: []updaterEntry{},
			drawers:  []drawerEntry{},
		},
	}
}

// AddUpdater adds an Updater to the ECS.
func (ecs *ECS) AddUpdater(u Updater) {
	ecs.updaters = append(ecs.updaters, updaterEntry{Updater: u})
}

// AddDrawer adds a Drawer to the ECS.
func (ecs *ECS) AddDrawer(d Drawer, opts *DrawerOpts) {
	ecs.drawers = append(ecs.drawers, drawerEntry{Drawer: d, Options: opts})
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
