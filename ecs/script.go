package ecs

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/query"
)

// EntryUpdater is a function that updates an entry.
type EntryUpdater interface {
	Update(entry *donburi.Entry)
}

// EntryDrawer is a function that draws an entry.
type EntryDrawer interface {
	Draw(entry *donburi.Entry, screen *ebiten.Image)
}

// ScriptSystem is a built-in system that manages scripts with queries.
type ScriptSystem struct {
	updaters []*entryUpdater
	drawers  []*entryDrawer
}

// ScriptOpts represents options for a script.
type ScriptOpts struct {
	// Image is the image that the script draws to.
	Image *ebiten.Image
	// Priority is the priority of the script. higher priority is executed first.
	Priority int
}

// NewScriptSystem creates a new ScriptSystem.
func NewScriptSystem() *ScriptSystem {
	return &ScriptSystem{
		updaters: []*entryUpdater{},
		drawers:  []*entryDrawer{},
	}
}

type entryUpdater struct {
	Query   *query.Query
	Updater EntryUpdater
	Options *ScriptOpts
}

type entryDrawer struct {
	Query   *query.Query
	Drawer  EntryDrawer
	Options *ScriptOpts
}

// AddScript adds a script to the system.
// Target entities are specified by the query.
func (ss *ScriptSystem) AddScript(q *query.Query, script interface{}, opts *ScriptOpts) {
	if opts == nil {
		opts = &ScriptOpts{}
	}
	flag := false
	if script, ok := script.(EntryUpdater); ok {
		ss.addUpdater(&entryUpdater{
			Query:   q,
			Updater: script,
			Options: opts,
		})
		flag = true
	}
	if script, ok := script.(EntryDrawer); ok {
		ss.addDrawer(&entryDrawer{
			Query:   q,
			Drawer:  script,
			Options: opts,
		})
		flag = true
	}
	if !flag {
		panic("script must be EntryUpdater or EntryDrawer at least")
	}
}

func (ss *ScriptSystem) Update(ecs *ECS) {
	for _, script := range ss.updaters {
		script.Query.EachEntity(ecs.World, script.Updater.Update)
	}
}

func (ss *ScriptSystem) Draw(ecs *ECS, screen *ebiten.Image) {
	for _, script := range ss.drawers {
		script.Query.EachEntity(ecs.World, func(entry *donburi.Entry) {
			if script.Options.Image != nil {
				script.Drawer.Draw(entry, script.Options.Image)
				return
			}
			script.Drawer.Draw(entry, screen)
		})
	}
}

func (ss *ScriptSystem) addUpdater(entry *entryUpdater) {
	ss.updaters = append(ss.updaters, entry)
	sortEntryUpdater(ss.updaters)
}

func (ss *ScriptSystem) addDrawer(entry *entryDrawer) {
	ss.drawers = append(ss.drawers, entry)
	sortEntryDrawer(ss.drawers)
}

func sortEntryUpdater(items []*entryUpdater) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Options.Priority > items[j].Options.Priority
	})
}

func sortEntryDrawer(items []*entryDrawer) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Options.Priority > items[j].Options.Priority
	})
}

var (
	_ = (Updater)((*ScriptSystem)(nil))
	_ = (Drawer)((*ScriptSystem)(nil))
)
